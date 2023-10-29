package requests

import (
	"app/src/models"
	"github.com/beego/beego/v2/core/validation"
)

type (
	TransportPostRequest struct {
		request
		CanBeRented   bool    `valid:"Required" json:"canBeRented"`
		TransportType string  `valid:"Required" json:"transportType"`
		Model         string  `valid:"Required" json:"model"`
		Color         string  `valid:"Required" json:"color"`
		Identifier    string  `valid:"Required" json:"identifier"`
		Description   string  `json:"description"`
		Latitude      float64 `valid:"Required" json:"latitude"`
		Longitude     float64 `valid:"Required" json:"longitude"`
		MinutePrice   float64 `json:"minutePrice"`
		DayPrice      float64 `json:"dayPrice"`
	}

	TransportPutRequest struct {
		request
		CanBeRented bool    `valid:"Required" json:"canBeRented"`
		Model       string  `valid:"Required" json:"model"`
		Color       string  `valid:"Required" json:"color"`
		Identifier  string  `valid:"Required" json:"identifier"`
		Description string  `json:"description"`
		Latitude    float64 `valid:"Required" json:"latitude"`
		Longitude   float64 `valid:"Required" json:"longitude"`
		MinutePrice float64 `json:"minutePrice"`
		DayPrice    float64 `json:"dayPrice"`
	}

	AdminTransportWriteRequest struct {
		request
		OwnerId       int64   `valid:"Required" json:"ownerId"`
		CanBeRented   bool    `valid:"Required" json:"canBeRented"`
		TransportType string  `valid:"Required" json:"transportType"`
		Model         string  `valid:"Required" json:"model"`
		Color         string  `valid:"Required" json:"color"`
		Identifier    string  `valid:"Required" json:"identifier"`
		Description   string  `json:"description"`
		Latitude      float64 `valid:"Required" json:"latitude"`
		Longitude     float64 `valid:"Required" json:"longitude"`
		MinutePrice   float64 `json:"minutePrice"`
		DayPrice      float64 `json:"dayPrice"`
	}
)

func (t *TransportPostRequest) Valid(v *validation.Validation) {
	transportTypeExists(v, t.TransportType)
}
func (t *AdminTransportWriteRequest) Valid(v *validation.Validation) {
	transportTypeExists(v, t.TransportType)
}

func transportTypeExists(v *validation.Validation, transportType string) {
	if transportType != "All" {
		transportType = models.GetTransportType(transportType)
		if transportType == "" {
			v.SetError("transportType", "transportType должен быть одним из [Car, Bike, All]")
			return
		}
	}
}
