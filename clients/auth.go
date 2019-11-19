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
		fmt.Sprintf("http://%s:%s%s/api/v1/credentials", authClient.AuthAPIHost, authClient.AuthAPIPort, authClient.AuthAPIPrefix),
		"application/json",
		bytes.NewBuffer(requestBody),
	)

	return res, err
}
