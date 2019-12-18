package models

import (
	"encoding/json"
	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
	"time"
)

type Credential struct {
	ID                          uuid.UUID    `json:"id" db:"id"`
	Password                    nulls.String `json:"password" db:"password"`
	FailedAttempts              int          `json:"failedAttempts" db:"failed_attempts"`
	LockedUntil                 nulls.Time   `json:"lockedUntil" db:"locked_until"`
	PasswordResetToken          nulls.String `json:"passwordResetToken" db:"password_reset_token"`
	PasswordResetTokenExpiresAt nulls.Time   `json:"passwordResetTokenExpiresAt" db:"password_reset_token_expires_at"`
	CreatedAt                   time.Time    `json:"createdAt" db:"created_at"`
	UpdatedAt                   time.Time    `json:"updatedAt" db:"updated_at"`
	User                        *User        `json:"user" belongs_to:"user"`
}

// String is not required by pop and may be deleted
func (c Credential) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

// Credentials is not required by pop and may be deleted
type Credentials []Credential

// String is not required by pop and may be deleted
func (c Credentials) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (c *Credential) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.IntIsPresent{Field: c.FailedAttempts, Name: "FailedAttempts"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (c *Credential) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (c *Credential) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
