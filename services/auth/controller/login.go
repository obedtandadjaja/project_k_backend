package controller

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/obedtandadjaja/project_k_backend/services/auth/helpers/hash"
	"github.com/obedtandadjaja/project_k_backend/services/auth/helpers/jwt"
	"github.com/obedtandadjaja/project_k_backend/services/auth/models/credential"
	"github.com/obedtandadjaja/project_k_backend/services/auth/models/session"
)

type LoginRequest struct {
	CredentialUuid string `json:"credential_uuid"`
	Email          string `json:"email"`
	Password       string `json:"password"`
}

type LoginResponse struct {
	Jwt            string `json:"jwt"`
	SessionJwt     string `json:"session"`
	CredentialUuid string `json:"credential_uuid"`
}

func Login(sr *SharedResources, w http.ResponseWriter, r *http.Request) error {
	request, err := parseLoginRequest(r)
	if err != nil {
		return HandlerError{400, err.Error(), err}
	}

	response, err := processLoginRequest(sr, request, r)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	return nil
}

func parseLoginRequest(r *http.Request) (*LoginRequest, error) {
	var request LoginRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	if request.Email == "" && request.CredentialUuid == "" {
		return &request, errors.New("Missing required field")
	}

	return &request, err
}

func processLoginRequest(sr *SharedResources, request *LoginRequest, r *http.Request) (*LoginResponse, error) {
	var response LoginResponse
	var cred *credential.Credential
	var err error

	if request.CredentialUuid != "" {
		cred, err = credential.FindBy(sr.DB, map[string]interface{}{
			"uuid": request.CredentialUuid,
		})
	} else {
		cred, err = credential.FindBy(sr.DB, map[string]interface{}{
			"email": request.Email,
		})
	}
	if err != nil {
		return &response, HandlerError{401, "Invalid credentials", err}
	}

	if cred.LockedUntil.Valid && cred.LockedUntil.Time.After(time.Now()) {
		return &response, HandlerError{
			401,
			fmt.Sprintf("Locked until %v", cred.LockedUntil.Time.Sub(time.Now())),
			nil,
		}
	}

	if hashValue := cred.Password.String; !hash.ValidatePasswordHash(request.Password, hashValue) {
		if cred.FailedAttempts == MAX_FAILED_ATTEMPTS {
			cred.Update(sr.DB, map[string]interface{}{
				"locked_until": time.Now().Add(time.Duration(cred.FailedAttempts*10) * time.Minute),
			})
		}
		cred.IncrementFailedAttempt(sr.DB)

		return &response, HandlerError{401, "Invalid credentials", nil}
	}

	go func() {
		cred.Update(sr.DB, map[string]interface{}{
			"failed_attempts": 0,
			"locked_until":    nil,
		})
	}()

	newSession := session.Session{
		CredentialId: cred.Id,
		ExpiresAt:    time.Now().Add(time.Duration(24 * 180 * time.Hour)),
		IpAddress:    sql.NullString{String: r.RemoteAddr, Valid: true},
		UserAgent:    sql.NullString{String: r.UserAgent(), Valid: true},
	}
	err = newSession.Create(sr.DB)
	if err != nil {
		return &response, HandlerError{500, "Internal Server Error", err}
	}

	sessionTokenChan := make(chan string)
	accessTokenChan := make(chan string)

	go func() {
		sessionTokenJwt, _ := jwt.GenerateSessionToken(cred.Uuid, newSession.Uuid)
		sessionTokenChan <- sessionTokenJwt
	}()

	go func() {
		accessTokenJwt, _ := jwt.GenerateAccessToken(cred.Uuid)
		accessTokenChan <- accessTokenJwt
	}()

	response.Jwt = <-accessTokenChan
	response.SessionJwt = <-sessionTokenChan
	response.CredentialUuid = cred.Uuid
	return &response, nil
}
