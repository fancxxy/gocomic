package routers

import (
	"net/http"
	"strconv"

	"github.com/fancxxy/gocomic/backend/models"
	"github.com/fancxxy/gocomic/backend/rpc"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// @Summary Update the chapter list of the comic
// @Produce json
// @Param id path int true "comic id"
// @Success 200 {object} routers.Response
// @Success 204 {object} routers.Response
// @Failure 400 {object} routers.Response
// @Failure 404 {object} routers.Response
// @Router /api/v1/comics/{id} [patch]
func update(c *gin.Context) {
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
	comic := models.GetComic(id, -1, -1)
	if comic == nil {
		c.JSON(http.StatusNotFound, Response{
			Code: 404,
			Msg:  "cannot find the comic",
			Data: nil,
		})
		return
	}

	response, err := rpc.UpdateComic(comic.URL, comic.Latest)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code: 400,
			Msg:  err.Error(),
			Data: nil,
		})
		return
	}

	session := sessions.Default(c)
	user := session.Get("user").(models.User)
	if comic.Update(response) {
		c.JSON(http.StatusOK, Response{
			Code: 200,
			Msg:  "update the comic succeed",
			Data: comicData{
				Time:  models.SubscribeDate(&user, comic),
				Comic: comic,
			},
		})
	} else {
		c.JSON(http.StatusNoContent, Response{
			Code: 204,
			Msg:  "comic not updated",
			Data: nil,
		})
	}
}
