package actions

import (

  "fmt"
  "net/http"
  "github.com/gobuffalo/buffalo"
  "github.com/gobuffalo/pop"
  "github.com/obedtandadjaja/project_k_backend/models"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (RoomOccupancy)
// DB Table: Plural (room_occupancies)
// Resource: Plural (RoomOccupancies)
// Path: Plural (/room_occupancies)
// View Template Folder: Plural (/templates/room_occupancies/)

// RoomOccupanciesResource is the resource for the RoomOccupancy model
type RoomOccupanciesResource struct{
  buffalo.Resource
}

// List gets all RoomOccupancies. This function is mapped to the path
// GET /room_occupancies
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

  return c.Render(http.StatusOK, r.Auto(c, roomOccupancies))
}

// Show gets the data for one RoomOccupancy. This function is mapped to
// the path GET /room_occupancies/{room_occupancy_id}
func (v RoomOccupanciesResource) Show(c buffalo.Context) error {
  // Get the DB connection from the context
  tx, ok := c.Value("tx").(*pop.Connection)
  if !ok {
    return fmt.Errorf("no transaction found")
  }

  // Allocate an empty RoomOccupancy
  roomOccupancy := &models.RoomOccupancy{}

  // To find the RoomOccupancy the parameter room_occupancy_id is used.
  if err := tx.Find(roomOccupancy, c.Param("room_occupancy_id")); err != nil {
    return c.Error(http.StatusNotFound, err)
  }

  return c.Render(http.StatusOK, r.Auto(c, roomOccupancy))
}

// Create adds a RoomOccupancy to the DB. This function is mapped to the
// path POST /room_occupancies
func (v RoomOccupanciesResource) Create(c buffalo.Context) error {
  // Allocate an empty RoomOccupancy
  roomOccupancy := &models.RoomOccupancy{}

  // Bind roomOccupancy to the html form elements
  if err := c.Bind(roomOccupancy); err != nil {
    return err
  }

  // Get the DB connection from the context
  tx, ok := c.Value("tx").(*pop.Connection)
  if !ok {
    return fmt.Errorf("no transaction found")
  }

  // Validate the data from the html form
  verrs, err := tx.ValidateAndCreate(roomOccupancy)
  if err != nil {
    return err
  }

  if verrs.HasAny() {
    // Make the errors available inside the html template
    c.Set("errors", verrs)

    // Render again the new.html template that the user can
    // correct the input.
    return c.Render(http.StatusUnprocessableEntity, r.Auto(c, roomOccupancy))
  }

  // and redirect to the room_occupancies index page
  return c.Render(http.StatusCreated, r.Auto(c, roomOccupancy))
}

// Update changes a RoomOccupancy in the DB. This function is mapped to
// the path PUT /room_occupancies/{room_occupancy_id}
func (v RoomOccupanciesResource) Update(c buffalo.Context) error {
  // Get the DB connection from the context
  tx, ok := c.Value("tx").(*pop.Connection)
  if !ok {
    return fmt.Errorf("no transaction found")
  }

  // Allocate an empty RoomOccupancy
  roomOccupancy := &models.RoomOccupancy{}

  if err := tx.Find(roomOccupancy, c.Param("room_occupancy_id")); err != nil {
    return c.Error(http.StatusNotFound, err)
  }

  // Bind RoomOccupancy to the html form elements
  if err := c.Bind(roomOccupancy); err != nil {
    return err
  }

  verrs, err := tx.ValidateAndUpdate(roomOccupancy)
  if err != nil {
    return err
  }

  if verrs.HasAny() {
    // Make the errors available inside the html template
    c.Set("errors", verrs)

    // Render again the edit.html template that the user can
    // correct the input.
    return c.Render(http.StatusUnprocessableEntity, r.Auto(c, roomOccupancy))
  }

  // and redirect to the room_occupancies index page
  return c.Render(http.StatusOK, r.Auto(c, roomOccupancy))
}

// Destroy deletes a RoomOccupancy from the DB. This function is mapped
// to the path DELETE /room_occupancies/{room_occupancy_id}
func (v RoomOccupanciesResource) Destroy(c buffalo.Context) error {
  // Get the DB connection from the context
  tx, ok := c.Value("tx").(*pop.Connection)
  if !ok {
    return fmt.Errorf("no transaction found")
  }

  // Allocate an empty RoomOccupancy
  roomOccupancy := &models.RoomOccupancy{}

  // To find the RoomOccupancy the parameter room_occupancy_id is used.
  if err := tx.Find(roomOccupancy, c.Param("room_occupancy_id")); err != nil {
    return c.Error(http.StatusNotFound, err)
  }

  if err := tx.Destroy(roomOccupancy); err != nil {
    return err
  }

  // Redirect to the room_occupancies index page
  return c.Render(http.StatusOK, r.Auto(c, roomOccupancy))
}
