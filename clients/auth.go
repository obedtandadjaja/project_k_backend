package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gobuffalo/envy"
)

type AuthClient struct {
	AuthAPIHost   string
	AuthAPIPort   string
	AuthAPIPrefix string
	AuthAPIUrl    string
}

var authClient *AuthClient

func NewAuthClient() *AuthClient {
	if authClient != nil {
		return authClient
	}

	authClient = &AuthClient{
		AuthAPIHost:   envy.Get("AUTH_API_HOST", "localhost"),
		AuthAPIPort:   envy.Get("AUTH_API_PORT", "3000"),
		AuthAPIPrefix: envy.Get("AUTH_API_PREFIX", ""),
		// AuthAPIUrl: fmt.Sprintf(
		// 	"http://%s:%s%s",
		// 	envy.Get("AUTH_API_HOST", "localhost"),
		// 	envy.Get("AUTH_API_PORT", "3000"),
		// 	envy.Get("AUTH_API_PREFIX", ""),
		// ),
		AuthAPIUrl: "http://api.stage.obedt.com/auth",
	}

	return authClient
}

// create credentials
type CreateCredentialRequest struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

func (authClient *AuthClient) CreateCredential(r *CreateCredentialRequest) (*http.Response, error) {
	requestBody, err := json.Marshal(r)

	fmt.Println("Calling auth service" + authClient.AuthAPIUrl)
	res, err := http.Post(
		authClient.AuthAPIUrl+"/api/v1/credentials",
		"application/json",
		bytes.NewBuffer(requestBody),
	)
	fmt.Println("Done calling auth service")

	return res, err
}

// login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (authClient *AuthClient) Login(r *LoginRequest) (*http.Response, error) {
	requestBody, err := json.Marshal(r)

	res, err := http.Post(
		authClient.AuthAPIUrl+"/api/v1/login",
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
		authClient.AuthAPIUrl+"/api/v1/verify_session_token",
		"application/json",
		bytes.NewBuffer(requestBody),
	)

	return res, err
}
