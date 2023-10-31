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
	id := c.GetIdFormPath()
	if id == 0 {
		return
	}
	transport := c.findModel(id)

	response := responses.New[*responses.TransportResponse](
		new(responses.TransportResponse), transport)
	c.Response(response)
}

// GetAll
// @Title GetAll
// @Security	api_key
// @Description Список транспорта, которым владеет пользователь
// @Success 200 {object} responses.TransportResponse	Список указанных объеков может быть получен по ключу data
// @Failure 401 unauthorized
// @Failure 404 not found
// @router / [get]
func (c *TransportController) GetAll() {
	accountId := c.GetIdentityId()
	if accountId == 0 {
		return
	}

	rowCount, list, err := models.TransportSearch(map[string]interface{}{
		"account_id": accountId,
	})
	if err != nil {
		c.ResponseError(err.Error(), 500)
		return
	}

	collection := responses.Collection[*responses.TransportResponse, *models.Transport](
		new(responses.TransportResponse), list)
	c.Response(collection, DataMap{
		"count": rowCount,
	})
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
	id := c.GetIdentityId()
	if id == 0 {
		return
	}

	var data requests.TransportPostRequest
	if !c.ParseAndValidateRequest(&data) {
		return
	}

	transport := &models.Transport{
		Account:     &models.Account{Id: id},
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

	if !transport.SetTransportType(data.TransportType) {
		c.ResponseError("Неверный тип транспорта", 400)
		return
	}

	_, err := models.Insert(transport)
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
	accountId := c.GetIdentityId()
	if accountId == 0 {
		return
	}
	id := c.GetIdFormPath()
	if id == 0 {
		return
	}

	var data requests.TransportPutRequest
	if !c.ParseAndValidateRequest(&data) {
		return
	}

	transport := c.findModel(id)

	if !transport.IsOwner(accountId) {
		c.ResponseError("Нет прав для изменения данных", 403)
		return
	}

	// TODO Нужно отрефакторить. Это не в контроллере надо делать
	transport.CanBeRented = data.CanBeRented
	transport.Model = data.Model
	transport.Color = data.Color
	transport.Identifier = data.Identifier
	transport.Description = data.Description
	transport.Latitude = data.Latitude
	transport.Longitude = data.Longitude
	transport.MinutePrice = data.MinutePrice
	transport.DayPrice = data.DayPrice

	_, err := models.Update(transport)
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
	accountId := c.GetIdentityId()
	if accountId == 0 {
		return
	}
	id := c.GetIdFormPath()
	if id == 0 {
		return
	}

	transport := c.findModel(id)

	if !transport.IsOwner(accountId) {
		c.ResponseError("Нет прав для изменения данных", 403)
		return
	}

	_, err := models.Delete(transport)
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
