package user

import (
	"strings"

	"github.com/gin-gonic/gin"

	"gin-template/internal/app/berr"
	appHTTP "gin-template/internal/app/http"
	"gin-template/internal/app/registry"
	"gin-template/internal/app/security"
	"gin-template/internal/middleware"
	sysuserStore "gin-template/internal/store/sysuser"
)

func init() {
	registry.RegisterRoute(registerRoutes)
}

func registerRoutes(group *gin.RouterGroup) {
	userGroup := group.Group("/users")
	userGroup.GET("/me", middleware.RequireAuth(), getMe)
	userGroup.PUT("/me", middleware.RequireAuth(), updateMe)
	userGroup.GET("", middleware.RequireAdmin(), listUsers)
	userGroup.GET("/search", middleware.RequireAdmin(), searchUsers)
	userGroup.POST("", middleware.RequireAdmin(), createUser)
	userGroup.GET("/:id", middleware.RequireAdmin(), getUser)
	userGroup.PUT("/:id", middleware.RequireAdmin(), updateUser)
	userGroup.DELETE("/:id", middleware.RequireAdmin(), deleteUser)
}

func getMe(c *gin.Context) {
	currentUser, _ := middleware.CurrentUser(c)
	appHTTP.OK(c, toUserPayload(currentUser))
}

func updateMe(c *gin.Context) {
	currentUser, _ := middleware.CurrentUser(c)
	var req updateMeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appHTTP.Abort(c, berr.ErrInvalidRequest)
		return
	}
	if req.DisplayName != "" {
		currentUser.DisplayName = req.DisplayName
	}
	if err := sysuserStore.Save(c.Request.Context(), currentUser); err != nil {
		appHTTP.Abort(c, berr.ErrUpdateProfileFailed.WithError(err))
		return
	}
	appHTTP.OK(c, toUserPayload(currentUser))
}

func listUsers(c *gin.Context) {
	params := parsePagination(c)
	items, total, err := sysuserStore.Search(c.Request.Context(), "", params.Limit, params.Offset())
	if err != nil {
		appHTTP.Abort(c, berr.ErrQueryUsersFailed.WithError(err))
		return
	}
	appHTTP.OK(c, usersListResponse{Items: itemsToPayload(items), Total: total})
}

func searchUsers(c *gin.Context) {
	params := parsePagination(c)
	items, total, err := sysuserStore.Search(c.Request.Context(), c.Query("q"), params.Limit, params.Offset())
	if err != nil {
		appHTTP.Abort(c, berr.ErrQueryUsersFailed.WithError(err))
		return
	}
	appHTTP.OK(c, usersListResponse{Items: itemsToPayload(items), Total: total})
}

func getUser(c *gin.Context) {
	uid, err := parseUID(c.Param("id"))
	if err != nil {
		appHTTP.Abort(c, berr.ErrInvalidUserID.WithError(err))
		return
	}
	item, err := sysuserStore.ByUID(c.Request.Context(), uid)
	if err != nil {
		appHTTP.Abort(c, berr.ErrUserNotFound.WithError(err))
		return
	}
	appHTTP.OK(c, toUserPayload(item))
}

func createUser(c *gin.Context) {
	currentUser, _ := middleware.CurrentUser(c)
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appHTTP.Abort(c, berr.ErrInvalidRequest)
		return
	}

	username := strings.TrimSpace(req.Username)
	email := strings.ToLower(strings.TrimSpace(req.Email))
	password := strings.TrimSpace(req.Password)
	displayName := strings.TrimSpace(req.DisplayName)
	role := normalizeRole(req.Role)
	status := normalizeStatus(req.Status)

	if username == "" || email == "" || password == "" {
		appHTTP.Abort(c, berr.ErrAuthFieldsRequired)
		return
	}
	if role == "" {
		appHTTP.Abort(c, berr.ErrInvalidRole)
		return
	}
	if status == "" {
		appHTTP.Abort(c, berr.ErrInvalidUserStatus)
		return
	}
	if role == sysuserStore.RoleRoot && currentUser.Role != sysuserStore.RoleRoot {
		appHTTP.Abort(c, berr.ErrOnlyRootCanCreateRootUser)
		return
	}
	if _, err := sysuserStore.ByUsername(c.Request.Context(), username); err == nil {
		appHTTP.Abort(c, berr.ErrUsernameExists)
		return
	} else if !sysuserStore.IsNotFound(err) {
		appHTTP.Abort(c, berr.ErrCheckUsernameFailed.WithError(err))
		return
	}
	if _, err := sysuserStore.ByEmail(c.Request.Context(), email); err == nil {
		appHTTP.Abort(c, berr.ErrEmailExists)
		return
	} else if !sysuserStore.IsNotFound(err) {
		appHTTP.Abort(c, berr.ErrCheckEmailFailed.WithError(err))
		return
	}

	hash, err := security.HashPassword(password)
	if err != nil {
		appHTTP.Abort(c, berr.ErrPasswordProcessFailed.WithError(err))
		return
	}

	item := sysuserStore.New(username, email, hash)
	if displayName != "" {
		item.DisplayName = displayName
	}
	item.Role = role
	item.Status = status
	item.EmailVerified = req.EmailVerified
	if err := sysuserStore.Create(c.Request.Context(), item); err != nil {
		appHTTP.Abort(c, berr.ErrCreateUserFailed.WithError(err))
		return
	}
	appHTTP.OK(c, toUserPayload(item))
}

