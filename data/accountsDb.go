package data

import (
	"database/sql"
	"errors"
	"log"
	"time"

	uuid "github.com/satori/go.uuid"
)

type AccountsDB struct {
	client *sql.DB
}

func (a Account) CanWithdraw(amount float64) bool {
	return a.Balance >= amount
}

func NewAccount(userId string, amount float64) Account {
	return Account{
		UserId:    userId,
		CreatedOn: time.Now().Format(dateLayout),
		Balance:   amount,
		Status:    1,
	}
}

func (a AccountsDB) CreateAccount(accnt Account) (*Account, error) {

	// TODO: return error if user is inactive

	accountCreateQuery := "INSERT INTO accounts (account_id, user_id, created_on, balance, status) values ($1, $2, $3, $4, $5) RETURNING account_id"
	var id string
	err := a.client.QueryRow(accountCreateQuery, uuid.NewV4(), accnt.UserId, accnt.CreatedOn, accnt.Balance, 1).Scan(&id)
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

	// TODO: return error if user/account is inactive

	row := tx.QueryRow(`INSERT INTO transactions (transaction_id, account_id, balance, transaction_type, transaction_date) 
											values ($1, $2, $3, $4, $5) RETURNING transaction_id`, uuid.NewV4(), t.AccountId, t.Amount, t.TransactionType, t.TransactionDate)

	var transactionId string
	err = row.Scan(&transactionId)
	if err != nil {
		log.Println("Error while getting the last transaction id: " + err.Error())
		return nil, errors.New("unexpected database error")
	}
	if t.TransactionType == "withdrawal" {
		_, err = tx.Exec(`UPDATE accounts SET balance = balance - $1 WHERE account_id = $2`, t.Amount, t.AccountId)
	} else {
		_, err = tx.Exec(`UPDATE accounts SET balance = balance + $1 WHERE account_id = $2`, t.Amount, t.AccountId)
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
	t.Amount = account.Balance
	return &t, nil
}

func (a AccountsDB) FindAccountById(accountId string) (*Account, error) {
	getAccountQuery := "SELECT account_id, user_id, created_on, balance from accounts where account_id = $1 AND status = $2"
	var account Account
	row := a.client.QueryRow(getAccountQuery, accountId, 1)
	err := row.Scan(
		&account.AccountId,
		&account.UserId,
		&account.CreatedOn,
		&account.Balance,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return &Account{}, errors.New("Account not found")
		}
		return &Account{}, err
	}
	return &account, nil
}

func NewAccountsDB(dbClient *sql.DB) AccountsDB {
	return AccountsDB{dbClient}
}
