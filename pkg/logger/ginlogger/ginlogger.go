package ginlogger

import (
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go-web-template/pkg/logger"
	"time"
)

func Logger() gin.HandlerFunc {
	return ginzap.Ginzap(logger.Log(), time.RFC3339, true)
}

func Recovery() gin.HandlerFunc {
	return ginzap.RecoveryWithZap(logger.Log(), true)
}
