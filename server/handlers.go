package server

import (
	"encoding/json"
	"finserv/domain"
	"finserv/services"
	"fmt"
	"net/http"
)

func greetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello!!")
}

func getAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	accountsHandlers := services.NewAccountService(domain.NewAccountRepository())

	account := accountsHandlers.GetAccount()
	json.NewEncoder(w).Encode(account)
}
