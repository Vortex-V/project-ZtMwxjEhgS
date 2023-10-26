package routes

import (
	"app/routers/routeHelpers"
	"app/src/controllers"
	"github.com/beego/beego/v2/server/web"
)

func Rent() *web.Namespace {
	controller := &controllers.RentController{}
	return web.NewNamespace("/Rent").
		Router("/Transport", controller, "get:Transport").
		Namespace(routeHelpers.RouteWithAuth("/:id:int", controller, "get:Get")).
		Namespace(routeHelpers.RouteWithAuth("/MyHistory", controller, "get:MyHistory")).
		Namespace(routeHelpers.RouteWithAuth("/TransportHistory/:id:int", controller, "get:TransportHistory")).
		Namespace(routeHelpers.RouteWithAuth("/New/:id:int", controller, "post:New")).
		Namespace(routeHelpers.RouteWithAuth("/End/:id:int", controller, "post:End"))
}
