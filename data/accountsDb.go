package data

import (
	"database/sql"

	"log"

	_ "github.com/lib/pq"
)

type Account struct {
	Name          string
	Location      string
	PAN           string
	Address       string
	ContactNumber int
	Gender        string
	Nationality   string
}

type AcccountsDB struct {
	client *sql.DB
}

func (a *AcccountsDB) FindBy(userId string) (Account, error) {
	findAccountQuery := "SELECT * FROM ACCOUNTS where id=userId"

	rows, err := a.client.Query(findAccountQuery)
	if err != nil {
		return Account{}, err
	}

	var accnt Account
	for rows.Next() {
		err := rows.Scan(&accnt)
		if err != nil {
			return Account{}, err
		}
	}

	return accnt, nil
}

func NewAccountsDB() *sql.DB {
	connStr := "postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	return db
}
