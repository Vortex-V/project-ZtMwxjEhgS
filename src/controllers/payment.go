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
// @Description create Transports
// @Param	body		body 	models.Transport	true		"body for Transports content"
// @Success 201 {int} models.Transport
// @Failure 403 body is empty
// @router / [post]
func (c *PaymentController) Hesoyam() {
	c.ServeJSON()
}
