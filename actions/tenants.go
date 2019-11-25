package actions

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/slices"
	"github.com/gofrs/uuid"
	"github.com/obedtandadjaja/project_k_backend/clients"
	"github.com/obedtandadjaja/project_k_backend/helpers"
	"github.com/obedtandadjaja/project_k_backend/models"
)

type TenantsResource struct {
	buffalo.Resource
}

func (v TenantsResource) List(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	users := &models.Users{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params()).
		InnerJoin("room_occupancies", "room_occupancies.user_id = users.id").
		InnerJoin("rooms", "rooms.id = room_occupancies.room_id").
		InnerJoin("properties", "properties.id = rooms.property_id").
		InnerJoin("user_property_relationships", "user_property_relationships.property_id = properties.id").
		Where("rooms.id = ?", c.Param("room_id")).
		Where("properties.id = ?", c.Param("property_id")).
		Where("user_property_relationships.user_id = ?", c.Value("current_user_id"))
	if err := q.All(users); err != nil {
		return err
	}

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

	return c.Render(http.StatusOK, r.JSON(users))
}

func (v TenantsResource) Show(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	user := &models.User{}
	q := tx.Q().
		InnerJoin("room_occupancies", "room_occupancies.user_id = users.id").
		InnerJoin("rooms", "rooms.id = room_occupancies.room_id").
		InnerJoin("properties", "properties.id = rooms.property_id").
		InnerJoin("user_property_relationships", "user_property_relationships.property_id = properties.id").
		Where("rooms.id = ?", c.Param("room_id")).
		Where("properties.id = ?", c.Param("property_id")).
		Where("user_property_relationships.user_id = ?", c.Value("current_user_id"))
	if err := q.Find(user, c.Param("tenant_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	return c.Render(http.StatusOK, r.JSON(user))
}

func (v TenantsResource) Create(c buffalo.Context) error {
	user := &models.User{
		Rooms: models.Rooms{
			models.Room{ID: helpers.ParseUUID(c.Param("room_id"))},
		},
		Data: slices.Map{},
	}

	// TODO: check user has ownership of property

	if err := c.Bind(user); err != nil {
		return err
	}

	// start of generating random credential on auth server
	dummyPassword, _ := helpers.GenerateRandomString(15)
	res, err := clients.NewAuthClient().CreateCredential(
		&clients.CreateCredentialRequest{
			Email:    user.Email,
			Phone:    user.Phone.String,
			Password: dummyPassword,
		},
	)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusCreated {
		return c.Render(http.StatusInternalServerError, r.JSON("Internal server error"))
	}
	// end of generating random credential on auth server

	var resBody map[string]interface{}
	json.NewDecoder(res.Body).Decode(&resBody)

	credentialUUID, _ := uuid.FromString(resBody["credential_uuid"].(string))
	user.CredentialUUID = nulls.NewUUID(credentialUUID)
	user.NotificationMethods = []string{"email"}

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
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	user := &models.User{}
	q := tx.Q().
		InnerJoin("room_occupancies", "room_occupancies.user_id = users.id").
		InnerJoin("rooms", "rooms.id = room_occupancies.room_id").
		InnerJoin("properties", "properties.id = rooms.property_id").
		InnerJoin("user_property_relationships", "user_property_relationships.property_id = properties.id").
		Where("rooms.id = ?", c.Param("room_id")).
		Where("properties.id = ?", c.Param("property_id")).
		Where("user_property_relationships.user_id = ?", c.Value("current_user_id"))
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
		return c.Render(http.StatusUnprocessableEntity, r.JSON(user))
	}

	return c.Render(http.StatusOK, r.JSON(user))
}

func (v TenantsResource) Destroy(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	user := &models.User{}

	if err := tx.Find(user, c.Param("tenant_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := tx.Destroy(user); err != nil {
		return err
	}

	return c.Render(http.StatusOK, r.JSON(user))
}
