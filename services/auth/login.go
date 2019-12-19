package actions

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"github.com/obedtandadjaja/project_k_backend/helpers"
	"github.com/obedtandadjaja/project_k_backend/models"
)

const (
	MAX_FAILED_ATTEMPTS = 5
)

type LoginRequest struct {
	CredentialID uuid.UUID `json:"credential_uuid"`
	Password     string    `json:"password"`
}

type LoginResponse struct {
	CredentialID uuid.UUID `json:"credential_uuid"`
	SessionJwt   string    `json:"session"`
}

func Login(tx *pop.Connection, request *LoginRequest, r http.Request) (*LoginResponse, error) {
	var response LoginResponse

	cred := &models.Credential{}
	if err := tx.Find(cred, request.CredentialID); err != nil {
		return &response, err
	}

	if cred.LockedUntil.Valid && cred.LockedUntil.Time.After(time.Now()) {
		return &response, errors.New(
			fmt.Sprintf("Locked until %v", cred.LockedUntil.Time.Sub(time.Now())),
		)
	}

	if hashValue := cred.Password.String; !helpers.ValidatePasswordHash(request.Password, hashValue) {
		if cred.FailedAttempts == MAX_FAILED_ATTEMPTS {
			cred.LockedUntil = nulls.Time{Time: time.Now().Add(time.Duration(cred.FailedAttempts*10) * time.Minute), Valid: true}
			tx.Update(cred)
		}

		cred.FailedAttempts += 1
		tx.Update(cred)

		return &response, errors.New("Account is locked")
	}

	go func() {
		cred.FailedAttempts = 0
		cred.LockedUntil = nulls.Time{Valid: false}
		tx.Update(cred)
	}()

	newSession := &models.Session{
		CredentialID: cred.ID,
		IpAddress:    nulls.String{String: r.RemoteAddr, Valid: true},
		UserAgent:    nulls.String{String: r.UserAgent(), Valid: true},
	}
	if err := tx.Create(newSession); err != nil {
		return &response, err
	}

	sessionTokenChan := make(chan string)
	go func() {
		sessionTokenJwt, _ := helpers.GenerateSessionToken(cred.ID.String(), newSession.ID.String())
		sessionTokenChan <- sessionTokenJwt
	}()

	response.SessionJwt = <-sessionTokenChan
	response.CredentialID = cred.ID
	return &response, nil
}
