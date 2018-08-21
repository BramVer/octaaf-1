package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type MessageCount struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Count     int       `json:"count" db:"count"`
	Diff      int       `json:"diff" db:"diff"`
}

// String is not required by pop and may be deleted
func (m MessageCount) String() string {
	jm, _ := json.Marshal(m)
	return string(jm)
}

// MessageCounts is not required by pop and may be deleted
type MessageCounts []MessageCount

// String is not required by pop and may be deleted
func (m MessageCounts) String() string {
	jm, _ := json.Marshal(m)
	return string(jm)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (m *MessageCount) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.IntIsPresent{Field: m.Count, Name: "Count"},
		&validators.IntIsPresent{Field: m.Diff, Name: "Diff"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (m *MessageCount) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (m *MessageCount) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func (m *MessageCount) BeforeCreate(tx *pop.Connection) {
	prevMC := MessageCount{}
	err := tx.Last(prevMC)

	m.Diff = 0

	if err == nil && prevMC.Count > 0 {
		m.Diff = (m.Count - prevMC.Count)
	} else {
		// This is the first message
		m.Diff = m.Count
	}
}
