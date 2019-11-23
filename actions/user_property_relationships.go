package actions

import (
	"fmt"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/obedtandadjaja/project_k_backend/models"
)

type UserPropertyRelationshipsResource struct {
	buffalo.Resource
}

func (v UserPropertyRelationshipsResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	userPropertyRelationships := &models.UserPropertyRelationships{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all UserPropertyRelationships from the DB
	if err := q.All(userPropertyRelationships); err != nil {
		return err
	}

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

	return c.Render(http.StatusOK, r.JSON(userPropertyRelationships))
}

func (v UserPropertyRelationshipsResource) Show(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	userPropertyRelationship := &models.UserPropertyRelationship{}

	if err := tx.Find(userPropertyRelationship, c.Param("user_property_relationship_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	return c.Render(http.StatusOK, r.JSON(userPropertyRelationship))
}

func (v UserPropertyRelationshipsResource) Create(c buffalo.Context) error {
	userPropertyRelationship := &models.UserPropertyRelationship{}

	if err := c.Bind(userPropertyRelationship); err != nil {
		return err
	}

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	verrs, err := tx.ValidateAndCreate(userPropertyRelationship)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		c.Set("errors", verrs)
		return c.Render(http.StatusUnprocessableEntity, r.JSON(userPropertyRelationship))
	}

	return c.Render(http.StatusCreated, r.JSON(userPropertyRelationship))
}

func (v UserPropertyRelationshipsResource) Update(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	userPropertyRelationship := &models.UserPropertyRelationship{}

	if err := tx.Find(userPropertyRelationship, c.Param("user_property_relationship_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := c.Bind(userPropertyRelationship); err != nil {
		return err
	}

	verrs, err := tx.ValidateAndUpdate(userPropertyRelationship)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		c.Set("errors", verrs)
		return c.Render(http.StatusUnprocessableEntity, r.JSON(userPropertyRelationship))
	}

	return c.Render(http.StatusOK, r.JSON(userPropertyRelationship))
}

func (v UserPropertyRelationshipsResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	userPropertyRelationship := &models.UserPropertyRelationship{}

	if err := tx.Find(userPropertyRelationship, c.Param("user_property_relationship_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := tx.Destroy(userPropertyRelationship); err != nil {
		return err
	}

	return c.Render(http.StatusOK, r.JSON(userPropertyRelationship))
}
