package models

import (
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

// Init
// Должен выполняться после регистрации всех моделей
func init() {
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

// Model
// Содаём тип интерфейса и структуру, указатель на которую его реализует.
//
// Можно было бы сделать
//
/*type Model interface {
	*Account | и т.д.
}*/
// но мне удобней указывать имплементацию непосредсвтенно в типе структуры
type Model interface {
	implement(_ Model)
}
type model struct{}

func (m *model) implement(_ Model) {}

// Функции Read, Insert, Update, Delete, Exists
// позволяют однозначно определить,
// что необходимо передавать указатель на Model в качестве параметра.
// Базовые функции orm.Ormer не защищают от этого и работают некорректно,
// если передавать не указатель.

func Read(model Model, by ...string) (err error) {
	return o.Read(model, by...)
}

func Insert(model Model) (int64, error) {
	return o.Insert(model)
}

func Update(model Model, cols ...string) (int64, error) {
	return o.Update(model, cols...)
}

func Delete(model Model) (int64, error) {
	num, err := o.Delete(model)
	return num, err
}

func Raw(build orm.QueryBuilder, args ...interface{}) orm.RawSeter {
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
