package requests

type (
	AdminRentWriteRequest struct {
		request
		TransportId int64   `valid:"Required" json:"transportId"`
		UserId      int64   `valid:"Required" json:"userId"`
		TimeStart   string  `valid:"Required" json:"timeStart"`
		TimeEnd     string  `json:"timeEnd"`
		PriceOfUnit float64 `valid:"Required" json:"priceOfUnit"`
		PriceType   string  `valid:"Required" json:"priceType"`
		FinalPrice  float64 `json:"finalPrice"`
	}
)
