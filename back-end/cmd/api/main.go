package main

import (
	"back-end/database/sqlite"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

const port = 8080

type application struct {
	database sqlite.SqliteDB // Mutex to protect concurrent access to connections
}

func main() {
	var app application

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
