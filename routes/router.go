package routes

import (
	"github.com/gin-gonic/gin"
	"go-web-template/pkg/conf"
)

func Init() *gin.Engine {
	if conf.AppConf.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	base := r.Group(conf.ServerConf.Prefix)
	{
		base.GET("/ping", func(c *gin.Context) {
			c.String(200, "pong")
		})
	}

	return r
}
