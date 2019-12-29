package actions

import (
	"fmt"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/obedtandadjaja/project_k_backend/models"
)

type UsersResource struct {
	buffalo.Resource
}

func (v UsersResource) Show(c buffalo.Context) error {
	if c.Value("current_user_id") != c.Param("user_id") {
		return c.Render(http.StatusUnauthorized, r.JSON("Unauthorized"))
	}

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	user := &models.User{}
	if err := tx.Find(user, c.Param("user_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	return c.Render(http.StatusOK, r.JSON(user))
}

func (v UsersResource) Update(c buffalo.Context) error {
	if c.Value("current_user_id") != c.Param("user_id") {
		return c.Render(http.StatusUnauthorized, r.JSON("Unauthorized"))
	}

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	user := &models.User{}

	if err := tx.Find(user, c.Value("current_user_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := c.Bind(user); err != nil {
		return err
	}

	verrs, err := tx.ValidateAndUpdate(user)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		c.Set("errors", verrs)
		return c.Render(http.StatusUnprocessableEntity, r.JSON(user))
	}

	return c.Render(http.StatusOK, r.JSON(user))
}
