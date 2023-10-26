package controllers

import (
	"app/src/components/requests"
	"app/src/components/responses"
	"app/src/models"
)

// AccountController operations for Account
type AccountController struct {
	Controller
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
		c.ResponseError(err.Error(), 500)
		return
	}
	account := c.findModel(id)
	if account == nil {
		return
	}

	c.ResponseMapTo(new(responses.AccountMeResponse), account)
}

// SignIn
// @Title SignIn
// @Param	body	body	requests.AccountSignInRequest	"sign in request"
// @Success 200 {"token":"string"}
// @Failure 400	user or password is incorrect
// @router /SignIn [post]
func (c *AccountController) SignIn() {
	var data requests.AccountSignInRequest
	if !c.LoadAndValidate(&data) {
		return
	}

	account := models.Account{
		Username: data.Username,
	}

	token, err := account.Login(data.Password)
	if err != nil {
		c.ResponseError(err.Error(), 500)
		return
	}

	c.ResponseJson(DataMap{"token": token})
}

// SignUp
// @Title SignUp
// @Param	body	body	requests.AccountSingUpRequest	"sign up request"
// @Success	200	{object}	responses.AccountSignUpResponse	Указанный объект может быть получен по ключу data
// @Failure 400	username already exists
// @router /SignUp [post]
func (c *AccountController) SignUp() {
	var data requests.AccountSingUpRequest
	if !c.LoadAndValidate(&data) {
		return
	}

	account := new(models.Account)
	err := account.Register(data.Username, data.Password)
	if err != nil {
		c.ResponseError(err.Error(), 500)
		return
	}

	c.ResponseMapTo(new(responses.AccountSignUpResponse), account, "Аккаунт успешно создан")
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
		c.ResponseError(err.Error(), 500)
		return
	}
	account := c.findModel(id)
	if account == nil {
		return
	}

	// Костыль. Старый jwt остаётся валидным в течение 1 часа.
	// Но jwt в базе хранить это кукож, товарищи, я так делать не буду.
	// Хотелось бы дополнительно использовать refresh token
	account.IsNeedRelogin = true
	_, err = models.Update(account, "IsNeedRelogin")
	if err != nil {
		c.ResponseError(err.Error(), 500)
		return
	}

	c.ResponseJson(DataMap{"message": "Выполнен выход из аккаунта"})
}

// Update
// @Title Update
// @Security	api_key
// @Param	body	body	requests.AccountUpdateRequest	"update request"
// @Success 200
// @Failure 400	username already exists
// @Failure 401	unauthorized
// @router /Update [put]
func (c *AccountController) Update() {
	id, err := c.GetInt64("accountId")
	if err != nil {
		c.ResponseError(err.Error(), 500)
		return
	}
	var data requests.AccountUpdateRequest
	if !c.LoadAndValidate(&data) {
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
		c.ResponseError(err.Error(), 500)
		return
	}

	c.ResponseJson(DataMap{"message": "Данные успешно изменены"})
}

func (c *AccountController) findModel(id int64) *models.Account {
	m := &models.Account{Id: id}
	if err := models.Get(m); err != nil {
		c.ResponseError(ErrorNotFound, 404)
		return nil
	}
	return m
}
