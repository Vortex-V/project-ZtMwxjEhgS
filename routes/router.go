// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"app/src/controllers"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/filter/cors"
)

func init() {
	s, _ := web.AppConfig.String("allowOrigins")
	web.InsertFilter("*", web.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{s},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	ns := web.NewNamespace("/api").
		Namespace(Account()) // TODO

	web.AddNamespace(ns)
}

func Account() *web.Namespace {
	controller := &controllers.AccountController{}
	return web.NewNamespace("/Account").
		// Router("/Me", controller, "get:Me").
		// Router("/SignIn", controller, "post:SignIn").
		Router("/SignUp", controller, "post:SignUp")
	// Router("/SignOut", controller, "post:SignOut").
	// Router("/Update", controller, "put:Update")
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
