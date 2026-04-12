package sysrefreshtoken

import (
	"context"
	"time"

	"gin-template/internal/app/db"
	"gin-template/pkg/errs"
)

const TableName = "sys_refresh_tokens"

type Model struct {
	ID        int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UID       int64      `gorm:"column:uid" json:"uid"`
	TokenHash string     `gorm:"column:token_hash" json:"token_hash"`
	ExpiresAt time.Time  `gorm:"column:expires_at" json:"expires_at"`
	RevokedAt *time.Time `gorm:"column:revoked_at" json:"revoked_at"`
	UserAgent string     `gorm:"column:user_agent" json:"user_agent"`
	ClientIP  string     `gorm:"column:client_ip" json:"client_ip"`
	CTime     time.Time  `gorm:"column:ctime;autoCreateTime" json:"ctime"`
	UTime     time.Time  `gorm:"column:utime;autoCreateTime;autoUpdateTime" json:"utime"`
}

func (Model) TableName() string {
	return TableName
}

func New(uid int64, tokenHash, userAgent, clientIP string, expiresAt time.Time) *Model {
	return &Model{
		UID:       uid,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
		UserAgent: userAgent,
		ClientIP:  clientIP,
	}
}

func Create(ctx context.Context, item *Model) error {
	if err := db.Get().WithContext(ctx).Create(item).Error; err != nil {
		return errs.Wrap(err, "创建刷新令牌失败")
	}
	return nil
}

func ByTokenHash(ctx context.Context, tokenHash string) (*Model, error) {
	var item Model
	err := db.Get().WithContext(ctx).First(&item, "token_hash = ?", tokenHash).Error
	if err != nil {
		return nil, errs.Wrap(err, "按 token hash 查询刷新令牌失败")
	}
	return &item, nil
}

func Revoke(ctx context.Context, id int64) error {
	now := time.Now()
	if err := db.Get().WithContext(ctx).Table(TableName).Where("id = ?", id).Updates(map[string]any{
		"revoked_at": &now,
		"utime":      now,
	}).Error; err != nil {
		return errs.Wrap(err, "撤销刷新令牌失败")
	}
	return nil
}

func DeleteExpired(ctx context.Context) error {
	if err := db.Get().WithContext(ctx).Where("expires_at < ?", time.Now()).Delete(&Model{}).Error; err != nil {
		return errs.Wrap(err, "删除过期刷新令牌失败")
	}
	return nil
}
