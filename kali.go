package main

import (
	"fmt"
	"octaaf/models"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize/english"
	log "github.com/sirupsen/logrus"
	"gopkg.in/telegram-bot-api.v4"
)

// KaliCount is an integer that holds the ID of the last sent message in the Kali group
var KaliCount int

func kaliHandler(message *tgbotapi.Message) {
	if message.Chat.ID == settings.Telegram.KaliID {
		log.Debug("Kalimember found")
		KaliCount = message.MessageID

		go kaliReport(message)

		if time.Now().Hour() == 13 && time.Now().Minute() == 37 {
			go addLeetBlazer(message, "1337")
		}

		if time.Now().Hour() == 16 && time.Now().Minute() == 20 {
			go addLeetBlazer(message, "420")
		}
	}
}

func kaliReport(message *tgbotapi.Message) {
	if message.From.ID == settings.Telegram.ReporterID {
		log.Debug("Reporter found")
		if strings.ToLower(message.Text) == "reported" ||
			(message.Sticker != nil && message.Sticker.FileID == "CAADBAAD5gEAAreTBA3s5qVy8bxHfAI") {
			DB.Save(&models.Report{})
		}
	}
}

func getLeetBlazers(event string) {
	log.Info("Getting blazers")
	participators := Redis.SMembers(event).Val()

	log.Info("Blazers count: ", len(participators))

	if len(participators) == 0 {
		sendGlobal(fmt.Sprintf("Nobody participated in today's %v", event))
		return
	}

	var usernames []string

	// Loop through the participators
	// Fetch their usernames
	// Store the kalivent in the DB
	for _, participator := range participators {
		userID, _ := strconv.Atoi(participator)
		user, err := getUsername(userID, settings.Telegram.KaliID)

		if err != nil {
			log.Error("Unable to fetch username for the kalivent %v; error: %v", event, err)
			continue
		}

		usernames = append(usernames, fmt.Sprintf("@%v", user.User.UserName))

		// Store this absolute unit in the database
		kali := models.Kalivent{
			UserID: userID,
			Date:   time.Now(),
			Type:   event}
		DB.Save(&kali)
	}

	reply := "Today "
	if len(usernames) == 1 {
		reply += "only "
	}

	reply += english.WordSeries(usernames, "and")

	reply += fmt.Sprintf(" participated in the %v.", event)
	sendGlobal(reply)
	Redis.Del(event)
}

func addLeetBlazer(message *tgbotapi.Message, event string) {
	if strings.Contains(message.Text, event) {
		log.Infof("Leetblazer found with id: %v!", message.From.ID)
		Redis.SAdd(event, message.From.ID)
	}
}

func setKaliCount() {
	count := models.MessageCount{
		Count: KaliCount,
		Diff:  0,
	}

	err := DB.Save(&count)
	if err != nil {
		log.Error("Unable to save today's kalicount: ", err)
	}
}
