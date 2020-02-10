package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"
)

type UserPropertyRelationship struct {
	ID         uuid.UUID `json:"id" db:"id"`
	UserID     uuid.UUID `json:"user_id" db:"user_id"`
	User       User      `json:"user,omitempty" belongs_to:"user"`
	PropertyID uuid.UUID `json:"property_id" db:"property_id"`
	Property   Property  `json:"property,omitempty" belongs_to:"property"`
	Type       string    `json:"type" db:"type"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt  time.Time `json:"updatedAt" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (u UserPropertyRelationship) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// UserPropertyRelationships is not required by pop and may be deleted
type UserPropertyRelationships []UserPropertyRelationship

// String is not required by pop and may be deleted
func (u UserPropertyRelationships) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (u *UserPropertyRelationship) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (u *UserPropertyRelationship) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (u *UserPropertyRelationship) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