func updateUser(c *gin.Context) {
	currentUser, _ := middleware.CurrentUser(c)
	uid, err := parseUID(c.Param("id"))
	if err != nil {
		appHTTP.Abort(c, berr.ErrInvalidUserID.WithError(err))
		return
	}

	item, err := sysuserStore.ByUID(c.Request.Context(), uid)
	if err != nil {
		appHTTP.Abort(c, berr.ErrUserNotFound.WithError(err))
		return
	}

	var req updateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appHTTP.Abort(c, berr.ErrInvalidRequest)
		return
	}

	if item.Role == sysuserStore.RoleRoot && currentUser.Role != sysuserStore.RoleRoot {
		appHTTP.Abort(c, berr.ErrOnlyRootCanModifyRootUser)
		return
	}

	if req.Email != "" {
		email := strings.ToLower(strings.TrimSpace(req.Email))
		if email == "" {
			appHTTP.Abort(c, berr.ErrEmailRequired)
			return
		}
		existing, err := sysuserStore.ByEmail(c.Request.Context(), email)
		if err == nil && existing.UID != item.UID {
			appHTTP.Abort(c, berr.ErrEmailExists)
			return
		} else if err != nil && !sysuserStore.IsNotFound(err) {
			appHTTP.Abort(c, berr.ErrCheckEmailFailed.WithError(err))
			return
		}
		item.Email = email
	}

	if req.DisplayName != "" {
		item.DisplayName = strings.TrimSpace(req.DisplayName)
	}
	if req.Role != "" {
		role := normalizeRole(req.Role)
		if role == "" {
			appHTTP.Abort(c, berr.ErrInvalidRole)
			return
		}
		if role == sysuserStore.RoleRoot && currentUser.Role != sysuserStore.RoleRoot {
			appHTTP.Abort(c, berr.ErrOnlyRootCanGrantRootRole)
			return
		}
		item.Role = role
	}
	if req.Status != "" {
		status := normalizeStatus(req.Status)
		if status == "" {
			appHTTP.Abort(c, berr.ErrInvalidUserStatus)
			return
		}
		item.Status = status
	}
	if password := strings.TrimSpace(req.Password); password != "" {
		hash, err := security.HashPassword(password)
		if err != nil {
			appHTTP.Abort(c, berr.ErrPasswordProcessFailed.WithError(err))
			return
		}
		item.PasswordHash = hash
	}
	if req.EmailVerified != nil {
		item.EmailVerified = *req.EmailVerified
	}

	if err := sysuserStore.Save(c.Request.Context(), item); err != nil {
		appHTTP.Abort(c, berr.ErrUpdateUserFailed.WithError(err))
		return
	}
	appHTTP.OK(c, toUserPayload(item))
}

func deleteUser(c *gin.Context) {
	currentUser, _ := middleware.CurrentUser(c)
	uid, err := parseUID(c.Param("id"))
	if err != nil {
		appHTTP.Abort(c, berr.ErrInvalidUserID.WithError(err))
		return
	}
	if uid == currentUser.UID {
		appHTTP.Abort(c, berr.ErrCannotDeleteCurrentUser)
		return
	}

	item, err := sysuserStore.ByUID(c.Request.Context(), uid)
	if err != nil {
		appHTTP.Abort(c, berr.ErrUserNotFound.WithError(err))
		return
	}
	if item.Role == sysuserStore.RoleRoot && currentUser.Role != sysuserStore.RoleRoot {
		appHTTP.Abort(c, berr.ErrOnlyRootCanDeleteRootUser)
		return
	}
	if err := sysuserStore.Delete(c.Request.Context(), uid); err != nil {
		appHTTP.Abort(c, berr.ErrDeleteUserFailed.WithError(err))
		return
	}
	appHTTP.OK(c, messageResponse{Message: "删除成功"})
}
