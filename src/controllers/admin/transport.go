package controllers

import (
	"app/src/components/requests"
	"app/src/components/responses"
	"app/src/controllers"
	"app/src/models"
)

// AdminTransportController operations for /Admin/Transport
type AdminTransportController struct {
	controllers.Controller
}

// GetAll
// @Title GetAll
// @Description Получение списка всех транспортных средств
// @Security	api_key
// @Param	start	query	int	1	false	"Начало выборки [применяет offset((start - 1) * count)]"
// @Param	count	query	int	10	false	"Размер выборки (по умолчанию 10)"
// @Param	transportType	query	string	"All"	false	"Тип транспорта [Car, Bike, Scooter, All]"
// @Success	200	{object}	responses.TransportResponse	Список из указанных объектов может быть получен по ключу data
// @Failure 401 unauthorized
// @router / [get]
func (c *AdminTransportController) GetAll() {
	start, _ := c.GetInt("start", 1)
	count, _ := c.GetInt("count", 10)
	transportType := c.GetString("transportType", "All")
	query := models.Find(new(models.Transport)).
		Offset((start - 1) * count).
		Limit(count)
	rawSetter := models.Raw(query)
	if transportType != "All" {
		query.Where("transport_type = ?")
		rawSetter = models.Raw(query).SetArgs(transportType)
	}
	collection := make([]responses.TransportResponse, 0)
	rowCount, err := rawSetter.QueryRows(&collection)
	if err != nil {
		c.ResponseError(err.Error(), 500)
		return
	}

	c.Response(collection, controllers.DataMap{
		"count": rowCount,
	})
}

// Get
// @Title Get
// @Description Получение информации о транспортном средстве по id
// @Security	api_key
// @Param	id	path 	int64	true	"transportId"
// @Success 200 {object}	responses.TransportResponse	"Указанный объект может быть получен по ключу data"
// @Failure 400 :id is empty
// @Failure 401 unauthorized
// @router /:id [get]
func (c *AdminTransportController) Get() {
	id, _ := c.GetInt64(":id", 0)
	if id == 0 {
		c.ResponseError(controllers.ErrorBadRequest, 400)
		return
	}
	transport := c.findModel(id)
	// TODO доработать responses.MapTo, чтобы так не писать
	response := &responses.TransportResponse{
		Id:            transport.Id,
		OwnerId:       transport.AccountId.Id,
		CanBeRented:   transport.CanBeRented,
		TransportType: transport.Type,
		Model:         transport.Model,
		Color:         transport.Color,
		Identifier:    transport.Identifier,
		Description:   transport.Description,
		Latitude:      transport.Latitude,
		Longitude:     transport.Longitude,
		MinutePrice:   transport.MinutePrice,
		DayPrice:      transport.DayPrice,
	}
	c.Response(response)
}

// Post
// @Title Post
// @Description Создание нового транспортного средства
// @Security	api_key
// @Param	body	body	requests.AdminTransportWriteRequest "transport info"
// @Success 200	{object}	models.Transport	Указанный объект может быть получен по ключу data
// @Failure 400 body is invalid
// @Failure 401 unauthorized
// @router / [post]
func (c *AdminTransportController) Post() {
	var data requests.AdminTransportWriteRequest
	if !c.LoadAndValidate(&data) {
		return
	}

	transport := &models.Transport{
		AccountId:   &models.Account{Id: data.OwnerId},
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

	_, err := models.Insert(transport)
	if err != nil {
		c.ResponseError(err.Error(), 500)
		return
	}
	response := &responses.TransportResponse{
		Id:            transport.Id,
		OwnerId:       transport.AccountId.Id,
		CanBeRented:   transport.CanBeRented,
		TransportType: transport.Type,
		Model:         transport.Model,
		Color:         transport.Color,
		Identifier:    transport.Identifier,
		Description:   transport.Description,
		Latitude:      transport.Latitude,
		Longitude:     transport.Longitude,
		MinutePrice:   transport.MinutePrice,
		DayPrice:      transport.DayPrice,
	}
	c.Response(response, "Транспорт добавлен")
}

// Put
// @Title Put
// @Description	Изменение транспортного средства по id
// @Security	api_key
// @Param	id	path 	int64	true	"id"
// @Param	body	body	requests.AdminTransportWriteRequest "transport info"
// @Success 200	{object}	responses.TransportResponse	Указанный объект может быть получен по ключу data
// @Failure	400	:id is empty
// @Failure	400	body is invalid
// @Failure 401 unauthorized
// @Failure 404 not found
// @router /:id [put]
func (c *AdminTransportController) Put() {
	id, _ := c.GetInt64(":id", 0)
	if id == 0 {
		c.ResponseError(controllers.ErrorBadRequest, 400)
		return
	}

	var data requests.AdminTransportWriteRequest
	if !c.LoadAndValidate(&data) {
		return
	}

	transport := &models.Transport{
		Id:          id,
		AccountId:   &models.Account{Id: data.OwnerId},
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

	_, err := models.Update(transport)
	if err != nil {
		c.ResponseError(err.Error(), 500)
		return
	}
	response := &responses.TransportResponse{
		Id:            transport.Id,
		OwnerId:       transport.AccountId.Id,
		CanBeRented:   transport.CanBeRented,
		TransportType: transport.Type,
		Model:         transport.Model,
		Color:         transport.Color,
		Identifier:    transport.Identifier,
		Description:   transport.Description,
		Latitude:      transport.Latitude,
		Longitude:     transport.Longitude,
		MinutePrice:   transport.MinutePrice,
		DayPrice:      transport.DayPrice,
	}
	c.Response(response, "Данные обновлены")
}

// Delete
// @Title Delete
// @Description Удаление транспорта по id
// @Security	api_key
// @Param	id	path 	int64	true	"id"
// @Success 201
// @Failure	400	:id is empty
// @Failure 401	unauthorized
// @router /:id [delete]
func (c *AdminTransportController) Delete() {
	id, _ := c.GetInt64(":id", 0)
	if id == 0 {
		c.ResponseError(controllers.ErrorBadRequest, 400)
		return
	}

	_, err := models.Delete(c.findModel(id))
	if err != nil {
		return
	}

	c.Response("Транспорт удален")
}

func (c *AdminTransportController) findModel(id int64) *models.Transport {
	m := &models.Transport{Id: id}
	if err := models.Get(m); err != nil {
		c.ResponseError(controllers.ErrorNotFound, 404)
		return nil
	}
	return m
}
