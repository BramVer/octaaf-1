package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/uuid"
)

type Reminder struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	ChatID    int64     `json:"chat_id" db:"chat_id"`
	MessageID int       `json:"message_id" db:"message_id"`
	Message   string    `json:"message" db:"message"`
	Deadline  time.Time `json:"deadline" db:"deadline"`
	Executed  bool      `json:"executed" db:"executed"`
}

// Alerts is not required by pop and may be deleted
type Reminders []Reminder

// String is not required by pop and may be deleted
func (r Reminders) String() string {
	jr, _ := json.Marshal(r)
	return string(jr)
}
