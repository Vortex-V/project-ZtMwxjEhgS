package controllers

import (
	"app/src/components/forms"
	"app/src/components/requests"
	"app/src/components/responses"
	"app/src/controllers"
	"app/src/models"
	"github.com/beego/beego/v2/client/orm/hints"
)

// AdminRentController operations for /Admin/Rent
type AdminRentController struct {
	controllers.Controller
}

// Get
// @Title Get
// @Description Получение информации об аренде по id
// @Security	api_key
// @Param	rentId	path 	int64	true	"rentId"
// @Success 200 {object}	responses.RentResponse
// @Failure	400	:id is empty
// @Failure 401 unauthorized
// @Failure 404 not found
// @router /Rent/:id [get]
func (c *AdminRentController) Get() {
	id := c.GetIdFormPath()
	if id == 0 {
		return
	}

	rent := c.findModel(id)

	response := responses.New[*responses.RentResponse](
		new(responses.RentResponse), rent)
	c.Response(response)
}

// UserHistory
// @Title UserHistory
// @Description Получение истории аренд пользователя с id={userId}
// @Security	api_key
// @Param	userId	path 	int64	true	"userId"
// @Success 200	{object}	responses.RentResponse	Список из указанных объектов может быть получен по ключу data
// @Failure	400	:id is empty
// @Failure 401 unauthorized
// @Failure 404 not found
// @router /UserHistory/:id [get]
func (c *AdminRentController) UserHistory() {
	id := c.GetIdFormPath()
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
	c.Response(collection, controllers.DataMap{
		"count": rowCount,
	})
}

// TransportHistory
// @Title	TransportHistory
// @Description Получение истории аренд транспорта с id={transportId}
// @Security	api_key
// @Param	transportId	path	int64	true	"transportId"
// @Success 200 {object}	responses.RentResponse	Список из указанных объектов может быть получен по ключу data
// @Failure	400	:id is empty
// @Failure 401 unauthorized
// @Failure 404 not found
// @router /TransportHistory/:id [get]
func (c *AdminRentController) TransportHistory() {
	transportId := c.GetIdFormPath()
	if transportId == 0 {
		return
	}

	rowCount, list, err := models.RentSearch(map[string]interface{}{
		"transport_id": transportId,
	})
	if err != nil {
		c.ResponseError(err.Error(), 500)
		return
	}

	collection := responses.Collection[*responses.RentResponse, *models.Rent](
		new(responses.RentResponse), list)
	c.Response(collection, controllers.DataMap{
		"count": rowCount,
	})
}

// Post
// @Title	Post
// @Description Создание новой аренды
// @Security	api_key
// @Param	body	body	requests.AdminRentWriteRequest "rent info"
// @Success 200 {object}	responses.RentResponse	Указанный объект может быть получен по ключу data
// @Failure	400	body is invalid
// @Failure 401 unauthorized
// @Failure 403 forbidden
// @Failure 404 not found
// @router /Rent [post]
func (c *AdminRentController) Post() {
	data := new(requests.AdminRentWriteRequest)
	if !c.ParseAndValidateRequest(data) {
		return
	}

	rent := &models.Rent{
		Account:     &models.Account{Id: data.UserId},
		Transport:   &models.Transport{Id: data.TransportId},
		PriceOfUnit: data.PriceOfUnit,
		FinalPrice:  data.FinalPrice,
	}
	if rent.SetType(data.PriceType) {
		c.ResponseError("rentType is invalid", 400)
		return
	}
	err := rent.SetTimeStart(data.TimeStart)
	if err != nil {
		c.ResponseError("timeStart is invalid", 400)
		return
	}
	err = rent.SetTimeEnd(data.TimeEnd)
	if err != nil {
		c.ResponseError("timeEnd is invalid", 400)
		return
	}
	err = models.Read(rent.Transport)
	if err != nil {
		c.ResponseError(controllers.ErrorNotFound, 404)
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
// @Title	End
// @Description Завершение аренды транспорта по id аренды
// @Security	api_key
// @Param	rentId	path 	int64	true	"rentId"
// @Param	lat	query	float64	true	"Географическая широта местонахождения транспорта"
// @Param	long	query	float64	true	"Географическая долгота местонахождения транспорта"
// @Success 200	{object}	responses.RentResponse	Указанный объект может быть получен по ключу data
// @Failure	400	:id is empty
// @Failure 401 unauthorized
// @Failure 404 not found
// @router /Rent/End/:id [post]
func (c *AdminRentController) End() {
	accountId := c.GetIdentityId()
	if accountId == 0 {
		return
	}
	rentId := c.GetIdFormPath()
	if rentId == 0 {
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
	err := rent.End(map[string]interface{}{
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

// Put
// @Title	Put
// @Description Изменение записи об аренде по id
// @Security	api_key
// @Param	id	path 	int64	true	"rentId"
// @Param	body	body	requests.AdminRentWriteRequest	"rent info"
// @Success 200	{object}	responses.RentResponse	Указанный объект может быть получен по ключу data
// @Failure	400	body is invalid
// @Failure 401 unauthorized
// @Failure 404 not found
// @router /Rent/:id [post]
func (c *AdminRentController) Put() {
	rentId := c.GetIdFormPath()
	if rentId == 0 {
		return
	}
	data := new(requests.AdminRentWriteRequest)
	if !c.ParseAndValidateRequest(data) {
		return
	}

	rent := &models.Rent{
		Account:     &models.Account{Id: data.UserId},
		Transport:   &models.Transport{Id: data.TransportId},
		PriceOfUnit: data.PriceOfUnit,
		FinalPrice:  data.FinalPrice,
	}
	if rent.SetType(data.PriceType) {
		c.ResponseError("rentType is invalid", 400)
		return
	}
	err := rent.SetTimeStart(data.TimeStart)
	if err != nil {
		c.ResponseError(err.Error(), 400)
		return
	}
	err = rent.SetTimeEnd(data.TimeEnd)
	if err != nil {
		c.ResponseError(err.Error(), 400)
		return
	}

	_, err = models.Update(rent)
	if err != nil {
		c.ResponseError(controllers.ErrorNotFound, 404)
		return
	}

	response := responses.New[*responses.RentResponse](
		new(responses.RentResponse), rent)
	c.Response(response)
}

// Delete
// @Title	Delete
// @Description	Удаление информации об аренде по id
// @Security	api_key
// @Param	id	path	int64	true	"rentId"
// @Success 201
// @Failure	400	:id is empty
// @Failure 401	unauthorized
// @Failure 404	not found
// @router /Rent/:id [post]
func (c *AdminRentController) Delete() {
	rentId := c.GetIdFormPath()
	if rentId == 0 {
		return
	}

	_, err := models.Delete(&models.Rent{Id: rentId})
	if err != nil {
		c.ResponseError(controllers.ErrorNotFound, 404)
		return
	}

	c.Response(201)
}

func (c *AdminRentController) findModel(id int64) *models.Rent {
	m := &models.Rent{Id: id}
	if err := models.Read(m); err != nil {
		c.ResponseError(controllers.ErrorNotFound, 404)
		return nil
	}
	return m
}
