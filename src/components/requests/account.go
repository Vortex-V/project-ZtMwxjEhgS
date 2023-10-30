package requests

import (
	"app/src/components/customValid"
	"github.com/beego/beego/v2/core/validation"
)

type (
	AccountSignUpRequest struct {
		request
		Username string `valid:"Required" json:"username"`
		Password string `valid:"Required" json:"password"`
	}
	AccountSignInRequest struct {
		request
		Username string `valid:"Required" json:"username"`
		Password string `valid:"Required" json:"password"`
	}
	AccountUpdateRequest struct {
		request
		Username string `json:"username"`
		Password string `json:"password"`
	}

	AdminAccountWriteRequest struct {
		request
		Username string  `valid:"Required" json:"username"`
		Password string  `valid:"Required" json:"password"`
		IsAdmin  bool    `json:"isAdmin"`
		Balance  float64 `valid:"Required" json:"balance"`
	}
)

func (r *AccountSignUpRequest) Valid(v *validation.Validation) {
	customValid.UsernameUniqueRequest(v, r.Username)
}
func (r *AccountUpdateRequest) Valid(v *validation.Validation) {
	customValid.UsernameUniqueRequest(v, r.Username)
}
func (r *AdminAccountWriteRequest) Valid(v *validation.Validation) {
	customValid.UsernameUniqueRequest(v, r.Username)
}
