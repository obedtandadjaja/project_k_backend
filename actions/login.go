package actions

import (
	"fmt"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
	"github.com/obedtandadjaja/project_k_backend/clients"
	"github.com/obedtandadjaja/project_k_backend/helpers"
	"github.com/obedtandadjaja/project_k_backend/models"
	authService "github.com/obedtandadjaja/project_k_backend/services/auth"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Jwt        string    `json:"jwt"`
	SessionJwt string    `json:"session"`
	UserID     uuid.UUID `json:"user_id"`
}

func Login(c buffalo.Context) error {
	req := &LoginRequest{}
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

	user := &models.User{}
	if err := tx.Select("id", "credential_uuid").Where("email = ?", req.Email).First(user); err != nil {
		return c.Error(http.StatusUnauthorized, nil)
	}

	res, err := clients.NewAuthClient().Login(
		tx,
		&authService.LoginRequest{
			CredentialID: user.CredentialUUID.UUID,
			Password:     req.Password,
		},
		c.Request(),
	)
	if err != nil {
		return c.Error(http.StatusUnauthorized, nil)
	}

	jwt, err := helpers.GenerateAccessToken(user.ID.String(), user.CredentialUUID.UUID.String())
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	return c.Render(http.StatusCreated, r.JSON(
		LoginResponse{
			Jwt:        jwt,
			SessionJwt: res.SessionJwt,
			UserID:     user.ID,
		},
	))
}
