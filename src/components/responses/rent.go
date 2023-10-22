package responses

import (
	"time"
)

type (
	RentTransportResponse struct {
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
	RentTransportResponseCollection struct {
		response
		Transports []RentTransportResponse
	}

	RentGetResponse struct {
		response
		Id          int64
		AccountId   int64
		Type        string
		TransportId int64
		TimeStart   time.Time
		TimeEnd     time.Time
		PriceOfUnit float64
		FinalPrice  float64
		Status      int
	}
	RentHistoryResponseCollection struct {
		response
		Rents []RentGetResponse
	}
)
