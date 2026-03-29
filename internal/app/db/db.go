package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"

	projectassets "gin-template"
	"gin-template/internal/app/config"
	appLog "gin-template/internal/app/log"
	"gin-template/pkg/errs"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

type SchemaInfo struct {
	SchemaVersion int    `gorm:"column:schema_version"`
	InitializedAt string `gorm:"column:initialized_at"`
}

var (
	once sync.Once
	db   *gorm.DB
	err  error
)

// ResetForTest 重置数据库单例，供测试场景切换临时数据库配置时使用。
// 如果当前连接已经建立，会先主动关闭底层连接，避免测试资源泄露。
func ResetForTest() {
	if db != nil {
		if sqlDB, sqlErr := db.DB(); sqlErr == nil {
			_ = sqlDB.Close()
		}
	}
	once = sync.Once{}
	db = nil
	err = nil
}

func Get() *gorm.DB {
	once.Do(func() {
		err = initDB()
	})
	if err != nil {
		panic(err)
	}
	return db
}

func VerifySchema(ctx context.Context) error {
	gormDB := Get()
	sqlDB, err := gormDB.DB()
	if err != nil {
		return errs.Wrap(err, "获取底层数据库连接失败")
	}

	cfg := config.Get()
	exists, err := tableExists(sqlDB, cfg.Database.Driver, "sys_schema_info")
	if err != nil {
		return errs.Wrap(err, "检查 schema 信息表是否存在失败")
	}
	if !exists {
		return errs.WithStack(errors.New("数据库尚未初始化"))
	}

	var info SchemaInfo
	if err := gormDB.WithContext(ctx).Table("sys_schema_info").First(&info).Error; err != nil {
		return errs.Wrap(err, "读取 schema 信息失败")
	}

	expected := cfg.Database.SchemaVersion
	if info.SchemaVersion < expected {
		return errs.WithStack(fmt.Errorf("schema version %d 低于当前要求 %d，请先手工执行 database/upgrade/%s 下的 SQL", info.SchemaVersion, expected, cfg.Database.Driver))
	}
	return nil
}

func initDB() error {
	cfg := config.Get()
	logger := newGormLogger()

	if cfg.Database.Driver == "mysql" && cfg.Database.DSN == "" {
		cfg.Database.Driver = "sqlite"
	}

	if cfg.Database.Driver == "sqlite" {
		if err := os.MkdirAll(filepath.Dir(cfg.Database.SQLitePath), 0o755); err != nil {
			return errs.Wrap(err, "创建 SQLite 数据目录失败")
		}
		db, err = gorm.Open(sqlite.Open(sqliteDSN(cfg.Database.SQLitePath)), &gorm.Config{Logger: logger})
	} else {
		db, err = gorm.Open(mysql.Open(cfg.Database.DSN), &gorm.Config{Logger: logger})
	}
	if err != nil {
		return errs.Wrap(err, "打开数据库失败")
	}

	if err := bootstrapIfNeeded(db, cfg.Database.Driver); err != nil {
		return errs.Wrap(err, "数据库初始化失败")
	}
	return nil
}

func bootstrapIfNeeded(gormDB *gorm.DB, driver string) error {
	sqlDB, err := gormDB.DB()
	if err != nil {
		return errs.Wrap(err, "获取底层数据库连接失败")
	}
	exists, err := tableExists(sqlDB, driver, "sys_schema_info")
	if err != nil {
		return errs.Wrap(err, "检查数据库初始化状态失败")
	}
	if exists {
		return nil
	}

	sqlBytes, err := projectassets.ReadBootstrapSQL(driver)
	if err != nil {
		return errs.Wrap(err, "读取 bootstrap SQL 失败")
	}

	if _, err := sqlDB.Exec(string(sqlBytes)); err != nil {
		return errs.Wrap(err, "执行 bootstrap SQL 失败")
	}
	appLog.Info("database bootstrapped", zap.String("driver", driver))
	return nil
}

func tableExists(sqlDB *sql.DB, driver, table string) (bool, error) {
	var query string
	switch driver {
	case "mysql":
		query = "SELECT COUNT(*) > 0 FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?"
	default:
		query = "SELECT COUNT(*) > 0 FROM sqlite_master WHERE type='table' AND name = ?"
	}
	var exists bool
	if err := sqlDB.QueryRow(query, table).Scan(&exists); err != nil {
		return false, errs.Wrap(err, "查询数据表是否存在失败")
	}
	return exists, nil
}

func sqliteDSN(path string) string {
	if strings.Contains(path, "_loc=") {
		return path
	}
	separator := "?"
	if strings.Contains(path, "?") {
		separator = "&"
	}
	return path + separator + "_loc=auto"
}

func newGormLogger() gormLogger.Interface {
	return &gormZapLogger{
		logger:                    appLog.Component("db"),
		level:                     gormLogger.Info,
		slowThreshold:             time.Second,
		ignoreRecordNotFoundError: true,
	}
}

type gormZapLogger struct {
	logger                    appLog.Logger
	level                     gormLogger.LogLevel
	slowThreshold             time.Duration
	ignoreRecordNotFoundError bool
}

func (l *gormZapLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	clone := *l
	clone.level = level
	return &clone
}

func (l *gormZapLogger) Info(ctx context.Context, msg string, args ...any) {
	if l.level < gormLogger.Info {
		return
	}
	l.logger.InfoCtx(ctx, "db info", zap.String("message", fmt.Sprintf(msg, args...)))
}

func (l *gormZapLogger) Warn(ctx context.Context, msg string, args ...any) {
	if l.level < gormLogger.Warn {
		return
	}
	l.logger.WarnCtx(ctx, "db warn", zap.String("message", fmt.Sprintf(msg, args...)))
}

func (l *gormZapLogger) Error(ctx context.Context, msg string, args ...any) {
	if l.level < gormLogger.Error {
		return
	}
	l.logger.ErrorCtx(ctx, "db error", zap.String("message", fmt.Sprintf(msg, args...)))
}

func (l *gormZapLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.level == gormLogger.Silent {
		return
	}
	execSQL, rows := fc()
	fields := []zap.Field{
		zap.String("sql", execSQL),
		zap.Int64("rows", rows),
		zap.Duration("cost", time.Since(begin)),
	}
	switch {
	case err != nil && (!l.ignoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)):
		l.logger.ErrorCtx(ctx, "db trace", append(fields, zap.Error(err))...)
	case l.slowThreshold > 0 && time.Since(begin) > l.slowThreshold && l.level >= gormLogger.Warn:
		l.logger.WarnCtx(ctx, "db slow query", fields...)
	case l.level >= gormLogger.Info:
		l.logger.InfoCtx(ctx, "db trace", fields...)
	}
}
