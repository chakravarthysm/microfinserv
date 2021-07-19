package server

import (
	"encoding/json"
	"finserv/common"
	"finserv/data"
	"finserv/service"
	"fmt"
	"net/http"
)

type AuthHandler struct {
	service service.AuthService
}

func (ah *AuthHandler) login(w http.ResponseWriter, r *http.Request) {
	var lr data.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&lr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	loginResponse, err := ah.service.Login(lr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(loginResponse)
}

func (ah *AuthHandler) logout(w http.ResponseWriter, r *http.Request) {
	token := common.GetTokenFromHeader(r.Header.Get("Authorization"))
	err := ah.service.Logout(token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Logged out")
}
