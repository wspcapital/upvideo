package main

import (
	"bitbucket.org/marketingx/upvideo/app/accounts"
	"bitbucket.org/marketingx/upvideo/app/domain/session"
	"bitbucket.org/marketingx/upvideo/app/domain/usr"
	"bitbucket.org/marketingx/upvideo/app/infrastructure"
	"bitbucket.org/marketingx/upvideo/app/infrastructure/web"
	"bitbucket.org/marketingx/upvideo/app/videos"
	"bitbucket.org/marketingx/upvideo/app/videos/titles"
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

	webServer := &web.WebServer{
		UserService:    usr.NewUserService(infrastructure.NewUserRepository(db)),
		SessionService: session.NewService(sessionRepository),
		VideoService:   videos.NewService(videos.NewRepository(db)),
		AccountService: accounts.NewService(accounts.NewRepository(db)),
		TitleService:   titles.NewService(titles.NewRepository(db)),
		Params:         cfg.WebServer,
		Config:         cfg,
	}

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
