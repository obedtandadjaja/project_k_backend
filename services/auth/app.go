package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/obedtandadjaja/project_k_backend/services/auth/controller"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
	Env    string
}

func (app *App) Initialize(env, host, port, user, password, dbName string) {
	app.Env = env

	err := app.initializeDB(host, port, user, password, dbName)
	if err != nil {
		log.Fatal(err)
		return
	}
	app.runMigration()

	app.Router = mux.NewRouter()
	sharedResources := &controller.SharedResources{DB: app.DB, Env: env}
	app.initializeRoutes(sharedResources)
}

func (app *App) initializeDB(host, port, user, password, dbName string) error {
	var connectionString string

	connectionString = fmt.Sprintf(
		"postgresql://%v:%v@%v:%v/%v?sslmode=disable",
		user, password, host, port, dbName,
	)

	fmt.Printf("Connecting to database... %v\n", connectionString)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		return err
	}

	app.DB = db
	return nil
}

func (app *App) initializeRoutes(sr *controller.SharedResources) {
	for _, route := range routes {
		app.Router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(logRequestMiddleware(controller.Handler{sr, route.HandlerFunc}))
	}
}

func (app *App) runMigration() {
	driver, err := postgres.WithInstance(app.DB, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
		return
	}

	migration, err := migrate.NewWithDatabaseInstance("file://./migrations", "postgres", driver)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println("Starting migration...")

	if err = migration.Up(); err != migrate.ErrNoChange && err != nil {
		log.Fatal(err)
	}
}

func (app *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, app.Router))
}

// TODO: move this out to a middleware package
func logRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestDump, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Println(err)
		}
		log.Println(string(requestDump))

		next.ServeHTTP(w, r)
	})
}
