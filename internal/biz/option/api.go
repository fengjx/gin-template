package option

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"gin-template/internal/app/berr"
	"gin-template/internal/app/config"
	appHTTP "gin-template/internal/app/http"
	"gin-template/internal/app/registry"
	"gin-template/internal/middleware"
	appService "gin-template/internal/service"
	sysoptionStore "gin-template/internal/store/sysoption"
)

func init() {
	registry.RegisterRoute(registerRoutes)
}

func registerRoutes(group *gin.RouterGroup) {
	group.GET("/system/status", status)
	group.GET("/system/about", about)
	group.GET("/system/notice", notice)
	group.GET("/system/pprof-url", middleware.RequireAuth(), pprofURL)

	optionGroup := group.Group("/options")
	optionGroup.GET("", middleware.RequireAdmin(), getOptions)
	optionGroup.PUT("", middleware.RequireRoot(), updateOption)
}

func status(c *gin.Context) {
	appHTTP.OK(c, statusResponse{Status: "ok"})
}

func about(c *gin.Context) {
	value, err := getOptionValue(c, "about")
	if err != nil {
		appHTTP.Abort(c, berr.ErrAboutNotFound.WithError(err))
		return
	}
	appHTTP.OK(c, optionValueResponse{Value: value})
}

func notice(c *gin.Context) {
	value, err := getOptionValue(c, "notice")
	if err != nil {
		appHTTP.Abort(c, berr.ErrNoticeNotFound.WithError(err))
		return
	}
	appHTTP.OK(c, optionValueResponse{Value: value})
}

func pprofURL(c *gin.Context) {
	url, err := buildPprofURL(c, config.Get())
	if err != nil {
		appHTTP.Abort(c, berr.ErrPprofURLUnavailable.WithError(err).WithDetail(err.Error()))
		return
	}
	appHTTP.OK(c, pprofURLResponse{URL: url})
}

func getOptions(c *gin.Context) {
	items, err := sysoptionStore.GetAll(c.Request.Context())
	if err != nil {
		appHTTP.Abort(c, berr.ErrGetOptionsFailed.WithError(err))
		return
	}
	appHTTP.OK(c, itemsToPayload(items))
}

func updateOption(c *gin.Context) {
	var req updateOptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appHTTP.Abort(c, berr.ErrInvalidRequest)
		return
	}
	item, err := appService.UpdateOption(c.Request.Context(), req.Key, appService.UpdateOptionRequest{
		Value:       req.Value,
		Description: req.Description,
		IsPublic:    req.IsPublic,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			appHTTP.Abort(c, berr.ErrOptionNotFound)
			return
		}
		appHTTP.Abort(c, berr.ErrUpdateOptionFailed.WithError(err))
		return
	}
	appHTTP.OK(c, toOptionPayload(item))
}
