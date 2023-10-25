package requests

type (
	TransportPostRequest struct {
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
