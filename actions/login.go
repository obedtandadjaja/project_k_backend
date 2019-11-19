package actions

import (
	"encoding/json"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/obedtandadjaja/project_k_backend/clients"
)

type LoginRequest struct {
	Email    string
	Password string
}

type LoginResponse struct {
	Jwt        string `json:"jwt"`
	SessionJwt string `json:"session"`
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

	res, err := clients.NewAuthClient().Login(
		&clients.LoginRequest{
			Email:    req.Email,
			Password: req.Password,
		},
	)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusCreated {
		return c.Render(http.StatusInternalServerError, r.JSON("Internal server error"))
	}

	var resBody map[string]interface{}
	json.NewDecoder(res.Body).Decode(&resBody)

	return c.Render(http.StatusCreated, r.JSON(
		LoginResponse{
			Jwt:        resBody["jwt"].(string),
			SessionJwt: resBody["session"].(string),
		},
	))
}
