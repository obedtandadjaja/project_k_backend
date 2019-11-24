package actions

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
	"github.com/obedtandadjaja/project_k_backend/clients"
	"github.com/obedtandadjaja/project_k_backend/models"
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

	res, err := clients.NewAuthClient().Login(
		&clients.LoginRequest{
			Email:    req.Email,
			Password: req.Password,
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

	return c.Render(http.StatusCreated, r.JSON(
		LoginResponse{
			Jwt:        resBody["jwt"].(string),
			SessionJwt: resBody["session"].(string),
			UserID:     user.ID,
		},
	))
}
