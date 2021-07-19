package data

import (
	"database/sql"
	"time"
)

func GetDbClient() (*sql.DB, error) {
	connStr := "user=postgres dbname=finserv password=postgres host=localhost sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute * 5)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db, nil
}
