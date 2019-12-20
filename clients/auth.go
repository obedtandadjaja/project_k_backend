package clients

import (
	"net/http"

	"github.com/gobuffalo/pop"
	authService "github.com/obedtandadjaja/project_k_backend/services/auth"
	"github.com/obedtandadjaja/project_k_backend/services/auth/credentials"
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
func (authClient *AuthClient) CreateCredential(tx *pop.Connection, request *credentials.CreateRequest, r http.Request) (*credentials.CreateResponse, error) {
	return credentials.Create(tx, request, r)
}

func (authClient *AuthClient) Login(tx *pop.Connection, request *authService.LoginRequest, r http.Request) (*authService.LoginResponse, error) {
	return authService.Login(tx, request, r)
}

func (authClient *AuthClient) VerifySessionToken(request *authService.VerifySessionTokenRequest) (*authService.VerifySessionTokenResponse, error) {
	return authService.VerifySessionToken(request)
}
