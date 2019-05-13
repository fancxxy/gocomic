package routers

import (
	"net/http"
	"time"

	"github.com/fancxxy/gocomic/backend/models"
	"github.com/fancxxy/gocomic/backend/rpc"
	"github.com/gin-contrib/sessions"

	"github.com/gin-gonic/gin"
)

// @Summary Subscribe a comic
// @Produce json
// @Param subscribe body routers.Subscribe true "website and comic"
// @Success 200 {object} routers.Response
// @Success 201 {object} routers.Response
// @Failure 400 {object} routers.Response
// @Router /api/v1/comics [post]
func subscribe(c *gin.Context) {
	var request Subscribe
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code: 400,
			Msg:  "website and comic are mandatory",
			Data: nil,
		})
		return
	}

	session := sessions.Default(c)
	user := session.Get("user").(models.User)
	comic := models.QueryComic(request.Website, request.Comic, -1, -1)
	if comic != nil {
		// 已经收录，直接添加一条订阅记录
		err := models.Subscribe(&user, comic)
		if err != nil {
			c.JSON(http.StatusBadRequest, Response{
				Code: 400,
				Msg:  err.Error(),
				Data: nil,
			})
		} else {
			c.JSON(http.StatusOK, Response{
				Code: 200,
				Msg:  "subscribe the comic succeed",
				Data: struct {
					models.Time   `json:"subscribed"`
					*models.Comic `json:"comic"`
				}{
					models.Time{time.Now()},
					comic,
				},
			})
		}
	} else {
		// 还没有收录，调用rpc接口，插入目录信息到数据库
		response, err := rpc.DownloadComic(request.Website, request.Comic)
		if err != nil {
			c.JSON(http.StatusBadRequest, Response{
				Code: 400,
				Msg:  err.Error(),
				Data: nil,
			})
		}

		models.AddComic(response, &user)
		c.JSON(http.StatusCreated, Response{
			Code: 201,
			Msg:  "create and subscribe the comic succeed",
			Data: struct {
				models.Time   `json:"subscribed"`
				*models.Comic `json:"comic"`
			}{
				models.Time{time.Now()},
				comic,
			},
		})
	}
}

// @Summary Unsubscribe the comic
// @Produce json
// @Param subscribe body routers.Unsubscribe true "comic id"
// @Success 200 {object} routers.Response
// @Failure 400 {object} routers.Response
// @Router /api/v1/comics [delete]
func unsubscribe(c *gin.Context) {
	var request Unsubscribe
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code: 400,
			Msg:  "id is mandatory",
			Data: nil,
		})
		return
	}

	comic := models.GetComic(request.ID, -1, -1)
	if comic == nil {
		c.JSON(http.StatusBadRequest, Response{
			Code: 400,
			Msg:  "the comic not collected",
			Data: nil,
		})
	}

	session := sessions.Default(c)
	user := session.Get("user").(models.User)
	err := models.Unsubscribe(&user, comic)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code: 400,
			Msg:  err.Error(),
			Data: nil,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code: 200,
		Msg:  "unsubscribe the comic succeed",
		Data: nil,
	})
}
