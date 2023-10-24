package routes

import (
	"app/routers/routeHelpers"
	adminControllers "app/src/controllers/admin"
	"github.com/beego/beego/v2/server/web"
)

func Admin() *web.Namespace {

	ns := web.NewNamespace("/Admin").
		Filter("before", routeHelpers.AuthFilter). // TODO фильтр админа
		Namespace(AdminAccount()).
		Namespace(AdminTransport())
	AdminRent(ns)
	return ns
}

func AdminAccount() *web.Namespace {
	controller := &adminControllers.AdminAccountController{}
	return web.NewNamespace("/Account").
		Router("/", controller,
			"get:GetAll;"+
				"post:Post",
		).
		Router("/:id:int", controller,
			"get:Get;"+
				"put:Put;"+
				"delete:Delete")
}

func AdminTransport() *web.Namespace {
	controller := &adminControllers.AdminTransportController{}
	return web.NewNamespace("/Transport").
		Router("/", controller,
			"get:GetAll;"+
				"post:Post").
		Router("/:id:int", controller,
			"get:Get;"+
				"put:Put;"+
				"delete:Delete")
}

func AdminRent(ns *web.Namespace) {
	controller := &adminControllers.AdminRentController{}
	ns.
		Router("/Rent/:id:int", controller,
			"get:Get;"+
				"put:Put;"+
				"delete:Delete").
		Router("/UserHistory/:id:int", controller, "get:UserHistory").
		Router("/TransportHistory/:id:int", controller, "get:TransportHistory").
		Router("/Rent", controller, "post:Post").
		Router("/Rent/End/:id:int", controller, "post:End")
}
