package logging

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger() *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(encoderCfg)
	consoleWriter := zapcore.Lock(os.Stdout)

	level := zapcore.InfoLevel
	consoleCore := zapcore.NewCore(consoleEncoder, consoleWriter, level)
	core := zapcore.NewTee(consoleCore)
	logger := zap.New(core, zap.AddCaller())
	return logger
}
