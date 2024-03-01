package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZapEchoLogger(options ...zap.Option) *zap.SugaredLogger {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "@timestamp"
	config.EncoderConfig.MessageKey = "message"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	logger, _ := config.Build(options...)
	return logger.Sugar()
}
