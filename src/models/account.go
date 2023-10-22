package models

import (
	"app/src/components/auth"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Account struct {
	model
	Id            int64     `orm:"column(id);auto;pk"`
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

func (m *Account) Login(password string) (string, error) {
	err := Get(m, "Username")
	if err != nil {
		return "", auth.ErrorUsernameOrPasswordIncorrect
	}

	if err = auth.CheckPasswordHash(password, m.Password); err != nil {
		return "", auth.ErrorUsernameOrPasswordIncorrect
	}

	token, err := auth.CreateAccessToken(m.Id)
	if err != nil {
		return "", err
	}
	m.IsNeedRelogin = false
	_, err = Update(m, "IsNeedRelogin")
	if err != nil {
		return "", err
	}
	return token, nil
}

func (m *Account) Register(username, password string) error {
	password, err := auth.HashPassword(password)
	if err != nil {
		return err
	}
	m.Username = username
	m.Password = password
	_, err = Insert(m)
	return err
}
