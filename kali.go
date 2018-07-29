package main

import (
	"fmt"
	"log"
	"octaaf/models"
	"strings"
	"time"

	"gopkg.in/telegram-bot-api.v4"
)

// ReporterID is the id of the user who reports everyone
var ReporterID int

// KaliCount is an integer that holds the ID of the last sent message in the Kali group
var KaliCount int

// KaliID is the ID of the kali group
var KaliID int64

func kaliHandler(message *tgbotapi.Message) {
	if message.Chat.ID == KaliID {
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
	if message.From.ID == ReporterID {
		if strings.ToLower(message.Text) == "reported" || (message.Sticker != nil && message.Sticker.FileID == "CAADBAAD5gEAAreTBA3s5qVy8bxHfAI") {
			DB.Save(&models.Report{})
		}
	}
}

func getLeetBlazers(event string) {
	log.Print("Getting blazers")
	participators := Redis.SMembers(event).Val()

	log.Printf("Blazers count: %v", len(participators))

	if len(participators) == 0 {
		sendGlobal(fmt.Sprintf("Nobody participated in today's %v", event))
		return
	}

	reply := "Today "

	if len(participators) == 1 {
		reply += fmt.Sprintf("only @%v participated in the %v.", participators[0], event)
		sendGlobal(reply)
		Redis.Del(event)
		return
	}

	for index, participator := range participators {
		if index == (len(participators) - 1) {
			reply += fmt.Sprintf(" and @%v", participator)
		} else {
			if index == 0 {
				reply += fmt.Sprintf("@%v", participator)
			} else {
				reply += fmt.Sprintf(", @%v", participator)
			}
		}
	}

	reply += fmt.Sprintf(" participated in the %v.", event)
	sendGlobal(reply)
	Redis.Del(event)
}

func addLeetBlazer(message *tgbotapi.Message, event string) {
	if strings.Contains(message.Text, event) {
		log.Printf("Leetblazer found!")
		Redis.SAdd(event, message.From.UserName)
	}
}

func setKaliCount() {
	lastCount := models.MessageCount{}

	err := DB.Last(&lastCount)

	count := models.MessageCount{
		Count: KaliCount,
		Diff:  0,
	}

	if err == nil && lastCount.Count > 0 {
		count.Diff = (KaliCount - lastCount.Count)
	}

	DB.Save(&count)
}
