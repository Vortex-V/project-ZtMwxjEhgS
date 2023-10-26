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

func (r *TransportResponse) CustomFields(m models.Model) interface{} {
	transport := m.(*models.Transport)
	return struct {
		OwnerId       int64
		TransportType string
	}{
		OwnerId:       transport.Account.Id,
		TransportType: transport.Type,
	}
}
