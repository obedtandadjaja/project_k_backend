package actions

import (
	"encoding/json"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/suite/fix"
	"github.com/obedtandadjaja/project_k_backend/helpers"
	"github.com/obedtandadjaja/project_k_backend/models"
)

func (as *ActionSuite) Test_TenantsResource_List() {
	as.LoadFixture("user with property with room with tenant")

	fixture, err := fix.Find("user with property with room with tenant")
	userID := fixture.Tables[0].Row[0]["id"]
	credentialUUID := fixture.Tables[0].Row[0]["credential_uuid"]
	propertyID := fixture.Tables[1].Row[0]["id"]
	roomID := fixture.Tables[3].Row[0]["id"]
	tenantID := fixture.Tables[0].Row[1]["id"]

	token, err := helpers.GenerateAccessToken(userID.(string), credentialUUID.(string))
	if err != nil {
		as.NoError(err)
	}

	req := as.JSON("/api/v1/properties/%s/rooms/%s/tenants", propertyID.(string), roomID.(string))
	req.Headers = map[string]string{
		"Authorization": token,
	}
	res := req.Get()
	as.Equal(200, res.Code)

	var responseBody []map[string]interface{}
	json.NewDecoder(res.Body).Decode(&responseBody)

	as.Equal(tenantID.(string), responseBody[0]["id"])
	as.Equal(1, len(responseBody))
}

func (as *ActionSuite) Test_TenantsResource_Show() {
	as.LoadFixture("user with property with room with tenant")

	fixture, err := fix.Find("user with property with room with tenant")
	userID := fixture.Tables[0].Row[0]["id"]
	credentialUUID := fixture.Tables[0].Row[0]["credential_uuid"]
	propertyID := fixture.Tables[1].Row[0]["id"]
	roomID := fixture.Tables[3].Row[0]["id"]
	tenantID := fixture.Tables[0].Row[1]["id"]

	token, err := helpers.GenerateAccessToken(userID.(string), credentialUUID.(string))
	if err != nil {
		as.NoError(err)
	}

	req := as.JSON("/api/v1/properties/%s/rooms/%s/tenants/%s", propertyID.(string), roomID.(string),
		tenantID.(string))
	req.Headers = map[string]string{
		"Authorization": token,
	}
	res := req.Get()
	as.Equal(200, res.Code)

	var responseBody map[string]interface{}
	json.NewDecoder(res.Body).Decode(&responseBody)

	as.Equal(tenantID.(string), responseBody["id"])
}

func (as *ActionSuite) Test_TenantsResource_Create() {
	as.LoadFixture("user with property with room")

	fixture, err := fix.Find("user with property with room")
	userID := fixture.Tables[0].Row[0]["id"]
	credentialUUID := fixture.Tables[0].Row[0]["credential_uuid"]
	propertyID := fixture.Tables[1].Row[0]["id"]
	roomID := fixture.Tables[3].Row[0]["id"]

	token, err := helpers.GenerateAccessToken(userID.(string), credentialUUID.(string))
	if err != nil {
		as.NoError(err)
	}

	propertyToCreate := &models.User{
		Name:  nulls.String{String: "tenant", Valid: true},
		Email: "tenant@example.com",
		Rooms: models.Rooms{models.Room{ID: helpers.ParseUUID(roomID.(string))}},
	}
	req := as.JSON("/api/v1/properties/%s/rooms/%s/tenants", propertyID.(string), roomID.(string))
	req.Headers = map[string]string{
		"Authorization": token,
	}
	res := req.Post(propertyToCreate)
	as.Equal(201, res.Code)

	var responseBody map[string]interface{}
	json.NewDecoder(res.Body).Decode(&responseBody)

	as.Equal("tenant", responseBody["name"])
	as.Equal("tenant@example.com", responseBody["email"])

	roomOccupancy := &models.RoomOccupancy{}
	err = as.DB.Where("room_id = ?", roomID.(string)).First(roomOccupancy)
	as.NoError(err)
}

func (as *ActionSuite) Test_TenantsResource_Update() {
	as.LoadFixture("user with property with room with tenant")

	fixture, err := fix.Find("user with property with room with tenant")
	userID := fixture.Tables[0].Row[0]["id"]
	credentialUUID := fixture.Tables[0].Row[0]["credential_uuid"]
	propertyID := fixture.Tables[1].Row[0]["id"]
	roomID := fixture.Tables[3].Row[0]["id"]
	tenantID := fixture.Tables[0].Row[1]["id"]

	token, err := helpers.GenerateAccessToken(userID.(string), credentialUUID.(string))
	if err != nil {
		as.NoError(err)
	}

	req := as.JSON("/api/v1/properties/%s/rooms/%s/tenants/%s", propertyID.(string), roomID.(string),
		tenantID.(string))
	req.Headers = map[string]string{
		"Authorization": token,
	}
	res := req.Put(
		&models.User{
			ID:    helpers.ParseUUID(tenantID.(string)),
			Rooms: models.Rooms{models.Room{ID: helpers.ParseUUID(roomID.(string))}},
			Name:  nulls.String{String: "Changed", Valid: true},
		},
	)
	as.Equal(200, res.Code)

	tenant := &models.User{}
	as.DB.Find(tenant, tenantID.(string))

	as.Equal(nulls.String{String: "Changed", Valid: true}, tenant.Name)
}

func (as *ActionSuite) Test_TenantsResource_Destroy() {
	as.LoadFixture("user with property with room with tenant")

	fixture, err := fix.Find("user with property with room with tenant")
	userID := fixture.Tables[0].Row[0]["id"]
	credentialUUID := fixture.Tables[0].Row[0]["credential_uuid"]
	propertyID := fixture.Tables[1].Row[0]["id"]
	roomID := fixture.Tables[3].Row[0]["id"]
	tenantID := fixture.Tables[0].Row[1]["id"]

	token, err := helpers.GenerateAccessToken(userID.(string), credentialUUID.(string))
	if err != nil {
		as.NoError(err)
	}

	req := as.JSON("/api/v1/properties/%s/rooms/%s/tenants/%s", propertyID.(string), roomID.(string),
		tenantID.(string))
	req.Headers = map[string]string{
		"Authorization": token,
	}
	res := req.Delete()
	as.Equal(200, res.Code)

	var responseBody map[string]interface{}
	json.NewDecoder(res.Body).Decode(&responseBody)

	as.Equal(tenantID.(string), responseBody["id"])

	roomOccupancy := &models.RoomOccupancy{}
	err = as.DB.Where("room_id = ?", roomID.(string)).First(roomOccupancy)
	as.NotNil(err)
}
