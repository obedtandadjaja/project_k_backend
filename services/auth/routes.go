package main

import (
	"net/http"

	"github.com/obedtandadjaja/project_k_backend/services/auth/controller"
	"github.com/obedtandadjaja/project_k_backend/services/auth/controller/credentials"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc func(sr *controller.SharedResources, w http.ResponseWriter, r *http.Request) error
}

type Routes []Route

var routes = Routes{
	Route{
		"Health",
		"GET",
		"/health",
		controller.Health,
	},
	Route{
		"Token",
		"POST",
		"/token",
		controller.Token,
	},
	Route{
		"Token",
		"POST",
		"/login",
		controller.Login,
	},
	Route{
		"VerifySessionToken",
		"POST",
		"/verify_session_token",
		controller.VerifySessionToken,
	},
	Route{
		"Verify",
		"POST",
		"/verify",
		controller.Verify,
	},
	Route{
		"CreateCredential",
		"POST",
		"/credentials",
		credentials.Create,
	},
	Route{
		"DeleteCredential",
		"DELETE",
		"/credentials",
		credentials.Delete,
	},
	Route{
		"UpdateCredential",
		"PUT",
		"/credentials/{uuid}",
		credentials.Update,
	},
	Route{
		"ResetPassword",
		"POST",
		"/credentials/reset_password",
		credentials.ResetPassword,
	},
	Route{
		"InitiateResetPassword",
		"POST",
		"/credentials/initiate_password_reset",
		credentials.InitiatePasswordReset,
	},
}
