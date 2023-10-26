package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Transport struct {
	model
	Id          int64     `orm:"auto;pk"`
	Account     *Account  `orm:"column(account_id);rel(one)"`
	CanBeRented bool      `orm:"" default:"false"`
	Type        string    `orm:""`
	Model       string    `orm:""`
	Color       string    `orm:""`
	Identifier  string    `orm:""`
	Description string    `orm:"null"`
	Latitude    float64   `orm:""`
	Longitude   float64   `orm:""`
	MinutePrice float64   `orm:"null"`
	DayPrice    float64   `orm:"null"`
	CreatedAt   time.Time `orm:"type(timestamp without time zone);auto_now_add"`
	UpdatedAt   time.Time `orm:"type(timestamp without time zone);auto_now"`
}

func (t *Transport) TableName() string {
	return "transports"
}

func init() {
	orm.RegisterModel(new(Transport))
}

func (t *Transport) IsOwner(id int64) bool {
	return t.Account.Id == id
}

func GetTransportType(label string) string {
	types := map[string]string{
		"Car":     "Car",
		"Bike":    "Bike",
		"Scooter": "Scooter",
	}
	if v, ok := types[label]; ok {
		return v
	} else {
		return ""
	}
}

func TransportSearch(params map[string]string, offset int, limit int) (int64, []*Transport, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Transport))

	for k, v := range params {
		if v != "" {
			qs = qs.Filter(k, v)
		}
	}

	var list []*Transport
	rowCount, err := qs.Limit(limit, (offset-1)*limit).All(&list)
	if err != nil {
		return 0, nil, err
	}
	return rowCount, list, nil
}
