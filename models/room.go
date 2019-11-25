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
	PropertyID      uuid.UUID  `json:"propertyId,omitempty" db:"property_id"`
	Property        *Property  `json:"property,omitempty" belongs_to:"property"`
	Name            string     `json:"name,omitempty" db:"name"`
	PriceAmount     int        `json:"priceAmount,omitempty" db:"price_amount"`
	PaymentSchedule string     `json:"paymentSchedule,omitempty" db:"payment_schedule"`
	Data            slices.Map `json:"data,omitempty" db:"data"`
	CreatedAt       time.Time  `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt       time.Time  `json:"updatedAt,omitempty" db:"updated_at"`
	Tenants         []User     `json:"tenants" many_to_many:"room_occupancies"`
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
