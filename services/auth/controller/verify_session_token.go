package controller

import (
	"encoding/json"
	"net/http"

	"github.com/obedtandadjaja/project_k_backend/services/auth/helpers/jwt"
)

type VerifySessionTokenRequest struct {
	SessionJwt string `json:"session"`
}

type VerifySessionTokenResponse struct {
	CredentialUuid string `json:"credential_uuid"`
}

func VerifySessionToken(sr *SharedResources, w http.ResponseWriter, r *http.Request) error {
	request, err := parseVerifySessionTokenRequest(r)
	if err != nil {
		return HandlerError{400, "", err}
	}

	response, err := processVerifySessionTokenRequest(sr, request)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	return nil
}

func parseVerifySessionTokenRequest(r *http.Request) (*VerifySessionTokenRequest, error) {
	var request VerifySessionTokenRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	return &request, err
}

func processVerifySessionTokenRequest(sr *SharedResources, request *VerifySessionTokenRequest) (*VerifySessionTokenResponse, error) {
	var response VerifySessionTokenResponse

	credentialUuid, _, err := jwt.VerifySessionToken(request.SessionJwt)
	if err != nil {
		return &response, HandlerError{401, "Invalid session token", err}
	}

	response.CredentialUuid = credentialUuid
	return &response, nil
}
