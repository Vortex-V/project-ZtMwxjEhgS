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
	RentTypeMinutes = "Minutes"
	RentTypeDays    = "Days"
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
	_ = Read(t.Transport)
	return t.Transport.IsOwner(id)
}

func (t *Rent) CanRent(accountId int64, transport *Transport) error {
	switch t.Type {
	case RentTypeDays:
		if transport.DayPrice == 0 {
			return errors.New("для данного транспорта недоступна аренда по дням")
		}
	case RentTypeMinutes:
		if transport.MinutePrice == 0 {
			return errors.New("для данного транспорта недоступна аренда по минутам")
		}
	}
	if !transport.CanBeRented /*TODO || t.Transport.Status == TransportStatusRented*/ {
		return errors.New("данный транспорт не может быть арендован")
	}
	if transport.IsOwner(accountId) {
		return errors.New("нельзя арендовать свой транспорт")
	}
	return nil
}

func (t *Rent) SetTimeStart(v string) (err error) {
	t.TimeStart, err = time.Parse(time.DateTime, v)
	return err
}

func (t *Rent) SetTimeEnd(v string) (err error) {
	t.TimeEnd, err = time.Parse(time.DateTime, v)
	return err
}

func (t *Rent) Create() error {
	tx, _ := o.Begin()
	defer tx.RollbackUnlessCommit()

	if t.PriceOfUnit == 0 {
		t.PriceOfUnit = t.GetPriceByRentType(t.Type)
	}
	// t.Transport.Status = TransportStatusRented TODO заместо изменения CanBeRented
	t.Transport.CanBeRented = false
	_, err := tx.Update(t.Transport)
	if err != nil {
		return err
	}

	if t.TimeStart.IsZero() {
		t.TimeStart = time.Now()
	}
	t.Status = RentStatusActive

	_, err = tx.Insert(t)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (t *Rent) End(lat, long float64) error {
	transport := t.Transport
	_ = Read(transport)
	ownerAccount := transport.Account
	_ = Read(ownerAccount)
	renterAccount := t.Account
	_ = Read(renterAccount)

	if t.Status != RentStatusActive {
		return errors.New("нельзя завершить неактивную аренду")
	}

	tx, _ := o.Begin()
	defer tx.RollbackUnlessCommit()

	t.TimeEnd = time.Now()
	t.calculateFinalPrice()
	if renterAccount.Balance < t.FinalPrice {
		return errors.New("недостаточно средств на счету")
	}
	// TODO хотелось бы учитывать комиссию сервиса
	ownerAccount.Balance += t.FinalPrice
	renterAccount.Balance -= t.FinalPrice
	_, err := tx.Update(ownerAccount, "Balance")
	if err != nil {
		return err
	}
	_, err = tx.Update(renterAccount, "Balance")
	if err != nil {
		return err
	}

	transport.Latitude = lat
	transport.Longitude = long
	transport.CanBeRented = true
	// t.Transport.Status = TransportStatusNotRented TODO заместо изменения CanBeRented
	_, err = tx.Update(transport)
	if err != nil {
		return err
	}

	t.Status = RentStatusCompleted

	_, err = tx.Update(t)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (t *Rent) calculateFinalPrice() {
	duration := t.TimeEnd.Sub(t.TimeStart)
	switch t.Type {
	case RentTypeDays:
		t.FinalPrice = t.PriceOfUnit * math.Ceil(duration.Hours()/24)
	case RentTypeMinutes:
		t.FinalPrice = t.PriceOfUnit * duration.Minutes()
	}
	// Сразу обновляем, чтобы пользователь мог посмотреть
	Update(t, "FinalPrice")
}

func RentSearch(params map[string]interface{}) (int64, []*Rent, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Rent))

	if params["owner_id"] != nil {
		ownerId := params["owner_id"].(int64)
		qs = qs.Filter("Transport__Account__Id", ownerId)
	}

	if params["account_id"] != nil {
		accountId := params["account_id"].(int64)
		qs = qs.Filter("Account__Id", accountId)
	}

	if params["transport_id"] != nil {
		transportId := params["transport_id"].(int64)
		qs = qs.Filter("Transport__Id", transportId)
	}

	if params["start"] != nil &&
		params["count"] != nil {
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
