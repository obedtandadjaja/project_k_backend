package actions

import (
	"fmt"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/slices"
	"github.com/obedtandadjaja/project_k_backend/helpers"
	"github.com/obedtandadjaja/project_k_backend/models"
)

type RoomsResource struct {
	buffalo.Resource
}

func (v RoomsResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	rooms := &models.Rooms{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params()).
		InnerJoin("properties", "properties.id = rooms.property_id").
		InnerJoin("user_property_relationships", "user_property_relationships.property_id = properties.id").
		Where("user_property_relationships.user_id = ?", c.Param("user_id"))
	if c.Param("eager") == "true" {
		if err := q.Eager().All(rooms); err != nil {
			return err
		}
	} else {
		if err := q.All(rooms); err != nil {
			return err
		}
	}

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

	return c.Render(http.StatusOK, r.JSON(rooms))
}

func (v RoomsResource) Show(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	room := &models.Room{}
	q := tx.Q().
		InnerJoin("properties", "properties.id = rooms.property_id").
		InnerJoin("user_property_relationships", "user_property_relationships.property_id = properties.id").
		Where("user_property_relationships.user_id = ?", c.Param("user_id"))
	if c.Param("eager") == "true" {
		if err := q.Eager().Find(room, c.Param("room_id")); err != nil {
			return err
		}
	} else {
		if err := q.Find(room, c.Param("room_id")); err != nil {
			return err
		}
	}

	return c.Render(http.StatusOK, r.JSON(room))
}

func (v RoomsResource) Create(c buffalo.Context) error {
	room := &models.Room{
		Property: &models.Property{ID: helpers.ParseUUID(c.Param("property_id"))},
		Data:     slices.Map{},
	}

	// TODO: check user has ownership of property

	if err := c.Bind(room); err != nil {
		return err
	}

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	verrs, err := tx.ValidateAndCreate(room)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		c.Set("errors", verrs)
		return c.Render(http.StatusUnprocessableEntity, r.JSON(room))
	}

	return c.Render(http.StatusCreated, r.JSON(room))
}

func (v RoomsResource) Update(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	room := &models.Room{}

	if err := tx.Find(room, c.Param("room_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := c.Bind(room); err != nil {
		return err
	}

	verrs, err := tx.ValidateAndUpdate(room)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		c.Set("errors", verrs)
		return c.Render(http.StatusUnprocessableEntity, r.JSON(room))
	}

	return c.Render(http.StatusOK, r.JSON(room))
}

func (v RoomsResource) Destroy(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	room := &models.Room{}

	if err := tx.Find(room, c.Param("room_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := tx.Destroy(room); err != nil {
		return err
	}

	return c.Render(http.StatusOK, r.JSON(room))
}
