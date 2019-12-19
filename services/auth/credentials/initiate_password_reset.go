package credentials

import (
	"fmt"
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"github.com/obedtandadjaja/project_k_backend/models"
)

type InitiatePasswordResetRequest struct {
	CredentialID uuid.UUID `json:"credential_uuid"`
}

type InitiatePasswordResetResponse struct {
	PasswordResetToken string `json:"password_reset_token"`
}

func InitiatePasswordReset(tx *pop.Connection, request *InitiatePasswordResetRequest) (*InitiatePasswordResetResponse, error) {
	var response InitiatePasswordResetResponse

	cred := models.Credential{}
	if err := tx.Find(cred, request.CredentialID); err != nil {
		return &response, err
	}

	// token is the last 6 digit of the unix nano second - should be unpredictable enough
	token := fmt.Sprintf("%v", time.Now().UnixNano())
	token = token[len(token)-6:]

	cred.PasswordResetToken = nulls.String{String: token, Valid: true}
	if err := tx.Update(cred); err != nil {
		return &response, err
	}

	return &response, nil
}
