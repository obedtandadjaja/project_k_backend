package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

type Payment struct {
	ID              uuid.UUID      `json:"id" db:"id"`
	Amount          int            `json:"amount" db:"amount"`
	Description     nulls.String   `json:"description" db:"description"`
	RoomOccupancyID uuid.UUID      `json:"roomOccupancyId" db:"room_occupancy_id"`
	RoomOccupancy   *RoomOccupancy `json:"roomOccupancy,omitempty" belongs_to:"room_occupancy"`
	PaymentDate     time.Time      `json:"paymentDate,omitempty" db:"payment_date"`
	CreatedAt       time.Time      `json:"createdAt" db:"created_at"`
	UpdatedAt       time.Time      `json:"updatedAt" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (p Payment) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Payments is not required by pop and may be deleted
type Payments []Payment

// String is not required by pop and may be deleted
func (p Payments) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (p *Payment) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.IntIsPresent{Field: p.Amount, Name: "Amount"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (p *Payment) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (p *Payment) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
