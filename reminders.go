package main

import (
	"fmt"
	"octaaf/models"

	log "github.com/sirupsen/logrus"
	"gopkg.in/telegram-bot-api.v4"
)

func startReminder(reminder models.Reminder) {
	log.Debugf("New reminder (%v) added for %v", reminder.Message, reminder.Deadline.String())
	err := DB.Save(&reminder)

	if err != nil {
		log.Errorf("reminder save error: %v", err)
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
	err = reminder.Complete(DB)
	if err != nil {
		log.Errorf("Unable to mark the reminder {%v} as completed: %v", reminder.ID, err)
	}
}

func loadReminders() {
	var reminders models.Reminders
	err := reminders.Pending(DB)

	if err != nil {
		log.Errorf("Unable to load pending reminders: %v", err)
		return
	}

	for _, reminder := range reminders {
		log.Debugf("Loaded reminder %v with message: %v", reminder.ID, reminder.Message)
		go startReminder(reminder)
	}
}
