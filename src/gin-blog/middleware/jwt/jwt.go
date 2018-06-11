package jwt

import (
	"github.com/gin-gonic/gin"
	"gin-blog/pkg/exception"
	"gin-blog/pkg/util"
	"time"
	"net/http"
)

func JWT() gin.HandlerFunc {
	return func(context *gin.Context) {
		var code int
		var data interface{}

		code = exception.SUCCESS
		token := context.GetHeader("token")

		if token == "" {
			code = exception.INVALID_PARAMS
		} else {
			claims, err := util.ParseToken(token)
			if err != nil {
				code = exception.ERROR_AUTH_CHECK_TOKEN_FAIL
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = exception.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			}
		}

		if code != exception.SUCCESS {
			context.JSON(http.StatusUnauthorized, gin.H{
				"code" : code,
				"msg" : exception.GetMsg(code),
				"data" : data,
			})

			context.Abort()

			return
		}

		context.Next()
	}
}
