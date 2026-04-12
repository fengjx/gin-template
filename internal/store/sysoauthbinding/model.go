package sysoauthbinding

import (
	"context"
	"time"

	"gin-template/internal/app/db"
	"gin-template/pkg/errs"
)

const TableName = "sys_oauth_bindings"

type Model struct {
	ID               int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UID              int64     `gorm:"column:uid" json:"uid"`
	Provider         string    `gorm:"column:provider" json:"provider"`
	ProviderUserID   string    `gorm:"column:provider_user_id" json:"provider_user_id"`
	ProviderUsername string    `gorm:"column:provider_username" json:"provider_username"`
	CTime            time.Time `gorm:"column:ctime;autoCreateTime" json:"ctime"`
	UTime            time.Time `gorm:"column:utime;autoCreateTime;autoUpdateTime" json:"utime"`
}

func (Model) TableName() string {
	return TableName
}

func Upsert(ctx context.Context, provider, providerUserID, providerUsername string, uid int64) error {
	now := time.Now()
	item := &Model{
		UID:              uid,
		Provider:         provider,
		ProviderUserID:   providerUserID,
		ProviderUsername: providerUsername,
	}
	if err := db.Get().WithContext(ctx).Where("provider = ? AND provider_user_id = ?", provider, providerUserID).Assign(map[string]any{
		"uid":               uid,
		"provider_username": providerUsername,
		"utime":             now,
	}).FirstOrCreate(item).Error; err != nil {
		return errs.Wrap(err, "保存 OAuth 绑定失败")
	}
	return nil
}

func ByProviderUserID(ctx context.Context, provider, providerUserID string) (*Model, error) {
	var item Model
	err := db.Get().WithContext(ctx).First(&item, "provider = ? AND provider_user_id = ?", provider, providerUserID).Error
	if err != nil {
		return nil, errs.Wrap(err, "查询 OAuth 绑定失败")
	}
	return &item, nil
}
