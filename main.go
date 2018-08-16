package main

import (
	"io/ioutil"

	"github.com/gobuffalo/envy"
	log "github.com/sirupsen/logrus"
	"gopkg.in/telegram-bot-api.v4"
	"octaaf/web"
)

// OctaafEnv is either development or production
var OctaafEnv string

// GitUri is the upstream development URL
const GitUri = "https://gitlab.com/BartWillems/octaaf"

var OctaafVersion string

func main() {
	envy.Load("config/.env")

	OctaafEnv = envy.Get("GO_ENV", "development")

	loadVersion()
	initDB()
	initRedis()
	initBot()
	go loadReminders()

	cron := initCrons()
	cron.Start()
	defer cron.Stop()

	router := web.New(web.Connections{
		Octaaf:   Octaaf,
		Postgres: DB,
		Redis:    Redis,
		KaliID:   KaliID,
	})

	go func() {
		err := router.Run()
		if err != nil {
			log.Errorf("Gin creation error: %v", err)
		}
	}()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := Octaaf.GetUpdatesChan(u)

	if err != nil {
		log.Fatalf("Failed to fetch updates: %v", err)
	}

	for update := range updates {

		if update.Message == nil {
			continue
		}

		go handle(update.Message)
	}
}

func loadVersion() {
	bytes, err := ioutil.ReadFile("assets/version")

	if err != nil {
		log.Errorf("Error while loading version string: %v", err)
		return
	}

	OctaafVersion = string(bytes)
	log.Infof("Loaded version %v", OctaafVersion)
}
