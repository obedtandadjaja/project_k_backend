package credentials

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/obedtandadjaja/project_k_backend/services/auth/controller"
	"github.com/obedtandadjaja/project_k_backend/services/auth/models/credential"
)

type InitiatePasswordResetRequest struct {
	CredentialUuid string `json:"credential_uuid"`
}

func InitiatePasswordReset(sr *controller.SharedResources, w http.ResponseWriter, r *http.Request) error {
	request, err := parseInitiatePasswordResetRequest(r)
	if err != nil {
		return controller.HandlerError{400, "", err}
	}

	err = processInitiatePasswordResetRequest(sr, request, r)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)

	return nil
}

func parseInitiatePasswordResetRequest(r *http.Request) (*InitiatePasswordResetRequest, error) {
	var request InitiatePasswordResetRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	return &request, err
}

func processInitiatePasswordResetRequest(sr *controller.SharedResources, request *InitiatePasswordResetRequest, r *http.Request) error {
	cred, err := credential.FindBy(sr.DB, map[string]interface{}{
		"id": request.CredentialUuid,
	})
	if err != nil {
		return controller.HandlerError{404, "Credential not found", err}
	}

	// token is the last 6 digit of the unix nano second - should be unpredictable enough
	token := fmt.Sprintf("%v", time.Now().UnixNano())
	token = token[len(token)-6:]

	err = cred.SetPasswordResetToken(sr.DB, token)
	if err != nil {
		return controller.HandlerError{500, "Failed to initiate password reset", err}
	}

	return nil
}
