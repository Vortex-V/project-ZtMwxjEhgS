package routes

import (
	"app/routers/routeHelpers"
	"app/src/controllers"
	"github.com/beego/beego/v2/server/web"
)

func Payment() *web.Namespace {
	controller := &controllers.PaymentController{}
	return web.NewNamespace("/Payment").
		Namespace(routeHelpers.RouteWithAuth("/Hesoyam/:id:int", controller, "post:Hesoyam"))
}
