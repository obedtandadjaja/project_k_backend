package actions

import (
	"encoding/json"
	"strings"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/suite/fix"
	"github.com/obedtandadjaja/project_k_backend/helpers"
	"github.com/obedtandadjaja/project_k_backend/models"
)

func (as *ActionSuite) Test_AdminMaintenanceRequestsResource_List() {
	as.LoadFixture("user with property with maintenance request")

	fixture, _ := fix.Find("user with property with maintenance request")
	maintenanceRequestID := fixture.Tables[3].Row[0]["id"]

	token := AccessTokenHelper(fixture.Tables[0].Row[0])

	req := as.JSON("/api/v1/maintenance_requests")
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

// TODO(obedt): This is not the cleanest way to test. We are only testing negative test
//              cases here. We should be adding positive test cases. The naming of the
//              test helper should also be changed to something more clear.
func (as *ActionSuite) Test_AdminMaintenanceRequestsResource_ListWithStatusParams() {
	CheckAdminMaintenanceRequestResourceCount(as, 0, "status=closed")
}

func (as *ActionSuite) Test_AdminMaintenanceRequestsResource_ListWithPropertyIDParams() {
	CheckAdminMaintenanceRequestResourceCount(as, 0, "property_id=8f255a27-fdfc-48a6-8092-d9db6d9347d7")
}

func (as *ActionSuite) Test_AdminMaintenanceRequestsResource_ListWithRoomIDParams() {
	CheckAdminMaintenanceRequestResourceCount(as, 0, "room_id=8f255a27-fdfc-48a6-8092-d9db6d9347d7")
}

func (as *ActionSuite) Test_AdminMaintenanceRequestsResource_ListWithOpenedStartDateParams() {
	CheckAdminMaintenanceRequestResourceCount(as, 0, "opened_start_date=2050-01-01")
}

func (as *ActionSuite) Test_AdminMaintenanceRequestsResource_ListWithClosedStartDateParams() {
	CheckAdminMaintenanceRequestResourceCount(as, 0, "closed_start_date=1970-01-01")
}

func CheckAdminMaintenanceRequestResourceCount(as *ActionSuite, count int, params ...string) {
	as.LoadFixture("user with property with maintenance request")

	fixture, _ := fix.Find("user with property with maintenance request")

	token := AccessTokenHelper(fixture.Tables[0].Row[0])

	req := as.JSON("/api/v1/maintenance_requests?" + strings.Join(params, "&"))
	req.Headers = map[string]string{
		"Authorization": token,
	}
	res := req.Get()
	as.Equal(200, res.Code)

	var responseBody []map[string]interface{}
	json.NewDecoder(res.Body).Decode(&responseBody)

	as.Equal(count, len(responseBody))
}

func (as *ActionSuite) Test_AdminMaintenanceRequestsResource_Show() {
	as.LoadFixture("user with property with maintenance request")

	fixture, _ := fix.Find("user with property with maintenance request")
	maintenanceRequestID := fixture.Tables[3].Row[0]["id"]

	token := AccessTokenHelper(fixture.Tables[0].Row[0])

	req := as.JSON("/api/v1/maintenance_requests/%s", maintenanceRequestID.(string))
	req.Headers = map[string]string{
		"Authorization": token,
	}
	res := req.Get()
	as.Equal(200, res.Code)

	var responseBody map[string]interface{}
	json.NewDecoder(res.Body).Decode(&responseBody)

	as.Equal(maintenanceRequestID.(string), responseBody["id"])
}

func (as *ActionSuite) Test_AdminMaintenanceRequestsResource_CreateRoomMaintenanceRequest() {
	as.LoadFixture("user with property with room")

	fixture, err := fix.Find("user with property with room")
	userID := fixture.Tables[0].Row[0]["id"]
	propertyID := fixture.Tables[1].Row[0]["id"]
	roomID := fixture.Tables[3].Row[0]["id"]

	token := AccessTokenHelper(fixture.Tables[0].Row[0])

	maintenanceRequestToCreate := &models.MaintenanceRequest{
		Title:       "title",
		Description: nulls.String{String: "description", Valid: true},
		ReporterID:  helpers.ParseUUID(userID.(string)),
		RoomID:      nulls.UUID{UUID: helpers.ParseUUID(roomID.(string)), Valid: true},
		PropertyID:  nulls.UUID{UUID: helpers.ParseUUID(propertyID.(string)), Valid: true},
	}
	req := as.JSON("/api/v1/maintenance_requests")
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

func (as *ActionSuite) Test_AdminMaintenanceRequestsResource_CreatePropertyMaintenanceRequest() {
	as.LoadFixture("user with property")

	fixture, err := fix.Find("user with property")
	userID := fixture.Tables[0].Row[0]["id"]
	propertyID := fixture.Tables[1].Row[0]["id"]

	token := AccessTokenHelper(fixture.Tables[0].Row[0])

	maintenanceRequestToCreate := &models.MaintenanceRequest{
		Title:       "title",
		Description: nulls.String{String: "description", Valid: true},
		ReporterID:  helpers.ParseUUID(userID.(string)),
		PropertyID:  nulls.UUID{UUID: helpers.ParseUUID(propertyID.(string)), Valid: true},
	}
	req := as.JSON("/api/v1/maintenance_requests")
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

func (as *ActionSuite) Test_AdminMaintenanceRequestsResource_Update() {
	as.LoadFixture("user with property with maintenance request")

	fixture, _ := fix.Find("user with property with maintenance request")
	userID := fixture.Tables[0].Row[0]["id"]
	propertyID := fixture.Tables[1].Row[0]["id"]
	maintenanceRequestID := fixture.Tables[3].Row[0]["id"]

	token := AccessTokenHelper(fixture.Tables[0].Row[0])

	req := as.JSON("/api/v1/maintenance_requests/%s", maintenanceRequestID.(string))
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

func (as *ActionSuite) Test_AdminMaintenanceRequestsResource_Destroy() {
	as.LoadFixture("user with property with maintenance request")

	fixture, err := fix.Find("user with property with maintenance request")
	maintenanceRequestID := fixture.Tables[3].Row[0]["id"]

	token := AccessTokenHelper(fixture.Tables[0].Row[0])

	req := as.JSON("/api/v1/maintenance_requests/%s", maintenanceRequestID.(string))
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

func (as *ActionSuite) Test_AdminMaintenanceRequestsResource_CompleteMaintenanceRequest() {
	as.LoadFixture("user with property with maintenance request")

	fixture, _ := fix.Find("user with property with maintenance request")
	maintenanceRequestID := fixture.Tables[3].Row[0]["id"]

	token := AccessTokenHelper(fixture.Tables[0].Row[0])

	req := as.JSON("/api/v1/maintenance_requests/%s/complete", maintenanceRequestID.(string))
	req.Headers = map[string]string{
		"Authorization": token,
	}
	res := req.Post("")
	as.Equal(200, res.Code)

	var responseBody map[string]interface{}
	json.NewDecoder(res.Body).Decode(&responseBody)

	as.Equal("closed", responseBody["status"])
	as.NotNil(responseBody["completedAt"])
}
