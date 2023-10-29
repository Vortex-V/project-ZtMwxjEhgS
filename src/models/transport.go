package models

import (
	"fmt"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Transport struct {
	model
	Id          int64     `orm:"auto;pk"`
	Account     *Account  `orm:"column(account_id);rel(fk)"`
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

const (
	TransportTypeCar     = "Car"
	TransportTypeBike    = "Bike"
	TransportTypeScooter = "Scooter"
)

var transportTypes = map[string]string{
	TransportTypeCar:     "Car",
	TransportTypeBike:    "Bike",
	TransportTypeScooter: "Scooter",
}

func GetTransportType(key string) string {
	if v, ok := transportTypes[key]; ok {
		return v
	} else {
		return ""
	}
}

func (t *Transport) SetTransportType(key string) bool {
	if v, ok := transportTypes[key]; ok {
		t.Type = v
		return true
	}

	return false
}

func (t *Transport) Create() error {

	if t.MinutePrice == 0 && t.DayPrice == 0 {
		t.CanBeRented = false
	}

	_, err := o.Insert(t)
	return err
}

func TransportSearch(params map[string]interface{}) (int64, []*Transport, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Transport))

	if params["lat"] != nil &&
		params["long"] != nil &&
		params["radius"] != nil {
		lat := params["lat"].(float64)
		long := params["long"].(float64)
		radius := params["radius"].(float64)
		// Делаем подзапрос с выборкой id, так как QuerySetter требует указать имя поля, для которого применить фильтр
		q := Find(new(Transport), "id").
			Where(fmt.Sprintf("pow(transports.latitude-%f, 2) + pow(transports.longitude-%f,2)<=pow(%f,2)", lat, long, radius)).
			String()
		cond := orm.NewCondition()
		// два раза in, потому что он сам не добавляет in (баг? || устаревшая документация?)
		cond = cond.Raw("id__in", fmt.Sprintf("in (%s)", q))
		qs = qs.SetCond(cond)

	}

	if params["type"] != nil {
		transportType := params["type"].(string)
		if transportType != "" {
			qs = qs.Filter("type", transportType)
		}
	}

	if params["can_be_rented"] != nil {
		qs = qs.Filter("can_be_rented", params["can_be_rented"] == "1")
	}

	if params["start"] != nil &&
		params["count"] != nil {
		start := params["start"].(int64)
		count := params["count"].(int64)
		qs = qs.Limit(count, (start-1)*count)
	}

	var list []*Transport
	rowCount, err := qs.All(&list)
	if err != nil {
		return 0, nil, err
	}
	return rowCount, list, nil
}
