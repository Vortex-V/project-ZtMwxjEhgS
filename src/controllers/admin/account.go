package controllers

import (
	"app/src/components/auth"
	"app/src/components/requests"
	"app/src/components/responses"
	"app/src/controllers"
	"app/src/models"
	"strconv"
)

// AdminAccountController operations for Admin/Account
type AdminAccountController struct {
	controllers.Controller
}

// GetAll
// @Title GetAll
// @Description	Получение списка всех аккаунтов
// @Security	api_key
// @Param	start	query	int	1	false	"Начало выборки [применяет offset((start - 1) * count)]"
// @Param	count	query	int	10	false	"Размер выборки (по умолчанию 10)"
// @Success	200	{object}	models.Account	Список из указанных объектов может быть получен по ключу data
// @Failure 401	unauthorized
// @router / [get]
func (c *AdminAccountController) GetAll() {
	start, _ := c.GetInt("start", 1)
	count, _ := c.GetInt("count", 10)
	query := models.Find(new(models.Account)).
		Offset((start - 1) * count).
		Limit(count)
	collection := make([]responses.AdminAccountResponse, 0)
	rowCount, err := models.Raw(query).QueryRows(&collection)
	if err != nil {
		c.ResponseError(err.Error(), 500)
		return
	}

	// TODO сделать Map func
	for _, response := range collection {
		typeKey, _ := strconv.Atoi(response.Type)
		response.Type = models.GetAccountTypeLabelByKey(typeKey)
	}

	c.Response(collection, controllers.DataMap{
		"count": rowCount,
	})
}

// Get
// @Title Get
// @Description	Получение информации об аккаунте по id
// @Security	api_key
// @Param	id	path 	int64	true	"accountId"
// @Success	200	{object}	responses.AdminAccountResponse	Указанный объект может быть получен по ключу data
// @Failure	400	:id is empty
// @Failure 401	unauthorized
// @Failure 404 not found
// @router /:id [get]
func (c *AdminAccountController) Get() {
	id, _ := c.GetInt64(":id", 0)
	if id == 0 {
		c.ResponseError(controllers.ErrorBadRequest, 400)
		return
	}
	account := c.findModel(id)
	response := &responses.AdminAccountResponse{
		Id:       account.Id,
		Username: account.Username,
		Balance:  account.Balance,
		Type:     account.GetTypeLabel(),
	}
	c.Response(response)
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
	var data requests.AdminAccountWriteRequest
	if !c.ParseAndValidateRequest(&data) {
		return
	}

	account := &models.Account{
		Username: data.Username,
		Balance:  data.Balance,
		Type: func() int {
			if data.IsAdmin {
				return models.AccountTypeAdmin
			}
			return models.AccountTypeUser
		}(),
	}
	account.Password, _ = auth.HashPassword(data.Password)
	_, err := models.Insert(account)
	if err != nil {
		c.ResponseError(err.Error(), 500)
		return
	}
	response := &responses.AdminAccountResponse{
		Id:       account.Id,
		Username: account.Username,
		Balance:  account.Balance,
		Type:     account.GetTypeLabel(),
	}
	c.Response(response, "Аккаунт успешно создан")
}

// Put
// @Title Put
// @Description	Изменение администратором аккаунта по id
// @Security	api_key
// @Param	id	path 	int64	true	"id"
// @Param	body	body	requests.AdminAccountWriteRequest	"account info"
// @Success	200	{object}	responses.AdminAccountResponse	Указанный объект может быть получен по ключу data
// @Failure	400	:id is empty
// @Failure 400 body is invalid
// @Failure 401	unauthorized
// @Failure 404 not found
// @router /:id [put]
func (c *AdminAccountController) Put() {
	id, _ := c.GetInt64(":id", 0)
	if id == 0 {
		c.ResponseError(controllers.ErrorBadRequest, 400)
		return
	}

	var data requests.AdminAccountWriteRequest
	if !c.ParseAndValidateRequest(&data) {
		return
	}

	account := c.findModel(id)

	// TODO Нужно отрефакторить. Это не в контроллере надо делать
	account.Username = data.Username
	account.Balance = data.Balance
	account.Type = func() int {
		if data.IsAdmin {
			return models.AccountTypeAdmin
		}
		return models.AccountTypeUser
	}()
	account.Password, _ = auth.HashPassword(data.Password)

	_, err := models.Update(account)
	if err != nil {
		c.ResponseError(err.Error(), 500)
		return
	}

	response := &responses.AdminAccountResponse{
		Id:       account.Id,
		Username: account.Username,
		Balance:  account.Balance,
		Type:     account.GetTypeLabel(),
	}
	c.Response(response, "Данные обновлены")
}

// Delete
// @Title Delete
// @Description	Удаление аккаунта по id
// @Security	api_key
// @Param	id	path 	int64	true	"id"
// @Success	200	{"message": "string"}
// @Failure	400	:id is empty
// @Failure 401	unauthorized
// @Failure 404 not found
// @router /:id [delete]
func (c *AdminAccountController) Delete() {
	id, _ := c.GetInt64(":id", 0)
	if id == 0 {
		c.ResponseError(controllers.ErrorBadRequest, 400)
		return
	}

	_, err := models.Delete(c.findModel(id))
	if err != nil {
		return
	}

	c.Response("Аккаунт удален")
}

func (c *AdminAccountController) findModel(id int64) *models.Account {
	m := &models.Account{Id: id}
	if err := models.Read(m); err != nil {
		c.ResponseError(controllers.ErrorNotFound, 404)
		return nil
	}
	return m
}
