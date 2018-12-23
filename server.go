package main

import (
	"bitbucket.org/marketingx/upvideo/app/storage/accounts"
	"bitbucket.org/marketingx/upvideo/app/storage/campaigns"
	"bitbucket.org/marketingx/upvideo/app/domain/session"
	"bitbucket.org/marketingx/upvideo/app/domain/usr"
	"bitbucket.org/marketingx/upvideo/app/domain/email"
	"bitbucket.org/marketingx/upvideo/app/httpserver"
	"bitbucket.org/marketingx/upvideo/app/httpserver/web"
	"bitbucket.org/marketingx/upvideo/app/storage/jobs"
	"bitbucket.org/marketingx/upvideo/app/utils/keywordtool"
	"bitbucket.org/marketingx/upvideo/app/utils/rapidtags"
	"bitbucket.org/marketingx/upvideo/app/storage/titles"
	"bitbucket.org/marketingx/upvideo/app/storage/videos"
	"bitbucket.org/marketingx/upvideo/app/storage/shortlinks"
	"bitbucket.org/marketingx/upvideo/app/storage/invites"
	"bitbucket.org/marketingx/upvideo/app/utils/aws"
	"bitbucket.org/marketingx/upvideo/app/utils/utils"
	"bitbucket.org/marketingx/upvideo/config"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Syntax: " + os.Args[0] + " <path-to-config.json>")
		return
	}

	cfg := config.ReadConfig(os.Args[1])
	db, err := sql.Open("mysql", cfg.Dsn)
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxIdleConns(100)
	defer db.Close()
	initDbTables(db)

	err = aws.AWSInitSession(&cfg.AWS)
	utils.InitProject()

	var sessionRepository session.Repository
	if cfg.Session.Storage == "db" {
		sessionRepository = httpserver.NewDbSession(db)
	} else if cfg.Session.Storage == "memory" {
		sessionRepository = httpserver.GetMemSession()
	} else {
		fmt.Println("Incorrect session storage specified")
		os.Exit(1)
	}
	fmt.Printf("config: %v\n", cfg)

	titlesService     := titles.NewService(titles.NewRepository(db))
	videoService      := videos.NewService(videos.NewRepository(db))
	jobsService       := jobs.NewService(jobs.NewRepository(db))
	accountService    := accounts.NewService(accounts.NewRepository(db))
	campaignService   := campaigns.NewService(campaigns.NewRepository(db))
	shortlinksService := shortlinks.NewService(shortlinks.NewRepository(db))
	invitesService    := invites.NewService(invites.NewRepository(db))

	webServer := &web.WebServer{
		UserService:        usr.NewUserService(httpserver.NewUserRepository(db)),
		SessionService:     session.NewService(sessionRepository),
		VideoService:       videoService,
		AccountService:     accountService,
		CampaignService:    campaignService,
		TitleService:       titlesService,
		JobService:         jobsService,
		KeywordtoolService: keywordtool.NewService(&cfg.Keywordtool),
		RapidtagsService:   rapidtags.NewService(),
		EmailService:       email.NewService(&cfg.Email),
		ShortlinksService:  shortlinksService,
		InviteService:      invitesService,
		Params:             cfg.WebServer,
		Config:             cfg,
	}

	webServer.Start()
}

func initDbTables(db *sql.DB) {
	var count int
	_ = db.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE();").Scan(&count)
	if count == 0 {
		file, err := ioutil.ReadFile("./db.sql")
		if err != nil {
			fmt.Println("Error: db.sql not found in current directory")
			os.Exit(1)
		}
		requests := strings.Split(string(file), ";")
		for _, request := range requests {
			if strings.TrimSpace(request) != "" {
				_, err := db.Exec(request)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}
