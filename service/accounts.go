package service

import "finserv/domain"

type AccountService interface {
	CreateAccount() domain.Account
	GetAccount() domain.Account
	UpdateAccount() domain.Account
	DeleteAccount() domain.Account
}

type AccountHandler struct {
	AccountDB domain.AcccountsDB
}

func (a *AccountHandler) CreateAccount() domain.Account {
	return a.AccountDB.FindBy()
}

func (a *AccountHandler) GetAccount(userId string) domain.Account {
	return a.AccountDB.FindBy(userId)
}

func (a *AccountHandler) UpdateAccount() domain.Account {
	return a.AccountDB.GetAccount()
}

func (a *AccountHandler) DeleteAccount() domain.Account {
	return a.AccountDB.GetAccount()
}

func NewAccountService(accountDb domain.AcccountsDB) AccountService {
	return &AccountHandler{
		AccountDB: accountDb,
	}
}
