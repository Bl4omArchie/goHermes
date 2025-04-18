package core

import (
	"database/sql"
	"fmt"
	"os"
)

type Document struct {
	Title string
	Authors []string
	Download_path string
	Category string
	Document_hash string
	Release_date string
}


type Author struct {
	First_name string
	Last_name string
}

func CreateDB() (*sql.DB, error) {
	dbPath := "sql/eprint.db"

	// Create db file
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		// Create the directory if it doesn't exist
		if err := os.MkdirAll("sql", os.ModePerm); err != nil {
			return nil, fmt.Errorf("failed to create directory: %v", err)
		}

		// Create the database file
		file, err := os.Create(dbPath)
		if err != nil {
			return nil, fmt.Errorf("failed to create database file: %v", err)
		}
		file.Close()
	}

	// Open a connection
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	return db, nil
}
