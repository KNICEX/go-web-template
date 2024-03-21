package ginlogger

import (
	"github.com/gin-gonic/gin"
	"my-bluebell/pkg/conf"
	"my-bluebell/pkg/logger"
	"testing"
)

func TestGinLogger(t *testing.T) {
	conf.Init("../../../config.yaml")
	logger.Init()

	r := gin.New()
	r.Use(Logger(), Recovery())
	r.GET("/ping", func(c *gin.Context) {
		logger.Log().Info("this is a test info")

		c.String(200, "pong")
	})
	r.GET("/panic", func(c *gin.Context) {
		logger.Log().Error("this is a test error")
		panic("this is a test panic")
	})

	r.Run(":8848")
}
