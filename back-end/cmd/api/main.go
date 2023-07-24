package main

import (
	"back-end/database/repository"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

const port = 8080

type application struct {
	database repository.SqliteDB
}

func main() {

	var app application

	conn, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	app.database = repository.SqliteDB{DB: conn}
	defer app.database.Connection().Close()

	//conn, err := openDB()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//app.DB = repository.SqliteDB{DB: conn}

	//defer app.DB.Connection().Close()

	log.Println("Starting application on port", port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
