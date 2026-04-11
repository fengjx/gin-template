package sysoption

import (
	"context"
	"errors"
	"time"

	"gin-template/internal/app/db"
	"gin-template/pkg/errs"
	"gorm.io/gorm"
)

const (
	TableName  = "sys_options"
	TypeString = "string"
	TypeJSON   = "json"

	StatusOnline  = "online"
	StatusOffline = "offline"
)

type Model struct {
	ID          string    `gorm:"column:id;primaryKey" json:"id"`
	OptionKey   string    `gorm:"column:option_key" json:"option_key"`
	OptionValue string    `gorm:"column:option_value" json:"option_value"`
	Description string    `gorm:"column:description" json:"description"`
	IsPublic    bool      `gorm:"column:is_public" json:"is_public"`
	Type        string    `gorm:"column:type" json:"type"`
	Status      string    `gorm:"column:status" json:"status"`
	CTime       time.Time `gorm:"column:ctime;autoCreateTime" json:"ctime"`
	UTime       time.Time `gorm:"column:utime;autoCreateTime;autoUpdateTime" json:"utime"`
}

func (Model) TableName() string {
	return TableName
}

func GetAll(ctx context.Context) ([]Model, error) {
	var items []Model
	err := db.Get().WithContext(ctx).Order("option_key ASC").Find(&items).Error
	if err != nil {
		return nil, errs.Wrap(err, "查询全部系统配置失败")
	}
	normalizeModels(items)
	return items, nil
}

func GetPublic(ctx context.Context) ([]Model, error) {
	var items []Model
	err := db.Get().WithContext(ctx).
		Where("is_public = ? AND status = ?", true, StatusOnline).
		Order("option_key ASC").
		Find(&items).Error
	if err != nil {
		return nil, errs.Wrap(err, "查询公开系统配置失败")
	}
	normalizeModels(items)
	return items, nil
}

func ByKey(ctx context.Context, key string) (*Model, error) {
	var item Model
	err := db.Get().WithContext(ctx).First(&item, "option_key = ?", key).Error
	if err != nil {
		return nil, errs.Wrap(err, "按 key 查询系统配置失败")
	}
	normalizeModel(&item)
	return &item, nil
}

func Create(ctx context.Context, item *Model) error {
	normalizeModel(item)
	if err := db.Get().WithContext(ctx).Create(item).Error; err != nil {
		return errs.Wrap(err, "创建系统配置失败")
	}
	return nil
}

func Save(ctx context.Context, item *Model) error {
	normalizeModel(item)
	if err := db.Get().WithContext(ctx).Save(item).Error; err != nil {
		return errs.Wrap(err, "保存系统配置失败")
	}
	return nil
}

func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func normalizeModels(items []Model) {
	for i := range items {
		normalizeModel(&items[i])
	}
}

func normalizeModel(item *Model) {
	if item == nil {
		return
	}
	if item.Type == "" {
		item.Type = TypeString
	}
	if item.Status == "" {
		item.Status = StatusOnline
	}
}
