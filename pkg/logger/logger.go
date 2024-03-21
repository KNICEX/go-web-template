package logger

import "go.uber.org/zap"

func Log() *zap.Logger {
	return globalLogger
}
func Sugar() *zap.SugaredLogger {
	return sugarLogger
}

func Info(v ...any) {
	sugarLogger.Info(v...)
}

func Debug(v ...any) {
	sugarLogger.Debug(v...)
}

func Error(v ...any) {
	sugarLogger.Error(v...)
}

func Warn(v ...any) {
	sugarLogger.Warn(v...)
}

func Panic(v ...any) {
	sugarLogger.DPanic(v...)
}
