package auth

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"gin-template/internal/app/config"
	appHTTP "gin-template/internal/app/http"
	"gin-template/internal/app/security"
	sysoauthbindingStore "gin-template/internal/store/sysoauthbinding"
	sysrefreshtokenStore "gin-template/internal/store/sysrefreshtoken"
	sysuserStore "gin-template/internal/store/sysuser"
	"github.com/gin-gonic/gin"
)

func finishOAuthLogin(c *gin.Context, provider, providerUserID, providerUsername, email, displayName string) {
	var currentUser *sysuserStore.Model
	if binding, err := sysoauthbindingStore.ByProviderUserID(c.Request.Context(), provider, providerUserID); err == nil {
		currentUser, err = sysuserStore.ByUID(c.Request.Context(), binding.UID)
		if err != nil {
			redirectOAuthResult(c, provider, "failed", err.Error())
			return
		}
	} else if existing, err := sysuserStore.ByEmail(c.Request.Context(), strings.ToLower(email)); err == nil {
		currentUser = existing
	} else {
		passwordHash, _ := security.HashPassword(security.NewOpaqueToken())
		currentUser = sysuserStore.New(providerUsername, strings.ToLower(email), passwordHash)
		currentUser.DisplayName = displayName
		currentUser.EmailVerified = true
		if err := sysuserStore.Create(c.Request.Context(), currentUser); err != nil {
			redirectOAuthResult(c, provider, "failed", err.Error())
			return
		}
	}
	if err := sysoauthbindingStore.Upsert(c.Request.Context(), provider, providerUserID, providerUsername, currentUser.UID); err != nil {
		redirectOAuthResult(c, provider, "failed", err.Error())
		return
	}
	setSession(c, currentUser)
	redirectOAuthResult(c, provider, "success", "")
}

func respondWithSession(c *gin.Context, currentUser *sysuserStore.Model) {
	setSession(c, currentUser)
	token, expireAt, _ := security.SignAccessToken(currentUser.UID, currentUser.Role)
	appHTTP.OK(c, authResponse{
		AccessToken: token,
		ExpiresAt:   expireAt.Unix(),
		User:        toUserPayload(currentUser),
	})
}

func setSession(c *gin.Context, currentUser *sysuserStore.Model) {
	cfg := config.Get()
	token := security.NewOpaqueToken()
	expiresAt := time.Now().Add(time.Duration(cfg.Auth.RefreshTokenTTLHours) * time.Hour)
	record := sysrefreshtokenStore.New(currentUser.UID, security.HashOpaqueToken(token), c.Request.UserAgent(), c.ClientIP(), expiresAt)
	_ = sysrefreshtokenStore.Create(c.Request.Context(), record)
	c.SetCookie(refreshCookieName, token, int(time.Until(expiresAt).Seconds()), "/", cfg.Auth.CookieDomain, cfg.Auth.SecureCookie, true)
}

func clearRefreshCookie(c *gin.Context) {
	cfg := config.Get()
	c.SetCookie(refreshCookieName, "", -1, "/", cfg.Auth.CookieDomain, cfg.Auth.SecureCookie, true)
}

func setShortCookie(c *gin.Context, name, value string, ttl time.Duration) {
	cfg := config.Get()
	c.SetCookie(name, value, int(ttl.Seconds()), "/", cfg.Auth.CookieDomain, cfg.Auth.SecureCookie, true)
}

func validateOAuthState(c *gin.Context, cookieName string) bool {
	state, err := c.Cookie(cookieName)
	if err != nil || state == "" || state != c.Query("state") {
		redirectOAuthResult(c, strings.TrimPrefix(cookieName, "oauth_state_"), "failed", "state 校验失败")
		return false
	}
	return true
}

func redirectOAuthResult(c *gin.Context, provider, status, message string) {
	cfg := config.Get()
	target, _ := url.Parse(strings.TrimRight(cfg.Server.FrontendURL, "/") + "/oauth/callback")
	q := target.Query()
	q.Set("provider", provider)
	q.Set("status", status)
	if message != "" {
		q.Set("message", message)
	}
	q.Set("trace_id", c.GetString("trace_id"))
	target.RawQuery = q.Encode()
	c.Redirect(http.StatusTemporaryRedirect, target.String())
}
