package actions

import (
	"errors"

	"github.com/gofrs/uuid"
	"github.com/obedtandadjaja/project_k_backend/helpers"
)

type VerifySessionTokenRequest struct {
	SessionJwt string `json:"session"`
}

type VerifySessionTokenResponse struct {
	CredentialID uuid.UUID `json:"credential_uuid"`
}

func VerifySessionToken(request *VerifySessionTokenRequest) (*VerifySessionTokenResponse, error) {
	var response VerifySessionTokenResponse

	credentialID, _, err := helpers.VerifySessionToken(request.SessionJwt)
	if err != nil {
		return &response, errors.New("Invalid session token")
	}

	credentialUUID, _ := uuid.FromString(credentialID)
	response.CredentialID = credentialUUID
	return &response, nil
}
