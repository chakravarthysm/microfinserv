package data

import (
	"database/sql"
	"errors"
	"log"
	"time"
)

const dateLayout = "2006-01-02 15:04:05"

type NewAccountRequest struct {
	UserId int     `json:"user_id"`
	Amount float64 `json:"amount"`
}

type NewAccountResponse struct {
	AccountId int `json:"account_id"`
}

type TransactionRequest struct {
	AccountId       int     `json:"account_id"`
	Amount          float64 `json:"amount"`
	TransactionType string  `json:"transaction_type"`
	TransactionDate string  `json:"transaction_date"`
	UserId          int     `json:"-"`
}

type TransactionResponse struct {
	TransactionId   int     `json:"transaction_id"`
	AccountId       int     `json:"account_id"`
	Amount          float64 `json:"new_balance"`
	TransactionType string  `json:"transaction_type"`
	TransactionDate string  `json:"transaction_date"`
}

type Account struct {
	AccountId int     `db:"account_id"`
	UserId    int     `db:"user_id"`
	CreatedOn string  `db:"created_on"`
	Amount    float64 `db:"amount"`
	Status    int     `db:"status"`
}

type Transaction struct {
	TransactionId   int     `db:"transaction_id"`
	AccountId       int     `db:"account_id"`
	Amount          float64 `db:"amount"`
	TransactionType string  `db:"transaction_type"`
	TransactionDate string  `db:"transaction_date"`
}

type AccountsDB struct {
	client *sql.DB
}

func (a Account) CanWithdraw(amount float64) bool {
	return a.Amount >= amount
}

func NewAccount(userId int, amount float64) Account {
	return Account{
		UserId:    userId,
		CreatedOn: time.Now().Format(dateLayout),
		Amount:    amount,
		Status:    1,
	}
}

func (a AccountsDB) CreateAccount(accnt Account) (*Account, error) {
	accountCreateQuery := "INSERT INTO accounts (user_id, created_on, amount, status) values ($1, $2, $3, $4) RETURNING account_id"
	var id int
	err := a.client.QueryRow(accountCreateQuery, accnt.UserId, accnt.CreatedOn, accnt.Amount, 1).Scan(&id)
	if err != nil {
		log.Println("Error while creating new account: " + err.Error())
		return nil, errors.New("error while creating new account")
	}

	accnt.AccountId = id
	return &accnt, nil
}

func (a AccountsDB) SaveTransaction(t Transaction) (*Transaction, error) {
	tx, err := a.client.Begin()
	if err != nil {
		log.Println("Error while starting a new transaction for bank account transaction: " + err.Error())
		return nil, errors.New("unexpected database error")
	}

	row := tx.QueryRow(`INSERT INTO transactions (account_id, amount, transaction_type, transaction_date) 
											values ($1, $2, $3, $4) RETURNING transaction_id`, t.AccountId, t.Amount, t.TransactionType, t.TransactionDate)

	var transactionId int
	err = row.Scan(&transactionId)
	if err != nil {
		log.Println("Error while getting the last transaction id: " + err.Error())
		return nil, errors.New("unexpected database error")
	}
	if t.TransactionType == "withdrawal" {
		_, err = tx.Exec(`UPDATE accounts SET amount = amount - $1 WHERE account_id = $2`, t.Amount, t.AccountId)
	} else {
		_, err = tx.Exec(`UPDATE accounts SET amount = amount + $1 WHERE account_id = $2`, t.Amount, t.AccountId)
	}

	if err != nil {
		tx.Rollback()
		log.Println("Error while saving transaction: " + err.Error())
		return nil, errors.New("unexpected database error")
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Println("Error while commiting transaction for bank account: " + err.Error())
		return nil, errors.New("unexpected database error")
	}

	account, err := a.FindAccountById(t.AccountId)
	if err != nil {
		return nil, err
	}

	t.TransactionId = transactionId
	t.Amount = account.Amount
	return &t, nil
}

func (a AccountsDB) FindAccountById(accountId int) (*Account, error) {
	getAccountQuery := "SELECT account_id, user_id, created_on, amount from accounts where account_id = $1 AND status = $2"
	var account Account
	row := a.client.QueryRow(getAccountQuery, accountId, 1)
	err := row.Scan(
		&account.AccountId,
		&account.UserId,
		&account.CreatedOn,
		&account.Amount,
	)
	if err != nil {
		return &Account{}, err
	}
	return &account, nil
}

func NewAccountsDB(dbClient *sql.DB) AccountsDB {
	return AccountsDB{dbClient}
}
