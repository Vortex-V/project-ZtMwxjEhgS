package requests

import (
	"app/src/models"
	"github.com/beego/beego/v2/core/validation"
)

// а в php для такого есть trait
type usernameUniqueRequest struct {
	request
	Username string
}

func (r *usernameUniqueRequest) Valid(v *validation.Validation) {
	query := models.Find(&models.Account{Username: r.Username}, "id").Where("username = ?")
	result, err := models.Raw(query, r.Username).Exec()
	if err != nil {
		return
	}
	count, _ := result.RowsAffected()
	if count > 0 {
		v.SetError("Username", "Указанное имя пользователя уже занято")
	}
}

type (
	AccountSingUpRequest struct {
		usernameUniqueRequest
		Username string `valid:"Required"`
		Password string `valid:"Required"`
	}
	AccountSignInRequest struct {
		request
		Username string `valid:"Required"`
		Password string `valid:"Required"`
	}
	AccountUpdateRequest struct {
		usernameUniqueRequest
		Username string
		Password string
	}

	AdminAccountWriteRequest struct {
		usernameUniqueRequest
		Username string `valid:"Required"`
		Password string `valid:"Required"`
		IsAdmin  bool
		Balance  float64
	}
)
