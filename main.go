package main

import (
	"io/ioutil"
	"log"

	"github.com/gobuffalo/envy"
	"gopkg.in/telegram-bot-api.v4"
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
	connectDB()
	migrateDB()
	initRedis()
	initBot()
	loadReminders()

	initCrons()

	go ginit()

	defer Cron.Stop()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := Octaaf.GetUpdatesChan(u)

	if err != nil {
		log.Panicf("Failed to fetch updates: %v", err)
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
		log.Printf("Error while loading version string: %v", err)
		return
	}

	OctaafVersion = string(bytes)
}
