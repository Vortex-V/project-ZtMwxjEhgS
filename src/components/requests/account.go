package requests

import (
	"app/src/models"
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
