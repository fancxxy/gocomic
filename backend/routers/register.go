package routers

import (
	"net/http"

	"github.com/fancxxy/gocomic/backend/models"
	"github.com/gin-gonic/gin"
)

// @Summary Register
// @Produce json
// @Param register body routers.Register true "username, email and password"
// @Success 201 {object} routers.Response
// @Failure 400 {object} routers.Response
// @Router /register [post]
func register(c *gin.Context) {
	var request Register
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code: 400,
			Msg:  "username, email and password are mandatory",
			Data: nil,
		})
		return
	}

	err := models.Register(request.Username, request.Email, request.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code: 400,
			Msg:  err.Error(),
			Data: nil,
		})
	} else {
		c.JSON(http.StatusCreated, Response{
			Code: 201,
			Msg:  "register succeed",
			Data: nil,
		})
	}
}
