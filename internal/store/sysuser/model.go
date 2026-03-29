package sysuser

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	"gin-template/pkg/errs"

	"gin-template/internal/app/db"
)

const (
	TableName = "sys_users"
	RoleRoot  = "root"
	RoleAdmin = "admin"
	RoleUser  = "user"

	StatusActive = "active"
	StatusLocked = "locked"
)

type Model struct {
	UID           int64     `gorm:"column:uid;primaryKey;autoIncrement" json:"uid"`
	Username      string    `gorm:"column:username" json:"username"`
	Email         string    `gorm:"column:email" json:"email"`
	PasswordHash  string    `gorm:"column:password_hash" json:"password_hash"`
	Role          string    `gorm:"column:role" json:"role"`
	Status        string    `gorm:"column:status" json:"status"`
	DisplayName   string    `gorm:"column:display_name" json:"display_name"`
	EmailVerified bool      `gorm:"column:email_verified" json:"email_verified"`
	CTime         time.Time `gorm:"column:ctime;autoCreateTime" json:"ctime"`
	UTime         time.Time `gorm:"column:utime;autoCreateTime;autoUpdateTime" json:"utime"`
}

func (Model) TableName() string {
	return TableName
}

func New(username, email, passwordHash string) *Model {
	return &Model{
		Username:      username,
		Email:         email,
		PasswordHash:  passwordHash,
		Role:          RoleUser,
		Status:        StatusActive,
		DisplayName:   username,
		EmailVerified: false,
	}
}

func Create(ctx context.Context, item *Model) error {
	if err := db.Get().WithContext(ctx).Create(item).Error; err != nil {
		return errs.Wrap(err, "创建用户失败")
	}
	return nil
}

func Delete(ctx context.Context, uid int64) error {
	if err := db.Get().WithContext(ctx).Delete(&Model{}, "uid = ?", uid).Error; err != nil {
		return errs.Wrap(err, "删除用户失败")
	}
	return nil
}

func Save(ctx context.Context, item *Model) error {
	if err := db.Get().WithContext(ctx).Save(item).Error; err != nil {
		return errs.Wrap(err, "保存用户失败")
	}
	return nil
}

func ByUID(ctx context.Context, uid int64) (*Model, error) {
	var item Model
	err := db.Get().WithContext(ctx).First(&item, "uid = ?", uid).Error
	if err != nil {
		return nil, errs.Wrap(err, "按 UID 查询用户失败")
	}
	return &item, nil
}

func ByEmail(ctx context.Context, email string) (*Model, error) {
	var item Model
	err := db.Get().WithContext(ctx).First(&item, "email = ?", email).Error
	if err != nil {
		return nil, errs.Wrap(err, "按邮箱查询用户失败")
	}
	return &item, nil
}

func ByUsername(ctx context.Context, username string) (*Model, error) {
	var item Model
	err := db.Get().WithContext(ctx).First(&item, "username = ?", username).Error
	if err != nil {
		return nil, errs.Wrap(err, "按用户名查询用户失败")
	}
	return &item, nil
}

func ByUsernameOrEmail(ctx context.Context, value string) (*Model, error) {
	var item Model
	err := db.Get().WithContext(ctx).Where("username = ? OR email = ?", value, value).First(&item).Error
	if err != nil {
		return nil, errs.Wrap(err, "按用户名或邮箱查询用户失败")
	}
	return &item, nil
}

func Search(ctx context.Context, keyword string, limit, offset int) ([]Model, int64, error) {
	var items []Model
	var total int64
	query := db.Get().WithContext(ctx).Model(&Model{})
	if trimmed := strings.TrimSpace(keyword); trimmed != "" {
		like := "%" + trimmed + "%"
		query = query.Where("username LIKE ? OR email LIKE ? OR display_name LIKE ?", like, like, like)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errs.Wrap(err, "统计用户列表失败")
	}
	if err := query.Order("ctime DESC").Limit(limit).Offset(offset).Find(&items).Error; err != nil {
		return nil, 0, errs.Wrap(err, "查询用户列表失败")
	}
	return items, total, nil
}

func Count(ctx context.Context) (int64, error) {
	var total int64
	if err := db.Get().WithContext(ctx).Model(&Model{}).Count(&total).Error; err != nil {
		return 0, errs.Wrap(err, "统计用户数量失败")
	}
	return total, nil
}

func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
