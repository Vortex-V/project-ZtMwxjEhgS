package controllers

import (
	"app/src/components/auth"
	"app/src/components/requests"
	"app/src/components/responses"
	"app/src/models"
	"errors"
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
	account := findModel(id)
	response, err := responses.MapTo(new(responses.AccountMeResponse), account)
	if err != nil {
		c.responseError(err, 500)
		return
	}

	c.response(response)
}

// SignIn
// @Title SignIn
// @Param	body	body	requests.AccountRequest	true
// @Success 200 {{"token": "token"}}
// @router /SignIn [post]
func (c *AccountController) SignIn() {
	var data requests.AccountRequest

	err := c.parseRequestBody(&data)
	if err != nil {
		c.responseError(err, 500)
		return
	}

	if validationErrors := validateRequest(&data); len(validationErrors) > 0 {
		c.responseValidationError(validationErrors, 400)
		return
	}

	account := models.Account{
		Username: data.Username,
	}

	query := models.Find(&account, "id", "password").Where("username = ?")
	err = models.Raw(query, data.Username).QueryRow(&account)
	if err != nil {
		c.responseError(err, 500)
		return
	}

	if err = auth.CheckPasswordHash(data.Password, account.Password); err != nil {
		c.responseError(errors.New("username or password is incorrect"), 400)
		return
	}

	token, err := auth.Login(account)
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
	err := c.parseRequestBody(&data)
	if err != nil {
		c.responseError(err, 500)
		return
	}

	if validationErrors := validateRequest(&data); len(validationErrors) > 0 {
		c.responseValidationError(validationErrors, 400)
		return
	}

	password, err := auth.HashPassword(data.Password)
	if err != nil {
		c.responseError(err, 500)
		return
	}

	account := models.Account{
		Username: data.Username,
		Password: password,
	}
	_, err = models.Insert(&account)
	if err != nil {
		c.responseError(err, 500)
		return
	}

	response, err := responses.MapTo(new(responses.AccountSignUpResponse), &account)
	if err != nil {
		c.responseError(err, 500)
		return
	}

	c.response(response)
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
	account := findModel(id)
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
	err = c.parseRequestBody(&data)
	if err != nil {
		c.responseError(err, 500)
		return
	}

	if validationErrors := validateRequest(&data); len(validationErrors) > 0 {
		c.responseValidationError(validationErrors, 400)
		return
	}

	account := findModel(id)
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

	c.response(dataMap{"message": "success"})
}

func findModel(id int) *models.Account {
	m := &models.Account{Id: id}
	if err := models.Get(m); err != nil {
		return nil
	}
	return m
}
