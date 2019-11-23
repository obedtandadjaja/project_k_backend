package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/slices"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

type Room struct {
	ID              uuid.UUID  `json:"id,omitempty" db:"id"`
	PropertyID      uuid.UUID  `json:"property_id,omitempty" db:"property_id"`
	Property        *Property  `json:"property,omitempty" belongs_to:property`
	Name            string     `json:"name,omitempty" db:"name"`
	PriceAmount     int        `json:"price_amount,omitempty" db:"price_amount"`
	PaymentSchedule string     `json:"payment_schedule,omitempty" db:"payment_schedule"`
	Data            slices.Map `json:"data,omitempty" db:"data"`
	CreatedAt       time.Time  `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at,omitempty" db:"updated_at"`
	Users           []User     `json:"users,omitempty" many_to_many:"room_occupancies"`
}

// String is not required by pop and may be deleted
func (r Room) String() string {
	jr, _ := json.Marshal(r)
	return string(jr)
}

// Rooms is not required by pop and may be deleted
type Rooms []Room

// String is not required by pop and may be deleted
func (r Rooms) String() string {
	jr, _ := json.Marshal(r)
	return string(jr)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (r *Room) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: r.Name, Name: "Name"},
		&validators.IntIsPresent{Field: r.PriceAmount, Name: "PriceAmount"},
		&validators.StringIsPresent{Field: r.PaymentSchedule, Name: "PaymentSchedule"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (r *Room) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (r *Room) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
