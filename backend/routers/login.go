package routers

import (
	"net/http"

	"github.com/fancxxy/gocomic/backend/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// @Summary Login
// @Produce json
// @Param login body routers.Login true "email and password"
// @Success 200 {object} routers.Response
// @Failure 400 {object} routers.Response
// @Failure 401 {object} routers.Response
// @Router /login [post]
func login(c *gin.Context) {
	session := sessions.Default(c)

	var request Login
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusUnauthorized, Response{
			Code: 401,
			Msg:  "email and password are mandatory",
			Data: nil,
		})
		return
	}

	user, err := models.Login(request.Email, request.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, Response{
			Code: 401,
			Msg:  err.Error(),
			Data: nil,
		})
		return
	}

	session.Set("user", user)
	session.Save()

	c.JSON(http.StatusOK, Response{
		Code: 200,
		Msg:  "login succeed",
		Data: nil,
	})
}
