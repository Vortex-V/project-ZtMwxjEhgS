package forms

import (
	"app/src/models"
	"github.com/beego/beego/v2/core/validation"
)

type TransportGetAllForm struct {
	form
	Start         int    `form:"start" validate:"Required;Min(0)"`
	Count         int    `form:"count" validate:"Required;Min(0)"`
	TransportType string `form:"transportType" validate:"Required;Alpha"`
}

func (f *TransportGetAllForm) Valid(v *validation.Validation) {
	transportTypeExists(v, f.TransportType)
}

func transportTypeExists(v *validation.Validation, transportType string) {
	if transportType != "All" {
		transportType = models.GetTransportType(transportType)
		if transportType == "" {
			v.SetError("transportType", "transportType должен быть одним из [Car, Bike, All]")
			return
		}
	}
}
