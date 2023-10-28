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
	Type          int       `orm:"column(type);" default:"2"`
	Status        int       `orm:"column(status)" default:"1"`
	IsNeedRelogin bool      `orm:"column(is_need_relogin);" default:"false"`
	Balance       float64   `orm:"column(balance);default(0)"`
	CreatedAt     time.Time `orm:"column(created_at);type(timestamp without time zone);null;auto_now_add"`
	UpdatedAt     time.Time `orm:"column(updated_at);type(timestamp without time zone);null;auto_now"`

	Transports []*Transport `orm:"reverse(many)"`
	Rents      []*Rent      `orm:"reverse(many)"`
}

func (m *Account) TableName() string {
	return "accounts"
}

func init() {
	orm.RegisterModel(new(Account))
}

const (
	AccountTypeAdmin = iota + 1
	AccountTypeUser
)

var typeLabels = map[int]string{
	AccountTypeAdmin: "admin",
	AccountTypeUser:  "user",
}

func (m *Account) GetTypeLabel() string {
	return typeLabels[m.Type]
}

func GetAccountTypeLabelByKey(key int) string {
	return typeLabels[key]
}

func (m *Account) GetTypeByLabel(label string) int {
	for k, v := range typeLabels {
		if v == label {
			return k
		}
	}
	return AccountTypeUser
}

func (m *Account) IsAdmin() bool {
	return m.Type == AccountTypeAdmin
}

const (
	AccountStatusActive = iota + 1
	AccountStatusDeleted
)

func (m *Account) GetStatusLabel(key int) string {
	types := map[int]string{
		AccountStatusActive:  "Active",
		AccountStatusDeleted: "Deleted",
	}
	return types[key]
}

func (m *Account) Login(password string) (string, error) {
	err := Read(m, "Username")
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

	// Костыль. Старый jwt остаётся валидным в течение 1 часа.
	// Но jwt в базе хранить это кукож, товарищи, я так делать не буду.
	// Хотелось бы дополнительно использовать refresh token
	m.IsNeedRelogin = false
	_, err = Update(m, "IsNeedRelogin")
	if err != nil {
		return "", err
	}
	return token, nil
}

// TODO сделать нормальный Create и перенести код из контроллеров
func (m *Account) Register(username, password string) error {
	password, err := auth.HashPassword(password)
	if err != nil {
		return err
	}
	m.Username = username
	m.Password = password
	m.Type = AccountTypeUser
	m.Status = AccountStatusActive
	_, err = Insert(m)
	return err
}
