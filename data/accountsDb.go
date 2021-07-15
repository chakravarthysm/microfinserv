package data

import (
	"database/sql"

	"log"

	_ "github.com/lib/pq"
)

type AcccountsDBImpl interface {
	Insert(accnt Account) (Account, error)
	FindById(id string) (Account, error)
	UpdateById(id string, accnt Account) (Account, error)
	DeleteById(id string) error
}

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

func (a *AcccountsDB) Insert(account Account) (Account, error) {
	findAccountQuery := "Insert INTO accounts"

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

func (a *AcccountsDB) FindById(id string) (Account, error) {
	findAccountQuery := "SELECT * FROM accounts where account_id=id"

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

func (a *AcccountsDB) UpdateById(id string, account Account) (Account, error) {
	findAccountQuery := "UPDATE accounts where account_id=id"

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

func (a *AcccountsDB) DeleteById(id string) error {
	deleteAccountQuery := "DELETE FROM accounts where id=id"

	_, err := a.client.Query(deleteAccountQuery)
	if err != nil {
		return err
	}

	return nil
}

func NewAccountsDB() AcccountsDBImpl {
	connStr := "postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	return &AcccountsDB{db}
}
