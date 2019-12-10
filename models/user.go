package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/slices"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

type User struct {
	ID                  uuid.UUID         `json:"id,omitempty" db:"id"`
	Name                nulls.String      `json:"name,omitempty" db:"name"`
	CredentialUUID      nulls.UUID        `json:"credentialUUID,omitmepty" db:"credential_uuid"`
	Email               string            `json:"email,omitempty" db:"email"`
	Phone               nulls.String      `json:"phone,omityempty" db:"phone"`
	NotificationMethods slices.String     `json:"notificationMethods,omitempty" db:"notification_methods"`
	DeactivatedAt       nulls.Time        `json:"deactivatedAt,omitempty" db:"deactivated_at"`
	Data                slices.Map        `json:"data,omitempty" db:"data"`
	CreatedAt           time.Time         `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt           time.Time         `json:"updatedAt,omitempty" db:"updated_at"`
	Properties          []Property        `json:"properties" many_to_many:"user_property_relationships"`
	Rooms               []Room            `json:"rooms" many_to_many:"room_occupancies"`
	RoomOccupancies     []RoomOccupancies `json:"room_occupancies" has_many:"room_occupancies"`
}

// String is not required by pop and may be deleted
func (u User) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Users is not required by pop and may be deleted
type Users []User

// String is not required by pop and may be deleted
func (u Users) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (u *User) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: u.Email, Name: "Email"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (u *User) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (u *User) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func (u User) TableName() string {
	return "users"
}
