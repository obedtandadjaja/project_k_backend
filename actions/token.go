package actions

import (
	"fmt"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/obedtandadjaja/project_k_backend/clients"
	"github.com/obedtandadjaja/project_k_backend/helpers"
	"github.com/obedtandadjaja/project_k_backend/models"
	authService "github.com/obedtandadjaja/project_k_backend/services/auth"
)

type TokenRequest struct {
	SessionJwt string `json:"session"`
}

type TokenResponse struct {
	Jwt string `json:"jwt"`
}

func Token(c buffalo.Context) error {
	req := &TokenRequest{}
	if err := c.Bind(req); err != nil {
		return c.Render(http.StatusBadRequest, r.JSON("Bad request"))
	}

	res, err := clients.NewAuthClient().VerifySessionToken(
		&authService.VerifySessionTokenRequest{
			SessionJwt: req.SessionJwt,
		},
	)
	if err != nil {
		c.Render(http.StatusUnauthorized, nil)
	}

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	user := &models.User{}
	if err := tx.Select("id").Where("credential_uuid = ?", res.CredentialID).First(user); err != nil {
		return c.Error(http.StatusUnauthorized, err)
	}

	token, err := helpers.GenerateAccessToken(user.ID.String(), res.CredentialID.String(),
		user.Type)
	return c.Render(http.StatusCreated, r.JSON(
		TokenResponse{
			Jwt: token,
		},
	))
}
