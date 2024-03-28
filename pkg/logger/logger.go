package logger

import "go.uber.org/zap"

func L() *zap.SugaredLogger {
	return sugarLogger
}
