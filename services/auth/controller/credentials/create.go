package credentials

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/lib/pq"
	"github.com/obedtandadjaja/project_k_backend/services/auth/helpers/jwt"
	"github.com/obedtandadjaja/project_k_backend/services/auth/controller"
	"github.com/obedtandadjaja/project_k_backend/services/auth/models/credential"
	"github.com/obedtandadjaja/project_k_backend/services/auth/models/session"
)

type CreateRequest struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type CreateResponse struct {
	CredentialUuid string `json:"credential_uuid"`
	Jwt            string `json:"jwt"`
	SessionJwt     string `json:"session"`
}

func Create(sr *controller.SharedResources, w http.ResponseWriter, r *http.Request) error {
	request, err := parseCreateRequest(r)
	if err != nil {
		return controller.HandlerError{400, err.Error(), err}
	}

	response, err := processCreateRequest(sr, request, r)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)

	return nil
}

func parseCreateRequest(r *http.Request) (*CreateRequest, error) {
	var request CreateRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	if request.Password == "" {
		return &request, errors.New("Missing required field")
	}

	return &request, err
}

func processCreateRequest(sr *controller.SharedResources, request *CreateRequest, r *http.Request) (*CreateResponse, error) {
	var response CreateResponse

	cred := credential.Credential{
		Email:    sql.NullString{String: request.Email, Valid: request.Email != ""},
		Phone:    sql.NullString{String: request.Phone, Valid: request.Phone != ""},
		Password: sql.NullString{String: request.Password, Valid: true},
	}

	err := cred.Create(sr.DB)
	if err != nil {
		if pgerr, ok := err.(*pq.Error); ok {
			if pgerr.Code == "23505" {
				return &response, controller.HandlerError{
					400,
					"Email has been taken",
					err,
				}
			}
		}

		return &response, controller.HandlerError{
			500,
			"Failed to create credential",
			err,
		}
	}

	newSession := session.Session{
		CredentialId: cred.Id,
		ExpiresAt:    time.Now().Add(time.Duration(24 * 180 * time.Hour)),
		IpAddress:    sql.NullString{String: r.RemoteAddr, Valid: true},
		UserAgent:    sql.NullString{String: r.UserAgent(), Valid: true},
	}
	err = newSession.Create(sr.DB)
	if err != nil {
		return &response, controller.HandlerError{500, "Failed to create session", err}
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
