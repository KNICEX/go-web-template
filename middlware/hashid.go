package middlware

import (
	"github.com/gin-gonic/gin"
	"go-web-template/pkg/hashid"
	"go-web-template/pkg/serializer"
	"net/http"
)

const SignIDKey = "object_id"

func HashID(IDType int) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Param("id") != "" {
			id, err := hashid.DecodeId(c.Param("id"), IDType)
			if err == nil {
				c.Set(SignIDKey, id)
				c.Next()
				return
			}
			c.JSON(http.StatusOK, serializer.ErrResponse(serializer.CodeInvalidSign, err))
			c.Abort()
			return
		}
		c.Next()
	}
}
