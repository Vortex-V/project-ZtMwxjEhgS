package responses

import "app/src/models"

type (
	TransportResponse struct {
		response
		Id          int64
		AccountId   int64
		CanBeRented bool
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

	AdminTransportResponseCollection struct {
		Transports []models.Transport
	}
)