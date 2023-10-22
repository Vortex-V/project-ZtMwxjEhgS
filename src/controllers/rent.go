package controllers

import (
	"app/src/models"
)

// RentController operations for /Rent
type RentController struct {
	Controller
}

// Transport
// @Title Transport
// @Description Получение транспорта доступного для аренды по параметрам
// @Param	lat	query	float64	false	Географическая широта местонахождения транспорта
// @Param	long	query	float64	false	Географическая долгота местонахождения транспорта
// @Param	radius	query	float64	false	Радиус круга поиска транспорта
// @Param	type	query	string	false	Тип транспорта [Car, Bike, Scooter, All]
// @Success 201 {object} responses.RentTransportResponseCollection
// @Failure 404 not found
// @router /Transport [get]
func (c *RentController) Transport() {

}

// Get
// @Title Get
// @Description Получение информации об аренде по id
// @Param	rentId	path 	int64	true	rentId
// @Success 201 {object} responses.RentGetResponse
// @Failure 401 unauthorized
// @Failure 404 not found
// @router /:rentId [get]
func (c *RentController) Get() {

}

// MyHistory
// @Title MyHistory
// @Description Получение истории аренд текущего аккаунта
// @Success 201 {object} responses.RentHistoryResponseCollection
// @Failure 401 unauthorized
// @router /MyHistory [get]
func (c *RentController) MyHistory() {

}

// TransportHistory
// @Title TransportHistory
// @Description Получение истории аренд транспорта
// @Param	transportId	path 	int64	true	transportId
// @Success 201 {object} responses.RentHistoryResponseCollection
// @Failure 401 unauthorized
// @Failure 404 not found
// @router /TransportHistory/:transportId [get]
func (c *RentController) TransportHistory() {

}

// New
// @Title New
// @Description Аренда транспорта в личное пользование
// @Param	transportId	path 	int64	true	transportId
// @Param	rentType	query 	string	true	Тип аренды [Minutes, Days]
// @Success 201 {object} responses.RentGetResponse
// @Failure 401 unauthorized
// @Failure 404 not found
// @router /New/:transportId [post]
func (c *RentController) New() {

}

// End
// @Title End
// @Description Завершение аренды транспорта по id аренды
// @Param	rentId	path 	int64	true	rentId
// @Param	lat	query	float64	false Географическая широта местонахождения транспорта
// @Param	long	query	float64	false Географическая долгота местонахождения транспорта
// @Success 201 {object}
// @Failure	400	:id is empty
// @Failure 401 unauthorized
// @Failure 404 not found
// @router /End/:rentId [post]
func (c *RentController) End() {

}

func (c *RentController) findModel(id int64) *models.Rent {
	m := &models.Rent{Id: id}
	if err := models.Get(m); err != nil {
		c.ResponseError(ErrorNotFound, 404)
		return nil
	}
	return m
}
