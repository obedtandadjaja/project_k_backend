package actions

import (
	"fmt"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/obedtandadjaja/project_k_backend/helpers"
	"github.com/obedtandadjaja/project_k_backend/models"
)

type PropertiesResource struct {
	buffalo.Resource
}

func (v PropertiesResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	properties := &models.Properties{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params()).
		InnerJoin("user_property_relationships", "user_property_relationships.property_id = properties.id").
		Where("user_property_relationships.user_id = ?", c.Param("user_id"))
	if c.Param("eager") == "true" {
		if err := q.Eager().All(properties); err != nil {
			return err
		}
	} else {
		if err := q.All(properties); err != nil {
			return err
		}
	}

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

	return c.Render(http.StatusOK, r.JSON(properties))
}

func (v PropertiesResource) Show(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	property := &models.Property{}

	err := tx.Q().
		InnerJoin("user_property_relationships", "user_property_relationships.property_id = properties.id").
		Where("user_property_relationships.user_id = ?", c.Param("user_id")).
		Find(property, c.Param("property_id"))
	if err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	return c.Render(http.StatusOK, r.JSON(property))
}

func (v PropertiesResource) Create(c buffalo.Context) error {
	property := &models.Property{}
	property.Users = models.Users{
		models.User{ID: helpers.ParseUUID(c.Param("user_id"))},
	}

	if err := c.Bind(property); err != nil {
		return err
	}
	fmt.Println(property)

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	verrs, err := tx.Eager().ValidateAndCreate(property)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		c.Set("errors", verrs)
		return c.Render(http.StatusUnprocessableEntity, r.JSON(property))
	}

	return c.Render(http.StatusCreated, r.JSON(property))
}

func (v PropertiesResource) Update(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	property := &models.Property{}
	err := tx.Q().
		InnerJoin("user_property_relationships", "user_property_relationships.property_id = properties.id").
		Where("user_property_relationships.user_id = ?", c.Param("user_id")).
		Find(property, c.Param("property_id"))
	if err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := c.Bind(property); err != nil {
		return err
	}

	fmt.Println(property)

	verrs, err := tx.ValidateAndUpdate(property)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		c.Set("errors", verrs)
		return c.Render(http.StatusUnprocessableEntity, r.JSON(property))
	}

	return c.Render(http.StatusOK, r.JSON(property))
}

func (v PropertiesResource) Destroy(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	property := &models.Property{}

	if err := tx.Find(property, c.Param("property_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := tx.Destroy(property); err != nil {
		return err
	}

	return c.Render(http.StatusOK, r.JSON(property))
}
