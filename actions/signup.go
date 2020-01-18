package actions

import (
	"fmt"
	"net/http"

	"github.com/lib/pq"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/obedtandadjaja/project_k_backend/clients"
	"github.com/obedtandadjaja/project_k_backend/helpers"
	"github.com/obedtandadjaja/project_k_backend/models"
	"github.com/obedtandadjaja/project_k_backend/services/auth/credentials"
)

type SignupRequest struct {
	Email    string
	Password string
}

type SignupResponse struct {
	Jwt        string `json:"jwt"`
	SessionJwt string `json:"session"`
	UserID     string `json:"userID"`
}

func Signup(c buffalo.Context) error {
	req := &SignupRequest{}
	if err := c.Bind(req); err != nil {
		return err
	}

	verrs := validate.Validate(
		&validators.StringIsPresent{Field: req.Email, Name: "Email"},
		&validators.StringIsPresent{Field: req.Password, Name: "Password"},
	)
	if verrs.HasAny() {
		c.Set("errors", verrs)
		return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
	}

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	res, err := clients.NewAuthClient().CreateCredential(
		tx,
		&credentials.CreateRequest{
			Password: req.Password,
		},
		c.Request(),
	)
	if err != nil {
		return c.Render(http.StatusInternalServerError, r.JSON("Internal server error"))
	}

	user := &models.User{
		Type:                models.USER_ADMIN,
		Email:               req.Email,
		CredentialUUID:      nulls.NewUUID(res.CredentialID),
		NotificationMethods: []string{"email"},
	}

	err = tx.Create(user)
	if err != nil {
		if pqerr, ok := err.(*pq.Error); ok && pqerr.Code == "23505" {
			return c.Render(http.StatusUnprocessableEntity, r.JSON("Email has been taken"))
		} else {
			return c.Render(http.StatusInternalServerError, r.JSON("Internal server error"))
		}
	}

	token, err := helpers.GenerateAccessToken(user.ID.String(), res.CredentialID.String())
	response := SignupResponse{
		Jwt:        token,
		SessionJwt: res.SessionJwt,
		UserID:     user.ID.String(),
	}
	return c.Render(http.StatusCreated, r.JSON(response))
}
