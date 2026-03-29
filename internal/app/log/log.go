package log

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	appEnv "gin-template/internal/app/env"
	"gin-template/internal/app/trace"
	"gin-template/pkg/rt"
)

type loggerSet struct {
	app    *zap.Logger
	access *zap.Logger
}

type encoderFormat string

type Logger struct {
	logger *zap.Logger
}

const (
	encoderFormatConsole encoderFormat = "console"
	encoderFormatJSON    encoderFormat = "json"
)

type encoderProfile struct {
	format     encoderFormat
	colorLevel bool
}

type sinkProfile struct {
	file   *encoderProfile
	stdout *encoderProfile
}

type runtimeConfig struct {
	AppLevel       zapcore.Level
	AppFilename    string
	AccessFilename string
	MaxBackups     int
	MaxAgeDays     int
	Compress       bool
	ConsoleColor   bool
}

var (
	once sync.Once
	set  loggerSet
)

const (
	defaultAppPath    = "runtime/logs/app.log"
	defaultAccessPath = "runtime/logs/access.log"
	defaultMaxBackups = 5
	defaultMaxAgeDays = 7
	defaultCompress   = false
)

func App() Logger {
	return Logger{logger: get().app}
}

func Access() Logger {
	return Logger{logger: get().access}
}

func Component(name string) Logger {
	return Logger{logger: get().app.Named(name)}
}

func Debug(msg string, fields ...zap.Field) {
	App().Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	App().Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	App().Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	App().Error(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	App().Panic(msg, fields...)
}

func DebugCtx(ctx context.Context, msg string, fields ...zap.Field) {
	App().Debug(msg, append(FieldsFromCtx(ctx), fields...)...)
}

func InfoCtx(ctx context.Context, msg string, fields ...zap.Field) {
	App().Info(msg, append(FieldsFromCtx(ctx), fields...)...)
}

func WarnCtx(ctx context.Context, msg string, fields ...zap.Field) {
	App().Warn(msg, append(FieldsFromCtx(ctx), fields...)...)
}

func ErrorCtx(ctx context.Context, msg string, fields ...zap.Field) {
	App().Error(msg, append(FieldsFromCtx(ctx), fields...)...)
}

func PanicCtx(ctx context.Context, msg string, fields ...zap.Field) {
	App().PanicCtx(ctx, msg, append(FieldsFromCtx(ctx), fields...)...)
}

func AccessInfo(msg string, fields ...zap.Field) {
	Access().Info(msg, fields...)
}

func AccessInfoCtx(ctx context.Context, msg string, fields ...zap.Field) {
	Access().Info(msg, append(FieldsFromCtx(ctx), fields...)...)
}

func Sync() {
	loggers := get()
	_ = loggers.app.Sync()
	_ = loggers.access.Sync()
}

func FieldsFromCtx(ctx context.Context) []zap.Field {
	fields := make([]zap.Field, 0, 3)
	if traceID := trace.IDFromCtx(ctx); traceID != "" {
		fields = append(fields, zap.String("trace_id", traceID))
	}
	if uid := trace.UIDFromCtx(ctx); uid > 0 {
		fields = append(fields, zap.Int64("uid", uid))
	}
	return fields
}

func FieldsFromContext(ctx context.Context) []zap.Field {
	return FieldsFromCtx(ctx)
}

func commFields() []zap.Field {
	var fields []zap.Field
	goID := rt.GetGoID()
	if goID > 0 {
		fields = append(fields, zap.Int64("goid", goID))
	}
	return fields
}

func (l Logger) Debug(msg string, fields ...zap.Field) {
	fields = append(fields, commFields()...)
	l.logWithSkip(1).Debug(msg, fields...)
}

func (l Logger) Info(msg string, fields ...zap.Field) {
	fields = append(fields, commFields()...)
	l.logWithSkip(1).Info(msg, fields...)
}

func (l Logger) Warn(msg string, fields ...zap.Field) {
	fields = append(fields, commFields()...)
	l.logWithSkip(1).Warn(msg, fields...)
}

func (l Logger) Error(msg string, fields ...zap.Field) {
	fields = append(fields, commFields()...)
	l.logWithSkip(1).Error(msg, fields...)
}

func (l Logger) Panic(msg string, fields ...zap.Field) {
	fields = append(fields, commFields()...)
	l.logWithSkip(1).Panic(msg, fields...)
}

func (l Logger) DebugCtx(ctx context.Context, msg string, fields ...zap.Field) {
	l.logWithSkip(1).Debug(msg, append(FieldsFromCtx(ctx), fields...)...)
}

func (l Logger) InfoCtx(ctx context.Context, msg string, fields ...zap.Field) {
	l.logWithSkip(1).Info(msg, append(FieldsFromCtx(ctx), fields...)...)
}

func (l Logger) WarnCtx(ctx context.Context, msg string, fields ...zap.Field) {
	l.logWithSkip(1).Warn(msg, append(FieldsFromCtx(ctx), fields...)...)
}

func (l Logger) ErrorCtx(ctx context.Context, msg string, fields ...zap.Field) {
	l.logWithSkip(1).Error(msg, append(FieldsFromCtx(ctx), fields...)...)
}

func (l Logger) PanicCtx(ctx context.Context, msg string, fields ...zap.Field) {
	l.logWithSkip(1).Panic(msg, append(FieldsFromCtx(ctx), fields...)...)
}

func (l Logger) logWithSkip(skip int) *zap.Logger {
	return l.logger.WithOptions(zap.AddCallerSkip(skip))
}

func get() loggerSet {
	once.Do(func() {
		cfg := loadRuntimeConfig()
		appFilename := resolveAppLogFilename(cfg)
		accessFilename := resolveAccessLogFilename(cfg)
		if appFilename != "" {
			_ = os.MkdirAll(filepath.Dir(appFilename), 0o755)
		}
		_ = os.MkdirAll(filepath.Dir(accessFilename), 0o755)

		set.app = buildAppLogger(
			cfg,
			appFilename,
			zapcore.DebugLevel,
		)
		set.access = buildAccessLogger(
			cfg,
			accessFilename,
			zapcore.InfoLevel,
		)
	})
	return set
}

func buildAppLogger(cfg runtimeConfig, filename string, level zapcore.Level) *zap.Logger {
	profile := resolveAppSinkProfile(cfg.ConsoleColor, filename != "")
	levelEnabler := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l >= level && l >= cfg.AppLevel
	})
	cores := make([]zapcore.Core, 0, 2)
	if filename != "" && profile.file != nil {
		fileWriter := zapcore.AddSync(newDailyFileWriter(
			filename,
			cfg.MaxBackups,
			cfg.MaxAgeDays,
			cfg.Compress,
		))
		cores = append(cores, zapcore.NewCore(buildEncoder(*profile.file), fileWriter, levelEnabler))
	}
	if filename == "" && profile.stdout != nil {
		consoleWriter := zapcore.AddSync(os.Stdout)
		cores = append(cores, zapcore.NewCore(buildEncoder(*profile.stdout), consoleWriter, levelEnabler))
	}
	if len(cores) == 0 {
		return zap.NewNop()
	}
	opts := []zap.Option{
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
		zap.AddCallerSkip(2),
	}
	return zap.New(zapcore.NewTee(cores...), opts...)
}

