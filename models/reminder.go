package models

import (
	"time"

	"github.com/gobuffalo/uuid"
)

type Reminder struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	ChatID    int64     `json:"chat_id" db:"chat_id"`
	UserID    int       `json:"user_id" db:"user_id"`
	MessageID int       `json:"message_id" db:"message_id"`
	Message   string    `json:"message" db:"message"`
	Deadline  time.Time `json:"deadline" db:"deadline"`
	Executed  bool      `json:"executed" db:"executed"`
}

type Reminders []Reminder

// Block the routine until the deadline is reached
func (r *Reminder) Wait() {
	if r.Deadline.After(time.Now()) {
		delay := r.Deadline.UnixNano() - time.Now().UnixNano()

		timer := time.NewTimer(time.Duration(delay))
		<-timer.C
	}
}
