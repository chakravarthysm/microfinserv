package service

import "finserv/data"

type AccountService interface {
	CreateAccount(account data.Account) (data.Account, error)
	GetAccount(id string) (data.Account, error)
	UpdateAccount(id string, account data.Account) (data.Account, error)
	DeleteAccount(id string) error
}

type AccountHandler struct {
	AccountDB data.AcccountsDBImpl
}

func (a *AccountHandler) CreateAccount(accnt data.Account) (data.Account, error) {
	return a.AccountDB.Insert(accnt)
}

func (a *AccountHandler) GetAccount(id string) (data.Account, error) {
	return a.AccountDB.FindById(id)
}

func (a *AccountHandler) UpdateAccount(id string, accnt data.Account) (data.Account, error) {
	return a.AccountDB.UpdateById(id, accnt)
}

func (a *AccountHandler) DeleteAccount(id string) error {
	return a.AccountDB.DeleteById(id)
}

func NewAccountService(accountDb data.AcccountsDBImpl) *AccountHandler {
	return &AccountHandler{
		AccountDB: accountDb,
	}
}
