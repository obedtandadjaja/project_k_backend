package actions

import (
	"encoding/json"

	"github.com/gobuffalo/suite/fix"
	"github.com/obedtandadjaja/project_k_backend/helpers"
	"github.com/obedtandadjaja/project_k_backend/models"
)

func (as *ActionSuite) Test_RoomsResource_List() {
	as.LoadFixture("user with property with room")

	fixture, _ := fix.Find("user with property with room")
	propertyID := fixture.Tables[1].Row[0]["id"]
	roomID := fixture.Tables[3].Row[0]["id"]

	token := AccessTokenHelper(fixture.Tables[0].Row[0])

	req := as.JSON("/api/v1/properties/%s/rooms", propertyID.(string))
	req.Headers = map[string]string{
		"Authorization": token,
	}
	res := req.Get()
	as.Equal(200, res.Code)

	var responseBody []map[string]interface{}
	json.NewDecoder(res.Body).Decode(&responseBody)

	as.Equal(roomID.(string), responseBody[0]["id"])
	as.Equal(1, len(responseBody))
}

func (as *ActionSuite) Test_RoomsResource_Show() {
	as.LoadFixture("user with property with room")

	fixture, _ := fix.Find("user with property with room")
	propertyID := fixture.Tables[1].Row[0]["id"]
	roomID := fixture.Tables[3].Row[0]["id"]

	token := AccessTokenHelper(fixture.Tables[0].Row[0])

	req := as.JSON("/api/v1/properties/%s/rooms/%s", propertyID.(string), roomID.(string))
	req.Headers = map[string]string{
		"Authorization": token,
	}
	res := req.Get()
	as.Equal(200, res.Code)

	var responseBody map[string]interface{}
	json.NewDecoder(res.Body).Decode(&responseBody)

	as.Equal(roomID.(string), responseBody["id"])
}

func (as *ActionSuite) Test_RoomsResource_Create() {
	as.LoadFixture("user with property")

	fixture, err := fix.Find("user with property")
	propertyID := fixture.Tables[1].Row[0]["id"]

	token := AccessTokenHelper(fixture.Tables[0].Row[0])

	propertyToCreate := &models.Room{
		Name:            "room",
		PriceAmount:     10000000,
		PaymentSchedule: "monthly",
		Property:        &models.Property{ID: helpers.ParseUUID(propertyID.(string))},
	}
	req := as.JSON("/api/v1/properties/%s/rooms", propertyID.(string))
	req.Headers = map[string]string{
		"Authorization": token,
	}
	res := req.Post(propertyToCreate)
	as.Equal(201, res.Code)

	var responseBody map[string]interface{}
	json.NewDecoder(res.Body).Decode(&responseBody)

	as.Equal("room", responseBody["name"])
	as.Equal(10000000.0, responseBody["priceAmount"])
	as.Equal("monthly", responseBody["paymentSchedule"])

	room := &models.Room{}
	err = as.DB.Where("property_id = ?", propertyID.(string)).First(room)
	as.NoError(err)
}

func (as *ActionSuite) Test_RoomsResource_Update() {
	as.LoadFixture("user with property with room")

	fixture, _ := fix.Find("user with property with room")
	propertyID := fixture.Tables[1].Row[0]["id"]
	roomID := fixture.Tables[3].Row[0]["id"]

	token := AccessTokenHelper(fixture.Tables[0].Row[0])

	req := as.JSON("/api/v1/properties/%s/rooms/%s", propertyID.(string), roomID.(string))
	req.Headers = map[string]string{
		"Authorization": token,
	}
	res := req.Put(
		&models.Room{
			ID:         helpers.ParseUUID(roomID.(string)),
			PropertyID: helpers.ParseUUID(propertyID.(string)),
			Name:       "Changed",
		},
	)
	as.Equal(200, res.Code)

	room := &models.Room{}
	as.DB.Find(room, roomID.(string))

	as.Equal("Changed", room.Name)
}

func (as *ActionSuite) Test_RoomsResource_Destroy() {
	as.LoadFixture("user with property with room")

	fixture, err := fix.Find("user with property with room")
	propertyID := fixture.Tables[1].Row[0]["id"]
	roomID := fixture.Tables[3].Row[0]["id"]

	token := AccessTokenHelper(fixture.Tables[0].Row[0])

	req := as.JSON("/api/v1/properties/%s/rooms/%s", propertyID.(string), roomID.(string))
	req.Headers = map[string]string{
		"Authorization": token,
	}
	res := req.Delete()
	as.Equal(200, res.Code)

	var responseBody map[string]interface{}
	json.NewDecoder(res.Body).Decode(&responseBody)

	as.Equal(roomID.(string), responseBody["id"])

	room := &models.Room{}
	err = as.DB.Where("property_id = ?", propertyID.(string)).First(room)
	as.NotNil(err)
}
