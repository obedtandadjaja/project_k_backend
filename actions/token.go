package actions

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/obedtandadjaja/project_k_backend/clients"
	"github.com/obedtandadjaja/project_k_backend/helpers"
	"github.com/obedtandadjaja/project_k_backend/models"
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
		&clients.VerifySessionTokenRequest{
			SessionJwt: req.SessionJwt,
		},
	)
	if err != nil {
		return err
	}

	if res.StatusCode == http.StatusUnauthorized {
		return c.Render(http.StatusUnauthorized, r.JSON("Unauthorized"))
	} else if res.StatusCode != http.StatusOK {
		return c.Render(http.StatusInternalServerError, r.JSON(res.StatusCode))
	}

	var resBody map[string]interface{}
	json.NewDecoder(res.Body).Decode(&resBody)

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	user := &models.User{}
	if err := tx.Select("id").Where("credential_uuid = ?", resBody["credential_uuid"].(string)).First(user); err != nil {
		return c.Error(http.StatusUnauthorized, err)
	}

	token, err := helpers.GenerateAccessToken(user.ID.String(), resBody["credential_uuid"].(string))
	return c.Render(http.StatusCreated, r.JSON(
		TokenResponse{
			Jwt: token,
		},
	))
}
