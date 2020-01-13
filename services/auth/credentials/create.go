package credentials

import (
	"net/http"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"github.com/obedtandadjaja/project_k_backend/helpers"
	"github.com/obedtandadjaja/project_k_backend/models"
)

type CreateRequest struct {
	Password string `json:"password"`
}

type CreateResponse struct {
	CredentialID uuid.UUID `json:"credential_id"`
	SessionJwt   string    `json:"session"`
}

func Create(tx *pop.Connection, request *CreateRequest, r *http.Request) (*CreateResponse, error) {
	var response CreateResponse

	hashedPassword, _ := helpers.HashPassword(request.Password)
	cred := &models.Credential{
		Password: nulls.String{String: hashedPassword, Valid: true},
	}

	err := tx.Create(cred)
	if err != nil {
		return &response, err
	}

	newSession := &models.Session{
		CredentialID: cred.ID,
		IpAddress:    nulls.String{String: r.RemoteAddr, Valid: true},
		UserAgent:    nulls.String{String: r.UserAgent(), Valid: true},
	}
	err = tx.Create(newSession)
	if err != nil {
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
