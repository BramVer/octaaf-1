package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type Quote struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Quote     string    `json:"quote" db:"quote"`
	UserID    int       `json:"user_id" db:"user_id"`
	ChatID    int64     `json:"chat_id" db:"chat_id"`
}

// String is not required by pop and may be deleted
func (q Quote) String() string {
	jq, _ := json.Marshal(q)
	return string(jq)
}

// Quotes is not required by pop and may be deleted
type Quotes []Quote

// String is not required by pop and may be deleted
func (q Quotes) String() string {
	jq, _ := json.Marshal(q)
	return string(jq)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (q *Quote) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: q.Quote, Name: "Quote"},
		&validators.IntIsPresent{Field: q.UserID, Name: "UserID"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (q *Quote) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (q *Quote) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func (q *Quote) Search(tx *pop.Connection, chatID int64, filter ...string) error {
	if len(filter) > 0 && filter[0] != "" {
		return tx.
			Where("chat_id = ?", chatID).
			Where("quote ilike '%' || ? || '%'", filter[0]).
			Order("random()").Limit(1).First(q)
	} else {
		return tx.Where("chat_id = ?", chatID).Order("random()").Limit(1).First(q)
	}
}
