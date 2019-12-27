package actions

import (
	"encoding/json"

	"github.com/obedtandadjaja/project_k_backend/models"
)

func (as *ActionSuite) Test_UsersResource_Show() {
	as.LoadFixture("user")
	res := as.JSON("/api/v1/users/%s", "user").Get()
	as.Equal(200, res.Code)

	var responseBody map[string]interface{}
	json.NewDecoder(res.Body).Decode(responseBody)

	as.Equal(responseBody["id"], "user")
	as.Equal(responseBody["email"], "user@example.com")
}

func (as *ActionSuite) Test_UsersResource_Update() {
	as.LoadFixture("user")
	res := as.JSON("/api/v1/users/%s", "user").Put(&models.User{Email: "asdf@example.com"})
	as.Equal(200, res.Code)

	user := &models.User{}
	as.DB.Find(user, "user")

	as.Equal(user.Email, "asdf@example.com")
}
