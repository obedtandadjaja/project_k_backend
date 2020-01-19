package actions

import (
	"encoding/json"

	"github.com/gobuffalo/suite/fix"
	"github.com/obedtandadjaja/project_k_backend/helpers"
	"github.com/obedtandadjaja/project_k_backend/models"
)

func (as *ActionSuite) Test_UsersResource_Show() {
	as.LoadFixture("user")

	fixture, _ := fix.Find("user")
	userID := fixture.Tables[0].Row[0]["id"]

	token := AccessTokenHelper(fixture.Tables[0].Row[0])

	req := as.JSON("/api/v1/users/%s", userID.(string))
	req.Headers = map[string]string{
		"Authorization": token,
	}
	res := req.Get()
	as.Equal(200, res.Code)

	var responseBody map[string]interface{}
	json.NewDecoder(res.Body).Decode(&responseBody)

	as.Equal(responseBody["id"], userID.(string))
	as.Equal(responseBody["email"], "user@example.com")
}

func (as *ActionSuite) Test_UsersResource_Update() {
	as.LoadFixture("user")

	fixture, _ := fix.Find("user")
	userID := fixture.Tables[0].Row[0]["id"]

	token := AccessTokenHelper(fixture.Tables[0].Row[0])

	req := as.JSON("/api/v1/users/%s", userID.(string))
	req.Headers = map[string]string{
		"Authorization": token,
	}
	res := req.Put(&models.User{
		ID:    helpers.ParseUUID(userID.(string)),
		Email: "asdf@example.com",
		Type:  models.USER_TENANT,
	})
	as.Equal(200, res.Code)

	user := &models.User{}
	as.DB.Find(user, userID.(string))

	as.Equal("asdf@example.com", user.Email)
}
