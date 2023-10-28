package controllers

import "app/src/models"

// PaymentController operations for Transports
type PaymentController struct {
	Controller
}

// Hesoyam Post ...
// @Title Post
// @Description Добавляет на баланс аккаунта с id={accountId} 250 000 денежных единиц.
// @Security	api_key
// @Param	id	path 	int64	true	"accountId"
// @Success 200
// @Failure	400	:id is empty
// @Failure 401 unauthorized
// @Failure 404 not found
// @router /Hesoyam/:id [post]
func (c *PaymentController) Hesoyam() {
	accountId := c.GetIdentityId()
	if accountId == 0 {
		return
	}
	id := c.GetIdFormPath()
	if id == 0 {
		return
	}

	account := &models.Account{Id: id}
	err := models.Read(account)
	if err != nil {
		c.ResponseError(ErrorNotFound, 404)
		return
	}

	if !account.IsAdmin() && accountId != id {
		c.ResponseError(ErrorNotFound, 404)
		return
	}

	account.Balance += 250000
	_, err = models.Update(account)
	if err != nil {
		c.ResponseError(err.Error(), 500)
		return
	}

	c.Response("Bekknqv")
}
