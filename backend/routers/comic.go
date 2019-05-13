package routers

import (
	"net/http"
	"strconv"

	"github.com/fancxxy/gocomic/backend/config"
	"github.com/fancxxy/gocomic/backend/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// @Summary Get the comic
// @Produce json
// @Param id path int true "comic id"
// @Param limit query int false "page limit"
// @Param page query int false "page number"
// @Success 200 {object} routers.Response
// @Failure 400 {object} routers.Response
// @Failure 404 {object} routers.Response
// @Router /api/v1/comics/{id} [get]
func comic(c *gin.Context) {
	pagination := config.Config.Pagination
	limitQuery := c.DefaultQuery("limit", pagination.Limit)
	pageQuery := c.DefaultQuery("page", pagination.Page)
	limit, err := strconv.Atoi(limitQuery)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code: 400,
			Msg:  "invalid limit query parameter",
			Data: nil,
		})
		return
	}
	page, err := strconv.Atoi(pageQuery)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code: 400,
			Msg:  "invalid page query parameter",
			Data: nil,
		})
		return
	}

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code: 400,
			Msg:  "invalid id parameter",
			Data: nil,
		})
		return
	}
	comic := models.GetComic(id, limit, page)
	if comic == nil {
		c.JSON(http.StatusNotFound, Response{
			Code: 404,
			Msg:  "cannot find the comic",
			Data: nil,
		})
		return
	}

	session := sessions.Default(c)
	user := session.Get("user").(models.User)

	c.JSON(http.StatusOK, Response{
		Code: 200,
		Msg:  "get the comic succeed",
		Data: comicData{
			Time:  models.SubscribeDate(&user, comic),
			Comic: comic,
		},
	})
}

// @Summary Get all subscribed comics by current user
// @Produce json
// @Success 200 {object} routers.Response
// @Router /api/v1/comics [get]
func list(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user").(models.User)
	c.JSON(http.StatusOK, Response{
		Code: 200,
		Msg:  "get the comics succeed",
		Data: user.ComicSlice(),
	})
}
