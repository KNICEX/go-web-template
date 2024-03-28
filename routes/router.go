package routes

import (
	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	r := gin.New()
	return r
}
