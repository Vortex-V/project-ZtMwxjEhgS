package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Transport struct {
	model
	Id          int64           `orm:"column(id);pk"`
	AccountId   *Account        `orm:"column(user_id);rel(fk)"`
	CanBeRented bool            `orm:"column(can_be_rented)"`
	TypeId      *TransportTypes `orm:"column(type_id);rel(fk)"`
	Model       string          `orm:"column(model)"`
	Color       string          `orm:"column(color)"`
	Identifier  string          `orm:"column(identifier)"`
	Description string          `orm:"column(description);null"`
	Latitude    float64         `orm:"column(latitude)"`
	Longitude   float64         `orm:"column(longitude)"`
	MinutePrice float64         `orm:"column(minute_price);null"`
	DayPrice    float64         `orm:"column(day_price);null"`
	CreatedAt   time.Time       `orm:"column(created_at);type(timestamp without time zone);auto_now_add"`
	UpdatedAt   time.Time       `orm:"column(updated_at);type(timestamp without time zone);auto_now_add"`
}

func (t *Transport) TableName() string {
	return "transports"
}

func init() {
	orm.RegisterModel(new(Transport))
}
