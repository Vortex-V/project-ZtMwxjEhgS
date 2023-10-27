package controllers

import (
	"app/src/components/requests"
	"app/src/components/responses"
	"app/src/models"
)

// TransportController operations for /Transport
type TransportController struct {
	Controller
}

// Get
// @Title Get
// @Description Получение информации о транспорте по id
// @Param	id	path 	int64	true	"transportId"
// @Success 200 {object} responses.TransportResponse	Указанный объект может быть получен по ключу data
// @Failure 400 :id is empty
// @Failure 404 not found
// @router /:id [get]
func (c *TransportController) Get() {
	id, _ := c.GetInt64(":id", 0)
	if id == 0 {
		c.ResponseError(ErrorBadRequest, 400)
		return
	}
	transport := c.findModel(id)

	response := responses.New[*responses.TransportResponse](
		new(responses.TransportResponse), transport)
	c.Response(response)
}

// Post
// @Title Post
// @Description Добавление нового транспорта
// @Security	api_key
// @Param	body	body	requests.TransportPostRequest	"transport info"
// @Success 200	{object}	responses.TransportResponse	Указанный объект может быть получен по ключу data
// @Failure 400 body is invalid
// @Failure 401 unauthorized
// @router / [post]
func (c *TransportController) Post() {
	id, err := c.GetInt64("accountId")
	if err != nil {
		c.ResponseError(ErrorBadRequest, 500)
		return
	}

	var data requests.TransportPostRequest
	if !c.LoadAndValidate(&data) {
		return
	}

	transport := &models.Transport{
		Account:     &models.Account{Id: id},
		CanBeRented: data.CanBeRented,
		Type:        data.TransportType,
		Model:       data.Model,
		Color:       data.Color,
		Identifier:  data.Identifier,
		Description: data.Description,
		Latitude:    data.Latitude,
		Longitude:   data.Longitude,
		MinutePrice: data.MinutePrice,
		DayPrice:    data.DayPrice,
	}

	_, err = models.Insert(transport)
	if err != nil {
		c.ResponseError(err.Error(), 500)
		return
	}

	response := responses.New[*responses.TransportResponse](
		new(responses.TransportResponse), transport)
	c.Response(response, "Транспорт добавлен")
}

// Put
// @Title Put
// @Description	Изменение транспорта по id
// @Security	api_key
// @Param	id	path 	int64	true	"transportId"
// @Param	body	body	requests.TransportPutRequest "transport info"
// @Success 200	{object}	responses.TransportResponse	Указанный объект может быть получен по ключу data
// @Failure	400	:id is empty
// @Failure	400	body is invalid
// @Failure 401 unauthorized
// @Failure	403	user is not owner
// @Failure 404 not found
// @router /:id [put]
func (c *TransportController) Put() {
	accountId, err := c.GetInt64("accountId")
	if err != nil {
		c.ResponseError(ErrorBadRequest, 500)
		return
	}
	id, _ := c.GetInt64(":id", 0)
	if id == 0 {
		c.ResponseError(ErrorBadRequest, 400)
		return
	}

	var data requests.TransportPutRequest
	if !c.LoadAndValidate(&data) {
		return
	}

	transport := c.findModel(id)

	if !transport.IsOwner(accountId) {
		c.ResponseError("Нет прав для изменения данных", 403)
		return
	}

	transport = &models.Transport{
		Id:          id,
		CanBeRented: data.CanBeRented,
		Model:       data.Model,
		Color:       data.Color,
		Identifier:  data.Identifier,
		Description: data.Description,
		Latitude:    data.Latitude,
		Longitude:   data.Longitude,
		MinutePrice: data.MinutePrice,
		DayPrice:    data.DayPrice,
	}

	_, err = models.Update(transport)
	if err != nil {
		c.ResponseError(err.Error(), 500)
		return
	}
	response := responses.New[*responses.TransportResponse](
		new(responses.TransportResponse), transport)
	c.Response(response, "Данные обновлены")
}

// Delete
// @Title Delete
// @Description Удаление транспорта по id
// @Security	api_key
// @Param	id	path 	int64	true	"transportId"
// @Success 201
// @Failure	400	:id is empty
// @Failure 401 unauthorized
// @Failure	403	user is not owner
// @Failure 404 not found
// @router /:id [delete]
func (c *TransportController) Delete() {
	accountId, err := c.GetInt64("accountId")
	if err != nil {
		c.ResponseError(ErrorBadRequest, 500)
		return
	}
	id, _ := c.GetInt64(":id", 0)
	if id == 0 {
		c.ResponseError(ErrorBadRequest, 400)
		return
	}

	transport := c.findModel(id)

	if !transport.IsOwner(accountId) {
		c.ResponseError("Нет прав для изменения данных", 403)
		return
	}

	_, err = models.Delete(transport)
	if err != nil {
		return
	}

	c.Response("Транспорт удален")
}

func (c *TransportController) findModel(id int64) *models.Transport {
	m := &models.Transport{Id: id}
	if err := models.Read(m); err != nil {
		c.ResponseError(ErrorNotFound, 404)
		return nil
	}
	return m
}
