package http

import (
	"net/http/pprof"

	"gin-template/internal/middleware"
	"github.com/gin-gonic/gin"
)

func registerPprof(engine *gin.Engine) {
	debugGroup := engine.Group("/debug/pprof")
	debugGroup.Use(middleware.RequireSessionPageAuth())
	debugGroup.GET("/", gin.WrapF(pprof.Index))
	debugGroup.GET("/cmdline", gin.WrapF(pprof.Cmdline))
	debugGroup.GET("/profile", gin.WrapF(pprof.Profile))
	debugGroup.GET("/symbol", gin.WrapF(pprof.Symbol))
	debugGroup.POST("/symbol", gin.WrapF(pprof.Symbol))
	debugGroup.GET("/trace", gin.WrapF(pprof.Trace))
	debugGroup.GET("/:name", func(c *gin.Context) {
		pprof.Handler(c.Param("name")).ServeHTTP(c.Writer, c.Request)
	})
}
