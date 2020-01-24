package actions

import (
	"fmt"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/slices"
	"github.com/gobuffalo/validate"
	"github.com/obedtandadjaja/project_k_backend/helpers"
	"github.com/obedtandadjaja/project_k_backend/models"
)

type TenantsResource struct {
	buffalo.Resource
}

func (v TenantsResource) getTransactionAndQueryContext(c buffalo.Context) (*pop.Connection, *pop.Query) {
	tx, _ := c.Value("tx").(*pop.Connection)

	return tx, tx.Q().
		InnerJoin("room_occupancies", "room_occupancies.user_id = users.id").
		InnerJoin("rooms", "rooms.id = room_occupancies.room_id").
		InnerJoin("properties", "properties.id = rooms.property_id").
		InnerJoin("user_property_relationships", "user_property_relationships.property_id = properties.id").
		Where("rooms.id = ?", c.Param("room_id")).
		Where("user_property_relationships.user_id = ?", c.Value("current_user_id"))
}

func (v TenantsResource) List(c buffalo.Context) error {
	_, q := v.getTransactionAndQueryContext(c)

	tenants := &models.Users{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q = q.PaginateFromParams(c.Params())
	if c.Param("eager") != "" {
		if err := q.Eager(c.Param("eager")).All(tenants); err != nil {
			return err
		}
	} else {
		if err := q.All(tenants); err != nil {
			return err
		}
	}

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

	return c.Render(http.StatusOK, r.JSON(tenants))
}

func (v TenantsResource) Show(c buffalo.Context) error {
	_, q := v.getTransactionAndQueryContext(c)

	tenant := &models.User{}
	if c.Param("eager") != "" {
		if err := q.Eager(c.Param("eager")).Find(tenant, c.Param("tenant_id")); err != nil {
			return c.Error(http.StatusNotFound, err)
		}
	} else {
		if err := q.Find(tenant, c.Param("tenant_id")); err != nil {
			return c.Error(http.StatusNotFound, err)
		}
	}

	return c.Render(http.StatusOK, r.JSON(tenant))
}

func (v TenantsResource) Create(c buffalo.Context) error {
	user := &models.User{
		Rooms: models.Rooms{
			models.Room{ID: helpers.ParseUUID(c.Param("room_id"))},
		},
		Data:                slices.Map{},
		NotificationMethods: []string{"email"},
		Type:                models.USER_TENANT,
	}
	if err := c.Bind(user); err != nil {
		return err
	}

	tx, _ := v.getTransactionAndQueryContext(c)

	q := tx.Q().
		InnerJoin("properties", "properties.id = rooms.property_id").
		InnerJoin("user_property_relationships", "user_property_relationships.property_id = properties.id").
		Where("user_property_relationships.user_id = ?", c.Value("current_user_id")).
		Where("rooms.id = ?", c.Param("room_id"))
	if err := q.First(&models.Room{}); err != nil {
		verrs := validate.NewErrors()
		verrs.Add("room", "Room does not exist")
		return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
	}

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	verrs, err := tx.ValidateAndCreate(user)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		c.Set("errors", verrs)
		return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
	}

	return c.Render(http.StatusCreated, r.JSON(user))
}

func (v TenantsResource) Update(c buffalo.Context) error {
	tx, q := v.getTransactionAndQueryContext(c)

	user := &models.User{}
	if err := q.Find(user, c.Param("tenant_id")); err != nil {
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
		return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
	}

	return c.Render(http.StatusOK, r.JSON(user))
}

func (v TenantsResource) Destroy(c buffalo.Context) error {
	tx, q := v.getTransactionAndQueryContext(c)

	user := &models.User{}

	if err := q.Find(user, c.Param("tenant_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := tx.Destroy(user); err != nil {
		return err
	}

	return c.Render(http.StatusOK, r.JSON(user))
}
