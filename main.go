package main

import (
	"log"

	"github.com/gobuffalo/pop"
	"github.com/obedtandadjaja/project_k_backend/actions"
	"github.com/obedtandadjaja/project_k_backend/models"
)

// main is the starting point for your Buffalo application.
// You can feel free and add to this `main` method, change
// what it does, etc...
// All we ask is that, at some point, you make sure to
// call `app.Serve()`, unless you don't want to start your
// application that is. :)
func main() {
	// Execute database migrations.
	// We need to do this here since it is not possible to run
	// `soda migrate up` in GAE. Remove this if we switch back to
	// GKE. see: https://golangtesting.com/posts/gobuffalo-app-engine
	migrator, err := pop.NewFileMigrator("./migrations", models.DB)
	if err != nil {
		panic(err)
	}
	migrator.Up()

	app := actions.App()
	if err := app.Serve(); err != nil {
		log.Fatal(err)
	}
}

/*
# Notes about `main.go`

## SSL Support

We recommend placing your application behind a proxy, such as
Apache or Nginx and letting them do the SSL heavy lifting
for you. https://gobuffalo.io/en/docs/proxy

## Buffalo Build

When `buffalo build` is run to compile your binary, this `main`
function will be at the heart of that binary. It is expected
that your `main` function will start your application using
the `app.Serve()` method.

*/
