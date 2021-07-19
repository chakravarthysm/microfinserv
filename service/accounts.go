package service

import (
	"errors"
	"finserv/data"
	"time"
)

const dateLayout = "2006-01-02 15:04:05"

type AccountServiceImpl interface {
	NewAccount(data.NewAccountRequest) (*data.NewAccountResponse, error)
	MakeTransaction(data.TransactionRequest) (*data.TransactionResponse, error)
}
type AccountService struct {
	source data.AccountsDB
}

func (d *AccountService) NewAccount(req data.NewAccountRequest) (*data.NewAccountResponse, error) {
	account := data.NewAccount(req.UserId, req.Amount)
	if newAccount, err := d.source.CreateAccount(account); err != nil {
		return nil, err
	} else {
		return &data.NewAccountResponse{AccountId: newAccount.AccountId}, nil
	}
}

func (d *AccountService) MakeTransaction(req data.TransactionRequest) (*data.TransactionResponse, error) {

	if req.TransactionType == "withdrawal" {
		account, err := d.source.FindAccountById(req.AccountId)
		if err != nil {
			return nil, err
		}
		if !account.CanWithdraw(req.Amount) {
			return nil, errors.New("insufficient balance in the account")
		}
	}

	t := data.Transaction{
		AccountId:       req.AccountId,
		Amount:          req.Amount,
		TransactionType: req.TransactionType,
		TransactionDate: time.Now().Format(dateLayout),
	}
	transaction, err := d.source.SaveTransaction(t)
	if err != nil {
		return nil, err
	}

	response := data.TransactionResponse{
		TransactionId:   transaction.TransactionId,
		AccountId:       transaction.AccountId,
		Balance:         transaction.Amount,
		TransactionType: transaction.TransactionType,
		TransactionDate: transaction.TransactionDate,
	}

	return &response, nil
}

func NewAccountService(repo data.AccountsDB) AccountServiceImpl {
	return &AccountService{repo}
}
