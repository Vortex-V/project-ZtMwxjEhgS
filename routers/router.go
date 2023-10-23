// Package routers
// @APIVersion 1.0.0
// @Title Simbir.GO
// @Description API сервиса для работы с арендой транспорта
// @Host localhost:8080
// @Schemes http
// @SecurityDefinition	api_key	apiKey	Authorization	header
package routers

import (
	"app/routers/routes"
	"app/src/controllers"
	adminControllers "app/src/controllers/admin"
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
		web.NSNamespace("/Transport",
			web.NSInclude(&controllers.TransportController{}),
		),
		web.NSNamespace("/Rent",
			web.NSInclude(&controllers.RentController{}),
		),
		web.NSNamespace("/Payment",
			web.NSInclude(&controllers.PaymentController{}),
		),
		web.NSNamespace("/Admin/Account",
			web.NSInclude(&adminControllers.AdminAccountController{}),
		),
		web.NSNamespace("/Admin/Transport",
			web.NSInclude(&adminControllers.AdminTransportController{}),
		),
		web.NSNamespace("/Admin/Rent",
			web.NSInclude(&adminControllers.AdminRentController{}),
		),
	)

	ns. // Регистрация маршрутов
		Namespace(routes.Account()).
		Namespace(routes.Transport()).
		Namespace(routes.Rent()).
		Namespace(routes.Payment()).
		Namespace(routes.Admin())

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
