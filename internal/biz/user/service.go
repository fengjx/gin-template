package user

import (
	"strconv"
	"strings"

	sysuserStore "gin-template/internal/store/sysuser"
	"gin-template/pkg/errs"
	"github.com/gin-gonic/gin"
)

type pagination struct {
	Limit int
	Page  int
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

func parseUID(value string) (int64, error) {
	uid, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, errs.Wrap(err, "解析用户 ID 失败")
	}
	return uid, nil
}

func normalizeRole(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", sysuserStore.RoleUser:
		return sysuserStore.RoleUser
	case sysuserStore.RoleAdmin:
		return sysuserStore.RoleAdmin
	case sysuserStore.RoleRoot:
		return sysuserStore.RoleRoot
	default:
		return ""
	}
}

func normalizeStatus(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", sysuserStore.StatusActive:
		return sysuserStore.StatusActive
	case sysuserStore.StatusLocked:
		return sysuserStore.StatusLocked
	default:
		return ""
	}
}

func (p pagination) Offset() int {
	return (p.Page - 1) * p.Limit
}
