package main

import (
	"fmt"
	"log"
	"octaaf/models"

	"gopkg.in/telegram-bot-api.v4"
)

func startReminder(reminder models.Reminder) {
	err := DB.Save(&reminder)

	if err != nil {
		log.Printf("reminder save error: %v", err)
		return
	}

	reminder.Wait()

	var username string
	user, err := getUsername(reminder.UserID, reminder.ChatID)
	if err == nil {
		username = "@" + user.User.UserName + ": "
	}

	msg := tgbotapi.NewMessage(reminder.ChatID, fmt.Sprintf("%v%v", username, reminder.Message))
	msg.ReplyToMessageID = reminder.MessageID
	go Octaaf.Send(msg)

	// Mark this reminder as completed
	reminder.Executed = true
	DB.Save(&reminder)
}

func loadReminders() {
	var reminders models.Reminders
	err := DB.Where("executed = false").Order("created_at").All(&reminders)

	if err != nil {
		log.Printf("Unable to load pending reminders: %v", err)
		return
	}

	for _, reminder := range reminders {
		go startReminder(reminder)
	}
}
