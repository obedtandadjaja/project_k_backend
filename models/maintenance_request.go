package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
	customValidators "github.com/obedtandadjaja/project_k_backend/helpers/validators"
)

type MaintenanceRequest struct {
	ID          uuid.UUID    `json:"id" db:"id"`
	PropertyID  nulls.UUID   `json:"propertyID,omitempty" db:"property_id"`
	Property    Property     `json:"property,omitempty" belongs_to:"property"`
	RoomID      nulls.UUID   `json:"roomID,omitempty" db:"room_id"`
	Room        *Room        `json:"room,omitempty" belongs_to:"room"`
	ReporterID  uuid.UUID    `json:"reporterID,omitempty" db:"reporter_id"`
	Reporter    *User        `json:"reporter,omitempty" belongs_to:"user"`
	Status      string       `json:"status,omitempty" db:"status"`
	Title       string       `json:"title" db:"title"`
	Description nulls.String `json:"description" db:"description"`
	CompletedAt nulls.Time   `json:"completedAt" db:"completed_at"`
	CreatedAt   time.Time    `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time    `json:"updatedAt" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (m MaintenanceRequest) String() string {
	jm, _ := json.Marshal(m)
	return string(jm)
}

// MaintenanceRequests is not required by pop and may be deleted
type MaintenanceRequests []MaintenanceRequest

// String is not required by pop and may be deleted
func (m MaintenanceRequests) String() string {
	jm, _ := json.Marshal(m)
	return string(jm)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (m *MaintenanceRequest) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: m.Status, Name: "Status"},
		&validators.StringIsPresent{Field: m.Title, Name: "Title"},
		&customValidators.AtLeastOne{
			Name: "BelongsTo",
			Validators: []validate.Validator{
				&validators.UUIDIsPresent{Field: m.PropertyID.UUID, Name: "PropertyID"},
				&validators.UUIDIsPresent{Field: m.RoomID.UUID, Name: "RoomID"},
			},
		},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (m *MaintenanceRequest) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (m *MaintenanceRequest) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
