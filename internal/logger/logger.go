// Package logger - отвечает за логирование в проекте
// Текстовое сообщение должно быть максимально лакончино, передавая суть.
// Если необходимо передать дополнительные данные - отправляем их отдельными полями
// (через zap.String, zap.Int, и т.д.).
// Не стоит добавлять новые поля если по ним не планируется делать отдельный поиск
package logger

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"protocol-adapter/internal/config"
	"protocol-adapter/internal/utils"
)

func zapToSyslog(l zapcore.Level) string {
	switch l {
	case zapcore.WarnLevel:
		return "warning"
	case zapcore.ErrorLevel:
		return "err"
	case zapcore.DPanicLevel:
		return "crit"
	case zapcore.PanicLevel:
		return "crit"
	case zapcore.FatalLevel:
		return "emerg"
	default:
		return l.String()
	}
}

func lowercaseLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(zapToSyslog(l))
}

func specCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	funcName := utils.RemoveBeforeFirst(caller.Function, ".")
	enc.AppendString(fmt.Sprintf("%s:%s", caller.TrimmedPath(), funcName))
}

func NewLogger(meta *config.Meta, conf *config.Logger) *zap.Logger {
	var loggerConfig zap.Config

	if strings.ToLower(conf.Mode) == "prod" {
		loggerConfig = zap.NewProductionConfig()
		loggerConfig.DisableStacktrace = true
		loggerConfig.Sampling = nil
	} else {
		loggerConfig = zap.NewDevelopmentConfig()
	}

	loggerConfig.EncoderConfig.TimeKey = "datetime"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339Nano)

	loggerConfig.EncoderConfig.LevelKey = "severity"
	loggerConfig.EncoderConfig.EncodeLevel = lowercaseLevelEncoder

	loggerConfig.EncoderConfig.CallerKey = "src"
	loggerConfig.EncoderConfig.EncodeCaller = specCallerEncoder

	loggerConfig.EncoderConfig.MessageKey = ""

	logger, err := loggerConfig.Build(
		zap.AddCallerSkip(2),
		zap.Fields(zap.String("svc-version", meta.SvcVersion)),
	)
	if err != nil {
		log.Fatal(err)
	}

	return logger
}

func logLevel(ctx context.Context, level zapcore.Level, log *zap.Logger, msg string, fields ...zap.Field) {
	if ctx == nil {
		ctx = context.Background()
	}
	var metaFields []zap.Field
	bodyFields := append([]zap.Field{zap.Namespace("body"), zap.String("message", msg)}, fields...)

	startTime, ok := GetRequestStartTimeFromCtx(ctx)
	if ok {
		metaFields = append(metaFields, zap.Duration("cutoff", time.Since(startTime)))
	}

	ce := log.Check(level, msg)
	if ce != nil {
		metaFields = append(metaFields, bodyFields...)
		ce.Write(metaFields...)
	}
}

func Debug(ctx context.Context, log *zap.Logger, msg string, fields ...zap.Field) {
	logLevel(ctx, zapcore.DebugLevel, log, msg, fields...)
}

func Info(ctx context.Context, log *zap.Logger, msg string, fields ...zap.Field) {
	logLevel(ctx, zapcore.InfoLevel, log, msg, fields...)
}

func Warn(ctx context.Context, log *zap.Logger, msg string, fields ...zap.Field) {
	logLevel(ctx, zapcore.WarnLevel, log, msg, fields...)
}

func Error(ctx context.Context, log *zap.Logger, msg string, fields ...zap.Field) {
	logLevel(ctx, zapcore.ErrorLevel, log, msg, fields...)
}

func DPanic(ctx context.Context, log *zap.Logger, msg string, fields ...zap.Field) {
	logLevel(ctx, zapcore.DPanicLevel, log, msg, fields...)
}

func Panic(ctx context.Context, log *zap.Logger, msg string, fields ...zap.Field) {
	logLevel(ctx, zapcore.PanicLevel, log, msg, fields...)
}

func Fatal(ctx context.Context, log *zap.Logger, msg string, fields ...zap.Field) {
	logLevel(ctx, zapcore.FatalLevel, log, msg, fields...)
}

func Sync(log *zap.Logger) {
	_ = log.Sync()
}
