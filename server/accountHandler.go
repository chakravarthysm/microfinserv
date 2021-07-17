package server

import (
	"encoding/json"
	"finserv/data"
	"finserv/service"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type AccountHandler struct {
	service service.DefaultAccountService
}

func (a AccountHandler) NewAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["user_id"]
	var request data.NewAccountRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		usrId, err := strconv.Atoi(userId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		request.UserId = usrId
		account, err := a.service.NewAccount(request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			if err := json.NewEncoder(w).Encode(account); err != nil {
				panic(err)
			}
		}
	}
}

func (a AccountHandler) MakeTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountId := vars["account_id"]
	userId := vars["user_id"]

	var request data.TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		accntId, err := strconv.Atoi(accountId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		usrId, err := strconv.Atoi(userId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		request.AccountId = accntId
		request.UserId = usrId

		account, err := a.service.MakeTransaction(request)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(account); err != nil {
				panic(err)
			}
		}
	}

}
