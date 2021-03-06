package web

import (
	"bitbucket.org/marketingx/upvideo/app/domain/email"
	"bitbucket.org/marketingx/upvideo/app/domain/session"
	"bitbucket.org/marketingx/upvideo/app/domain/usr"
	"bitbucket.org/marketingx/upvideo/app/storage/accounts"
	"bitbucket.org/marketingx/upvideo/app/storage/campaigns"
	"bitbucket.org/marketingx/upvideo/app/storage/invites"
	"bitbucket.org/marketingx/upvideo/app/storage/jobs"
	"bitbucket.org/marketingx/upvideo/app/storage/shortlinks"
	"bitbucket.org/marketingx/upvideo/app/storage/titles"
	"bitbucket.org/marketingx/upvideo/app/storage/topusers"
	"bitbucket.org/marketingx/upvideo/app/storage/videos"
	"bitbucket.org/marketingx/upvideo/app/utils/keywordtool"
	"bitbucket.org/marketingx/upvideo/app/utils/rapidtags"
	"bitbucket.org/marketingx/upvideo/config"
	"fmt"
	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net/http"
)

type WebServer struct {
	Config             config.Config
	Params             config.WebServerParams
	UserService        usr.UserService
	SessionService     session.Service
	VideoService       *videos.Service
	AccountService     *accounts.Service
	CampaignService    *campaigns.Service
	TitleService       *titles.Service
	TopUserService     *topusers.Service
	ShortlinksService  *shortlinks.Service
	JobService         *jobs.Service
	KeywordtoolService *keywordtool.Service
	RapidtagsService   *rapidtags.Service
	InviteService      *invites.Service
	EmailService       email.Service
}

func (this *WebServer) Start() {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	if this.Params.Production {
		gin.SetMode(gin.ReleaseMode)
	}

	// cross Origin helper
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Cache-Control, Authorization")
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
