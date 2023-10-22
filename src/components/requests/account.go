package requests

import (
	"app/src/models"
	"github.com/beego/beego/v2/core/validation"
)

type AccountSingUpRequest struct {
	request
	Username string `valid:"Required"`
	Password string `valid:"Required"`
}

func (r *AccountSingUpRequest) Valid(v *validation.Validation) {
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

type AccountSignInRequest struct {
	request
	Username string `valid:"Required"`
	Password string `valid:"Required"`
}

type AccountUpdateRequest struct {
	AccountSingUpRequest
	Username string
	Password string
}
