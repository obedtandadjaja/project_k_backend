package clients

import (
)

type AuthClient struct {
}

var authClient *AuthClient

func NewAuthClient() *AuthClient {
    if authClient != nil {
        return authClient
    }

    return &AuthClient{}
}

// create credentials
func (authClient *AuthClient) CreateCredential(r *CreateCredentialRequest) (*http.Response, error) {
    services.auth.controller.credentials.Create(r)
}

// update credentials
type UpdateCredentialRequest struct {
	CredentialUUID uuid.UUID `json:"uuid"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
}

func (authClient *AuthClient) UpdateCredential(r *UpdateCredentialRequest) (*http.Response, error) {
	requestBody, err := json.Marshal(r)

	client := &http.Client{}
	request, err := http.NewRequest(
		"PUT",
		authClient.AuthAPIUrl+"/credentials/"+r.CredentialUUID.String(),
		bytes.NewBuffer(requestBody),
	)
	res, err := client.Do(request)

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
