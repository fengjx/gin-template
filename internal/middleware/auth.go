package middleware

import (
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"gin-template/internal/app/berr"
	"gin-template/internal/app/security"
	"gin-template/internal/app/trace"
	sysrefreshtokenStore "gin-template/internal/store/sysrefreshtoken"
	sysuserStore "gin-template/internal/store/sysuser"
	"gin-template/pkg/errs"
)

const currentUserKey = "current_user"
const refreshCookieName = "refresh_token"

func AuthOptional() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := authenticate(c, false); err != nil {
			berr.Abort(c, berr.ErrInvalidToken.WithError(err).WithDetail(err.Error()))
			return
		}
		c.Next()
	}
}

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := authenticate(c, true); err != nil {
			berr.Abort(c, berr.ErrRequireLogin)
			return
		}
		c.Next()
	}
}

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := authenticate(c, true); err != nil {
			berr.Abort(c, berr.ErrRequireLogin)
			return
		}
		currentUser, _ := CurrentUser(c)
		roles := []string{sysuserStore.RoleAdmin, sysuserStore.RoleRoot}
		if !slices.Contains(roles, currentUser.Role) {
			berr.Abort(c, berr.ErrRequireAdmin)
			return
		}
		c.Next()
	}
}

func RequireSessionPageAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := authenticateByRefreshCookie(c); err != nil {
			target := "/login?redirect=" + url.QueryEscape(c.Request.URL.RequestURI())
			c.Redirect(http.StatusTemporaryRedirect, target)
			c.Abort()
			return
		}
		c.Next()
	}
}

func CurrentUser(c *gin.Context) (*sysuserStore.Model, bool) {
	item, ok := c.Get(currentUserKey)
	if !ok {
		return nil, false
	}
	user, ok := item.(*sysuserStore.Model)
	return user, ok
}

func authenticate(c *gin.Context, required bool) error {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
		if required {
			return errs.WithStack(http.ErrNoCookie)
		}
		return nil
	}

	tokenString := strings.TrimSpace(authHeader[7:])
	claims, err := security.ParseAccessToken(tokenString)
	if err != nil {
		return errs.Wrap(err, "解析 Authorization 访问令牌失败")
	}

	currentUser, err := sysuserStore.ByUID(c.Request.Context(), claims.UID)
	if err != nil {
		return errs.Wrap(err, "加载当前登录用户失败")
	}

	c.Set(currentUserKey, currentUser)
	c.Set("uid", currentUser.UID)
	c.Request = c.Request.WithContext(trace.WithUID(c.Request.Context(), currentUser.UID))
	return nil
}

func authenticateByRefreshCookie(c *gin.Context) error {
	token, err := c.Cookie(refreshCookieName)
	if err != nil || token == "" {
		return errs.WithStack(http.ErrNoCookie)
	}

	record, err := sysrefreshtokenStore.ByTokenHash(c.Request.Context(), security.HashOpaqueToken(token))
	if err != nil {
		return errs.Wrap(err, "通过 refresh cookie 查询会话失败")
	}
	if record.RevokedAt != nil || record.ExpiresAt.Before(time.Now()) {
		return errs.WithStack(http.ErrNoCookie)
	}

	currentUser, err := sysuserStore.ByUID(c.Request.Context(), record.UID)
	if err != nil {
		return errs.Wrap(err, "通过 refresh cookie 加载用户失败")
	}

	c.Set(currentUserKey, currentUser)
	c.Set("uid", currentUser.UID)
	c.Request = c.Request.WithContext(trace.WithUID(c.Request.Context(), currentUser.UID))
	return nil
}
