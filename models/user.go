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

var USER_ADMIN string = "admin"
var USER_TENANT string = "tenant"
var USER_VALID_TYPES []string = []string{USER_ADMIN, USER_TENANT}

type User struct {
	ID                  uuid.UUID             `json:"id" db:"id"`
	Type                string                `json:"type" db:"type"`
	Name                nulls.String          `json:"name,omitempty" db:"name"`
	CredentialUUID      nulls.UUID            `json:"credentialUUID" db:"credential_uuid"`
	Email               string                `json:"email" db:"email"`
	Phone               nulls.String          `json:"phone" db:"phone"`
	NotificationMethods slices.String         `json:"notificationMethods,omitempty" db:"notification_methods"`
	DeactivatedAt       nulls.Time            `json:"deactivatedAt,omitempty" db:"deactivated_at"`
	Data                slices.Map            `json:"data" db:"data"`
	CreatedAt           time.Time             `json:"createdAt" db:"created_at"`
	UpdatedAt           time.Time             `json:"updatedAt" db:"updated_at"`
	Properties          []Property            `json:"properties" many_to_many:"user_property_relationships"`
	Rooms               []Room                `json:"rooms" many_to_many:"room_occupancies"`
	RoomOccupancies     []RoomOccupancies     `json:"room_occupancies" has_many:"room_occupancies"`
	MaintenanceRequests []MaintenanceRequests `json:"maintenance_requests" has_many:"maintenance_requests"`
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
		&validators.StringInclusion{Field: u.Type, Name: "Type", List: USER_VALID_TYPES},
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
