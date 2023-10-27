package models

import (
	"errors"
	"fmt"
	"strconv"
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

// TODO minute day price -> CanBeRented при создании

func TransportSearch(params map[string]string) (int64, []*Transport, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Transport))

	if params["lat"] != "" &&
		params["long"] != "" &&
		params["radius"] != "" {
		lat, err := strconv.ParseFloat(params["lat"], 64)
		long, err := strconv.ParseFloat(params["long"], 64)
		radius, err := strconv.ParseFloat(params["radius"], 64)
		if err != nil {
			return 0, nil, errors.New("lat, long, radius must be float")
		}
		// Делаем подзапрос с выборкой id, так как QuerySetter требует указать имя поля, для которого применить фильтр
		q := Find(new(Transport), "id").
			Where(fmt.Sprintf("pow(transports.latitude-%f, 2) + pow(transports.longitude-%f,2)<=pow(%f,2)", lat, long, radius)).
			String()
		cond := orm.NewCondition()
		// два раза in, потому что он сам не добавляет in (баг? || устаревшая документация?)
		cond = cond.Raw("id__in", fmt.Sprintf("in (%s)", q))
		qs = qs.SetCond(cond)

	}

	if params["type"] != "" && params["type"] != "All" {
		qs = qs.Filter("type", params["type"])
	}

	if params["can_be_rented"] != "" && params["can_be_rented"] != "All" {
		qs = qs.Filter("can_be_rented", params["can_be_rented"] == "1")
	}

	if params["start"] != "" &&
		params["count"] != "" {
		start, err := strconv.ParseInt(params["start"], 10, 64)
		count, err := strconv.ParseInt(params["count"], 10, 64)
		if err != nil {
			return 0, nil, errors.New("start, count must be int")
		}
		qs = qs.Limit(count, (start-1)*count)
	}

	var list []*Transport
	rowCount, err := qs.All(&list)
	if err != nil {
		return 0, nil, err
	}
	return rowCount, list, nil
}
