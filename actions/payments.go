package actions

import (
	"fmt"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/obedtandadjaja/project_k_backend/models"
)

type PaymentsResource struct {
	buffalo.Resource
}

func (v PaymentsResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	payments := &models.Payments{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all Payments from the DB
	if err := q.All(payments); err != nil {
		return err
	}

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

	return c.Render(http.StatusOK, r.JSON(payments))
}

func (v PaymentsResource) Show(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	payment := &models.Payment{}
	q := tx.Q().
		InnerJoin("room_occupancies", "room_occupancies.id = payments.room_occupancy_id").
		InnerJoin("rooms", "rooms.id = room_occupancies.room_id").
		InnerJoin("properties", "properties.id = rooms.property_id").
		InnerJoin("user_property_relationships", "user_property_relationships.property_id = properties.id").
		Where("room_occupancies.user_id = ?", c.Param("tenant_id")).
		Where("rooms.id = ?", c.Param("room_id")).
		Where("properties.id = ?", c.Param("property_id")).
		Where("user_property_relationships.user_id = ?", c.Value("current_user_id"))
	if err := q.Find(payment, c.Param("payment_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	return c.Render(http.StatusOK, r.JSON(payment))
}

func (v PaymentsResource) Create(c buffalo.Context) error {
	payment := &models.Payment{}

	if err := c.Bind(payment); err != nil {
		return err
	}

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	verrs, err := tx.ValidateAndCreate(payment)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		c.Set("errors", verrs)
		return c.Render(http.StatusUnprocessableEntity, r.JSON(payment))
	}

	return c.Render(http.StatusCreated, r.JSON(payment))
}

func (v PaymentsResource) Update(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	payment := &models.Payment{}

	if err := tx.Find(payment, c.Param("payment_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := c.Bind(payment); err != nil {
		return err
	}

	verrs, err := tx.ValidateAndUpdate(payment)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		c.Set("errors", verrs)
		return c.Render(http.StatusUnprocessableEntity, r.JSON(payment))
	}

	return c.Render(http.StatusOK, r.JSON(payment))
}

func (v PaymentsResource) Destroy(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	payment := &models.Payment{}

	if err := tx.Find(payment, c.Param("payment_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := tx.Destroy(payment); err != nil {
		return err
	}

	return c.Render(http.StatusOK, r.JSON(payment))
}
