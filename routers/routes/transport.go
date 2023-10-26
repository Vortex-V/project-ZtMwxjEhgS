package routes

import (
	"app/routers/routeHelpers"
	"app/src/controllers"
	"github.com/beego/beego/v2/server/web"
)

func Transport() *web.Namespace {
	controller := &controllers.TransportController{}
	web.InsertFilter("/api/Transport/", web.BeforeRouter, routeHelpers.AuthFilter)
	return web.NewNamespace("/Transport").
		Router("/:id:int", controller,
			"get:Get;"+
				"put:Put;"+
				"delete:Delete",
		).
		Router("/", controller, "post:Post")
}
