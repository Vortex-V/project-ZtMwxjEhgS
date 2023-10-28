package forms

import (
	"app/src/models"
	"github.com/beego/beego/v2/core/validation"
)

type RentEndForm struct {
	form
	Long float64 `form:"long" valid:"Required"`
	Lat  float64 `form:"lat" valid:"Required"`
}

type RentNewForm struct {
	form
	RentType string `form:"rentType" valid:"Required"`
}

func (f *RentNewForm) Valid(v *validation.Validation) {
	rentType := models.GetRentTypeKeyByLabel(f.RentType)
	if rentType == "" {
		v.SetError("rentType", "rentType должен быть одним из [Days, Minutes]")
		return
	}
}
