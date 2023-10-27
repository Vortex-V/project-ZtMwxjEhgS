package controllers

import (
	"app/src/components/responses"
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
	lat := c.GetString("lat", "")
	long := c.GetString("long", "")
	radius := c.GetString("radius", "")
	transportType := c.GetString("type", "All")

	rowCount, list, err := models.TransportSearch(map[string]string{
		"type":          models.GetTransportType(transportType),
		"lat":           lat,
		"long":          long,
		"radius":        radius,
		"can_be_rented": "1",
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

// Get
// @Title Get
// @Description Получение информации об аренде по id
// @Security	api_key
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
// @Security	api_key
// @Success 201 {object} responses.RentGetResponse	Список из указанных объектов может быть получен по ключу data
// @Failure 401 unauthorized
// @router /MyHistory [get]
func (c *RentController) MyHistory() {

}

// TransportHistory
// @Title TransportHistory
// @Description Получение истории аренд транспорта
// @Security	api_key
// @Param	transportId	path 	int64	true	transportId
// @Success 201 {object} responses.RentGetResponse	Список из указанных объектов может быть получен по ключу data
// @Failure 401 unauthorized
// @Failure 404 not found
// @router /TransportHistory/:transportId [get]
func (c *RentController) TransportHistory() {

}

// New
// @Title New
// @Description Аренда транспорта в личное пользование
// @Security	api_key
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
// @Security	api_key
// @Param	rentId	path 	int64	true	rentId
// @Param	lat	query	float64	false Географическая широта местонахождения транспорта
// @Param	long	query	float64	false Географическая долгота местонахождения транспорта
// @Success 201
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
