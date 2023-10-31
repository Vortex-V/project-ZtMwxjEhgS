package controllers

import (
	"app/src/components/auth"
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
// @Description	получение данных о текущем аккаунте
// @Security	api_key
// @Success	200	{object}	responses.AccountMeResponse Указанный объект может быть получен по ключу data
// @Failure 401	unauthorized
// @router /Me [get]
func (c *AccountController) Me() {
	id := c.GetIdentityId()
	if id == 0 {
		return
	}

	account := c.findModel(id)
	if account == nil {
		return
	}

	response := responses.New[*responses.AccountMeResponse](
		new(responses.AccountMeResponse), account)
	c.Response(response)
}

// SignIn
// @Title SignIn
// @Description	получение нового jwt токена пользователя
// @Param	body	body	requests.AccountSignInRequest	"sign in request"
// @Success 200 {"token":"string"}
// @Failure 400	user or password is incorrect
// @router /SignIn [post]
func (c *AccountController) SignIn() {
	var data requests.AccountSignInRequest
	if !c.ParseAndValidateRequest(&data) {
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
// @Description	регистрация нового аккаунта
// @Param	body	body	requests.AccountSignUpRequest	"sign up request"
// @Success	200	{object}	responses.AccountResponse	Указанный объект может быть получен по ключу data
// @Failure 400	username already exists
// @router /SignUp [post]
func (c *AccountController) SignUp() {
	var data requests.AccountSignUpRequest
	if !c.ParseAndValidateRequest(&data) {
		return
	}

	account := new(models.Account)
	err := account.Register(data.Username, data.Password)
	if err != nil {
		c.ResponseError(err.Error(), 500)
		return
	}

	response := responses.New[*responses.AccountResponse](
		new(responses.AccountResponse), account)
	c.Response(response, "Аккаунт успешно создан")
}

// SignOut
// @Title SignOut
// @Description	выход из аккаунта
// @Security	api_key
// @Success 200
// @Failure 401	unauthorized
// @router /SignOut [post]
func (c *AccountController) SignOut() {
	id := c.GetIdentityId()
	if id == 0 {
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
	_, err := models.Update(account, "IsNeedRelogin")
	if err != nil {
		c.ResponseError(err.Error(), 500)
		return
	}

	c.Response("Выполнен выход из аккаунта")
}

// Update
// @Title Update
// @Description	обновление своего аккаунта
// @Security	api_key
// @Param	body	body	requests.AccountUpdateRequest	"update request"
// @Success 200	{object}	responses.AccountResponse
// @Failure 400	username already exists
// @Failure 401	unauthorized
// @router /Update [put]
func (c *AccountController) Update() {
	id := c.GetIdentityId()
	if id == 0 {
		return
	}
	var data requests.AccountUpdateRequest
	if !c.ParseAndValidateRequest(&data) {
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
		account.Password, _ = auth.HashPassword(data.Password)
	}

	_, err := models.Update(account)
	if err != nil {
		c.ResponseError(err.Error(), 500)
		return
	}
	response := responses.New[*responses.AccountResponse](
		new(responses.AccountResponse), account)
	c.Response(response, "Данные успешно изменены")
}

func (c *AccountController) findModel(id int64) *models.Account {
	m := &models.Account{Id: id}
	if err := models.Read(m); err != nil {
		c.ResponseError(ErrorNotFound, 404)
		return nil
	}
	return m
}
