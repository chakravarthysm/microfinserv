package server

import (
	"encoding/json"
	"finserv/data"
	"finserv/service"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type AccountHandlers struct {
	service service.AccountService
}

func (ac *AccountHandlers) getAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["account_id"]
	w.Header().Add("Content-Type", "application/json")
	account, err := ac.service.GetAccount(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
	}
	json.NewEncoder(w).Encode(account)
}

func (ac *AccountHandlers) createAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	//json.NewEncoder(w).Encode(account)
}

func (ac *AccountHandlers) updateAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["account_id"]
	w.Header().Add("Content-Type", "application/json")
	var account data.Account
	account, err := ac.service.UpdateAccount(id, account)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
	}

	json.NewEncoder(w).Encode(account)
}

func (ac *AccountHandlers) deleteAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["account_id"]
	w.Header().Add("Content-Type", "application/json")
	err := ac.service.DeleteAccount(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
	}

	w.WriteHeader(http.StatusNoContent)
}

func IsAuthenticated() {

}
