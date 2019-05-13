package routers

import (
	"net/http"

	"github.com/fancxxy/gocomic/backend/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// LoginRequired funcion check user login or not
func loginRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		if user == nil {
			c.JSON(http.StatusForbidden, Response{
				Code: 403,
				Msg:  "you must login first",
				Data: nil,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// AdminRequired funcion check user login or not
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		if user == nil || user.(models.User).Admin == false {
			c.JSON(http.StatusForbidden, Response{
				Code: 403,
				Msg:  "not administrator",
				Data: nil,
			})
			c.Abort()
		}
		c.Next()
	}
}
