package actions

import (
	"fmt"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/obedtandadjaja/project_k_backend/models"
)

type RoomOccupanciesResource struct {
	buffalo.Resource
}

func (v RoomOccupanciesResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	roomOccupancies := &models.RoomOccupancies{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all RoomOccupancies from the DB
	if err := q.All(roomOccupancies); err != nil {
		return err
	}

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

	return c.Render(http.StatusOK, r.JSON(roomOccupancies))
}

func (v RoomOccupanciesResource) Show(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	roomOccupancy := &models.RoomOccupancy{}

	if err := tx.Find(roomOccupancy, c.Param("room_occupancy_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	return c.Render(http.StatusOK, r.JSON(roomOccupancy))
}

func (v RoomOccupanciesResource) Create(c buffalo.Context) error {
	roomOccupancy := &models.RoomOccupancy{}

	if err := c.Bind(roomOccupancy); err != nil {
		return err
	}

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	verrs, err := tx.ValidateAndCreate(roomOccupancy)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		c.Set("errors", verrs)
		return c.Render(http.StatusUnprocessableEntity, r.JSON(roomOccupancy))
	}

	return c.Render(http.StatusCreated, r.JSON(roomOccupancy))
}

func (v RoomOccupanciesResource) Update(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	roomOccupancy := &models.RoomOccupancy{}

	if err := tx.Find(roomOccupancy, c.Param("room_occupancy_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := c.Bind(roomOccupancy); err != nil {
		return err
	}

	verrs, err := tx.ValidateAndUpdate(roomOccupancy)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		c.Set("errors", verrs)
		return c.Render(http.StatusUnprocessableEntity, r.JSON(roomOccupancy))
	}

	return c.Render(http.StatusOK, r.JSON(roomOccupancy))
}

func (v RoomOccupanciesResource) Destroy(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	roomOccupancy := &models.RoomOccupancy{}

	if err := tx.Find(roomOccupancy, c.Param("room_occupancy_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := tx.Destroy(roomOccupancy); err != nil {
		return err
	}

	return c.Render(http.StatusOK, r.JSON(roomOccupancy))
}
