package data

import (
	"database/sql"
)

func GetDbClient() (*sql.DB, error) {
	connStr := "user=postgres dbname=finserv password=postgres host=localhost sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return db, nil
}
