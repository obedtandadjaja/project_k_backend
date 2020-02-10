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

type Session struct {
	ID             uuid.UUID    `json:"id" db:"id"`
	CredentialID   uuid.UUID    `json:"credentialId" db:"credential_id"`
	Credential     Credential   `json:"credential" belongs_to:"credential"`
	IpAddress      nulls.String `json:"ipAddress,omitempty" db:"ip_address"`
	UserAgent      nulls.String `json:"userAgent,omitempty" db:"user_agent"`
	LastAccessedAt time.Time    `json:"lastAccessedAt" db:"last_accessed_at"`
	CreatedAt      time.Time    `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time    `json:"updatedAt" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (s Session) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Sessions is not required by pop and may be deleted
type Sessions []Session

// String is not required by pop and may be deleted
func (s Sessions) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (s *Session) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.TimeIsPresent{Field: s.LastAccessedAt, Name: "LastAccessedAt"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (s *Session) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (s *Session) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
