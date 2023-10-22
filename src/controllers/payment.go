package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
)

// PaymentController operations for Transports
type PaymentController struct {
	beego.Controller
}

// Hesoyam Post ...
// @Title Post
// @Description Добавляет на баланс аккаунта с id={accountId} 250 000 денежных единиц.
// @Param	id	path 	int64	true	"accountId"
// @Success 201
// @Failure 401 unauthorized
// @Failure 404 not found
// @router /Hesoyam/:id [post]
func (c *PaymentController) Hesoyam() {}
