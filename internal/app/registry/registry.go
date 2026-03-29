package registry

import "github.com/gin-gonic/gin"

type RouteRegistrar func(group *gin.RouterGroup)

var routeRegistrars []RouteRegistrar

func RegisterRoute(registrar RouteRegistrar) {
	routeRegistrars = append(routeRegistrars, registrar)
}

func Routes() []RouteRegistrar {
	return append([]RouteRegistrar(nil), routeRegistrars...)
}
