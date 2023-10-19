package app

import (
	_ "app/routes"
	"fmt"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	_ "github.com/lib/pq"
)

func registerDb() error {
	var dbConf, _ = web.AppConfig.GetSection("database")
	cfg := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbConf["host"],
		dbConf["port"],
		dbConf["user"],
		dbConf["password"],
		dbConf["name"],
		dbConf["ssl"])
	return orm.RegisterDataBase("default", c("database::driver"), cfg)
}

func main() {
	if err := registerDb(); err != nil {
		logs.Error("%s", err)
		return
	}

	if web.BConfig.RunMode == web.DEV {
		logs.Warn("Run mode DEV")
		orm.Debug = true
		web.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	var runConf, _ = web.AppConfig.GetSection("run")
	web.Run(fmt.Sprintf("%v:%v", runConf["host"], runConf["port"]))
}
