// Package routers
// @APIVersion 1.0.0
// @Title Simbir.GO
// @Host localhost:8080
package routers

import (
	"app/src/controllers"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/filter/cors"
)

func init() {
	insertCorsFilter()

	var ns *web.Namespace

	// Выполняет Include контроллеров для автогенерации в swagger
	// Видимо он парсит файл и ожидает, что они будут записаны именно внутри init и строго в таком формате.
	// Если попробовать вынести эту этажерку в функцию или записать попроще, то оно будет генерировать пустой path.
	ns = web.NewNamespace("/api",
		web.NSNamespace("/Account",
			web.NSInclude(&controllers.AccountController{}),
		),
	)

	ns.Namespace(account()) // TODO

	web.AddNamespace(ns)
}

func insertCorsFilter() {
	s, _ := web.AppConfig.String("allowOrigins")
	web.InsertFilter("*", web.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{s},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
}

func routeWithAuth(rootpath string, controller web.ControllerInterface, methods string) *web.Namespace {
	return web.NewNamespace(rootpath).
		Filter("before", authFilter).
		Router("/", controller, methods)
}

func account() *web.Namespace {
	controller := &controllers.AccountController{}
	return web.NewNamespace("/Account").
		Namespace(routeWithAuth("/Me", controller, "get:Me")).
		Router("/SignIn", controller, "post:SignIn").
		Router("/SignUp", controller, "post:SignUp").
		Namespace(routeWithAuth("/SignOut", controller, "post:SignOut")).
		Namespace(routeWithAuth("/Update", controller, "put:Update"))
}

func Transport() *web.Namespace {
	controller := &controllers.TransportController{}
	return web.NewNamespace("/Transport").
		Router("/:id:int", controller,
			"get:Get;"+
				"put:Put;"+
				"delete:Delete",
		).
		Router("/", controller, "post:Post")
}

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

func Payment() *web.Namespace {
	controller := &controllers.PaymentController{}
	return web.NewNamespace("/Payment").
		Router("/Hesoyam/:id:int", controller, "post:Hesoyam")
}

func AdminAccount() *web.Namespace {
	controller := &controllers.AccountController{}
	return web.NewNamespace("/Admin/Account").
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
	controller := &controllers.TransportController{}
	return web.NewNamespace("/Admin/Transport").
		Router("/", controller,
			"get:GetAll;"+
				"post:Post").
		Router("/:id:int", controller,
			"get:Get;"+
				"put:Put;"+
				"delete:Delete")
}

func AdminRent() *web.Namespace {
	controller := &controllers.RentController{}
	return web.NewNamespace("/Admin").
		Router("/Rent/:id:int", controller,
			"get:Get;"+
				"put:Put;"+
				"delete:PutRent").
		Router("/UserHistory/:id:int", controller, "get:UserHistory").
		Router("/TransportHistory/:id:int", controller, "get:TransportHistory").
		Router("/Rent", controller, "post:PostRent").
		Router("/Rent/End/:id:int", controller, "post:RentEnd")
}
