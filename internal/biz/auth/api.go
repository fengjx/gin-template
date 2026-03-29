package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"gin-template/internal/app/berr"
	"gin-template/internal/app/config"
	appEnv "gin-template/internal/app/env"
	appHTTP "gin-template/internal/app/http"
	appLog "gin-template/internal/app/log"
	"gin-template/internal/app/oauth"
	"gin-template/internal/app/registry"
	"gin-template/internal/app/security"
	emailModule "gin-template/internal/biz/email"
	turnstileModule "gin-template/internal/biz/turnstile"
	"gin-template/internal/middleware"
	sysemailverificationStore "gin-template/internal/store/sysemailverification"
	syspasswordresetStore "gin-template/internal/store/syspasswordreset"
	sysrefreshtokenStore "gin-template/internal/store/sysrefreshtoken"
	sysuserStore "gin-template/internal/store/sysuser"
)

func init() {
	registry.RegisterRoute(registerRoutes)
}

func registerRoutes(group *gin.RouterGroup) {
	authGroup := group.Group("/auth")
	authGroup.POST("/register", middleware.CriticalRateLimit(), register)
	authGroup.POST("/login", middleware.CriticalRateLimit(), login)
	authGroup.POST("/logout", logout)
	authGroup.POST("/refresh", refresh)
	authGroup.POST("/password/request-reset", middleware.CriticalRateLimit(), requestPasswordReset)
	authGroup.POST("/password/reset", middleware.CriticalRateLimit(), resetPassword)
	authGroup.POST("/email/send-verification", middleware.RequireAuth(), sendEmailVerification)
	authGroup.POST("/email/verify", verifyEmail)
	authGroup.GET("/oauth/github/start", startGitHubOAuth)
	authGroup.GET("/oauth/github/callback", handleGitHubCallback)
	authGroup.GET("/oauth/wechat/start", startWeChatOAuth)
	authGroup.GET("/oauth/wechat/callback", handleWeChatCallback)
}

func register(c *gin.Context) {
	var req authRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appHTTP.Abort(c, berr.ErrInvalidRequest)
		return
	}
	if err := turnstileModule.Verify(c.Request.Context(), req.TurnstileToken, c.ClientIP()); err != nil {
		appHTTP.Abort(c, berr.ErrTurnstileVerifyFailed.WithError(err).WithDetail(err.Error()))
		return
	}
	if req.Username == "" || req.Email == "" || req.Password == "" {
		appHTTP.Abort(c, berr.ErrAuthFieldsRequired)
		return
	}
	if _, err := sysuserStore.ByUsernameOrEmail(c.Request.Context(), req.Username); err == nil {
		appHTTP.Abort(c, berr.ErrUsernameExists)
		return
	}
	if _, err := sysuserStore.ByEmail(c.Request.Context(), strings.ToLower(req.Email)); err == nil {
		appHTTP.Abort(c, berr.ErrEmailExists)
		return
	}

	hash, err := security.HashPassword(req.Password)
	if err != nil {
		appHTTP.Abort(c, berr.ErrPasswordProcessFailed.WithError(err))
		return
	}
	item := sysuserStore.New(req.Username, strings.ToLower(req.Email), hash)
	if req.DisplayName != "" {
		item.DisplayName = req.DisplayName
	}
	if err := sysuserStore.Create(c.Request.Context(), item); err != nil {
		appHTTP.Abort(c, berr.ErrCreateUserFailed.WithError(err))
		return
	}
	respondWithSession(c, item)
}

func login(c *gin.Context) {
	var req authRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appHTTP.Abort(c, berr.ErrInvalidRequest)
		return
	}
	currentUser, err := sysuserStore.ByUsernameOrEmail(c.Request.Context(), req.Identifier)
	if err != nil {
		appHTTP.Abort(c, berr.ErrInvalidCredentials.WithError(err))
		return
	}
	if currentUser.Status != sysuserStore.StatusActive {
		appHTTP.Abort(c, berr.ErrUserDisabled)
		return
	}
	if err := security.ComparePassword(currentUser.PasswordHash, req.Password); err != nil {
		appHTTP.Abort(c, berr.ErrInvalidCredentials.WithError(err))
		return
	}
	appLog.InfoCtx(c, "用户登录", zap.Any("email", currentUser.Email))
	respondWithSession(c, currentUser)
}

