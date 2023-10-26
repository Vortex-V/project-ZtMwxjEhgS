package requests

import (
	"app/src/models"
	"github.com/beego/beego/v2/core/validation"
)

type (
	AccountSingUpRequest struct {
		request
		Username string `valid:"Required"`
		Password string `valid:"Required"`
	}
	AccountSignInRequest struct {
		request
		Username string `valid:"Required"`
		Password string `valid:"Required"`
	}
	AccountUpdateRequest struct {
		request
		Username string
		Password string
	}

	AdminAccountWriteRequest struct {
		request
		Username string `valid:"Required"`
		Password string `valid:"Required"`
		IsAdmin  bool
		Balance  float64 `valid:"Required"`
	}
)

func (r *AccountSingUpRequest) Valid(v *validation.Validation) {
	usernameUniqueRequest(v, r.Username)
}
func (r *AccountUpdateRequest) Valid(v *validation.Validation) {
	usernameUniqueRequest(v, r.Username)
}
func (r *AdminAccountWriteRequest) Valid(v *validation.Validation) {
	usernameUniqueRequest(v, r.Username)
}

func usernameUniqueRequest(v *validation.Validation, username string) {
	query := models.Find(&models.Account{Username: username}, "id").Where("username = ?")
	result, err := models.Raw(query, username).Exec()
	if err != nil {
		return
	}
	count, _ := result.RowsAffected()
	if count > 0 {
		v.SetError("Username", "Указанное имя пользователя уже занято")
	}
}
