package sysoption

import (
	"context"
	"time"

	"gin-template/internal/app/db"
	"gin-template/pkg/errs"
)

const TableName = "sys_options"

type Model struct {
	ID          string    `gorm:"column:id;primaryKey" json:"id"`
	OptionKey   string    `gorm:"column:option_key" json:"option_key"`
	OptionValue string    `gorm:"column:option_value" json:"option_value"`
	Description string    `gorm:"column:description" json:"description"`
	IsPublic    bool      `gorm:"column:is_public" json:"is_public"`
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
	return items, nil
}

func GetPublic(ctx context.Context) ([]Model, error) {
	var items []Model
	err := db.Get().WithContext(ctx).Where("is_public = ?", true).Order("option_key ASC").Find(&items).Error
	if err != nil {
		return nil, errs.Wrap(err, "查询公开系统配置失败")
	}
	return items, nil
}

func ByKey(ctx context.Context, key string) (*Model, error) {
	var item Model
	err := db.Get().WithContext(ctx).First(&item, "option_key = ?", key).Error
	if err != nil {
		return nil, errs.Wrap(err, "按 key 查询系统配置失败")
	}
	return &item, nil
}

func Create(ctx context.Context, item *Model) error {
	if err := db.Get().WithContext(ctx).Create(item).Error; err != nil {
		return errs.Wrap(err, "创建系统配置失败")
	}
	return nil
}

func Save(ctx context.Context, item *Model) error {
	if err := db.Get().WithContext(ctx).Save(item).Error; err != nil {
		return errs.Wrap(err, "保存系统配置失败")
	}
	return nil
}
