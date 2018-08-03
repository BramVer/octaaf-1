package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type Kalivent struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	UserID    int       `json:"user_id" db:"user_id"`
	Date      time.Time `json:"date" db:"date"`
	Type      string    `json:"type" db:"type"`
}

// String is not required by pop and may be deleted
func (k Kalivent) String() string {
	jk, _ := json.Marshal(k)
	return string(jk)
}

// Kalivents is not required by pop and may be deleted
type Kalivents []Kalivent

// String is not required by pop and may be deleted
func (k Kalivents) String() string {
	jk, _ := json.Marshal(k)
	return string(jk)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (k *Kalivent) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.IntIsPresent{Field: k.UserID, Name: "UserID"},
		&validators.TimeIsPresent{Field: k.Date, Name: "Date"},
		&validators.StringIsPresent{Field: k.Type, Name: "Type"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (k *Kalivent) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (k *Kalivent) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
