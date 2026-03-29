package user

import (
	sysuserStore "gin-template/internal/store/sysuser"
	"gin-template/pkg/timex"
)

type updateMeRequest struct {
	DisplayName string `json:"display_name"`
}

type createUserRequest struct {
	Username      string `json:"username"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	DisplayName   string `json:"display_name"`
	Role          string `json:"role"`
	Status        string `json:"status"`
	EmailVerified bool   `json:"email_verified"`
}

type updateUserRequest struct {
	Email         string `json:"email"`
	Password      string `json:"password"`
	DisplayName   string `json:"display_name"`
	Role          string `json:"role"`
	Status        string `json:"status"`
	EmailVerified *bool  `json:"email_verified"`
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

type usersListResponse struct {
	Items []userPayload `json:"items"`
	Total int64         `json:"total"`
}

type messageResponse struct {
	Message string `json:"message"`
}

func itemsToPayload(items []sysuserStore.Model) []userPayload {
	resp := make([]userPayload, 0, len(items))
	for _, item := range items {
		resp = append(resp, toUserPayload(&item))
	}
	return resp
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
