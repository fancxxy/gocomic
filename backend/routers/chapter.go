package routers

import (
	"net/http"
	"strconv"

	"github.com/fancxxy/gocomic/backend/config"
	"github.com/fancxxy/gocomic/backend/models"
	"github.com/fancxxy/gocomic/backend/rpc"
	"github.com/gin-gonic/gin"
)

// @Summary Get the chapter
// @Produce json
// @Param id path int true "comic id"
// @Param cid path int true "chapter index"
// @Success 200 {object} routers.Response
// @Success 200 {object} routers.Response
// @Failure 400 {object} routers.Response
// @Failure 404 {object} routers.Response
// @Router /api/v1/comics/{id}/{cid} [get]
func chapter(c *gin.Context) {
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

	cidParam := c.Param("cid")
	cid, err := strconv.Atoi(cidParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code: 400,
			Msg:  "invalid cid parameter",
			Data: nil,
		})
		return
	}

	chapter := models.GetChapter(id, cid)
	if chapter == nil {
		c.JSON(http.StatusNotFound, Response{
			Code: 404,
			Msg:  "cannot find the chapter",
			Data: nil,
		})
		return
	}

	if chapter.Cached {
		c.JSON(http.StatusOK, Response{
			Code: 200,
			Msg:  "get the chapter succeed",
			Data: struct {
				*models.Chapter `json:"chapter"`
			}{chapter},
		})
		return
	}

	response, err := rpc.DownloadChapter(chapter.URL, config.Config.Download.Path)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code: 400,
			Msg:  err.Error(),
			Data: nil,
		})
		return
	}

	chapter.Save(response)
	c.JSON(http.StatusOK, Response{
		Code: 200,
		Msg:  "get the chapter succeed",
		Data: chapterData{Chapter: chapter},
	})
}
