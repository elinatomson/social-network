package main

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var dbPath = "./database/database.db"

func openDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (app *application) connectToDB() (*sql.DB, error) {
	connection, err := openDB()
	if err != nil {
		return nil, err
	}
	log.Println("Connected to database")
	return connection, nil
}

func (app *application) applyMigrations() error {
	// Create a new SQLite driver for Golang Migrate
	db, err := openDB()
	if err != nil {
		return err
	}

	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		return err
	}

	// Create a new migration instance
	m, err := migrate.NewWithDatabaseInstance("file://./database/migrations", "sqlite3", driver)
	if err != nil {
		return err
	}

	// Migrate the database to the latest version
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
