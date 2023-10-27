package models

import (
	"errors"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Rent struct {
	model
	Id          int64      `orm:"auto;pk"`
	Account     *Account   `orm:"rel(one)"`
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

func GetRentType(label string) string {
	label = strings.ToLower(label)
	if v, ok := rentTypeLabels[label]; ok {
		return v
	} else {
		return ""
	}
}

func (t *Rent) IsRenter(id int64) bool {
	return t.Account.Id == id
}

func (t *Rent) IsOwner(id int64) bool {
	return t.Transport.Account.Id == id
}

func (t *Rent) Create() error {
	if t.IsOwner(t.Account.Id) {
		return errors.New("нельзя арендовать свой транспорт")
	}
	if !t.Transport.CanBeRented /*TODO || t.Transport.Status == TransportStatusRented*/ {
		return errors.New("транспорт уже арендован")
	}

	t.PriceOfUnit = t.Transport.GetPriceByRentType(t.Type)
	// t.Transport.Status = TransportStatusRented TODO заместо изменения CanBeRented
	t.Transport.CanBeRented = false
	t.TimeStart = time.Now()
	t.Status = RentStatusActive

	_, err := o.Insert(t)

	return err
}

func (t *Rent) End() error {
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

func RentSearch(params map[string]string) (int64, []*Rent, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Rent))

	if params["owner_id"] != "" {
		ownerId, err := strconv.ParseInt(params["owner_id"], 10, 64)
		if err != nil {
			return 0, nil, errors.New("owner_id must be int")
		}
		qs.Filter("transport__account__id", ownerId)
	}

	if params["account_id"] != "" {
		accountId, err := strconv.ParseInt(params["account_id"], 10, 64)
		if err != nil {
			return 0, nil, errors.New("account_id must be int")
		}
		qs.Filter("account__id", accountId)
	}

	if params["transport_id"] != "" {
		transportId, err := strconv.ParseInt(params["transport_id"], 10, 64)
		if err != nil {
			return 0, nil, errors.New("transport_id must be int")
		}
		qs.Filter("transport__id", transportId)
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

	var list []*Rent
	rowCount, err := qs.All(&list)
	if err != nil {
		return 0, nil, err
	}
	return rowCount, list, nil
}
