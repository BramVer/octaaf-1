package main

import (
	"io/ioutil"

	"octaaf/web"

	"github.com/go-redis/cache"
	"github.com/go-redis/redis"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop"
	log "github.com/sirupsen/logrus"
	"gopkg.in/telegram-bot-api.v4"
)

type State struct {
	Environment string
	Version     string
	GitUri      string
	DB          *pop.Connection
	Redis       *redis.Client
	Codec       *cache.Codec
}

var state *State

func main() {
	envy.Load("config/.env")

	state = &State{
		Environment: envy.Get("GO_ENV", "development"),
		Version:     loadVersion(),
		GitUri:      "https://gitlab.com/BartWillems/octaaf",
		DB:          getDB(),
		Redis:       getRedis(),
		Codec:       getCodec(),
	}

	migrateDB()

	initBot()
	go loadReminders()

	cron := initCrons()
	cron.Start()
	defer cron.Stop()

	router := web.New(web.Connections{
		Octaaf:   Octaaf,
		Postgres: state.DB,
		Redis:    state.Redis,
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

func loadVersion() string {
	bytes, err := ioutil.ReadFile("assets/version")

	if err != nil {
		log.Errorf("Error while loading version string: %v", err)
		return ""
	}
	log.Infof("Loaded version %v", string(bytes))
	return string(bytes)
}
