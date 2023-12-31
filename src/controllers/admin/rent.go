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
// @router /Rent/:rentId [get]
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
// @router /UserHistory/:userId [get]
func (c *AdminRentController) UserHistory() {
	id := c.GetIdFormPath()
	if id == 0 {
		return
	}

	account := &models.Account{Id: id}
	err := models.Read(account)
	if err != nil {
		c.ResponseError(controllers.ErrorNotFound, 404)
		return
	}

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
// @router /TransportHistory/:transportId [get]
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
		Type:        models.GetRentType(data.PriceType),
		Transport:   &models.Transport{Id: data.TransportId},
		PriceOfUnit: data.PriceOfUnit,
		FinalPrice:  data.FinalPrice,
		Status:      models.RentStatusActive,
	}
	err := rent.SetTimeStart(data.TimeStart)
	if err != nil {
		c.ResponseError("timeStart is invalid (example: 2021-01-01 00:00:00)", 400)
		return
	}
	if data.TimeEnd != "" {
		err = rent.SetTimeEnd(data.TimeEnd)
		if err != nil {
			c.ResponseError("timeEnd is invalid (example: 2021-01-01 00:00:00)", 400)
			return
		}
	}
	err = models.Read(rent.Transport)
	if err != nil {
		c.ResponseError(controllers.ErrorNotFound, 404)
		return
	}
	// Передаём rent.Account.Id, так как создаётся аренда от лица этого пользователя
	/*
		Админ не должен иметь ограничений, поэтому отключил проверку
		if err := rent.CanRent(rent.Account.Id, rent.Transport); err != nil {
			c.ResponseError(err.Error(), 403)
			return
		}*/
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
// @router /Rent/End/:rentId [post]
func (c *AdminRentController) End() {
	rentId := c.GetIdFormPath()
	if rentId == 0 {
		return
	}

	rent := c.findModel(rentId)

	form := new(forms.RentEndForm)
	if !c.ParseAndValidateQuery(form) {
		return
	}
	err := rent.End(form.Lat, form.Long)
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
// @Param	rentId	path 	int64	true	"rentId"
// @Param	body	body	requests.AdminRentWriteRequest	"rent info"
// @Success 200	{object}	responses.RentResponse	Указанный объект может быть получен по ключу data
// @Failure	400	body is invalid
// @Failure 401 unauthorized
// @Failure 404 not found
// @router /Rent/:rentId [put]
func (c *AdminRentController) Put() {
	rentId := c.GetIdFormPath()
	if rentId == 0 {
		return
	}
	data := new(requests.AdminRentWriteRequest)
	if !c.ParseAndValidateRequest(data) {
		return
	}

	rent := c.findModel(rentId)

	/*
		Админ не должен иметь ограничений, поэтому отключил проверку
		if rent.IsOwner(data.UserId) {
			c.ResponseError("нельзя арендовать свой транспорт", 403)
		}*/

	// TODO Нужно отрефакторить. Это не в контроллере надо делать
	rent.Account = &models.Account{Id: data.UserId}
	rent.Type = models.GetRentType(data.PriceType)
	rent.Transport = &models.Transport{Id: data.TransportId}
	rent.PriceOfUnit = data.PriceOfUnit
	rent.FinalPrice = data.FinalPrice
	err := rent.SetTimeStart(data.TimeStart)
	if err != nil {
		c.ResponseError("timeStart is invalid (example: 2021-01-01 00:00:00)", 400)
		return
	}
	if data.TimeEnd != "" {
		err = rent.SetTimeEnd(data.TimeEnd)
		if err != nil {
			c.ResponseError("timeEnd is invalid (example: 2021-01-01 00:00:00)", 400)
			return
		}
	}

	_, err = models.Update(rent)
	if err != nil {
		c.ResponseError(err.Error(), 500)
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
// @Param	rentId	path	int64	true	"rentId"
// @Success 201
// @Failure	400	:id is empty
// @Failure 401	unauthorized
// @Failure 404	not found
// @router /Rent/:rentId [delete]
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
