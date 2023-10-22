package requests

type (
	TransportPostRequest struct {
		TransportPutRequest
		TransportType string `valid:"Required"`
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
)
