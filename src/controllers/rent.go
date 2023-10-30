package controllers

import (
	"app/src/components/forms"
	"app/src/components/responses"
	"app/src/models"
	"github.com/beego/beego/v2/client/orm/hints"
)

// RentController operations for /Rent
type RentController struct {
	Controller
}

// Transport
// @Title Transport
// @Description Получение транспорта доступного для аренды по параметрам
// @Param	lat	query	float64	true	"Географическая широта местонахождения транспорта"
// @Param	long	query	float64	true	"Географическая долгота местонахождения транспорта"
// @Param	radius	query	float64	true	"Радиус круга поиска транспорта"
// @Param	type	query	string	true	"Тип транспорта [Car, Bike, Scooter, All]"
// @Success 200	{object}	responses.TransportResponse
// @Failure 400 invalid params
// @Failure 404 not found
// @router /Transport [get]
func (c *RentController) Transport() {
	form := new(forms.RentTransportForm)
	if !c.ParseAndValidateQuery(form) {
		return
	}

	rowCount, list, err := models.TransportSearch(map[string]interface{}{
		"type":          models.GetTransportType(form.Type),
		"lat":           form.Lat,
		"long":          form.Long,
		"radius":        form.Radius,
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
// @Param	rentId	path 	int64	true	"rentId"
// @Success 200 {object} responses.RentResponse
// @Failure 400	:id is empty
// @Failure 401 unauthorized
// @Failure 404 not found
// @router /:rentId [get]
func (c *RentController) Get() {
	accountId := c.GetIdentityId()
	if accountId == 0 {
		return
	}
	id := c.GetIdFormPath()
	if id == 0 {
		return
	}

	rent := c.findModel(id)

	if !rent.IsOwner(accountId) || !rent.IsRenter(accountId) {
		c.ResponseError("Нет прав для получения данных", 403)
		return
	}

	response := responses.New[*responses.RentResponse](
		new(responses.RentResponse), rent)
	c.Response(response)
}

// MyHistory
// @Title MyHistory
// @Description Получение истории аренд текущего аккаунта
// @Security	api_key
// @Success 200 {object} responses.RentResponse	Список из указанных объектов может быть получен по ключу data
// @Failure 401 unauthorized
// @router /MyHistory [get]
func (c *RentController) MyHistory() {
	id := c.GetIdentityId()
	if id == 0 {
		return
	}

	account := &models.Account{Id: id}
	rowCount, err := models.LoadRelated(account, "Rents",
		hints.OrderBy("-Id"))
	if err != nil {
		c.ResponseError(err.Error(), 500)
		return
	}

	collection := responses.Collection[*responses.RentResponse, *models.Rent](
		new(responses.RentResponse), account.Rents)
	c.Response(collection, DataMap{
		"count": rowCount,
	})
}

// TransportHistory
// @Title TransportHistory
// @Description Получение истории аренд транспорта
// @Security	api_key
// @Param	transportId	path	int64	true	"transportId"
// @Success 200 {object} responses.RentResponse	Список из указанных объектов может быть получен по ключу data
// @Failure 400 :id is empty
// @Failure 401 unauthorized
// @Failure 404 not found
// @router /TransportHistory/:transportId [get]
func (c *RentController) TransportHistory() {
	accountId := c.GetIdentityId()
	if accountId == 0 {
		return
	}
	transportId := c.GetIdFormPath()
	if transportId == 0 {
		return
	}

	var transport = &models.Transport{Id: transportId}
	err := models.Read(transport)
	if err != nil {
		c.ResponseError(ErrorNotFound, 404)
		return
	}

	if !transport.IsOwner(accountId) {
		c.ResponseError("Нет прав для получения данных", 403)
		return
	}

	rowCount, list, err := models.RentSearch(map[string]interface{}{
		"owner_id":     accountId,
		"transport_id": transportId,
	})
	if err != nil {
		c.ResponseError(err.Error(), 500)
		return
	}

	collection := responses.Collection[*responses.RentResponse, *models.Rent](
		new(responses.RentResponse), list)
	c.Response(collection, DataMap{
		"count": rowCount,
	})
}

// New
// @Title New
// @Description Аренда транспорта в личное пользование
// @Security	api_key
// @Param	transportId	path 	int64	true	"transportId"
// @Param	rentType	query 	string	true	"Тип аренды [Minutes, Days]"
// @Success 200 {object}	responses.RentResponse	Указанный объект может быть получен по ключу data
// @Failure 400 :id is empty
// @Failure 400 invalid params
// @Failure 401 unauthorized
// @Failure 403 forbidden
// @Failure 404 not found
// @router /New/:transportId [post]
func (c *RentController) New() {
	accountId := c.GetIdentityId()
	if accountId == 0 {
		return
	}
	transportId := c.GetIdFormPath()
	if transportId == 0 {
		return
	}

	form := new(forms.RentNewForm)
	if !c.ParseAndValidateQuery(form) {
		return
	}

	rentType := models.GetRentType(form.RentType)
	if rentType == "" {
		c.ResponseError("rentType is invalid", 400)
		return
	}

	rent := &models.Rent{
		Account:     &models.Account{Id: accountId},
		Type:        rentType,
		Transport:   &models.Transport{Id: transportId},
		PriceOfUnit: 0,
		FinalPrice:  0,
	}
	err := models.Read(rent.Transport)
	if err != nil {
		c.ResponseError(ErrorNotFound, 404)
		return
	}
	if rent.IsOwner(rent.Account.Id) {
		c.ResponseError("Нельзя арендовать свой транспорт", 403)
		return
	}
	if !rent.Transport.CanBeRented /*TODO || t.Transport.Status == TransportStatusRented*/ {
		c.ResponseError("Транспорт уже арендован", 403)
		return
	}
	err = rent.Create()
	if err != nil {
		c.ResponseError(err.Error(), 500)
		return
	}

	response := responses.New[*responses.RentResponse](
		new(responses.RentResponse), rent)
	c.Response(response)
}

// End
// @Title End
// @Description Завершение аренды транспорта по id аренды
// @Security	api_key
// @Param	rentId	path 	int64	true	"rentId"
// @Param	lat	query	float64	true "Географическая широта местонахождения транспорта"
// @Param	long	query	float64	true "Географическая долгота местонахождения транспорта"
// @Success 200	{object}	responses.RentResponse	Указанный объект может быть получен по ключу data
// @Failure	400	:id is empty
// @Failure 401 unauthorized
// @Failure 404 not found
// @router /End/:rentId [post]
func (c *RentController) End() {
	accountId, err := c.GetInt64("accountId")
	if err != nil {
		c.ResponseError(err.Error(), 500)
		return
	}
	rentId, _ := c.GetInt64(":id", 0)
	if rentId == 0 {
		c.ResponseError(ErrorBadRequest, 400)
		return
	}

	rent := c.findModel(rentId)
	if !rent.IsRenter(accountId) {
		c.ResponseError("Нет прав для завершения аренды", 403)
	}

	form := new(forms.RentEndForm)
	if !c.ParseAndValidateQuery(form) {
		return
	}
	err = rent.End(map[string]interface{}{
		"lat":  form.Lat,
		"long": form.Long,
	})
	if err != nil {
		c.ResponseError(err.Error(), 500)
		return
	}

	response := responses.New[*responses.RentResponse](
		new(responses.RentResponse), rent)
	c.Response(response, "Аренда завершена")
}

func (c *RentController) findModel(id int64) *models.Rent {
	m := &models.Rent{Id: id}
	if err := models.Read(m); err != nil {
		c.ResponseError(ErrorNotFound, 404)
		return nil
	}
	return m
}
