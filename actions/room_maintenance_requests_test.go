package actions

import (
	"encoding/json"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/suite/fix"
	"github.com/obedtandadjaja/project_k_backend/helpers"
	"github.com/obedtandadjaja/project_k_backend/models"
)

func (as *ActionSuite) Test_RoomMaintenanceRequestsResource_List() {
	as.LoadFixture("user with property with room with maintenance request")

	fixture, _ := fix.Find("user with property with room with maintenance request")
	propertyID := fixture.Tables[1].Row[0]["id"]
	roomID := fixture.Tables[3].Row[0]["id"]
	maintenanceRequestID := fixture.Tables[4].Row[0]["id"]

	token := AccessTokenHelper(fixture.Tables[0].Row[0])

	req := as.JSON("/api/v1/properties/%s/rooms/%s/maintenance_requests", propertyID.(string),
		roomID.(string))
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

func (as *ActionSuite) Test_RoomMaintenanceRequestsResource_Show() {
	as.LoadFixture("user with property with room with maintenance request")

	fixture, _ := fix.Find("user with property with room with maintenance request")
	propertyID := fixture.Tables[1].Row[0]["id"]
	roomID := fixture.Tables[3].Row[0]["id"]
	maintenanceRequestID := fixture.Tables[4].Row[0]["id"]

	token := AccessTokenHelper(fixture.Tables[0].Row[0])

	req := as.JSON("/api/v1/properties/%s/rooms/%s/maintenance_requests/%s", propertyID.(string),
		roomID.(string), maintenanceRequestID.(string))
	req.Headers = map[string]string{
		"Authorization": token,
	}
	res := req.Get()
	as.Equal(200, res.Code)

	var responseBody map[string]interface{}
	json.NewDecoder(res.Body).Decode(&responseBody)

	as.Equal(maintenanceRequestID.(string), responseBody["id"])
}

func (as *ActionSuite) Test_RoomMaintenanceRequestsResource_Create() {
	as.LoadFixture("user with property with room")

	fixture, err := fix.Find("user with property with room")
	userID := fixture.Tables[0].Row[0]["id"]
	propertyID := fixture.Tables[1].Row[0]["id"]
	roomID := fixture.Tables[3].Row[0]["id"]

	token := AccessTokenHelper(fixture.Tables[0].Row[0])

	maintenanceRequestToCreate := &models.MaintenanceRequest{
		Title:       "title",
		Description: nulls.String{String: "description", Valid: true},
		Status:      "pending",
		ReporterID:  helpers.ParseUUID(userID.(string)),
		RoomID:      nulls.UUID{UUID: helpers.ParseUUID(roomID.(string)), Valid: true},
	}
	req := as.JSON("/api/v1/properties/%s/rooms/%s/maintenance_requests", propertyID.(string),
		roomID.(string))
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

func (as *ActionSuite) Test_RoomMaintenanceRequestsResource_Update() {
	as.LoadFixture("user with property with room with maintenance request")

	fixture, _ := fix.Find("user with property with room with maintenance request")
	userID := fixture.Tables[0].Row[0]["id"]
	propertyID := fixture.Tables[1].Row[0]["id"]
	roomID := fixture.Tables[3].Row[0]["id"]
	maintenanceRequestID := fixture.Tables[4].Row[0]["id"]

	token := AccessTokenHelper(fixture.Tables[0].Row[0])

	req := as.JSON("/api/v1/properties/%s/rooms/%s/maintenance_requests/%s", propertyID.(string),
		roomID.(string), maintenanceRequestID.(string))
	req.Headers = map[string]string{
		"Authorization": token,
	}
	res := req.Put(
		&models.MaintenanceRequest{
			ID:         helpers.ParseUUID(maintenanceRequestID.(string)),
			RoomID:     nulls.UUID{UUID: helpers.ParseUUID(roomID.(string)), Valid: true},
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

func (as *ActionSuite) Test_RoomMaintenanceRequestsResource_Destroy() {
	as.LoadFixture("user with property with room with maintenance request")

	fixture, err := fix.Find("user with property with room with maintenance request")
	propertyID := fixture.Tables[1].Row[0]["id"]
	roomID := fixture.Tables[3].Row[0]["id"]
	maintenanceRequestID := fixture.Tables[4].Row[0]["id"]

	token := AccessTokenHelper(fixture.Tables[0].Row[0])

	req := as.JSON("/api/v1/properties/%s/rooms/%s/maintenance_requests/%s", propertyID.(string),
		roomID.(string), maintenanceRequestID.(string))
	req.Headers = map[string]string{
		"Authorization": token,
	}
	res := req.Delete()
	as.Equal(200, res.Code)

	var responseBody map[string]interface{}
	json.NewDecoder(res.Body).Decode(&responseBody)

	as.Equal(maintenanceRequestID.(string), responseBody["id"])

	maintenanceRequest := &models.MaintenanceRequest{}
	err = as.DB.Where("room_id = ?", roomID.(string)).First(maintenanceRequest)

	as.NotNil(err)
}
