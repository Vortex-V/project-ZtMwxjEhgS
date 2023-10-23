package controllers

import (
	"app/src/controllers"
	"app/src/models"
)

// AdminAccountController operations for Admin/Account
type AdminAccountController struct {
	controllers.Controller
}

// GetAll
// @Title GetAll
// @Description	Получение списка всех аккаунтов
// @Security	api_key
// @Param	start	query	int	false	Начало выборки
// @Param	count	query	int	false	Размер выборки
// @Success	200	{object}	responses.AdminAccountResponseCollection	Указанный объект может быть получен по ключу data
// @Failure 401	unauthorized
// @router /GetAll [get]
func (c *AdminAccountController) GetAll() {

}

// Get
// @Title Get
// @Description	Получение информации об аккаунте по id
// @Security	api_key
// @Param	id	path 	int64	true	"id"
// @Success	200	{object}	models.Account	Указанный объект может быть получен по ключу data
// @Failure 401	unauthorized
// @Failure 404 not found
// @router /:id [get]
func (c *AdminAccountController) Get() {

}

// Post
// @Title Post
// @Description	Создание администратором нового аккаунта
// @Security	api_key
// @Param	body	body	requests.AdminAccountWriteRequest	"account info"
// @Success	200	{object}	models.Account	Указанный объект может быть получен по ключу data
// @Failure 400 body is invalid
// @Failure 401	unauthorized
// @router / [post]
func (c *AdminAccountController) Post() {

}

// Put
// @Title Put
// @Description	Изменение администратором аккаунта по id
// @Security	api_key
// @Param	id	path 	int64	true	"id"
// @Param	body	body	requests.AdminAccountWriteRequest	"account info"
// @Success	200	{object}	models.Account	Указанный объект может быть получен по ключу data
// @Failure 400 body is invalid
// @Failure 401	unauthorized
// @Failure 404 not found
// @router /:id [put]
func (c *AdminAccountController) Put() {

}

// Delete
// @Title Delete
// @Description	Удаление аккаунта по id
// @Security	api_key
// @Param	id	path 	int64	true	"id"
// @Success	200	{"message": "string"}
// @Failure 401	unauthorized
// @Failure 404 not found
// @router /:id [delete]
func (c *AdminAccountController) Delete() {

}

func (c *AdminAccountController) findModel(id int64) *models.Account {
	m := &models.Account{Id: id}
	if err := models.Get(m); err != nil {
		c.ResponseError(controllers.ErrorNotFound, 404)
		return nil
	}
	return m
}
