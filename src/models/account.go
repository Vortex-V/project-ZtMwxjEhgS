package models

import (
	"app/src/components/auth"
	"app/src/components/requests"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Account struct {
	model
	Id            int       `orm:"column(id);auto;pk"`
	Username      string    `orm:"column(username);unique"`
	Password      string    `orm:"column(password)"`
	Type          int       `orm:"column(type);" default:"1"`
	Status        int       `orm:"column(status)" default:"1"`
	IsNeedRelogin bool      `orm:"column(is_need_relogin);" default:"false"`
	Balance       float64   `orm:"column(balance);default(0)"`
	CreatedAt     time.Time `orm:"column(created_at);type(timestamp without time zone);null;auto_now_add"`
	UpdatedAt     time.Time `orm:"column(updated_at);type(timestamp without time zone);null;auto_now"`
}

func (m *Account) TableName() string {
	return "accounts"
}

func init() {
	orm.RegisterModel(new(Account))
}

func (m *Account) Register(data requests.AccountRequest) error {
	password, err := auth.HashPassword(data.Password)
	if err != nil {
		return err
	}
	m.Username = data.Username
	m.Password = password
	_, err = Insert(m)
	return err
}
