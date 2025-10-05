package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func Init(level string) error {
	var zapLevel zapcore.Level
	err := zapLevel.UnmarshalText([]byte(level))
	if err != nil {
		return err
	}

	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapLevel),
		Development:      false,
		Encoding:         "json",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	if os.Getenv("ENV") == "development" {
		config.Development = true
		config.Encoding = "console"

	}

	logger, err := config.Build()
	if err != nil {
		return err
	}

	log = logger
	return nil
}

func GetLogger() *zap.Logger {
	if log == nil {
		log, _ = zap.NewProduction()
	}
	return log
}

func Info(msg string, fields ...zap.Field) {
	GetLogger().Info(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	GetLogger().Debug(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	GetLogger().Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	GetLogger().Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	GetLogger().Fatal(msg, fields...)
}

func With(fields ...zap.Field) *zap.Logger {
	return GetLogger().With(fields...)
}

func Sync() error {
	return GetLogger().Sync()
}
