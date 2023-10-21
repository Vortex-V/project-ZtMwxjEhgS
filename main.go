package main

import (
	_ "app/routes"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	_ "github.com/lib/pq"
)

func main() {
	if web.BConfig.RunMode == web.DEV {
		logs.Warn("Run mode DEV")
		orm.Debug = true
		web.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	var runConf, _ = web.AppConfig.GetSection("run")
	web.Run(fmt.Sprintf("%v:%v", runConf["host"], runConf["port"]))
}
