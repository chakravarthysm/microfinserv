package service

import (
	"errors"
	"finserv/data"
	"time"
)

const dateLayout = "2006-01-02 15:04:05"

type AccountService interface {
	NewAccount(request data.NewAccountRequest) (*data.NewAccountResponse, error)
	MakeTransaction(request data.TransactionRequest) (*data.TransactionResponse, error)
}

type DefaultAccountService struct {
	source data.AccountsDB
}

func (d DefaultAccountService) NewAccount(req data.NewAccountRequest) (*data.NewAccountResponse, error) {
	account := data.NewAccount(req.UserId, req.Amount)
	if newAccount, err := d.source.CreateAccount(account); err != nil {
		return nil, err
	} else {
		return &data.NewAccountResponse{AccountId: newAccount.AccountId}, nil
	}
}

func (d DefaultAccountService) MakeTransaction(req data.TransactionRequest) (*data.TransactionResponse, error) {

	if req.TransactionType == "withdrawal" {
		account, err := d.source.FindAccountById(req.AccountId)
		if err != nil {
			return nil, err
		}
		if !account.CanWithdraw(req.Amount) {
			return nil, errors.New("Insufficient balance in the account")
		}
	}

	t := data.Transaction{
		AccountId:       req.AccountId,
		Amount:          req.Amount,
		TransactionType: req.TransactionType,
		TransactionDate: time.Now().Format(dateLayout),
	}
	transaction, appError := d.source.SaveTransaction(t)
	if appError != nil {
		return nil, appError
	}

	response := data.TransactionResponse{
		TransactionId:   transaction.TransactionId,
		AccountId:       transaction.AccountId,
		Amount:          transaction.Amount,
		TransactionType: transaction.TransactionType,
		TransactionDate: transaction.TransactionDate,
	}

	return &response, nil
}

func NewAccountService(repo data.AccountsDB) DefaultAccountService {
	return DefaultAccountService{repo}
}
