package main

import (
	"log"

	"github.com/go-redis/redis"
	"github.com/gobuffalo/envy"
	"gopkg.in/robfig/cron.v2"
	"gopkg.in/telegram-bot-api.v4"
)

// KaliCount is an integer that holds the ID of the last send message in the Kali group
var KaliCount int

// KaliID is the ID of the kali group
var KaliID int64

// ReporterID is the id of the user who reports everyone
var ReporterID int

// Cron executes all the cron jobs
var Cron *cron.Cron

// Redis client
var Redis *redis.Client

func main() {
	envy.Load("config/.env")

	connectDB()
	migrateDB()
	initBot()

	Redis = redis.NewClient(&redis.Options{
		Addr:     envy.Get("REDIS_URI", "localhost:6379"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	initCrons()

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
