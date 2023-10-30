package forms

import (
	"app/src/components/customValid"
	"github.com/beego/beego/v2/core/validation"
)

type TransportGetAllForm struct {
	form
	Start         int    `form:"start" validate:"Required;Min(0)"`
	Count         int    `form:"count" validate:"Required;Min(0)"`
	TransportType string `form:"transportType" validate:"Required;Alpha"`
}

func (f *TransportGetAllForm) Valid(v *validation.Validation) {
	customValid.TransportTypeExists(v, f.TransportType)
}
