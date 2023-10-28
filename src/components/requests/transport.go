package requests

import (
	"app/src/models"
	"github.com/beego/beego/v2/core/validation"
)

type (
	TransportPostRequest struct {
		request
		CanBeRented   bool   `valid:"Required"`
		TransportType string `valid:"Required"`
		Model         string `valid:"Required"`
		Color         string `valid:"Required"`
		Identifier    string `valid:"Required"`
		Description   string
		Latitude      float64 `valid:"Required"`
		Longitude     float64 `valid:"Required"`
		MinutePrice   float64
		DayPrice      float64
	}

	TransportPutRequest struct {
		request
		CanBeRented bool   `valid:"Required"`
		Model       string `valid:"Required"`
		Color       string `valid:"Required"`
		Identifier  string `valid:"Required"`
		Description string
		Latitude    float64 `valid:"Required"`
		Longitude   float64 `valid:"Required"`
		MinutePrice float64
		DayPrice    float64
	}

	AdminTransportWriteRequest struct {
		request
		OwnerId       int64  `valid:"Required"`
		CanBeRented   bool   `valid:"Required"`
		TransportType string `valid:"Required"`
		Model         string `valid:"Required"`
		Color         string `valid:"Required"`
		Identifier    string `valid:"Required"`
		Description   string
		Latitude      float64 `valid:"Required"`
		Longitude     float64 `valid:"Required"`
		MinutePrice   float64
		DayPrice      float64
	}
)

func (t *TransportPostRequest) Valid(v *validation.Validation) {
	transportTypeExists(v, t.TransportType)
}
func (t *AdminTransportWriteRequest) Valid(v *validation.Validation) {
	transportTypeExists(v, t.TransportType)
}

func transportTypeExists(v *validation.Validation, transportType string) {
	if models.GetTransportTypeLabel(transportType) == "" {
		v.SetError("TransportType", "Тип транспорта должен быть одним из [Car, Bike, Scooter]")
	}
}
