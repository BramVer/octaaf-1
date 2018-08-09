package main

import (
	"fmt"
	"log"
	"octaaf/models"
	"time"

	"gopkg.in/telegram-bot-api.v4"
)

func startReminder(reminder models.Reminder) {
	err := DB.Save(&reminder)

	if err != nil {
		log.Printf("reminder save error: %v", err)
		return
	}

	// Only wait for the deadline if it's actually in the future
	if reminder.Deadline.After(time.Now()) {
		delay := reminder.Deadline.UnixNano() - time.Now().UnixNano()

		timer := time.NewTimer(time.Duration(delay))
		<-timer.C
	}

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
