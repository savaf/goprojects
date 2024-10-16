package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// Create an exported global variable to hold the database connection pool.
var DB *sql.DB

// InitDB initializes the SQLite connection and stores it in the package-level variable
func ConnectToDB(dataSourceName string) (*sql.DB, error) {
	if DB != nil {
		return DB, nil
	}

	var err error
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	DB = db

	return DB, nil
}

// CloseDB closes the database connection
func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}

func GetDB() (*sql.DB, error) {
	if DB != nil {
		return DB, nil
	}

	return nil, fmt.Errorf("error opening database")
}
