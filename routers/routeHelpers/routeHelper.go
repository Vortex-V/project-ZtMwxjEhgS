package routeHelpers

import (
	"github.com/beego/beego/v2/server/web"
	"net/http"
)

func RouteWithAuth(prefix, route string, controller web.ControllerInterface, methods string) *web.Namespace {
	return web.NewNamespace(prefix).
		Filter("before", AuthFilter).
		Router(route, controller, methods)
}

func ErrorNotFound(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("Not found"))
}

func ErrorInternalServerError(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("Internal server error"))
}
