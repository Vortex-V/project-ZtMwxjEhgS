package controllers

import (
	"app/src/controllers"
	"app/src/models"
)

// AdminRentController operations for /Admin/Rent
type AdminRentController struct {
	controllers.Controller
}

// Get
// @Title Get
// @Description Получение информации об аренде по id
// @Param	rentId	path 	int64	true	"rentId"
// @Success 200 {object} models.Rent
// @Failure	400	:id is empty
// @Failure 401 unauthorized
// @Failure 404 not found
// @router /Rent/:rentId [get]
func (c *AdminRentController) Get() {

}

// UserHistory
// @Title UserHistory
// @Description Получение истории аренд пользователя с id={userId}
// @Param	userId	path 	int64	true	userId
// @Success 200 {object} responses.RentHistoryResponseCollection
// @Failure	400	:id is empty
// @Failure 401 unauthorized
// @Failure 404 not found
// @router /UserHistory/:userId [get]
func (c *AdminRentController) UserHistory() {

}

// TransportHistory
// @Title TransportHistory
// @Description Получение истории аренд транспорта с id={transportId}
// @Param	transportId	path 	int64	true	transportId
// @Success 200 {object} responses.RentHistoryResponseCollection
// @Failure	400	:id is empty
// @Failure 401 unauthorized
// @Failure 404 not found
// @router /TransportHistory/:transportId [get]
func (c *AdminRentController) TransportHistory() {

}

// Post
// @Title Post
// @Description Создание новой аренды
// @Param	body	body	requests.AdminRentWriteRequest rent info
// @Success 200	{object}	models.Rent	Указанный объект может быть получен по ключу data
// @Failure	400	body is invalid
// @Failure 401 unauthorized
// @router /Rent [post]
func (c *AdminRentController) Post() {

}

// End
// @Title End
// @Description Завершение аренды транспорта по id аренды
// @Param	rentId	path 	int64	true	rentId
// @Param	lat	query	float64	false	Географическая широта местонахождения транспорта
// @Param	long	query	float64	false	Географическая долгота местонахождения транспорта
// @Success 200	{object}	models.Rent	Указанный объект может быть получен по ключу data
// @Failure	400	:id is empty
// @Failure 401 unauthorized
// @Failure 404 not found
// @router /Rent/End/:rentId [post]
func (c *AdminRentController) End() {

}

// Put
// @Title Put
// @Description Изменение записи об аренде по id
// @Param	id	path 	int64	true	rentId
// @Param	body	body	requests.AdminRentWriteRequest	rent info
// @Success 200	{object}	models.Rent	Указанный объект может быть получен по ключу data
// @Failure	400	body is invalid
// @Failure 401 unauthorized
// @router /Rent/:id [post]
func (c *AdminRentController) Put() {

}

// Delete
// @Title Delete
// @Description Удаление информации об аренде по id
// @Param	id	path 	int64	true	"rentId"
// @Success 201
// @Failure	400	:id is empty
// @Failure 401 unauthorized
// @router /Rent/:id [post]
func (c *AdminRentController) Delete() {

}

func (c *AdminRentController) findModel(id int64) *models.Rent {
	m := &models.Rent{Id: id}
	if err := models.Get(m); err != nil {
		c.ResponseError(controllers.ErrorNotFound, 404)
		return nil
	}
	return m
}