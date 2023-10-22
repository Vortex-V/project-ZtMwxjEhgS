package controllers

import (
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
// @Security	api_key
// @Success	200	{object}	responses.AccountMeResponse Указанный объект может быть получен по ключу data
// @Failure 401	unauthorized
// @router /Me [get]
func (c *AccountController) Me() {
	id, err := c.GetInt64("accountId")
	if err != nil {
		c.responseError(err.Error(), 500)
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
// @Param	body	body	requests.AccountSignInRequest	"sign in request"
// @Success 200 {"token":"string"}
// @Failure 400	user or password is incorrect
// @router /SignIn [post]
func (c *AccountController) SignIn() {
	var data requests.AccountSignInRequest
	if !c.load(&data) {
		return
	}

	account := models.Account{
		Username: data.Username,
	}

	token, err := account.Login(data.Password)
	if err != nil {
		c.responseError(err.Error(), 500)
		return
	}

	c.response(dataMap{"token": token})
}

// SignUp
// @Title SignUp
// @Param	body	body	requests.AccountSingUpRequest	"sign up request"
// @Success	200	{object}	responses.AccountSignUpResponse	Указанный объект может быть получен по ключу data
// @Failure 400	username already exists
// @router /SignUp [post]
func (c *AccountController) SignUp() {
	var data requests.AccountSingUpRequest
	if !c.load(&data) {
		return
	}

	account := new(models.Account)
	err := account.Register(data.Username, data.Password)
	if err != nil {
		c.responseError(err.Error(), 500)
		return
	}

	c.responseMapTo(new(responses.AccountSignUpResponse), account, "Аккаунт успешно создан")
}

// SignOut
// @Title SignOut
// @Security	api_key
// @Success 200
// @Failure 401	unauthorized
// @router /SignOut [post]
func (c *AccountController) SignOut() {
	id, err := c.GetInt64("accountId")
	if err != nil {
		c.responseError(err.Error(), 500)
		return
	}
	account := c.findModel(id)
	if account == nil {
		return
	}

	account.IsNeedRelogin = true

	_, err = models.Update(account, "IsNeedRelogin")
	if err != nil {
		c.responseError(err.Error(), 500)
		return
	}

	c.response(dataMap{"message": "success"})
}

// Update
// @Title Update
// @Security	api_key
// @Param	body	body	requests.AccountSingUpRequest "update request"
// @Success 200
// @Failure 400	username already exists
// @Failure 401	unauthorized
// @router /Update [put]
func (c *AccountController) Update() {
	id, err := c.GetInt64("accountId")
	if err != nil {
		c.responseError(err.Error(), 500)
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
		c.responseError(err.Error(), 500)
		return
	}

	c.response(dataMap{"message": "Данные успешно изменены"})
}

func (c *AccountController) findModel(id int64) *models.Account {
	m := &models.Account{Id: id}
	if err := models.Get(m); err != nil {
		c.responseError(ErrorNotFound, 404)
		return nil
	}
	return m
}
