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

	r.POST("/auth/signin", this.signin)
	r.POST("/auth/signup", this.signup)
	r.POST("/auth/forgot-password", this.userForgotPassword)
	r.POST("/auth/restore-password", this.userRestorePassword)

	r.POST("/auth/change-password", this.requireAuth, this.userChangePassword)
	r.POST("/auth/reset-apikey", this.requireAuth, this.userResetApikey)

	//dashboard
	dashboard := r.Group("/api/dashboard")
	dashboard.Use(this.requireAuth)
	{
		dashboard.GET("/:period", this.dashboardView)
	}

	profile := r.Group("/api/profile")
	profile.Use(this.requireAuth)
	{
		profile.GET("/", this.profile)
	}

	// video group api
	videos := r.Group("/api/videos")
	videos.Use(this.requireAuth)
	{
		videos.GET("", this.videoIndex)
		videos.GET("/:id", this.videoView)
		videos.POST("", this.videoCreate)
		videos.PUT("/:id", this.videoUpdate)
		videos.DELETE("/:id", this.videoDelete)
	}

	// accounts group api
	accounts := r.Group("/api/accounts")
	accounts.Use(this.requireAuth)
	{
		accounts.GET("", this.accountIndex)
		accounts.GET("/:id", this.accountView)
		accounts.GET("/:id/select", this.accountSelect)
		accounts.POST("", this.accountCreate)
		accounts.POST("/confirm", this.accountConfirm)
		accounts.PUT("/:id", this.accountUpdate)
		accounts.DELETE("/:id", this.accountDelete)
	}

	// campaigns group api
	campaigns := r.Group("/api/campaigns")
	campaigns.Use(this.requireAuth)
	{
		campaigns.GET("", this.campaignIndex)
		campaigns.GET("/:id", this.campaignView)
		campaigns.POST("", this.campaignCreate)
		campaigns.PUT("/:id", this.campaignUpdate)
		campaigns.DELETE("/:id", this.campaignDelete)
		campaigns.POST("/:id/gen_titles", this.campaignGenerateTitles)
		campaigns.GET("/:id/get_titles", this.campaignGetTitles)
		campaigns.POST("/:id/add_titles", this.campaignAddTitles)
	}

	// titles group api
	_titles := r.Group("/api/titles")
	_titles.Use(this.requireAuth)
	{
		_titles.GET("", this.titleIndex)
		_titles.GET("/:id", this.titleView)
		_titles.POST("", this.titleCreate)
		_titles.GET("/:id/suggest", this.titleSuggest)
		_titles.PUT("/:id", this.titleUpdate)
		_titles.DELETE("/:id", this.titleDelete)

		_titles.POST("/:id/convert", this.titleConvert)
		_titles.POST("/:id/publish", this.titlePublish)
		_titles.GET("/:id/status", this.titleStatus)
	}

	// shortlinks group apis.
	r.GET("/go/:uniq_id", this.ShortLinkRedirect)
	_shortlinks := r.Group("/api/shortlinks")
	_shortlinks.Use(this.requireAuth)
	{
		_shortlinks.GET("", this.shortlinksIndex)
		_shortlinks.GET("/:id", this.shortlinkView)
		_shortlinks.POST("", this.shortlinkCreate)
		_shortlinks.POST("/:id", this.shortlinkUpdate)
		_shortlinks.PUT("/:uniqid/update", this.shortlinkUpdatebyUniqId)
		_shortlinks.DELETE("/:id", this.shortlinkDelete)
	}

	// inites group apis.
	invite := r.Group("/api/invites")
	invite.Use(this.requireAuth)
	{
		invite.GET("/", this.inviteIndex)
		invite.GET("/:id", this.inviteView)
		invite.POST("/", this.inviteCreate)
		invite.PUT("/:id", this.inviteUpdate)
		invite.DELETE("/:id", this.inviteDelete)
	}


}
