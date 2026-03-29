package file

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	sysfileStore "gin-template/internal/store/sysfile"
	"gin-template/pkg/errs"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	errOpenUploadFile   = errors.New("open upload file")
	errCreateUploadDir  = errors.New("create upload dir")
	errCreateTargetFile = errors.New("create target file")
	errWriteTargetFile  = errors.New("write target file")
	errSaveFileRecord   = errors.New("save file record")
)

func saveUploadedFile(c *gin.Context, uid int64, formFile *multipart.FileHeader, storageDir string) (*sysfileStore.Model, error) {
	fileHandle, err := formFile.Open()
	if err != nil {
		return nil, errs.WithStack(fmt.Errorf("%w: %w", errOpenUploadFile, err))
	}
	defer fileHandle.Close()

	if err := os.MkdirAll(storageDir, 0o755); err != nil {
		return nil, errs.WithStack(fmt.Errorf("%w: %w", errCreateUploadDir, err))
	}

	storageName := uuid.NewString() + filepath.Ext(formFile.Filename)
	targetPath := filepath.Join(storageDir, storageName)
	targetFile, err := os.Create(targetPath)
	if err != nil {
		return nil, errs.WithStack(fmt.Errorf("%w: %w", errCreateTargetFile, err))
	}
	defer targetFile.Close()

	if _, err := io.Copy(targetFile, fileHandle); err != nil {
		return nil, errs.WithStack(fmt.Errorf("%w: %w", errWriteTargetFile, err))
	}

	item := sysfileStore.New(uid, formFile.Filename, storageName, formFile.Header.Get("Content-Type"), targetPath, formFile.Size)
	if err := sysfileStore.Create(c.Request.Context(), item); err != nil {
		return nil, errs.WithStack(fmt.Errorf("%w: %w", errSaveFileRecord, err))
	}

	return item, nil
}
