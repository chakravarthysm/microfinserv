package services

import "finserv/domain"

type AccountService interface {
	GetAccount() domain.Account
}

type AccountHandler struct {
	AccountRepo domain.AccountsRepository
}

func (a *AccountHandler) GetAccount() domain.Account {
	return a.AccountRepo.GetAccount()
}

func NewAccountService(accountRepo domain.AccountsRepository) AccountService {
	return &AccountHandler{
		AccountRepo: accountRepo,
	}
}
