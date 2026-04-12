package sysfile

import (
	"context"
	"strings"
	"time"

	"gin-template/pkg/errs"

	"gin-template/internal/app/db"
)

const TableName = "sys_files"

type Model struct {
	ID           int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UID          int64     `gorm:"column:uid" json:"uid"`
	OriginalName string    `gorm:"column:original_name" json:"original_name"`
	StorageName  string    `gorm:"column:storage_name" json:"storage_name"`
	ContentType  string    `gorm:"column:content_type" json:"content_type"`
	Size         int64     `gorm:"column:size" json:"size"`
	Path         string    `gorm:"column:path" json:"path"`
	CTime        time.Time `gorm:"column:ctime;autoCreateTime" json:"ctime"`
	UTime        time.Time `gorm:"column:utime;autoCreateTime;autoUpdateTime" json:"utime"`
}

func (Model) TableName() string {
	return TableName
}

func New(uid int64, originalName, storageName, contentType, path string, size int64) *Model {
	return &Model{
		UID:          uid,
		OriginalName: originalName,
		StorageName:  storageName,
		ContentType:  contentType,
		Size:         size,
		Path:         path,
	}
}

func Create(ctx context.Context, item *Model) error {
	if err := db.Get().WithContext(ctx).Create(item).Error; err != nil {
		return errs.Wrap(err, "创建文件记录失败")
	}
	return nil
}

func ByID(ctx context.Context, id int64) (*Model, error) {
	var item Model
	err := db.Get().WithContext(ctx).First(&item, "id = ?", id).Error
	if err != nil {
		return nil, errs.Wrap(err, "按 ID 查询文件失败")
	}
	return &item, nil
}

func Search(ctx context.Context, keyword string, limit, offset int) ([]Model, int64, error) {
	var items []Model
	var total int64
	query := db.Get().WithContext(ctx).Model(&Model{})
	if trimmed := strings.TrimSpace(keyword); trimmed != "" {
		like := "%" + trimmed + "%"
		query = query.Where("original_name LIKE ? OR content_type LIKE ?", like, like)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errs.Wrap(err, "统计文件列表失败")
	}
	if err := query.Order("ctime DESC").Limit(limit).Offset(offset).Find(&items).Error; err != nil {
		return nil, 0, errs.Wrap(err, "查询文件列表失败")
	}
	return items, total, nil
}

func Delete(ctx context.Context, id int64) error {
	if err := db.Get().WithContext(ctx).Delete(&Model{}, "id = ?", id).Error; err != nil {
		return errs.Wrap(err, "删除文件记录失败")
	}
	return nil
}
