package requests

import (
	"app/src/models"
	"github.com/beego/beego/v2/core/validation"
)

type AccountRequest struct {
	request
	Username string `valid:"Required"`
	Password string `valid:"Required"`
}

type AccountUpdateRequest struct {
	AccountRequest
	Username string
	Password string
}

func (r *AccountRequest) Valid(v *validation.Validation) {
	query := models.Find(&models.Account{Username: r.Username}, "").Where("username = ?")
	result, _ := models.Raw(query).Exec()
	count, _ := result.RowsAffected()
	if count > 0 {
		v.SetError("Username", "Указанное имя пользователя уже занято")
	}
}
