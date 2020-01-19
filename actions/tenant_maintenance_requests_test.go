package actions

import (
	"encoding/json"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/suite/fix"
	"github.com/obedtandadjaja/project_k_backend/helpers"
	"github.com/obedtandadjaja/project_k_backend/models"
)

func (as *ActionSuite) Test_TenantMaintenanceRequestsResource_List() {
	as.LoadFixture("user with maintenance request")

	fixture, err := fix.Find("user with maintenance request")
	userID := fixture.Tables[0].Row[0]["id"]
	credentialUUID := fixture.Tables[0].Row[0]["credential_uuid"]
	maintenanceRequestID := fixture.Tables[1].Row[0]["id"]

	token, err := helpers.GenerateAccessToken(userID.(string), credentialUUID.(string))
	if err != nil {
		as.NoError(err)
	}

	req := as.JSON("/api/tenant/v1/maintenance_requests")
	req.Headers = map[string]string{
		"Authorization": token,
	}
	res := req.Get()
	as.Equal(200, res.Code)

	var responseBody []map[string]interface{}
	json.NewDecoder(res.Body).Decode(&responseBody)

	as.Equal(maintenanceRequestID.(string), responseBody[0]["id"])
	as.Equal(1, len(responseBody))
}

func (as *ActionSuite) Test_TenantMaintenanceRequestsResource_Show() {
	as.LoadFixture("user with maintenance request")

	fixture, err := fix.Find("user with maintenance request")
	userID := fixture.Tables[0].Row[0]["id"]
	credentialUUID := fixture.Tables[0].Row[0]["credential_uuid"]
	maintenanceRequestID := fixture.Tables[1].Row[0]["id"]

	token, err := helpers.GenerateAccessToken(userID.(string), credentialUUID.(string))
	if err != nil {
		as.NoError(err)
	}

	req := as.JSON("/api/tenant/v1/maintenance_requests/%s", maintenanceRequestID.(string))
	req.Headers = map[string]string{
		"Authorization": token,
	}
	res := req.Get()
	as.Equal(200, res.Code)

	var responseBody map[string]interface{}
	json.NewDecoder(res.Body).Decode(&responseBody)

	as.Equal(maintenanceRequestID.(string), responseBody["id"])
}

func (as *ActionSuite) Test_TenantMaintenanceRequestsResource_CreateRoomMaintenanceRequest() {
	as.LoadFixture("tenant with property with room")

	fixture, err := fix.Find("tenant with property with room")
	userID := fixture.Tables[0].Row[0]["id"]
	credentialUUID := fixture.Tables[0].Row[0]["credential_uuid"]
	propertyID := fixture.Tables[1].Row[0]["id"]
	roomID := fixture.Tables[2].Row[0]["id"]

	token, err := helpers.GenerateAccessToken(userID.(string), credentialUUID.(string))
	if err != nil {
		as.NoError(err)
	}

	maintenanceRequestToCreate := &models.MaintenanceRequest{
		Title:       "title",
		Description: nulls.String{String: "description", Valid: true},
		Status:      "pending",
		ReporterID:  helpers.ParseUUID(userID.(string)),
		RoomID:      nulls.UUID{UUID: helpers.ParseUUID(roomID.(string)), Valid: true},
		PropertyID:  nulls.UUID{UUID: helpers.ParseUUID(propertyID.(string)), Valid: true},
	}
	req := as.JSON("/api/tenant/v1/maintenance_requests")
	req.Headers = map[string]string{
		"Authorization": token,
	}
	res := req.Post(maintenanceRequestToCreate)
	as.Equal(201, res.Code)

	var responseBody map[string]interface{}
	json.NewDecoder(res.Body).Decode(&responseBody)

	as.Equal("title", responseBody["title"])
	as.Equal("description", responseBody["description"])
	as.Equal("pending", responseBody["status"])

	maintenanceRequest := &models.MaintenanceRequest{}
	err = as.DB.Where("room_id = ?", roomID.(string)).First(maintenanceRequest)
	as.NoError(err)
}

