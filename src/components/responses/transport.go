package responses

type (
	TransportResponse struct {
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
)
