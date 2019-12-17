package controller

import (
	"encoding/json"
	"net/http"

	"github.com/obedtandadjaja/project_k_backend/services/auth/helpers/jwt"
)

type VerifyRequest struct {
	Jwt string `json:"jwt"`
}

type VerifyResponse struct {
	Verified bool `json:"verified"`
}

func Verify(sr *SharedResources, w http.ResponseWriter, r *http.Request) error {
	request, err := parseVerifyRequest(r)
	if err != nil {
		return HandlerError{400, "", err}
	}

	response, err := processVerifyRequest(sr, request)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	return nil
}

func parseVerifyRequest(r *http.Request) (*VerifyRequest, error) {
	var request VerifyRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	return &request, err
}

func processVerifyRequest(sr *SharedResources, request *VerifyRequest) (*VerifyResponse, error) {
	var response VerifyResponse

	_, err := jwt.VerifyAccessToken(request.Jwt)
	if err != nil {
		return &response, HandlerError{401, "Invalid JWT token", err}
	}

	response.Verified = true
	return &response, nil
}
