package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gobuffalo/envy"
	"github.com/gofrs/uuid"

	"github.com/obedtandadjaja/project_k_backend/services/auth/controller/credentials/create"
	"github.com/obedtandadjaja/project_k_backend/services/auth/controller/credentials/update"
)

type AuthClient struct{}

var authClient *AuthClient

func NewAuthClient() *AuthClient {
	if authClient != nil {
		return authClient
	}

	return &AuthClient{}
}

// create credentials
func (authClient *AuthClient) CreateCredential(r *create.CreateRequest) (*create.CreateResponse, error) {
	return create.processCreateRequest(r)
}

// update credentials
func (authClient *AuthClient) UpdateCredential(r *UpdateCredentialRequest) (*http.Response, error) {
}

// login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (authClient *AuthClient) Login(r *LoginRequest) (*http.Response, error) {
	requestBody, err := json.Marshal(r)

	res, err := http.Post(
		authClient.AuthAPIUrl+"/login",
		"application/json",
		bytes.NewBuffer(requestBody),
	)

	return res, err
}

// Verify session token
type VerifySessionTokenRequest struct {
	SessionJwt string `json:"session"`
}

func (authClient *AuthClient) VerifySessionToken(r *VerifySessionTokenRequest) (*http.Response, error) {
	requestBody, err := json.Marshal(r)

	res, err := http.Post(
		authClient.AuthAPIUrl+"/verify_session_token",
		"application/json",
		bytes.NewBuffer(requestBody),
	)

	return res, err
}