func logout(c *gin.Context) {
	if token, err := c.Cookie(refreshCookieName); err == nil && token != "" {
		if item, err := sysrefreshtokenStore.ByTokenHash(c.Request.Context(), security.HashOpaqueToken(token)); err == nil {
			_ = sysrefreshtokenStore.Revoke(c.Request.Context(), item.ID)
		}
	}
	clearRefreshCookie(c)
	appHTTP.OK(c, messageResponse{Message: "已退出登录"})
}

func refresh(c *gin.Context) {
	token, err := c.Cookie(refreshCookieName)
	if err != nil || token == "" {
		c.Status(http.StatusNoContent)
		return
	}
	tokenHash := security.HashOpaqueToken(token)
	record, err := sysrefreshtokenStore.ByTokenHash(c.Request.Context(), tokenHash)
	if err != nil || record.RevokedAt != nil || record.ExpiresAt.Before(time.Now()) {
		clearRefreshCookie(c)
		c.Status(http.StatusNoContent)
		return
	}
	currentUser, err := sysuserStore.ByUID(c.Request.Context(), record.UID)
	if err != nil {
		clearRefreshCookie(c)
		c.Status(http.StatusNoContent)
		return
	}
	_ = sysrefreshtokenStore.Revoke(c.Request.Context(), record.ID)
	respondWithSession(c, currentUser)
}

func requestPasswordReset(c *gin.Context) {
	var req tokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appHTTP.Abort(c, berr.ErrInvalidRequest)
		return
	}
	if err := turnstileModule.Verify(c.Request.Context(), req.TurnstileToken, c.ClientIP()); err != nil {
		appHTTP.Abort(c, berr.ErrTurnstileVerifyFailed.WithError(err).WithDetail(err.Error()))
		return
	}

	currentUser, err := sysuserStore.ByEmail(c.Request.Context(), strings.ToLower(req.Email))
	if err != nil {
		appHTTP.OK(c, messageResponse{Message: "如果邮箱存在，将发送重置邮件"})
		return
	}

	token := security.NewOpaqueToken()
	record := syspasswordresetStore.New(currentUser.UID, currentUser.Email, security.HashOpaqueToken(token), time.Now().Add(30*time.Minute))
	if err := syspasswordresetStore.Create(c.Request.Context(), record); err != nil {
		appHTTP.Abort(c, berr.ErrCreatePasswordResetTokenFailed.WithError(err))
		return
	}

	resp := messageResponse{Message: "如果邮箱存在，将发送重置邮件"}
	if err := emailModule.Send(c.Request.Context(), currentUser.Email, "密码重置", "重置令牌: "+token); errors.Is(err, emailModule.ErrProviderDisabled) && appEnv.IsDev() {
		resp.DebugToken = token
	}
	appHTTP.OK(c, resp)
}

func resetPassword(c *gin.Context) {
	var req tokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appHTTP.Abort(c, berr.ErrInvalidRequest)
		return
	}
	record, err := syspasswordresetStore.ByTokenHash(c.Request.Context(), security.HashOpaqueToken(req.Token))
	if err != nil || record.ExpiresAt.Before(time.Now()) {
		appHTTP.Abort(c, berr.ErrInvalidPasswordResetToken)
		return
	}
	currentUser, err := sysuserStore.ByUID(c.Request.Context(), record.UID)
	if err != nil {
		appHTTP.Abort(c, berr.ErrUserNotFound.WithError(err))
		return
	}
	hash, err := security.HashPassword(req.NewPassword)
	if err != nil {
		appHTTP.Abort(c, berr.ErrPasswordProcessFailed.WithError(err))
		return
	}
	currentUser.PasswordHash = hash
	if err := sysuserStore.Save(c.Request.Context(), currentUser); err != nil {
		appHTTP.Abort(c, berr.ErrPasswordUpdateFailed.WithError(err))
		return
	}
	_ = syspasswordresetStore.MarkUsed(c.Request.Context(), record.ID)
	appHTTP.OK(c, messageResponse{Message: "密码已重置"})
}

func sendEmailVerification(c *gin.Context) {
	currentUser, _ := middleware.CurrentUser(c)
	token := security.NewOpaqueToken()
	record := sysemailverificationStore.New(currentUser.UID, currentUser.Email, security.HashOpaqueToken(token), time.Now().Add(30*time.Minute))
	if err := sysemailverificationStore.Create(c.Request.Context(), record); err != nil {
		appHTTP.Abort(c, berr.ErrCreateEmailVerificationTokenFailed.WithError(err))
		return
	}
	resp := messageResponse{Message: "验证邮件已发送"}
	if err := emailModule.Send(c.Request.Context(), currentUser.Email, "邮箱验证", "验证令牌: "+token); errors.Is(err, emailModule.ErrProviderDisabled) && appEnv.IsDev() {
		resp.DebugToken = token
	}
	appHTTP.OK(c, resp)
}

