package main

import (
	"io/ioutil"

	"octaaf/web"

	"github.com/gobuffalo/envy"
	log "github.com/sirupsen/logrus"
	"gopkg.in/telegram-bot-api.v4"
)

type State struct {
	Environment string
	Version     string
	GitUri      string
}

var state *State

func main() {
	envy.Load("config/.env")

	state = &State{
		Environment: envy.Get("GO_ENV", "development"),
		Version:     loadVersion(),
		GitUri:      "https://gitlab.com/BartWillems/octaaf",
	}

	initRedis()

	initDB()
	migrateDB()

	initBot()

	go loadReminders()

	cron := initCrons()
	cron.Start()
	defer cron.Stop()

	go func() {
		router := web.New(web.Connections{
			Octaaf:   Octaaf,
			Postgres: DB,
			Redis:    Redis,
			KaliID:   KaliID,
		})
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

func loadVersion() string {
	bytes, err := ioutil.ReadFile("assets/version")

	if err != nil {
		log.Errorf("Error while loading version string: %v", err)
		return ""
	}
	log.Infof("Loaded version %v", string(bytes))
	return string(bytes)
}
