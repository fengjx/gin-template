package file

import (
	"errors"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"

	"gin-template/internal/app/berr"
	"gin-template/internal/app/config"
	appHTTP "gin-template/internal/app/http"
	"gin-template/internal/app/registry"
	"gin-template/internal/middleware"
	sysfileStore "gin-template/internal/store/sysfile"
)

func init() {
	registry.RegisterRoute(registerRoutes)
}

func registerRoutes(group *gin.RouterGroup) {
	fileGroup := group.Group("/files")
	fileGroup.Use(middleware.RequireAdmin())
	fileGroup.GET("", listFiles)
	fileGroup.GET("/search", searchFiles)
	fileGroup.GET("/:id", getFile)
	fileGroup.POST("/upload", middleware.UploadRateLimit(), uploadFile)
	fileGroup.DELETE("/:id", deleteFile)
}

func listFiles(c *gin.Context) {
	params := parsePagination(c)
	items, total, err := sysfileStore.Search(c.Request.Context(), "", params.Limit, params.Offset())
	if err != nil {
		appHTTP.Abort(c, berr.ErrQueryFilesFailed.WithError(err))
		return
	}
	appHTTP.OK(c, listFilesResponse{Items: itemsToPayload(items), Total: total})
}

func searchFiles(c *gin.Context) {
	params := parsePagination(c)
	items, total, err := sysfileStore.Search(c.Request.Context(), c.Query("q"), params.Limit, params.Offset())
	if err != nil {
		appHTTP.Abort(c, berr.ErrQueryFilesFailed.WithError(err))
		return
	}
	appHTTP.OK(c, listFilesResponse{Items: itemsToPayload(items), Total: total})
}

func getFile(c *gin.Context) {
	id, err := parseFileID(c.Param("id"))
	if err != nil {
		appHTTP.Abort(c, berr.ErrInvalidRequest.WithDetail(err.Error()))
		return
	}
	item, err := sysfileStore.ByID(c.Request.Context(), id)
	if err != nil {
		appHTTP.Abort(c, berr.ErrFileNotFound.WithError(err))
		return
	}
	appHTTP.OK(c, toFilePayload(item))
}

func uploadFile(c *gin.Context) {
	currentUser, _ := middleware.CurrentUser(c)
	formFile, err := c.FormFile("file")
	if err != nil {
		appHTTP.Abort(c, berr.ErrMissingUploadFile)
		return
	}

	item, err := saveUploadedFile(c, currentUser.UID, formFile, config.Get().Storage.LocalDir)
	if err != nil {
		switch {
		case errors.Is(err, errOpenUploadFile):
			appHTTP.Abort(c, berr.ErrOpenUploadFileFailed.WithError(err))
		case errors.Is(err, errCreateUploadDir):
			appHTTP.Abort(c, berr.ErrCreateUploadDirFailed.WithError(err))
		case errors.Is(err, errCreateTargetFile):
			appHTTP.Abort(c, berr.ErrCreateTargetFileFailed.WithError(err))
		case errors.Is(err, errWriteTargetFile):
			appHTTP.Abort(c, berr.ErrWriteTargetFileFailed.WithError(err))
		case errors.Is(err, errSaveFileRecord):
			appHTTP.Abort(c, berr.ErrSaveFileRecordFailed.WithError(err))
		default:
			appHTTP.Abort(c, berr.ErrUploadFileFailed.WithError(err))
		}
		return
	}

	appHTTP.OK(c, toFilePayload(item))
}

func deleteFile(c *gin.Context) {
	id, err := parseFileID(c.Param("id"))
	if err != nil {
		appHTTP.Abort(c, berr.ErrInvalidRequest.WithDetail(err.Error()))
		return
	}
	item, err := sysfileStore.ByID(c.Request.Context(), id)
	if err != nil {
		appHTTP.Abort(c, berr.ErrFileNotFound.WithError(err))
		return
	}
	_ = os.Remove(item.Path)
	if err := sysfileStore.Delete(c.Request.Context(), item.ID); err != nil {
		appHTTP.Abort(c, berr.ErrDeleteFileFailed.WithError(err))
		return
	}
	appHTTP.OK(c, messageResponse{Message: fmt.Sprintf("文件 %d 已删除", item.ID)})
}
