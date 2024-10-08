package main

import (
	"fmt"
	"log"
	"net/http"
	"social-network/database/sqlite"

)

const port = 8080

type application struct {
	database sqlite.SqliteDB
}

func main() {
	var app application
	err := app.applyMigrations()
	if err != nil {
		log.Fatal(err)
	}

	conn, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	app.database = sqlite.SqliteDB{DB: conn}
	defer app.database.Connection().Close()

	log.Println("Starting application on port", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
