package web

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func (this *WebServer) initRoutes(r *gin.Engine) {

	r.Use(static.Serve("/", static.LocalFile(this.Params.StaticDir, false)))
	r.NoRoute(func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.Status(200)
		} else {
			c.File(this.Params.StaticDir + "/index.html")
		}
	})

	r.POST("/api/login", this.todoLogin)
	profile := r.Group("/api/profile")
	profile.Use(this.requireAuth)
	{
		profile.GET("/", this.profile)
	}
	r.POST("/api/register", this.register)

	// video group api
	videos := r.Group("/api/videos")
	videos.Use(this.requireAuth)
	{
		videos.GET("/", this.videoIndex)
		videos.GET("/:id", this.videoView)
		videos.POST("/", this.videoCreate)
		videos.PUT("/:id", this.videoUpdate)
		videos.DELETE("/:id", this.videoDelete)
		videos.POST("/:id/gen_titles", this.VideoGenerateTitles)
		videos.GET("/:id/get_titles", this.VideoGetTitles)
	}

	// accounts group api
	accounts := r.Group("/api/accounts")
	accounts.Use(this.requireAuth)
	{
		accounts.GET("/", this.accountIndex)
		accounts.GET("/:id", this.accountView)
		accounts.POST("/", this.accountCreate)
		accounts.POST("/confirm", this.accountConfirm)
		accounts.PUT("/:id", this.accountUpdate)
		accounts.DELETE("/:id", this.accountDelete)
	}

	// titles group api
	_titles := r.Group("/api/titles")
	_titles.Use(this.requireAuth)
	{
		_titles.GET("/", this.titleIndex)
		_titles.GET("/:id", this.titleView)
		_titles.POST("/", this.titleCreate)
		_titles.PUT("/:id", this.titleUpdate)
		_titles.DELETE("/:id", this.titleDelete)

		_titles.POST("/:id/convert", this.titleConvert)
		_titles.POST("/:id/publish", this.titlePublish)
		_titles.GET("/:id/status", this.titleStatus)
	}

}
