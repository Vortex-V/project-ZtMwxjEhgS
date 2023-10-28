package models

import (
	"errors"
	"math"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Rent struct {
	model
	Id          int64      `orm:"auto;pk"`
	Account     *Account   `orm:"rel(fk)"`
	Type        string     `orm:""`
	Transport   *Transport `orm:"column(transport_id);rel(one)"`
	TimeStart   time.Time  `orm:"column(time_start);type(timestamp without time zone);auto_now_add"`
	TimeEnd     time.Time  `orm:"column(time_end);type(timestamp without time zone);null"`
	PriceOfUnit float64    `orm:"column(price_of_unit)"`
	FinalPrice  float64    `orm:"column(final_price);null"`
	Status      string     `orm:"column(status)"`
	CreatedAt   time.Time  `orm:"column(created_at);type(timestamp without time zone);auto_now_add"`
	UpdatedAt   time.Time  `orm:"column(updated_at);type(timestamp without time zone);auto_now"`
}

func (t *Rent) TableName() string {
	return "rental"
}

func init() {
	orm.RegisterModel(new(Rent))
}

const (
	RentStatusActive    = "active"
	RentStatusCanceled  = "canceled"
	RentStatusCompleted = "completed"
)

var statusLabels = map[string]string{
	RentStatusActive:    "Active",
	RentStatusCanceled:  "Canceled",
	RentStatusCompleted: "Completed",
}

func (t *Rent) GetStatusLabel() string {
	return statusLabels[t.Status]
}

// TODO сейчас типы сильно перегружены. Если единственным различием останется написание с заглавной, то убрать все эти проверки и хранить только лейблы
const (
	RentTypeMinutes = "minutes"
	RentTypeDays    = "days"
)

var rentTypeLabels = map[string]string{
	RentTypeMinutes: "Minutes",
	RentTypeDays:    "Days",
}

func (t *Rent) GetRentTypeLabel() string {
	return rentTypeLabels[t.Type]
}

func GetRentType(key string) string {
	if v, ok := rentTypeLabels[key]; ok {
		return v
	} else {
		return ""
	}
}

func GetRentTypeKeyByLabel(label string) string {
	for k, v := range rentTypeLabels {
		if v == label {
			return k
		}
	}
	return ""
}

func (t *Rent) GetPriceByRentType(key string) float64 {
	if key == RentTypeMinutes {
		return t.Transport.MinutePrice
	} else if key == RentTypeDays {
		return t.Transport.DayPrice
	}
	return 0
}

func (t *Rent) IsRenter(id int64) bool {
	return t.Account.Id == id
}

func (t *Rent) IsOwner(id int64) bool {
	return t.Transport.Account.Id == id
}

func (t *Rent) SetTimeStart(v string) (err error) {
	t.TimeStart, err = time.Parse(time.DateTime, v)
	return err
}

func (t *Rent) SetTimeEnd(v string) (err error) {
	t.TimeEnd, err = time.Parse(time.DateTime, v)
	return err
}

func (t *Rent) SetType(v string) bool {
	t.Type = GetRentTypeKeyByLabel(v)
	return !(t.Type == "")
}

func (t *Rent) Create() error {
	if t.IsOwner(t.Account.Id) {
		return errors.New("нельзя арендовать свой транспорт")
	}
	if !t.Transport.CanBeRented /*TODO || t.Transport.Status == TransportStatusRented*/ {
		return errors.New("транспорт уже арендован")
	}

	if t.PriceOfUnit == 0 {
		t.PriceOfUnit = t.GetPriceByRentType(t.Type)
	}
	// t.Transport.Status = TransportStatusRented TODO заместо изменения CanBeRented
	t.Transport.CanBeRented = false
	if t.TimeStart.IsZero() {
		t.TimeStart = time.Now()
	}
	t.Status = RentStatusActive

	_, err := o.Insert(t)

	return err
}

func (t *Rent) End(params map[string]interface{}) error {
	t.Transport.Latitude = (params["lat"]).(float64)
	t.Transport.Longitude = (params["long"]).(float64)

	if t.Status != RentStatusActive {
		return errors.New("нельзя завершить неактивную аренду")
	}

	t.TimeEnd = time.Now()
	t.calculateFinalPrice()
	t.Status = RentStatusCompleted
	t.Transport.CanBeRented = true
	// t.Transport.Status = TransportStatusNotRented TODO заместо изменения CanBeRented

	_, err := o.Update(t)

	return err
}

func (t *Rent) calculateFinalPrice() {
	duration := t.TimeEnd.Sub(t.TimeStart)
	switch t.Type {
	case RentTypeDays:
		t.FinalPrice = t.PriceOfUnit * math.Ceil(duration.Hours()/24)
	case RentTypeMinutes:
		t.FinalPrice = t.PriceOfUnit * duration.Minutes()
	}
}

func RentSearch(params map[string]interface{}) (int64, []*Rent, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Rent))

	if params["owner_id"] != "" {
		ownerId := params["owner_id"].(int64)
		qs.Filter("transport__account__id", ownerId)
	}

	if params["account_id"] != "" {
		accountId := params["account_id"].(int64)
		qs.Filter("account__id", accountId)
	}

	if params["transport_id"] != "" {
		transportId := params["transport_id"].(int64)
		qs.Filter("transport__id", transportId)
	}

	if params["start"] != "" &&
		params["count"] != "" {
		start := params["start"].(int)
		count := params["count"].(int)
		qs = qs.Limit(count, (start-1)*count)
	}

	var list []*Rent
	rowCount, err := qs.All(&list)
	if err != nil {
		return 0, nil, err
	}
	return rowCount, list, nil
}
