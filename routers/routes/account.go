package routes

import (
	"app/routers/routeHelpers"
	"app/src/controllers"
	"github.com/beego/beego/v2/server/web"
)

func Account() *web.Namespace {
	controller := &controllers.AccountController{}
	return web.NewNamespace("/Account").
		Namespace(routeHelpers.RouteWithAuth("/Me", "/", controller, "get:Me")).
		Router("/SignIn", controller, "post:SignIn").
		Router("/SignUp", controller, "post:SignUp").
		Namespace(routeHelpers.RouteWithAuth("/SignOut", "/", controller, "post:SignOut")).
		Namespace(routeHelpers.RouteWithAuth("/Update", "/", controller, "put:Update"))
}
