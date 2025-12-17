package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "golang-event.db")

	if err != nil {
		panic("Failed to connect to database")
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	// Create tables
	createTables()
}

func createTables() {
	createEventTable := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		name TEXT NOT NULL, 
		description TEXT NOT NULL, 
		location TEXT NOT NULL, 
		datetime DATETIME NOT NULL, 
		user_id INTEGER
	)
	`

	_, err := DB.Exec(createEventTable)

	if err != nil {
		// panic("Failed to create event table")
		panic(err)
	}
}
