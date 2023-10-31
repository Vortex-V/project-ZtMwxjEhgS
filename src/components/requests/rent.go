package requests

import (
	"app/src/components/customValid"
	"github.com/beego/beego/v2/core/validation"
)

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

func (r *AdminRentWriteRequest) Valid(v *validation.Validation) {
	customValid.RentTypeExists(v, r.PriceType)
}
