package data

import "database/sql"

type AuthDB struct {
	client *sql.DB
}
