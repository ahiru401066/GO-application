package logger

import (
	"os"

	"go.uber.org/zap"
)

var (
	ZapLogger        *zap.Logger
	zapSugaredLogger *zap.SugaredLogger
)

func init() {
	cfg := zap.NewProductionConfig()
	logFile := os.Getenv("APP_LOG_FILE")
	if logFile != "" {
		cfg.OutputPaths = []string{"stderr", logFile}
	}

	ZapLogger = zap.Must(cfg.Build())
	if os.Getenv("APP_ENV") == "development" {
		ZapLogger = zap.Must(zap.NewDevelopment())
	}
	zapSugaredLogger = ZapLogger.Sugar()
}

func Sync() {
	err := zapSugaredLogger.Sync()
	if err != nil {
		zap.Error(err)
	}
}

func Info(msg string, keysAndValues ...interface{}) {
	zapSugaredLogger.Infow(msg, keysAndValues...)
}

func Debug(msg string, keysAndValues ...interface{}) {
	zapSugaredLogger.Debugw(msg, keysAndValues...)
}

func Warn(msg string, keysAndValues ...interface{}) {
	zapSugaredLogger.Warnw(msg, keysAndValues...)
}

func Error(msg string, keysAndValues ...interface{}) {
	zapSugaredLogger.Errorw(msg, keysAndValues...)
}

func Fatal(msg string, keysAndValues ...interface{}) {
	zapSugaredLogger.Fatalw(msg, keysAndValues...)
}

func Panic(msg string, keysAndValues ...interface{}) {
	zapSugaredLogger.Panicw(msg, keysAndValues...)
}
