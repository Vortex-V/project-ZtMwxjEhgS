package responses

import (
	"app/src/models"
	"time"
)

type RentTransportResponse struct {
	response
	Id          int64
	AccountId   int64
	Type        string
	Model       string
	Color       string
	Identifier  string
	Description string
	Latitude    float64
	Longitude   float64
	MinutePrice float64
	DayPrice    float64
}

type RentResponse struct {
	response
	Id          int64
	AccountId   int64
	Type        string
	TransportId int64
	TimeStart   string
	TimeEnd     string
	PriceOfUnit float64
	FinalPrice  float64
	Status      string
}

func (r *RentResponse) Construct(m models.Model) interface{} {
	rent := m.(*models.Rent)
	return &RentResponse{
		Id:          rent.Id,
		AccountId:   rent.Account.Id,
		Type:        models.GetRentType(rent.Type),
		TransportId: rent.Transport.Id,
		TimeStart:   rent.TimeStart.Format(time.DateTime),
		TimeEnd:     rent.TimeEnd.Format(time.DateTime),
		PriceOfUnit: rent.PriceOfUnit,
		FinalPrice:  rent.FinalPrice,
		Status:      rent.GetStatusLabel(),
	}
}
