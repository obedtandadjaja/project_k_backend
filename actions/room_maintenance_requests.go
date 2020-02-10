package actions

import (
	"net/http"
	"strings"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/obedtandadjaja/project_k_backend/helpers"
	"github.com/obedtandadjaja/project_k_backend/models"
)

type RoomMaintenanceRequestsResource struct {
	buffalo.Resource
}

func (v RoomMaintenanceRequestsResource) getTransactionAndQueryContext(c buffalo.Context) (*pop.Connection, *pop.Query) {
	tx, _ := c.Value("tx").(*pop.Connection)

	return tx, tx.Q().
		InnerJoin("rooms", "rooms.id = maintenance_requests.room_id").
		InnerJoin("properties", "properties.id = rooms.property_id").
		InnerJoin("user_property_relationships", "user_property_relationships.property_id = properties.id").
		Where("user_property_relationships.user_id = ?", c.Value("current_user_id"))
}

func (v RoomMaintenanceRequestsResource) List(c buffalo.Context) error {
	_, q := v.getTransactionAndQueryContext(c)

	maintenanceRequests := &models.MaintenanceRequests{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q = q.PaginateFromParams(c.Params())
	if c.Param("eager") != "" {
		if err := q.Eager(strings.Split(c.Param("eager"), ",")...).All(maintenanceRequests); err != nil {
			return err
		}
	} else {
		if err := q.All(maintenanceRequests); err != nil {
			return err
		}
	}

	c.Set("pagination", q.Paginator)

	return c.Render(http.StatusOK, r.JSON(maintenanceRequests))
}

func (v RoomMaintenanceRequestsResource) Show(c buffalo.Context) error {
	_, q := v.getTransactionAndQueryContext(c)

	maintenanceRequest := &models.MaintenanceRequest{}

	if c.Param("eager") != "" {
		if err := q.Eager(strings.Split(c.Param("eager"), ",")...).Find(maintenanceRequest, c.Param("room_maintenance_request_id")); err != nil {
			return c.Error(http.StatusNotFound, err)
		}
	} else {
		if err := q.Find(maintenanceRequest, c.Param("room_maintenance_request_id")); err != nil {
			return c.Error(http.StatusNotFound, err)
		}
	}

	return c.Render(http.StatusOK, r.JSON(maintenanceRequest))
}

func (v RoomMaintenanceRequestsResource) Create(c buffalo.Context) error {
	maintenanceRequest := &models.MaintenanceRequest{
		RoomID:     nulls.UUID{UUID: helpers.ParseUUID(c.Param("room_id")), Valid: true},
		PropertyID: nulls.UUID{UUID: helpers.ParseUUID(c.Param("property_id")), Valid: true},
		ReporterID: helpers.ParseUUID(c.Value("current_user_id").(string)),
		Status:     "pending",
	}
	if err := c.Bind(maintenanceRequest); err != nil {
		return err
	}

	tx, _ := v.getTransactionAndQueryContext(c)
	q := tx.Q().
		InnerJoin("properties", "properties.id = rooms.property_id").
		InnerJoin("user_property_relationships", "user_property_relationships.property_id = properties.id").
		Where("user_property_relationships.user_id = ?", c.Value("current_user_id")).
		Where("properties.id = ?", c.Param("property_id")).
		Where("rooms.id = ?", c.Param("room_id"))
	if err := q.First(&models.Room{}); err != nil {
		verrs := validate.NewErrors()
		verrs.Add("room", "Room does not exist")
		return c.Render(http.StatusNotFound, r.JSON(verrs))
	}

	verrs, err := tx.ValidateAndCreate(maintenanceRequest)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		c.Set("errors", verrs)
		return c.Render(http.StatusUnprocessableEntity, r.JSON(maintenanceRequest))
	}

	return c.Render(http.StatusCreated, r.JSON(maintenanceRequest))
}

func (v RoomMaintenanceRequestsResource) Update(c buffalo.Context) error {
	tx, q := v.getTransactionAndQueryContext(c)

	maintenanceRequest := &models.MaintenanceRequest{}
	if err := q.Find(maintenanceRequest, c.Param("room_maintenance_request_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := c.Bind(maintenanceRequest); err != nil {
		return err
	}

	verrs, err := tx.ValidateAndUpdate(maintenanceRequest)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		c.Set("errors", verrs)

		return c.Render(http.StatusUnprocessableEntity, r.JSON(maintenanceRequest))
	}

	return c.Render(http.StatusOK, r.JSON(maintenanceRequest))
}

func (v RoomMaintenanceRequestsResource) Destroy(c buffalo.Context) error {
	tx, q := v.getTransactionAndQueryContext(c)

	maintenanceRequest := &models.MaintenanceRequest{}

	if err := q.Find(maintenanceRequest, c.Param("room_maintenance_request_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := tx.Destroy(maintenanceRequest); err != nil {
		return err
	}

	return c.Render(http.StatusOK, r.JSON(maintenanceRequest))
}
