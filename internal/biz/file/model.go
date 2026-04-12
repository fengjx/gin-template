package file

import (
	"strconv"

	sysfileStore "gin-template/internal/store/sysfile"
	"gin-template/pkg/timex"

	"github.com/gin-gonic/gin"
)

type pagination struct {
	Limit int
	Page  int
}

type listFilesResponse struct {
	Items []filePayload `json:"items"`
	Total int64         `json:"total"`
}

type filePayload struct {
	ID           int64  `json:"id"`
	UID          int64  `json:"uid"`
	OriginalName string `json:"original_name"`
	StorageName  string `json:"storage_name"`
	ContentType  string `json:"content_type"`
	Size         int64  `json:"size"`
	Path         string `json:"path"`
	CTime        int64  `json:"ctime"`
	UTime        int64  `json:"utime"`
}

type messageResponse struct {
	Message string `json:"message"`
}

func parsePagination(c *gin.Context) pagination {
	limit, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}
	return pagination{
		Limit: limit,
		Page:  page,
	}
}

func (p pagination) Offset() int {
	return (p.Page - 1) * p.Limit
}

func itemsToPayload(items []sysfileStore.Model) []filePayload {
	resp := make([]filePayload, 0, len(items))
	for _, item := range items {
		resp = append(resp, toFilePayload(&item))
	}
	return resp
}

func toFilePayload(item *sysfileStore.Model) filePayload {
	return filePayload{
		ID:           item.ID,
		UID:          item.UID,
		OriginalName: item.OriginalName,
		StorageName:  item.StorageName,
		ContentType:  item.ContentType,
		Size:         item.Size,
		Path:         item.Path,
		CTime:        timex.ToUnixSeconds(item.CTime),
		UTime:        timex.ToUnixSeconds(item.UTime),
	}
}
