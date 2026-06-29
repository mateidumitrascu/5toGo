// Package database defines access to the database, drivers, connections
package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	_ "modernc.org/sqlite"
)

var MaxOpenCons = 5

func GetDatabaseConnection() (*sql.DB, error) {
	db, err := sql.Open("sqlite", os.Getenv("DB_PATH"))
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}
	db.SetMaxOpenConns(MaxOpenCons)
	return db, nil
}
