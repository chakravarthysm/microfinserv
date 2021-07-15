package server

import (
	"encoding/json"
	"finserv/data"
	"finserv/service"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type UserHandlers struct {
	service service.UserService
}

func (ac *UserHandlers) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["user_id"]
	user, err := ac.service.GetUser(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (ac *UserHandlers) createUser(w http.ResponseWriter, r *http.Request) {
	var usr data.User

	err := json.NewDecoder(r.Body).Decode(&usr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userCreated, err := ac.service.CreateUser(usr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userCreated)
}

func (ac *UserHandlers) updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["user_id"]
	w.Header().Add("Content-Type", "application/json")
	var user data.User
	user, err := ac.service.UpdateUser(id, user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
	}

	json.NewEncoder(w).Encode(user)
}

func (ac *UserHandlers) deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["user_id"]
	w.Header().Add("Content-Type", "application/json")
	err := ac.service.DeleteUser(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
	}

	w.WriteHeader(http.StatusNoContent)
}

func IsAuthenticated() {

}
