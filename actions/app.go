package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	forcessl "github.com/gobuffalo/mw-forcessl"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
	"github.com/rs/cors"
	"github.com/unrolled/secure"

	"github.com/gobuffalo/buffalo-pop/pop/popmw"
	contenttype "github.com/gobuffalo/mw-contenttype"
	"github.com/gobuffalo/x/sessions"
	"github.com/obedtandadjaja/project_k_backend/helpers"
	"github.com/obedtandadjaja/project_k_backend/models"
)

var ENV = envy.Get("ENV", "development")
var app *buffalo.App

func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:          ENV,
			SessionStore: sessions.Null{},
			SessionName:  "_project_k_backend_session",
		})

		if app.Env == "development" {
			app.PreWares = []buffalo.PreWare{cors.New(cors.Options{
				AllowedOrigins:   []string{"*"},
				AllowedMethods:   []string{"OPTIONS", "GET", "POST", "PUT", "DELETE"},
				AllowedHeaders:   []string{"Content-Type", "Cookie", "Authorization"},
				AllowCredentials: true,
			}).Handler}
		}

		// Automatically redirect to SSL
		// app.Use(forceSSL())

		// parse access token and set current_user
		app.Use(parseAccessToken)

		// Log request parameters (filters apply).
		app.Use(paramlogger.ParameterLogger)

		// Set the request content type to JSON
		app.Use(contenttype.Set("application/json"))

		app.Use(popmw.Transaction(models.DB))

		app.GET("/api/health", Health)
		app.POST("/api/v1/token", Token)
		app.POST("/api/v1/signup", Signup)
		app.POST("/api/v1/login", Login)
		app.Resource("/api/v1/users", UsersResource{})
		app.Resource("/api/v1/properties", PropertiesResource{})
		app.Resource("/api/v1/properties/{property_id}/rooms", RoomsResource{})
		app.Resource("/api/v1/properties/{property_id}/rooms/{room_id}/tenants", TenantsResource{})
		app.Resource("/api/v1/properties/{property_id}/rooms/{room_id}/tenants/{tenant_id}/payments", PaymentsResource{})
	}

	return app
}

// forceSSL will return a middleware that will redirect an incoming request
// if it is not HTTPS. "http://example.com" => "https://example.com".
// This middleware does **not** enable SSL. for your application. To do that
// we recommend using a proxy: https://gobuffalo.io/en/docs/proxy
// for more information: https://github.com/unrolled/secure/
func forceSSL() buffalo.MiddlewareFunc {
	return forcessl.Middleware(secure.Options{
		SSLRedirect:     ENV == "production" || ENV == "stage",
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
	})
}

func parseAccessToken(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if jwt := c.Request().Header.Get("Authorization"); jwt != "" {
			userID, _ := helpers.VerifyAccessToken(jwt)
			c.Set("current_user_id", userID)
		}

		return next(c)
	}
}
