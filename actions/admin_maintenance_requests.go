package actions

import (
	"net/http"
	"strings"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/obedtandadjaja/project_k_backend/helpers"
	"github.com/obedtandadjaja/project_k_backend/models"
)

type AdminMaintenanceRequestsResource struct {
	buffalo.Resource
}

func (v AdminMaintenanceRequestsResource) getTransactionAndQueryContext(c buffalo.Context) (*pop.Connection, *pop.Query) {
	tx, _ := c.Value("tx").(*pop.Connection)

	return tx, tx.Q().
		InnerJoin("user_property_relationships", "user_property_relationships.property_id = maintenance_requests.property_id").
		Where("user_property_relationships.user_id = ?", c.Value("current_user_id"))
}

func (v AdminMaintenanceRequestsResource) List(c buffalo.Context) error {
	_, q := v.getTransactionAndQueryContext(c)

	maintenanceRequests := &models.MaintenanceRequests{}

	if c.Param("status") != "" {
		q.Where("status = ?", c.Param("status"))
	}

	if c.Param("category") != "" {
		q.Where("category = ?", c.Param("category"))
	}

	if c.Param("property_id") != "" {
		q.Where("maintenance_requests.property_id = ?", c.Param("property_id"))
	}

	if c.Param("room_id") != "" {
		q.Where("maintenance_requests.room_id = ?", c.Param("room_id"))
	}

	if c.Param("opened_start_date") != "" {
		q.Where("maintenance_requests.created_at >= ?", c.Param("opened_start_date"))
	}

	if c.Param("opened_end_date") != "" {
		q.Where("maintenance_requests.created_at <= ?", c.Param("opened_end_date"))
	}

	if c.Param("closed_start_date") != "" {
		q.Where("maintenance_requests.updated_at >= ? AND status='closed'",
			c.Param("closed_start_date"))
	}

	if c.Param("closed_end_date") != "" {
		q.Where("maintenance_requests.updated_at <= ? AND status='closed'",
			c.Param("closed_end_date"))
	}

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

func (v AdminMaintenanceRequestsResource) Show(c buffalo.Context) error {
	_, q := v.getTransactionAndQueryContext(c)

	maintenanceRequest := &models.MaintenanceRequest{}

	if c.Param("eager") != "" {
		if err := q.Eager(c.Param("eager")).Find(maintenanceRequest, c.Param("admin_maintenance_request_id")); err != nil {
			return c.Error(http.StatusNotFound, err)
		}
	} else {
		if err := q.Find(maintenanceRequest, c.Param("admin_maintenance_request_id")); err != nil {
			return c.Error(http.StatusNotFound, err)
		}
	}

	return c.Render(http.StatusOK, r.JSON(maintenanceRequest))
}

func (v AdminMaintenanceRequestsResource) Create(c buffalo.Context) error {
	maintenanceRequest := &models.MaintenanceRequest{
		ReporterID: helpers.ParseUUID(c.Value("current_user_id").(string)),
		Status:     "pending",
	}
	if err := c.Bind(maintenanceRequest); err != nil {
		return err
	}

	tx, _ := v.getTransactionAndQueryContext(c)

	verrs, err := maintenanceRequest.Validate(tx)
	if verrs.HasAny() {
		c.Set("errors", verrs)
		return c.Render(http.StatusUnprocessableEntity, r.JSON(maintenanceRequest))
	}

	// double check that either the property or room exists
	// double check that user also has access to those resources
	if !maintenanceRequest.RoomID.Valid {
		q := tx.Q().
			InnerJoin("user_property_relationships", "user_property_relationships.property_id = properties.id").
			Where("user_property_relationships.user_id = ?", c.Value("current_user_id")).
			Where("properties.id = ?", maintenanceRequest.PropertyID)
		if err := q.First(&models.Property{}); err != nil {
			verrs.Add("property", "Property does not exist")
			return c.Render(http.StatusNotFound, r.JSON(verrs))
		}
	} else {
		q := tx.Q().
			InnerJoin("properties", "properties.id = rooms.property_id").
			InnerJoin("user_property_relationships", "user_property_relationships.property_id = properties.id").
			Where("user_property_relationships.user_id = ?", c.Value("current_user_id")).
			Where("properties.id = ?", maintenanceRequest.PropertyID).
			Where("rooms.id = ?", maintenanceRequest.RoomID)
		if err := q.First(&models.Room{}); err != nil {
			verrs.Add("room", "Room does not exist")
			return c.Render(http.StatusNotFound, r.JSON(verrs))
		}
	}

	err = tx.Create(maintenanceRequest)
	if err != nil {
		return err
	}

	return c.Render(http.StatusCreated, r.JSON(maintenanceRequest))
}

func (v AdminMaintenanceRequestsResource) Update(c buffalo.Context) error {
	tx, q := v.getTransactionAndQueryContext(c)

	maintenanceRequest := &models.MaintenanceRequest{}
	if err := q.Find(maintenanceRequest, c.Param("admin_maintenance_request_id")); err != nil {
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

func (v AdminMaintenanceRequestsResource) Destroy(c buffalo.Context) error {
	tx, q := v.getTransactionAndQueryContext(c)

	maintenanceRequest := &models.MaintenanceRequest{}

	if err := q.Find(maintenanceRequest, c.Param("admin_maintenance_request_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := tx.Destroy(maintenanceRequest); err != nil {
		return err
	}

	return c.Render(http.StatusOK, r.JSON(maintenanceRequest))
}

func (v AdminMaintenanceRequestsResource) Complete(c buffalo.Context) error {
	tx, q := v.getTransactionAndQueryContext(c)

	maintenanceRequest := &models.MaintenanceRequest{}
	if err := q.Find(maintenanceRequest, c.Param("admin_maintenance_request_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if maintenanceRequest.Status != "pending" {
		verrs := validate.NewErrors()
		verrs.Add("status", "Status must be pending")
		return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
	}

	maintenanceRequest.Status = "closed"
	maintenanceRequest.CompletedAt = nulls.Time{Time: time.Now(), Valid: true}

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
