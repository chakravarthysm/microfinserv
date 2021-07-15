package server

import (
	"encoding/json"
	"finserv/data"
	"finserv/service"
	"fmt"
	"net/http"
)

type AuthHandlers struct {
	service service.DefaultAuthService
}

func (au *AuthHandlers) login(w http.ResponseWriter, r *http.Request) {
	var lr data.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&lr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	loginResponse, err := au.service.Login(lr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(loginResponse)
}

func (au *AuthHandlers) logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
