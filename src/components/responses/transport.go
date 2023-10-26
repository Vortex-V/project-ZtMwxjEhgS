package responses

import "app/src/models"

type TransportResponse struct {
	response
	Id            int64
	OwnerId       int64
	CanBeRented   bool
	TransportType string
	Model         string
	Color         string
	Identifier    string
	Description   string
	Latitude      float64
	Longitude     float64
	MinutePrice   float64
	DayPrice      float64
}

func (r TransportResponse) Construct(m models.Model) interface{} {
	transport := m.(*models.Transport)
	return &TransportResponse{
		Id:            transport.Id,
		OwnerId:       transport.Account.Id,
		CanBeRented:   transport.CanBeRented,
		TransportType: transport.Type,
		Model:         transport.Model,
		Color:         transport.Color,
		Identifier:    transport.Identifier,
		Description:   transport.Description,
		Latitude:      transport.Latitude,
		Longitude:     transport.Longitude,
		MinutePrice:   transport.MinutePrice,
		DayPrice:      transport.DayPrice,
	}
}