func (as *ActionSuite) Test_TenantMaintenanceRequestsResource_CreatePropertyMaintenanceRequest() {
	as.LoadFixture("tenant with property with room")

	fixture, err := fix.Find("tenant with property with room")
	userID := fixture.Tables[0].Row[0]["id"]
	credentialUUID := fixture.Tables[0].Row[0]["credential_uuid"]
	propertyID := fixture.Tables[1].Row[0]["id"]

	token, err := helpers.GenerateAccessToken(userID.(string), credentialUUID.(string))
	if err != nil {
		as.NoError(err)
	}

	maintenanceRequestToCreate := &models.MaintenanceRequest{
		Title:       "title",
		Description: nulls.String{String: "description", Valid: true},
		Status:      "pending",
		ReporterID:  helpers.ParseUUID(userID.(string)),
		PropertyID:  nulls.UUID{UUID: helpers.ParseUUID(propertyID.(string)), Valid: true},
	}
	req := as.JSON("/api/tenant/v1/maintenance_requests")
	req.Headers = map[string]string{
		"Authorization": token,
	}
	res := req.Post(maintenanceRequestToCreate)
	as.Equal(201, res.Code)

	var responseBody map[string]interface{}
	json.NewDecoder(res.Body).Decode(&responseBody)

	as.Equal("title", responseBody["title"])
	as.Equal("description", responseBody["description"])
	as.Equal("pending", responseBody["status"])

	maintenanceRequest := &models.MaintenanceRequest{}
	err = as.DB.Where("property_id = ?", propertyID.(string)).First(maintenanceRequest)
	as.NoError(err)
}

func (as *ActionSuite) Test_TenantMaintenanceRequestsResource_Update() {
	as.LoadFixture("user with property with maintenance request")

	fixture, err := fix.Find("user with maintenance request")
	userID := fixture.Tables[0].Row[0]["id"]
	credentialUUID := fixture.Tables[0].Row[0]["credential_uuid"]
	propertyID := fixture.Tables[1].Row[0]["id"]
	maintenanceRequestID := fixture.Tables[1].Row[0]["id"]

	token, err := helpers.GenerateAccessToken(userID.(string), credentialUUID.(string))
	if err != nil {
		as.NoError(err)
	}

	req := as.JSON("/api/tenant/v1/maintenance_requests/%s", maintenanceRequestID.(string))
	req.Headers = map[string]string{
		"Authorization": token,
	}
	res := req.Put(
		&models.MaintenanceRequest{
			ID:         helpers.ParseUUID(maintenanceRequestID.(string)),
			PropertyID: nulls.UUID{UUID: helpers.ParseUUID(propertyID.(string)), Valid: true},
			ReporterID: helpers.ParseUUID(userID.(string)),
			Title:      "Changed",
			Status:     "Changed",
		},
	)
	as.Equal(200, res.Code)

	maintenanceRequest := &models.MaintenanceRequest{}
	as.DB.Find(maintenanceRequest, maintenanceRequestID.(string))

	as.Equal("Changed", maintenanceRequest.Title)
}

func (as *ActionSuite) Test_TenantMaintenanceRequestsResource_Destroy() {
	as.LoadFixture("user with maintenance request")

	fixture, err := fix.Find("user with maintenance request")
	userID := fixture.Tables[0].Row[0]["id"]
	credentialUUID := fixture.Tables[0].Row[0]["credential_uuid"]
	maintenanceRequestID := fixture.Tables[1].Row[0]["id"]

	token, err := helpers.GenerateAccessToken(userID.(string), credentialUUID.(string))
	if err != nil {
		as.NoError(err)
	}

	req := as.JSON("/api/tenant/v1/maintenance_requests/%s", maintenanceRequestID.(string))
	req.Headers = map[string]string{
		"Authorization": token,
	}
	res := req.Delete()
	as.Equal(200, res.Code)

	var responseBody map[string]interface{}
	json.NewDecoder(res.Body).Decode(&responseBody)

	as.Equal(maintenanceRequestID.(string), responseBody["id"])

	maintenanceRequest := &models.MaintenanceRequest{}
	err = as.DB.Where("id = ?", maintenanceRequestID.(string)).First(maintenanceRequest)

	as.NotNil(err)
}
