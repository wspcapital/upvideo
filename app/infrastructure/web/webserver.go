package web

import (
	"fmt"
	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"bitbucket.org/marketingx/upvideo/app/domain/session"
	"bitbucket.org/marketingx/upvideo/app/domain/usr"
	"bitbucket.org/marketingx/upvideo/app/videos"
	"bitbucket.org/marketingx/upvideo/app/videos/titles"
	"bitbucket.org/marketingx/upvideo/app/accounts"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net/http"
)

type WebServerParams struct {
	Bind           string
	StaticDir      string
	DebugMode      bool
	SSL            bool
	CertCache      string
	HostsWhitelist []string
}

type WebServer struct {
	Params           WebServerParams
	UserService      usr.UserService
	SessionService   session.Service
	VideoService     *videos.Service
	AccountService   *accounts.Service
	TitleService     *titles.Service
}

func (this *WebServer) Start() {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// cross Origin helper
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "content-type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		if c.Request.Method == "OPTIONS" {
			c.Status(200)
		} else {
			c.Next()
		}
	})

	this.initRoutes(r)

	fmt.Println("Starting server on " + this.Params.Bind)
	if this.Params.SSL {
		fmt.Println("SSL Enabled")
		m := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(this.Params.HostsWhitelist...),
			Cache:      autocert.DirCache(this.Params.CertCache),
		}

		log.Fatal(autotls.RunWithManager(r, &m))
	} else {
		r.Run(this.Params.Bind)
	}

}

// statics
func (this *WebServer) homepage(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "./index.html")
}

func (this *WebServer) assets(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, req.URL.Path[1:])
}
