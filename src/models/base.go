package models

import (
	"database/sql"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"reflect"
)

var (
	o  orm.Ormer
	qb orm.QueryBuilder
)

func init() {

	// TODO перенести
	var dbConf, _ = web.AppConfig.GetSection("database")
	cfg := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbConf["host"],
		dbConf["port"],
		dbConf["user"],
		dbConf["password"],
		dbConf["name"],
		dbConf["ssl"])
	if err := orm.RegisterDataBase("default", dbConf["driver"], cfg); err != nil {
		logs.Error("%s", err)
		return
	}

	o = orm.NewOrm()
	driver, _ := web.AppConfig.String("database::driver")
	qb, _ = orm.NewQueryBuilder(driver)
}

type Model interface {
	implement()
}

type model struct {
	Model
}

func Insert(model Model) (int64, error) {
	return o.Insert(model)
}

func Delete(model Model) error {
	_, err := o.Delete(model)
	return err
}

func Update(model Model) (int64, error) {
	return o.Update(model)
}

func Exists(model Model) (bool, error) {
	var (
		result sql.Result
		count  int64
		err    error
	)
	id := reflect.ValueOf(model).Elem().FieldByName("Id").Int()
	query := Find(model).Where("id = ?")
	result, err = SetArgs(query, id).Exec()
	count, err = result.RowsAffected()
	return count > 0, err
}

func Get(model Model) (err error) {
	id := reflect.ValueOf(model).Elem().FieldByName("Id").Int()
	if id == 0 {
		return fmt.Errorf("invalid id: %v", id)
	}
	return o.Read(model)
}

func SetArgs(build orm.QueryBuilder, args ...interface{}) orm.RawSeter {
	return o.Raw(build.String(), args...)
}

func Find(model Model, selectFields ...string) orm.QueryBuilder {
	if len(selectFields) == 0 {
		selectFields = append(selectFields, "*")
	}
	tableName := reflect.ValueOf(model).
		MethodByName("TableName").
		Call([]reflect.Value{})[0].
		String()
	return qb.Select(selectFields...).From(tableName)
}
