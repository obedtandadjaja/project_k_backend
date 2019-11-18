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
// Model: Singular (Property)
// DB Table: Plural (properties)
// Resource: Plural (Properties)
// Path: Plural (/properties)
// View Template Folder: Plural (/templates/properties/)

// PropertiesResource is the resource for the Property model
type PropertiesResource struct{
  buffalo.Resource
}

// List gets all Properties. This function is mapped to the path
// GET /properties
func (v PropertiesResource) List(c buffalo.Context) error {
  // Get the DB connection from the context
  tx, ok := c.Value("tx").(*pop.Connection)
  if !ok {
    return fmt.Errorf("no transaction found")
  }

  properties := &models.Properties{}

  // Paginate results. Params "page" and "per_page" control pagination.
  // Default values are "page=1" and "per_page=20".
  q := tx.PaginateFromParams(c.Params())

  // Retrieve all Properties from the DB
  if err := q.All(properties); err != nil {
    return err
  }

  // Add the paginator to the context so it can be used in the template.
  c.Set("pagination", q.Paginator)

  return c.Render(http.StatusOK, r.Auto(c, properties))
}

// Show gets the data for one Property. This function is mapped to
// the path GET /properties/{property_id}
func (v PropertiesResource) Show(c buffalo.Context) error {
  // Get the DB connection from the context
  tx, ok := c.Value("tx").(*pop.Connection)
  if !ok {
    return fmt.Errorf("no transaction found")
  }

  // Allocate an empty Property
  property := &models.Property{}

  // To find the Property the parameter property_id is used.
  if err := tx.Find(property, c.Param("property_id")); err != nil {
    return c.Error(http.StatusNotFound, err)
  }

  return c.Render(http.StatusOK, r.Auto(c, property))
}

// Create adds a Property to the DB. This function is mapped to the
// path POST /properties
func (v PropertiesResource) Create(c buffalo.Context) error {
  // Allocate an empty Property
  property := &models.Property{}

  // Bind property to the html form elements
  if err := c.Bind(property); err != nil {
    return err
  }

  // Get the DB connection from the context
  tx, ok := c.Value("tx").(*pop.Connection)
  if !ok {
    return fmt.Errorf("no transaction found")
  }

  // Validate the data from the html form
  verrs, err := tx.ValidateAndCreate(property)
  if err != nil {
    return err
  }

  if verrs.HasAny() {
    // Make the errors available inside the html template
    c.Set("errors", verrs)

    // Render again the new.html template that the user can
    // correct the input.
    return c.Render(http.StatusUnprocessableEntity, r.Auto(c, property))
  }

  // and redirect to the properties index page
  return c.Render(http.StatusCreated, r.Auto(c, property))
}

// Update changes a Property in the DB. This function is mapped to
// the path PUT /properties/{property_id}
func (v PropertiesResource) Update(c buffalo.Context) error {
  // Get the DB connection from the context
  tx, ok := c.Value("tx").(*pop.Connection)
  if !ok {
    return fmt.Errorf("no transaction found")
  }

  // Allocate an empty Property
  property := &models.Property{}

  if err := tx.Find(property, c.Param("property_id")); err != nil {
    return c.Error(http.StatusNotFound, err)
  }

  // Bind Property to the html form elements
  if err := c.Bind(property); err != nil {
    return err
  }

  verrs, err := tx.ValidateAndUpdate(property)
  if err != nil {
    return err
  }

  if verrs.HasAny() {
    // Make the errors available inside the html template
    c.Set("errors", verrs)

    // Render again the edit.html template that the user can
    // correct the input.
    return c.Render(http.StatusUnprocessableEntity, r.Auto(c, property))
  }

  // and redirect to the properties index page
  return c.Render(http.StatusOK, r.Auto(c, property))
}

// Destroy deletes a Property from the DB. This function is mapped
// to the path DELETE /properties/{property_id}
func (v PropertiesResource) Destroy(c buffalo.Context) error {
  // Get the DB connection from the context
  tx, ok := c.Value("tx").(*pop.Connection)
  if !ok {
    return fmt.Errorf("no transaction found")
  }

  // Allocate an empty Property
  property := &models.Property{}

  // To find the Property the parameter property_id is used.
  if err := tx.Find(property, c.Param("property_id")); err != nil {
    return c.Error(http.StatusNotFound, err)
  }

  if err := tx.Destroy(property); err != nil {
    return err
  }

  // Redirect to the properties index page
  return c.Render(http.StatusOK, r.Auto(c, property))
}