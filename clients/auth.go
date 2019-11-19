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
		AuthAPIUrl: fmt.Sprintf(
			"https://%s:%s%s/api/v1/verify_token",
			envy.Get("AUTH_API_HOST", "localhost"),
			envy.Get("AUTH_API_PORT", "3000"),
			envy.Get("AUTH_API_PREFIX", ""),
		),
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

	res, err := http.Post(
		authClient.AuthAPIUrl,
		"application/json",
		bytes.NewBuffer(requestBody),
	)

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
		authClient.AuthAPIUrl,
		"application/json",
		bytes.NewBuffer(requestBody),
	)

	return res, err
}

// Token
type TokenRequest struct {
	SessionJwt string `json:"session"`
}

func (authClient *AuthClient) Token(r *LoginRequest) (*http.Response, error) {
	requestBody, err := json.Marshal(r)

	res, err := http.Post(
		authClient.AuthAPIUrl,
		"application/json",
		bytes.NewBuffer(requestBody),
	)

	return res, err
}

// Verify token
type VerifyTokenRequest struct {
	Jwt string `json:"jwt"`
}

func (authClient *AuthClient) VerifyToken(r *LoginRequest) (*http.Response, error) {
	requestBody, err := json.Marshal(r)

	res, err := http.Post(
		authClient.AuthAPIUrl,
		"application/json",
		bytes.NewBuffer(requestBody),
	)

	return res, err
}
