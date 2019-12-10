package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/slices"
	"github.com/obedtandadjaja/project_k_backend/helpers"
	"github.com/obedtandadjaja/project_k_backend/models"
)

type PropertiesResource struct {
	buffalo.Resource
}

func (v PropertiesResource) getTransactionAndQueryContext(c buffalo.Context) (*pop.Connection, *pop.Query) {
	tx, _ := c.Value("tx").(*pop.Connection)

	return tx, tx.Q().
		InnerJoin("user_property_relationships", "user_property_relationships.property_id = properties.id").
		Where("user_property_relationships.user_id = ?", c.Value("current_user_id"))
}

func (v PropertiesResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	_, q := v.getTransactionAndQueryContext(c)

	properties := &models.Properties{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q = q.PaginateFromParams(c.Params())
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
	_, q := v.getTransactionAndQueryContext(c)

	property := &models.Property{}

	if c.Param("eager") == "true" {
		if err := q.Eager().Find(property, c.Param("property_id")); err != nil {
			return c.Error(http.StatusNotFound, err)
		}
	} else {
		if err := q.Find(property, c.Param("property_id")); err != nil {
			return c.Error(http.StatusNotFound, err)
		}
	}

	return c.Render(http.StatusOK, r.JSON(property))
}

func (v PropertiesResource) Create(c buffalo.Context) error {
	property := &models.Property{
		Users: models.Users{
			models.User{ID: helpers.ParseUUID(c.Value("current_user_id").(string))},
		},
		Data: slices.Map{},
	}

	if err := c.Bind(property); err != nil {
		return err
	}

	tx, _ := v.getTransactionAndQueryContext(c)

	verrs, err := tx.ValidateAndCreate(property)
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
	tx, q := v.getTransactionAndQueryContext(c)

	property := &models.Property{}
	err := q.Find(property, c.Param("property_id"))
	if err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := c.Bind(property); err != nil {
		return err
	}

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
	tx, q := v.getTransactionAndQueryContext(c)

	property := &models.Property{}

	if err := q.Find(property, c.Param("property_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := tx.Destroy(property); err != nil {
		return err
	}

	return c.Render(http.StatusOK, r.JSON(property))
}
