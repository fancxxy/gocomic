package routers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// @Summary Logout
// @Produce json
// @Success 200 {object} routers.Response
// @Failure 400 {object} routers.Response
// @Router /logout [get]
func logout(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")
	if user == nil {
		c.JSON(http.StatusBadRequest, Response{
			Code: 400,
			Msg:  "not login",
			Data: nil,
		})
		return
	}

	session.Delete("user")
	session.Save()

	c.JSON(http.StatusOK, Response{
		Code: 200,
		Msg:  "logout succeed",
		Data: nil,
	})
}
