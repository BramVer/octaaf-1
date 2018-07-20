package main

import (
	"fmt"
	"log"
	"octaaf/models"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gobuffalo/envy"
	"gopkg.in/telegram-bot-api.v4"
)

// Octaaf is the global bot endpoint
var Octaaf *tgbotapi.BotAPI

func initBot() {
	var err error
	Octaaf, err = tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_API_TOKEN"))

	if err != nil {
		log.Panic(err)
	}

	env := envy.Get("GO_ENV", "development")
	Octaaf.Debug = env == "development"

	log.Printf("Authorized on account %s", Octaaf.Self.UserName)

	KaliID, err = strconv.ParseInt(os.Getenv("TELEGRAM_ROOM_ID"), 10, 64)

	if err != nil {
		log.Panic(err)
	}

	ReporterID, err = strconv.Atoi(envy.Get("REPORTER_ID", "-1"))
	if err != nil {
		log.Panic(err)
	}

	if env != "development" {
		sendGlobal("I'm up and running! 👌")

		c := make(chan os.Signal, 2)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			sendGlobal("I'm going to sleep! 💤💤")
			DB.Close()
			os.Exit(0)
		}()
	}
}

func handle(message *tgbotapi.Message) {
	if message.Chat.ID == KaliID {

		go func() {
			KaliCount = message.MessageID

			if message.From.ID == ReporterID {
				if strings.ToLower(message.Text) == "reported" ||
					(message.Sticker != nil && message.Sticker.FileID == "CAADBAAD5gEAAreTBA3s5qVy8bxHfAI") {
					DB.Save(&models.Report{})
				}
			}
		}()
	}

	if message.IsCommand() {
		switch message.Command() {
		case "roll":
			sendRoll(message)
		case "m8ball":
			m8Ball(message)
		case "bodegem":
			sendBodegem(message)
		case "img", "img_sfw", "more":
			sendImage(message)
		case "stallman":
			sendStallman(message)
		case "avatar":
			sendAvatar(message)
		case "search", "search_nsfw":
			search(message)
		case "where":
			where(message)
		case "count":
			count(message)
		case "weather":
			weather(message)
		case "what":
			what(message)
		case "bol":
			bol(message)
		case "xkcd":
			xkcd(message)
		case "quote":
			quote(message)
		case "next_launch":
			nextLaunch(message)
		case "doubt":
			doubt(message)
		case "issues":
			issues(message)
		case "kalirank":
			kaliRank(message)
		case "iasip":
			iasip(message)
		}
	}

	if message.MessageID%100000 == 0 {
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("💯💯💯💯 YOU HAVE MESSAGE %v 💯💯💯💯", message.MessageID))
		msg.ReplyToMessageID = message.MessageID
		msg.ParseMode = "markdown"
		Octaaf.Send(msg)
	}

}

func sendGlobal(message string) {
	// Wait 1.5 seconds because Telegram has bad NTP
	time.Sleep(1500)

	msg := tgbotapi.NewMessage(KaliID, message)
	_, err := Octaaf.Send(msg)

	if err != nil {
		log.Printf("Error while sending '%s': %s", message, err)
	}
}

func reply(message *tgbotapi.Message, text string, markdown ...bool) {
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ReplyToMessageID = message.MessageID

	if len(markdown) > 0 {
		if markdown[0] {
			msg.ParseMode = "markdown"
		}
	} else {
		msg.ParseMode = "markdown"
	}

	Octaaf.Send(msg)
}
