// Package database defines access to the database, drivers, connections
package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	_ "modernc.org/sqlite"
)

func GetDatabasConnection() *sql.DB {
	db, err := sql.Open("sqlite", os.Getenv("DB_PATH"))
	if err != nil {
		fmt.Printf("Error opening database: %v\n", err)
	}
	if err := db.Ping(); err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		return nil
	}

	return db
}
