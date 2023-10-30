package forms

import (
	"app/src/components/customValid"
	"github.com/beego/beego/v2/core/validation"
)

type RentTransportForm struct {
	form
	Long   float64 `form:"long" valid:"Required"`
	Lat    float64 `form:"lat" valid:"Required"`
	Radius float64 `form:"radius" valid:"Required"`
	Type   string  `form:"type" valid:"Required;Alpha"`
}

func (f *RentTransportForm) Valid(v *validation.Validation) {
	customValid.TransportTypeExists(v, f.Type)
}

type RentEndForm struct {
	form
	Long float64 `form:"long" valid:"Required"`
	Lat  float64 `form:"lat" valid:"Required"`
}

type RentNewForm struct {
	form
	RentType string `form:"rentType" valid:"Required;Alpha"`
}

func (f *RentNewForm) Valid(v *validation.Validation) {
	customValid.RentTypeExists(v, f.RentType)
}
