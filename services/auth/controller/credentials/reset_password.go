package credentials

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/obedtandadjaja/project_k_backend/services/auth/helpers/hash"
	"github.com/obedtandadjaja/project_k_backend/services/auth/controller"
	"github.com/obedtandadjaja/project_k_backend/services/auth/models/credential"
)

type ResetPasswordRequest struct {
	CredentialId       string `json:"credential_uuid"`
	PasswordResetToken string `json:"password_reset_token"`
	NewPassword        string `json:"new_password"`
}

type ResetPasswordResponse struct {
	Uuid string `json:"uuid"`
}

func ResetPassword(sr *controller.SharedResources, w http.ResponseWriter, r *http.Request) error {
	request, err := parseResetPasswordRequest(r)
	if err != nil {
		return controller.HandlerError{400, "", err}
	}

	response, err := processResetPasswordRequest(sr, request, r)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	return nil
}

func parseResetPasswordRequest(r *http.Request) (*ResetPasswordRequest, error) {
	var request ResetPasswordRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	return &request, err
}

func processResetPasswordRequest(sr *controller.SharedResources, request *ResetPasswordRequest, r *http.Request) (*ResetPasswordResponse, error) {
	var response ResetPasswordResponse

	cred, err := credential.FindBy(sr.DB, map[string]interface{}{
		"id": request.CredentialId,
	})
	if err != nil {
		return &response, controller.HandlerError{404, "Credential not found", err}
	}

	if !cred.PasswordResetToken.Valid {
		return &response, controller.HandlerError{400, "Credential did not apply for password reset", err}
	}

	if cred.PasswordResetTokenExpiresAt.Valid && cred.PasswordResetTokenExpiresAt.Time.Before(time.Now()) {
		return &response, controller.HandlerError{401, "Wrong password reset token", err}
	}

	if !hash.ValidatePasswordHash(request.PasswordResetToken, cred.Password.String) {
		return &response, controller.HandlerError{401, "Wrong password reset token", err}
	}

	response.Uuid = cred.Uuid
	return &response, nil
}
