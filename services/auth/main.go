package main

import (
	"log"
	"os"
)

func main() {
	app := App{}
	app.Initialize(
		os.Getenv("ENV"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	log.Println("App running on port " + os.Getenv("APP_PORT"))
	app.Run(os.Getenv("APP_HOST") + ":" + os.Getenv("APP_PORT"))
}
