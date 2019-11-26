package actions

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
	"github.com/obedtandadjaja/project_k_backend/clients"
	"github.com/obedtandadjaja/project_k_backend/models"
)

type SignupRequest struct {
	Email    string
	Password string
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

	res, err := clients.NewAuthClient().CreateCredential(
		&clients.CreateCredentialRequest{
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

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	credentialUUID, _ := uuid.FromString(resBody["credential_uuid"].(string))
	user := &models.User{
		Email:               req.Email,
		CredentialUUID:      nulls.NewUUID(credentialUUID),
		NotificationMethods: []string{"email"},
	}

	err = tx.Create(user)
	if err != nil {
		return err
	}

	return c.Render(http.StatusCreated, r.JSON(user))
}
