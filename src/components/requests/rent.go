package requests

type (
	AdminRentWriteRequest struct {
		request
		TransportId int64  `valid:"Required"`
		UserId      int64  `valid:"Required"`
		TimeStart   string `valid:"Required"`
		TimeEnd     string
		PriceOfUnit float64 `valid:"Required"`
		PriceType   string  `valid:"Required"`
		FinalPrice  float64
	}
)
