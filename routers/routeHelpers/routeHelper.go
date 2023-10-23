package routeHelpers

import "github.com/beego/beego/v2/server/web"

func RouteWithAuth(rootpath string, controller web.ControllerInterface, methods string) *web.Namespace {
	return web.NewNamespace(rootpath).
		Filter("before", AuthFilter).
		Router("/", controller, methods)
}
