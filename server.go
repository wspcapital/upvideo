package main

import (
	"bitbucket.org/marketingx/upvideo/app/accounts"
	"bitbucket.org/marketingx/upvideo/app/domain/session"
	"bitbucket.org/marketingx/upvideo/app/domain/usr"
	"bitbucket.org/marketingx/upvideo/app/infrastructure"
	"bitbucket.org/marketingx/upvideo/app/infrastructure/web"
	"bitbucket.org/marketingx/upvideo/app/jobs"
	"bitbucket.org/marketingx/upvideo/app/jobworker"
	"bitbucket.org/marketingx/upvideo/app/services/keywordtool"
	"bitbucket.org/marketingx/upvideo/app/services/rapidtags"
	"bitbucket.org/marketingx/upvideo/app/videos"
	"bitbucket.org/marketingx/upvideo/app/videos/titles"
	"bitbucket.org/marketingx/upvideo/aws"
	"bitbucket.org/marketingx/upvideo/config"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"os"
	_ "strconv"
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

	err = aws.AWSInitSession(cfg)

	var sessionRepository session.Repository
	if cfg.Session.Storage == "db" {
		sessionRepository = infrastructure.NewDbSession(db)
	} else if cfg.Session.Storage == "memory" {
		sessionRepository = infrastructure.GetMemSession()
	} else {
		fmt.Println("Incorrect session storage specified")
		os.Exit(1)
	}
	fmt.Printf("config: %v\n", cfg)

	titlesService := titles.NewService(titles.NewRepository(db))
	videoService := videos.NewService(videos.NewRepository(db))
	jobsService := jobs.NewService(jobs.NewRepository(db))
	accountService := accounts.NewService(accounts.NewRepository(db))

	webServer := &web.WebServer{
		UserService:        usr.NewUserService(infrastructure.NewUserRepository(db)),
		SessionService:     session.NewService(sessionRepository),
		VideoService:       videoService,
		AccountService:     accountService,
		TitleService:       titlesService,
		JobService:         jobsService,
		KeywordtoolService: keywordtool.NewService(&cfg.Keywordtool),
		RapidtagsService:   rapidtags.NewService(),
		Params:             cfg.WebServer,
		Config:             cfg,
	}

	jobWorkerService := &jobworker.Service{
		VideoService:   videoService,
		TitleService:   titlesService,
		JobService:     jobsService,
		AccountService: accountService,
		Config:         &cfg,
	}
	jobWorkerService.Start()

	webServer.Start()
}

func initDbTables(db *sql.DB) {
	var count int
	db.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE();").Scan(&count)
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
