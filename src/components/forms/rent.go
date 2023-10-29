package forms

import (
	"app/src/models"
	"github.com/beego/beego/v2/core/validation"
)

type RentTransportForm struct {
	form
	Long   float64 `form:"long" valid:"Required;Match(/^([0-9]*[.])?[0-9]+$/)"`
	Lat    float64 `form:"lat" valid:"Required;Match(/^([0-9]*[.])?[0-9]+$/)"`
	Radius float64 `form:"radius" valid:"Required;Match(/^([0-9]*[.])?[0-9]+$/)"`
	Type   string  `form:"type" valid:"Required;Alpha"`
}

func (f *RentTransportForm) Valid(v *validation.Validation) {
	transportTypeExists(v, f.Type)
}

type RentEndForm struct {
	form
	Long float64 `form:"long" valid:"Required;Match(/^([0-9]*[.])?[0-9]+$/)"`
	Lat  float64 `form:"lat" valid:"Required;Match(/^([0-9]*[.])?[0-9]+$/)"`
}

type RentNewForm struct {
	form
	RentType string `form:"rentType" valid:"Required;Alpha"`
}

func (f *RentNewForm) Valid(v *validation.Validation) {
	rentTypeExists(v, f.RentType)
}

func rentTypeExists(v *validation.Validation, rentType string) {
	if rentType != "All" {
		rentType = models.GetRentType(rentType)
		if rentType == "" {
			v.SetError("rentType", "rentType должен быть одним из [Days, Minutes, All]")
			return
		}
	}
}
