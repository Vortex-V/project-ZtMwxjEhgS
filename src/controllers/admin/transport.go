package controllers

import (
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
// @Param	id	path 	int64	true	"id"
// @Param	start	query	int	false	Начало выборки
// @Param	count	query	int	false	Размер выборки
// @Param	transportType	query	string	false	Тип транспорта [Car, Bike, Scooter, All]
// @Success 200 {object} responses.AdminTransportResponseCollection	Указанный объект может быть получен по ключу data
// @Failure 401 unauthorized
// @router / [get]
func (c *AdminTransportController) GetAll() {

}

// Get
// @Title Get
// @Description Получение информации о транспортном средстве по id
// @Param	id	path 	int64	true	"id"
// @Success 200 {object} models.Transport	Указанный объект может быть получен по ключу data
// @Failure 400 :id is empty
// @Failure 401 unauthorized
// @router /:id [get]
func (c *AdminTransportController) Get() {

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

}

// Put
// @Title Put
// @Description	Изменение транспортного средства по id
// @Security	api_key
// @Param	id	path 	int64	true	"id"
// @Param	body	body	requests.AdminTransportWriteRequest "transport info"
// @Success 200	{object}	models.Transport	Указанный объект может быть получен по ключу data
// @Failure	400	:id is empty
// @Failure	400	body is invalid
// @Failure 401 unauthorized
// @Failure 404 not found
// @router /:id [put]
func (c *AdminTransportController) Put() {

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

}

func (c *AdminTransportController) findModel(id int64) *models.Transport {
	m := &models.Transport{Id: id}
	if err := models.Get(m); err != nil {
		c.ResponseError(controllers.ErrorNotFound, 404)
		return nil
	}
	return m
}
