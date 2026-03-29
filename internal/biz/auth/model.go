package auth

import (
	sysuserStore "gin-template/internal/store/sysuser"
	"gin-template/pkg/timex"
)

const refreshCookieName = "refresh_token"

type authRequest struct {
	Username       string `json:"username"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	DisplayName    string `json:"display_name"`
	Identifier     string `json:"identifier"`
	TurnstileToken string `json:"turnstile_token"`
}

type tokenRequest struct {
	Token          string `json:"token"`
	NewPassword    string `json:"new_password"`
	TurnstileToken string `json:"turnstile_token"`
	Email          string `json:"email"`
}

type authResponse struct {
	AccessToken string      `json:"access_token"`
	ExpiresAt   int64       `json:"expires_at"`
	User        userPayload `json:"user"`
}

type userPayload struct {
	UID           int64  `json:"uid"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	DisplayName   string `json:"display_name"`
	Role          string `json:"role"`
	Status        string `json:"status"`
	EmailVerified bool   `json:"email_verified"`
	CTime         int64  `json:"ctime"`
	UTime         int64  `json:"utime"`
}

type messageResponse struct {
	Message    string `json:"message"`
	DebugToken string `json:"debug_token,omitempty"`
}

func toUserPayload(item *sysuserStore.Model) userPayload {
	return userPayload{
		UID:           item.UID,
		Username:      item.Username,
		Email:         item.Email,
		DisplayName:   item.DisplayName,
		Role:          item.Role,
		Status:        item.Status,
		EmailVerified: item.EmailVerified,
		CTime:         timex.ToUnixSeconds(item.CTime),
		UTime:         timex.ToUnixSeconds(item.UTime),
	}
}
