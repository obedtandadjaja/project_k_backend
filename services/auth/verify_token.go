package actions

import (
	"github.com/obedtandadjaja/project_k_backend/helpers"
)

type VerifyTokenRequest struct {
	Jwt string `json:"jwt"`
}

type VerifyTokenResponse struct {
	Verified bool `json:"verified"`
}

func VerifyToken(request *VerifyTokenRequest) (*VerifyTokenResponse, error) {
	var response VerifyTokenResponse

	_, err := helpers.VerifyAccessToken(request.Jwt)
	if err != nil {
		response.Verified = false
		return &response, nil
	}

	response.Verified = true
	return &response, nil
}
