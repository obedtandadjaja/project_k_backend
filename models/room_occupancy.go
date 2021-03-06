package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"
)

type RoomOccupancy struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	UserID       uuid.UUID  `json:"user_id" db:"user_id"`
	User         *User      `json:"user,omitempty" belongs_to:"user"`
	RoomID       uuid.UUID  `json:"room_id" db:"room_id"`
	Room         *Room      `json:"room,omitempty" belongs_to:"room"`
	TerminatedAt nulls.Time `json:"terminatedAt" db:"terminated_at"`
	Type         string     `json:"type" db:"type"`
	CreatedAt    time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time  `json:"updatedAt" db:"updated_at"`
	Payments     []Payment  `json:"payments" has_many:"payments"`
}

// String is not required by pop and may be deleted
func (r RoomOccupancy) String() string {
	jr, _ := json.Marshal(r)
	return string(jr)
}

// RoomOccupancies is not required by pop and may be deleted
type RoomOccupancies []RoomOccupancy

// String is not required by pop and may be deleted
func (r RoomOccupancies) String() string {
	jr, _ := json.Marshal(r)
	return string(jr)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (r *RoomOccupancy) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (r *RoomOccupancy) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (r *RoomOccupancy) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
