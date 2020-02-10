package actions

import (
	"fmt"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/slices"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/obedtandadjaja/project_k_backend/helpers"
	"github.com/obedtandadjaja/project_k_backend/models"
)

type RoomsResource struct {
	buffalo.Resource
}

type BatchCreateRooms struct {
	PriceAmount     int        `json:"priceAmount"`
	PaymentSchedule string     `json:"paymentSchedule"`
	Data            slices.Map `json:"data"`
	Type            string     `json:"type"`
	Quantity        int        `json:"quantity"`
}

func (v RoomsResource) getTransactionAndQueryContext(c buffalo.Context) (*pop.Connection, *pop.Query) {
	tx, _ := c.Value("tx").(*pop.Connection)

	return tx, tx.Q().
		InnerJoin("properties", "properties.id = rooms.property_id").
		InnerJoin("user_property_relationships", "user_property_relationships.property_id = properties.id").
		Where("user_property_relationships.user_id = ?", c.Value("current_user_id")).
		Where("properties.id = ?", c.Param("property_id"))
}

func (v RoomsResource) List(c buffalo.Context) error {
	_, q := v.getTransactionAndQueryContext(c)

	rooms := &models.Rooms{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q = q.PaginateFromParams(c.Params())
	if c.Param("eager") != "" {
		if err := q.Eager(c.Param("eager")).All(rooms); err != nil {
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
	_, q := v.getTransactionAndQueryContext(c)

	room := &models.Room{}
	if c.Param("eager") != "" {
		if err := q.Eager(c.Param("eager")).Find(room, c.Param("room_id")); err != nil {
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
		Property: models.Property{ID: helpers.ParseUUID(c.Param("property_id"))},
		Data:     slices.Map{},
	}
	if err := c.Bind(room); err != nil {
		return err
	}

	tx, _ := v.getTransactionAndQueryContext(c)

	// Check that user has ownership of property and that property exists
	q := tx.Q().
		InnerJoin("user_property_relationships", "user_property_relationships.property_id = properties.id").
		Where("user_property_relationships.user_id = ?", c.Value("current_user_id")).
		Where("properties.id = ?", c.Param("property_id"))
	err := q.First(&models.Property{})
	if err != nil {
		verrs := validate.NewErrors()
		verrs.Add("property", "Property does not exist")
		return c.Render(http.StatusNotFound, r.JSON(verrs))
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

func (v RoomsResource) BatchCreate(c buffalo.Context) error {
	req := BatchCreateRooms{Data: slices.Map{}}
	if err := c.Bind(req); err != nil {
		return err
	}

	verrs := validate.Validate(
		&validators.IntIsPresent{Field: req.Quantity, Name: "Quantity"},
		&validators.IntIsPresent{Field: req.PriceAmount, Name: "PriceAmount"},
		&validators.StringIsPresent{Field: req.PaymentSchedule, Name: "PaymentSchedule"},
		&validators.StringIsPresent{Field: req.Type, Name: "Type"},
	)
	if verrs.HasAny() {
		c.Set("errors", verrs)
		return c.Render(http.StatusBadRequest, r.JSON(verrs))
	}

	tx, _ := v.getTransactionAndQueryContext(c)

	rooms := make([]models.Room, req.Quantity)
	for i := 0; i < req.Quantity; i++ {
		rooms[i] = models.Room{
			PropertyID:      helpers.ParseUUID(c.Param("property_id")),
			Name:            fmt.Sprintf("Room %d", i+1),
			PriceAmount:     req.PriceAmount,
			PaymentSchedule: req.PaymentSchedule,
			Data:            req.Data,
		}
		tx.Create(&rooms[i])
	}

	return c.Render(http.StatusCreated, r.JSON(rooms))
}

func (v RoomsResource) Update(c buffalo.Context) error {
	tx, q := v.getTransactionAndQueryContext(c)

	room := &models.Room{}

	if err := q.Find(room, c.Param("room_id")); err != nil {
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
		return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
	}

	return c.Render(http.StatusOK, r.JSON(room))
}

func (v RoomsResource) Destroy(c buffalo.Context) error {
	tx, q := v.getTransactionAndQueryContext(c)

	room := &models.Room{}

	if err := q.Find(room, c.Param("room_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := tx.Destroy(room); err != nil {
		return err
	}

	return c.Render(http.StatusOK, r.JSON(room))
}
