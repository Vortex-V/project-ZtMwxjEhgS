package controllers

import (
	"app/src/components/auth"
	"app/src/components/requests"
	"app/src/components/responses"
	"app/src/models"
)

// AccountController operations for Account
type AccountController struct {
	controller
}

// Me
// @Title Me
// @Success	200	{object}	responses.AccountMeResponse
// @router /Me [get]
func (c *AccountController) Me() {
	id, err := c.GetInt(":id")
	if err != nil {
		c.responseError(err, 500)
		return
	}
	account := c.findModel(id)
	if account == nil {
		return
	}

	c.responseMapTo(new(responses.AccountMeResponse), account)
}

// SignIn
// @Title SignIn
// @Param	body	body	requests.AccountRequest	true
// @Success 200 {{"token": "token"}}
// @router /SignIn [post]
func (c *AccountController) SignIn() {
	var data requests.AccountRequest

	if !c.load(&data) {
		return
	}

	account := models.Account{
		Username: data.Username,
	}

	token, err := auth.Login(&account, data.Password)
	if err != nil {
		c.responseError(err, 500)
		return
	}

	c.response(dataMap{"token": token})
}

// SignUp
// @Title SignUp
// @Param	body	body	requests.AccountRequest	true
// @Success	200	{object}	responses.AccountSignUpResponse
// @router /SignUp [post]
func (c *AccountController) SignUp() {
	var data requests.AccountRequest
	if !c.load(&data) {
		return
	}

	var account models.Account
	err := account.Register(data)
	if err != nil {
		c.responseError(err, 500)
		return
	}

	c.responseMapTo(new(responses.AccountSignUpResponse), account, "Аккаунт успешно создан")
}

// SignOut
// @Title SignOut
// @Success 200
// @router /SignOut [post]
func (c *AccountController) SignOut() {
	id, err := c.GetInt(":id")
	if err != nil {
		c.responseError(err, 500)
		return
	}
	account := c.findModel(id)
	if account == nil {
		return
	}

	account.IsNeedRelogin = true

	_, err = models.Update(account, "IsNeedRelogin")
	if err != nil {
		c.responseError(err, 500)
		return
	}

	c.response(dataMap{"message": "success"})
}

// Update
// @Title Update
// @Param	body	body	requests.AccountRequest	true
// @Success 200
// @router /Update [put]
func (c *AccountController) Update() {
	id, err := c.GetInt(":id")
	if err != nil {
		c.responseError(err, 500)
		return
	}
	var data requests.AccountUpdateRequest
	if !c.load(&data) {
		return
	}

	account := c.findModel(id)
	if account == nil {
		return
	}
	if data.Username != "" {
		account.Username = data.Username
	}
	if data.Password != "" {
		account.Password = data.Password
	}

	_, err = models.Update(account)
	if err != nil {
		c.responseError(err, 500)
		return
	}

	c.response(dataMap{"message": "Данные успешно изменены"})
}

func (c *AccountController) findModel(id int) *models.Account {
	m := &models.Account{Id: id}
	if err := models.Get(m); err != nil {
		c.responseError(ErrorNotFound, 404)
		return nil
	}
	return m
}
