package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"sort"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var dbPath = "./back-end/database/database.db"

func openDB() (*sql.DB, error) {
	exists := true

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		exists = false
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	if !exists {
		path := "./database/migrations/"
		fileNames, err := os.ReadDir(path)
		if err != nil {
			return nil, err
		}

		sort.SliceStable(fileNames, func(i, j int) bool {
			return fileNames[i].Name() < fileNames[j].Name()
		})

		tx, err := db.Begin()
		if err != nil {
			return nil, err
		}

		for _, name := range fileNames {
			fileName := name.Name()

			if !strings.Contains(fileName, ".down") {
				readFile, err := os.ReadFile(path + fileName)
				if err != nil {
					tx.Rollback()
					return nil, err
				}

				_, err = tx.ExecContext(context.Background(), string(readFile))
				if err != nil {
					tx.Rollback()
					return nil, err
				}
			}
		}

		err = tx.Commit()
		if err != nil {
			return nil, err
		}
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