func verifyEmail(c *gin.Context) {
	var req tokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appHTTP.Abort(c, berr.ErrInvalidRequest)
		return
	}
	record, err := sysemailverificationStore.ByTokenHash(c.Request.Context(), security.HashOpaqueToken(req.Token))
	if err != nil || record.ExpiresAt.Before(time.Now()) {
		appHTTP.Abort(c, berr.ErrInvalidEmailVerificationToken)
		return
	}
	currentUser, err := sysuserStore.ByUID(c.Request.Context(), record.UID)
	if err != nil {
		appHTTP.Abort(c, berr.ErrUserNotFound.WithError(err))
		return
	}
	currentUser.EmailVerified = true
	if err := sysuserStore.Save(c.Request.Context(), currentUser); err != nil {
		appHTTP.Abort(c, berr.ErrUpdateEmailVerificationFailed.WithError(err))
		return
	}
	_ = sysemailverificationStore.MarkUsed(c.Request.Context(), record.ID)
	appHTTP.OK(c, messageResponse{Message: "邮箱验证成功"})
}

func startGitHubOAuth(c *gin.Context) {
	cfg := config.Get()
	if !cfg.OAuth.GitHub.Enabled {
		appHTTP.Abort(c, berr.ErrGitHubOAuthDisabled)
		return
	}
	state := security.NewOpaqueToken()
	setShortCookie(c, "oauth_state_github", state, 10*time.Minute)
	c.Redirect(http.StatusTemporaryRedirect, oauth.GitHubAuthURL(cfg.OAuth.GitHub.ClientID, cfg.OAuth.GitHub.RedirectURL, state))
}

func handleGitHubCallback(c *gin.Context) {
	cfg := config.Get()
	if !validateOAuthState(c, "oauth_state_github") {
		return
	}
	accessToken, err := oauth.ExchangeGitHubToken(c.Request.Context(), cfg.OAuth.GitHub.ClientID, cfg.OAuth.GitHub.ClientSecret, c.Query("code"))
	if err != nil {
		redirectOAuthResult(c, "github", "failed", err.Error())
		return
	}
	ghUser, err := oauth.FetchGitHubUser(c.Request.Context(), accessToken)
	if err != nil {
		redirectOAuthResult(c, "github", "failed", err.Error())
		return
	}
	email := ghUser.Email
	if email == "" {
		email = fmt.Sprintf("%s@users.noreply.github.com", ghUser.Login)
	}
	finishOAuthLogin(c, "github", fmt.Sprintf("%d", ghUser.ID), ghUser.Login, email, ghUser.Name)
}

func startWeChatOAuth(c *gin.Context) {
	cfg := config.Get()
	if !cfg.OAuth.WeChat.Enabled {
		appHTTP.Abort(c, berr.ErrWeChatOAuthDisabled)
		return
	}
	state := security.NewOpaqueToken()
	setShortCookie(c, "oauth_state_wechat", state, 10*time.Minute)
	c.Redirect(http.StatusTemporaryRedirect, oauth.WeChatAuthURL(cfg.OAuth.WeChat.AppID, cfg.OAuth.WeChat.RedirectURL, state))
}

func handleWeChatCallback(c *gin.Context) {
	cfg := config.Get()
	if !validateOAuthState(c, "oauth_state_wechat") {
		return
	}
	accessToken, openID, err := oauth.ExchangeWeChatToken(c.Request.Context(), cfg.OAuth.WeChat.AppID, cfg.OAuth.WeChat.AppSecret, c.Query("code"))
	if err != nil {
		redirectOAuthResult(c, "wechat", "failed", err.Error())
		return
	}
	wxUser, err := oauth.FetchWeChatUser(c.Request.Context(), accessToken, openID)
	if err != nil {
		redirectOAuthResult(c, "wechat", "failed", err.Error())
		return
	}
	email := fmt.Sprintf("wechat_%s@oauth.local", openID)
	finishOAuthLogin(c, "wechat", openID, wxUser.Name, email, wxUser.Name)
}
