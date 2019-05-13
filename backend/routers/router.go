package routers

import (
	"encoding/gob"
	"net/http"

	"github.com/fancxxy/gocomic/backend/config"
	"github.com/fancxxy/gocomic/backend/cron"
	"github.com/fancxxy/gocomic/backend/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	ginSwagger "github.com/swaggo/gin-swagger"   // gin-swagger middleware
	"github.com/swaggo/gin-swagger/swaggerFiles" // swagger embed files

	// docs
	_ "github.com/fancxxy/gocomic/backend/docs"
)

// Init function initialize and return gin engine
// @title gocomic api
// @version 1.0
// @description This is a comic server.
// termsOfService https://github.com/fancxxy/gocomic/backend
// @licenese.name MIT
// @licenese.url https://github.com/fancxxy/gocomic/blob/master/LICENSE
// @host 127.0.0.1:8080
func Init() error {
	server := config.Config.Server
	gin.SetMode(server.Mode)
	gobRegister()
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	store := cookie.NewStore([]byte(server.Secret))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 24 * 7,
		HttpOnly: true,
	})
	r.Use(sessions.Sessions("session", store))

	r.POST("/register", register)
	r.POST("/login", login)
	r.GET("/logout", logout)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api1 := r.Group("/api/v1")
	api1.Use(loginRequired())
	{
		api1.GET("/comics", list)
		api1.POST("/comics", subscribe)

		api1.DELETE("/comics", unsubscribe)
		api1.GET("/comics/:id", comic)
		api1.PATCH("/comics/:id", update)
		api1.GET("/comics/:id/:cid", chapter)
	}

	api2 := r.Group("/admin")
	api2.Use(loginRequired(), AdminRequired())
	{
		api2.GET("/update", func(c *gin.Context) {
			c.JSON(http.StatusOK, Response{
				Code: 200,
				Msg:  "ok",
				Data: cron.Update(),
			})
		})
		api2.GET("/delete", func(c *gin.Context) {
			c.JSON(http.StatusOK, Response{
				Code: 200,
				Msg:  "ok",
				Data: cron.Clean(),
			})
		})
	}

	return r.Run(server.Address)
}

func gobRegister() {
	gob.Register(models.User{})
}
