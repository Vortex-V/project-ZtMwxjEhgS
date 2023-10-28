package routeHelpers

import "github.com/beego/beego/v2/server/web"

func RouteWithAuth(prefix, route string, controller web.ControllerInterface, methods string) *web.Namespace {
	return web.NewNamespace(prefix).
		Filter("before", AuthFilter).
		Router(route, controller, methods)
}
