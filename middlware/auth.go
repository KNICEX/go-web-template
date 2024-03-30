package middlware

import (
	"github.com/gin-gonic/gin"
	"go-web-template/models"
	"go-web-template/pkg/auth"
	"go-web-template/pkg/authstorage"
	"go-web-template/pkg/serializer"
	"net/http"
	"strings"
)

func SignRequired(authInstance auth.Auth) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		switch c.Request.Method {
		case http.MethodPut, http.MethodPost, http.MethodPatch:
		//TODO check request
		default:
			err = auth.CheckURL(authInstance, c.Request.URL)
		}

		if err != nil {
			c.JSON(http.StatusOK, serializer.ErrResponse(serializer.CodeInvalidSign, err))
			c.Abort()
			return
		}
		c.Next()
	}
}

const (
	SessionUserKey = "SessionUser"
)

func LoginRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		// try to get token from header
		token := c.GetHeader("Authorization")
		if token == "" {
			// try to get token from cookie
			token, err = c.Cookie("token")
			if err != nil {
				// try to get token from query
				token = c.Query("token")
				if token == "" {
					c.JSON(http.StatusOK, serializer.ErrResponse(serializer.CodeUnauthorized, nil))
					c.Abort()
					return
				}
			}
		} else {
			token = strings.TrimPrefix(token, "Bearer ")
		}
		authItem, err := authstorage.Get(token)
		if err != nil {
			c.JSON(http.StatusOK, serializer.ErrResponse(serializer.CodeUnauthorized, nil))
			c.Abort()
			return
		}
		c.Set(SessionUserKey, authItem.Value.(models.SessionUser))
		c.Next()
	}
}
