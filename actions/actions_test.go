package actions

import (
	"testing"

	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/suite"
	"github.com/obedtandadjaja/project_k_backend/helpers"
)

type ActionSuite struct {
	*suite.Action
}

func Test_ActionSuite(t *testing.T) {
	action, err := suite.NewActionWithFixtures(App(), packr.New("Test_ActionSuite", "../fixtures"))
	if err != nil {
		t.Fatal(err)
	}

	as := &ActionSuite{
		Action: action,
	}
	suite.Run(t, as)
}

func AccessTokenHelper(user map[string]interface{}) string {
	token, _ := helpers.GenerateAccessToken(
		user["id"].(string),
		user["credential_uuid"].(string),
		user["type"].(string),
	)
	return token
}
