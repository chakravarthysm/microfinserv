package server

import (
	"encoding/json"
	"finserv/data"
	"finserv/service"
	"net/http"
)

type AuthHandlers struct {
	service service.AuthService
}

type Creds struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (au *AuthHandlers) login(w http.ResponseWriter, r *http.Request) {
	err := json.NewDecoder(r.Body).Decode(&Creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

func (au *AuthHandlers) logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	accountsHandlers := service.NewAccountService(data.NewAccountsDB())

	account := accountsHandlers.GetAccount()
	json.NewEncoder(w).Encode(account)
}

type AccountHandlers struct {
	service service.AccountService
}

func (ac *AccountHandlers) getAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	accountsHandlers := service.NewAccountService(data.NewAccountsDB())

	account := accountsHandlers.GetAccount()
	json.NewEncoder(w).Encode(account)
}

func (ac *AccountHandlers) createAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	accountsHandlers := service.NewAccountService(data.NewAccountsDB())

	account := accountsHandlers.GetAccount()
	json.NewEncoder(w).Encode(account)
}

func (ac *AccountHandlers) updateAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	accountsHandlers := service.NewAccountService(data.NewAccountsDB())

	account := accountsHandlers.GetAccount()
	json.NewEncoder(w).Encode(account)
}

func (ac *AccountHandlers) deleteAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	accountsHandlers := service.NewAccountService(data.NewAccountsDB())

	account := accountsHandlers.GetAccount()
	json.NewEncoder(w).Encode(account)
}
