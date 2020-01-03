package actions

import (
	"encoding/json"

	"github.com/gobuffalo/suite/fix"
	"github.com/obedtandadjaja/project_k_backend/helpers"
	"github.com/obedtandadjaja/project_k_backend/models"
)

func (as *ActionSuite) Test_PropertiesResource_List() {
	as.LoadFixture("user with property")

	fixture, err := fix.Find("user with property")
	userID := fixture.Tables[0].Row[0]["id"]
	credentialUUID := fixture.Tables[0].Row[0]["credential_uuid"]
	propertyID := fixture.Tables[1].Row[0]["id"]

	token, err := helpers.GenerateAccessToken(userID.(string), credentialUUID.(string))
	if err != nil {
		as.NoError(err)
	}

	req := as.JSON("/api/v1/properties")
	req.Headers = map[string]string{
		"Authorization": token,
	}
	res := req.Get()
	as.Equal(200, res.Code)

	var responseBody []map[string]interface{}
	json.NewDecoder(res.Body).Decode(&responseBody)

	as.Equal(propertyID.(string), responseBody[0]["id"])
	as.Equal(1, len(responseBody))
}

func (as *ActionSuite) Test_PropertiesResource_Show() {
	as.LoadFixture("user with property")

	fixture, err := fix.Find("user with property")
	userID := fixture.Tables[0].Row[0]["id"]
	credentialUUID := fixture.Tables[0].Row[0]["credential_uuid"]
	propertyID := fixture.Tables[1].Row[0]["id"]

	token, err := helpers.GenerateAccessToken(userID.(string), credentialUUID.(string))
	if err != nil {
		as.NoError(err)
	}

	req := as.JSON("/api/v1/properties/%s", propertyID.(string))
	req.Headers = map[string]string{
		"Authorization": token,
	}
	res := req.Get()
	as.Equal(200, res.Code)

	var responseBody map[string]interface{}
	json.NewDecoder(res.Body).Decode(&responseBody)

	as.Equal(propertyID.(string), responseBody["id"])
}

func (as *ActionSuite) Test_PropertiesResource_Create() {
	as.LoadFixture("user")

	fixture, err := fix.Find("user")
	userID := fixture.Tables[0].Row[0]["id"]
	credentialUUID := fixture.Tables[0].Row[0]["credential_uuid"]

	token, err := helpers.GenerateAccessToken(userID.(string), credentialUUID.(string))
	if err != nil {
		as.NoError(err)
	}

	propertyToCreate := &models.Property{
		Name:    "property",
		Address: "address",
		Type:    "house",
		Users:   models.Users{models.User{ID: helpers.ParseUUID(userID.(string))}},
	}
	req := as.JSON("/api/v1/properties")
	req.Headers = map[string]string{
		"Authorization": token,
	}
	res := req.Post(propertyToCreate)
	as.Equal(201, res.Code)

	var responseBody map[string]interface{}
	json.NewDecoder(res.Body).Decode(&responseBody)

	as.Equal("property", responseBody["name"])
	as.Equal("address", responseBody["address"])
	as.Equal("house", responseBody["type"])

	userPropertyRelationship := &models.UserPropertyRelationship{}
	err = as.DB.Where("user_id = ?", userID.(string)).First(userPropertyRelationship)
	as.NoError(err)
}

func (as *ActionSuite) Test_PropertiesResource_Update() {
	as.LoadFixture("user with property")

	fixture, err := fix.Find("user with property")
	userID := fixture.Tables[0].Row[0]["id"]
	credentialUUID := fixture.Tables[0].Row[0]["credential_uuid"]
	propertyID := fixture.Tables[1].Row[0]["id"]

	token, err := helpers.GenerateAccessToken(userID.(string), credentialUUID.(string))
	if err != nil {
		as.NoError(err)
	}

	req := as.JSON("/api/v1/properties/%s", propertyID.(string))
	req.Headers = map[string]string{
		"Authorization": token,
	}
	res := req.Put(&models.Property{ID: helpers.ParseUUID(propertyID.(string)), Name: "Changed"})
	as.Equal(200, res.Code)

	property := &models.Property{}
	as.DB.Find(property, propertyID.(string))

	as.Equal("Changed", property.Name)
}

func (as *ActionSuite) Test_PropertiesResource_Destroy() {
	as.LoadFixture("user with property")

	fixture, err := fix.Find("user with property")
	userID := fixture.Tables[0].Row[0]["id"]
	credentialUUID := fixture.Tables[0].Row[0]["credential_uuid"]
	propertyID := fixture.Tables[1].Row[0]["id"]

	token, err := helpers.GenerateAccessToken(userID.(string), credentialUUID.(string))
	if err != nil {
		as.NoError(err)
	}

	req := as.JSON("/api/v1/properties/%s", propertyID.(string))
	req.Headers = map[string]string{
		"Authorization": token,
	}
	res := req.Delete()
	as.Equal(200, res.Code)

	var responseBody map[string]interface{}
	json.NewDecoder(res.Body).Decode(&responseBody)

	as.Equal(propertyID.(string), responseBody["id"])

	userPropertyRelationship := &models.UserPropertyRelationship{}
	err = as.DB.Where("user_id = ?", userID.(string)).First(userPropertyRelationship)

	as.NotNil(err)
}
