package server

import (
	"encoding/json"
	"finserv/service"
	"net/http"
)

type AuthHandlers struct {
	service service.DefaultAuthService
}

type Creds struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (au *AuthHandlers) login(w http.ResponseWriter, r *http.Request) {
	var c Creds
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

func (au *AuthHandlers) logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
