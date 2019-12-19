package credentials

import (
	"errors"
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"github.com/obedtandadjaja/project_k_backend/helpers"
	"github.com/obedtandadjaja/project_k_backend/models"
)

type ResetPasswordRequest struct {
	CredentialID       uuid.UUID `json:"credential_id"`
	PasswordResetToken string    `json:"password_reset_token"`
	NewPassword        string    `json:"new_password"`
}

type ResetPasswordResponse struct {
	CredentialID uuid.UUID `json:"credential_id"`
}

func ResetPassword(tx *pop.Connection, request *ResetPasswordRequest) (*ResetPasswordResponse, error) {
	var response ResetPasswordResponse

	cred := models.Credential{}
	if err := tx.Eager().Find(cred, request.CredentialID); err != nil {
		return &response, err
	}

	if cred.PasswordResetTokenExpiresAt.Valid && cred.PasswordResetTokenExpiresAt.Time.Before(time.Now()) {
		return &response, errors.New("Invalid password reset token")
	}

	if !helpers.ValidatePasswordHash(request.PasswordResetToken, cred.Password.String) {
		return &response, errors.New("Invalid password reset token")
	}

	newPasswordHashed, _ := helpers.HashPassword(request.NewPassword)
	cred.Password = nulls.String{String: newPasswordHashed, Valid: true}
	if err := tx.Update(cred); err != nil {
		return &response, err
	}

	// delete all sessions to force everyone to login again
	for i := 0; i < len(cred.Sessions); i++ {
		if err := tx.Destroy(cred.Sessions); err != nil {
			return &response, err
		}
	}

	response.CredentialID = cred.ID
	return &response, nil
}
