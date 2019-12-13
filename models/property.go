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

type Property struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	Name      string     `json:"name,omitempty" db:"name"`
	Type      string     `json:"type" db:"type"`
	Address   string     `json:"address" db:"address"`
	Data      slices.Map `json:"data,omitempty" db:"data"`
	CreatedAt time.Time  `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt time.Time  `json:"updatedAt,omitempty" db:"updated_at"`
	Users     []User     `json:"users" many_to_many:"user_property_relationships"`
	Rooms     []Room     `json:"rooms" has_many:"rooms"`
}

// String is not required by pop and may be deleted
func (p Property) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Properties is not required by pop and may be deleted
type Properties []Property

// String is not required by pop and may be deleted
func (p Properties) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (p *Property) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: p.Name, Name: "Name"},
		&validators.StringIsPresent{Field: p.Name, Name: "Type"},
		&validators.StringIsPresent{Field: p.Name, Name: "Address"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (p *Property) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (p *Property) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
