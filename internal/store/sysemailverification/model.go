package sysemailverification

import (
	"context"
	"time"

	"gin-template/internal/app/db"
	"gin-template/pkg/errs"
	"github.com/google/uuid"
)

const TableName = "sys_email_verifications"

type Model struct {
	ID        string     `gorm:"column:id;primaryKey" json:"id"`
	UID       int64      `gorm:"column:uid" json:"uid"`
	Email     string     `gorm:"column:email" json:"email"`
	TokenHash string     `gorm:"column:token_hash" json:"token_hash"`
	ExpiresAt time.Time  `gorm:"column:expires_at" json:"expires_at"`
	UsedAt    *time.Time `gorm:"column:used_at" json:"used_at"`
	CTime     time.Time  `gorm:"column:ctime;autoCreateTime" json:"ctime"`
	UTime     time.Time  `gorm:"column:utime;autoCreateTime;autoUpdateTime" json:"utime"`
}

func (Model) TableName() string {
	return TableName
}

func New(uid int64, email, tokenHash string, expiresAt time.Time) *Model {
	return &Model{
		ID:        uuid.NewString(),
		UID:       uid,
		Email:     email,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	}
}

func Create(ctx context.Context, item *Model) error {
	if err := db.Get().WithContext(ctx).Create(item).Error; err != nil {
		return errs.Wrap(err, "创建邮箱验证记录失败")
	}
	return nil
}

func ByTokenHash(ctx context.Context, tokenHash string) (*Model, error) {
	var item Model
	err := db.Get().WithContext(ctx).Where("token_hash = ? AND used_at IS NULL", tokenHash).First(&item).Error
	if err != nil {
		return nil, errs.Wrap(err, "按 token hash 查询邮箱验证记录失败")
	}
	return &item, nil
}

func MarkUsed(ctx context.Context, id string) error {
	now := time.Now()
	if err := db.Get().WithContext(ctx).Table(TableName).Where("id = ?", id).Updates(map[string]any{
		"used_at": &now,
		"utime":   now,
	}).Error; err != nil {
		return errs.Wrap(err, "标记邮箱验证记录已使用失败")
	}
	return nil
}