func buildAccessLogger(cfg runtimeConfig, filename string, level zapcore.Level) *zap.Logger {
	profile := encoderProfile{format: encoderFormatJSON}
	levelEnabler := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l >= level
	})
	fileWriter := zapcore.AddSync(newDailyFileWriter(
		filename,
		cfg.MaxBackups,
		cfg.MaxAgeDays,
		cfg.Compress,
	))
	opts := []zap.Option{
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
		zap.AddCallerSkip(2),
	}
	return zap.New(zapcore.NewCore(buildEncoder(profile), fileWriter, levelEnabler), opts...)
}

func resolveAppSinkProfile(consoleColor bool, hasFile bool) sinkProfile {
	if hasFile {
		return sinkProfile{
			file: &encoderProfile{
				format: encoderFormatJSON,
			},
		}
	}
	return sinkProfile{
		stdout: &encoderProfile{
			format:     encoderFormatConsole,
			colorLevel: consoleColor,
		},
	}
}

func resolveAppLogFilename(cfg runtimeConfig) string {
	return strings.TrimSpace(cfg.AppFilename)
}

func resolveAccessLogFilename(cfg runtimeConfig) string {
	return strings.TrimSpace(cfg.AccessFilename)
}

func buildEncoder(profile encoderProfile) zapcore.Encoder {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "time"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeLevel = zapcore.LowercaseLevelEncoder
	if profile.colorLevel {
		encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	switch profile.format {
	case encoderFormatConsole:
		return zapcore.NewConsoleEncoder(encoderCfg)
	default:
		return zapcore.NewJSONEncoder(encoderCfg)
	}
}

func ResetForTest() {
	once = sync.Once{}
	set = loggerSet{}
}

func loadRuntimeConfig() runtimeConfig {
	if appEnv.IsDev() {
		return runtimeConfig{
			AppLevel:       zapcore.DebugLevel,
			AppFilename:    "",
			AccessFilename: defaultAccessPath,
			MaxBackups:     defaultMaxBackups,
			MaxAgeDays:     defaultMaxAgeDays,
			Compress:       defaultCompress,
			ConsoleColor:   true,
		}
	}
	return runtimeConfig{
		AppLevel:       zapcore.InfoLevel,
		AppFilename:    defaultAppPath,
		AccessFilename: defaultAccessPath,
		MaxBackups:     defaultMaxBackups,
		MaxAgeDays:     defaultMaxAgeDays,
		Compress:       defaultCompress,
		ConsoleColor:   false,
	}
}
