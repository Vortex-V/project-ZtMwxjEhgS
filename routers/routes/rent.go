package routes

import (
	"app/src/controllers"
	"github.com/beego/beego/v2/server/web"
)

func Rent() *web.Namespace {
	controller := &controllers.RentController{}
	return web.NewNamespace("/Rent").
		Router("/:id:int", controller, "get:Get").
		Router("/Transport", controller, "get:Transport").
		Router("/MyHistory", controller, "get:MyHistory").
		Router("/TransportHistory/:id:int", controller, "get:TransportHistory").
		Router("/New/:id:int", controller, "post:New").
		Router("/End/:id:int", controller, "post:End")
}
