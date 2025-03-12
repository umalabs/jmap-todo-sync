package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func SetupDatabase() *sql.DB {
	dbPath := "./todos.db"
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		log.Println("Database file does not exist, creating...")
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS todos (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			is_completed INTEGER NOT NULL DEFAULT 0
		);
	`)
	if err != nil {
		log.Fatal(err)
	}

	return db
}
