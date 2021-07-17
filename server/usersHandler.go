package server

import (
	"encoding/json"
	"finserv/data"
	"finserv/service"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type UserHandlers struct {
	service service.UserService
}

func (ac *UserHandlers) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["userId"]
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
	err = validateUserPayload(usr)
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
	id := vars["userId"]

	var usr data.User
	err := json.NewDecoder(r.Body).Decode(&usr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := ac.service.UpdateUser(id, usr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (ac *UserHandlers) deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["userId"]
	w.Header().Add("Content-Type", "application/json")
	err := ac.service.DeleteUser(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
	}

	w.WriteHeader(http.StatusNoContent)
}

func validateUserPayload(usr data.User) error {
	var user data.User
	var errs []string
	if usr.Username == user.Username {
		errs = append(errs, "Username")
	}
	if usr.Password == user.Password {
		errs = append(errs, "Password")
	}
	if usr.Name == user.Name {
		errs = append(errs, "Name")
	}
	if usr.Location == user.Location {
		errs = append(errs, "Location")
	}
	if usr.PAN == user.PAN {
		errs = append(errs, "PAN")
	}
	if usr.Address == user.Address {
		errs = append(errs, "Address")
	}
	if usr.ContactNumber == user.ContactNumber {
		errs = append(errs, "ContactNumber")
	}
	if usr.Gender == user.Gender {
		errs = append(errs, "Gender")
	}
	if usr.Nationality == user.Nationality {
		errs = append(errs, "Nationality")
	}

	return fmt.Errorf("%s is/are required", strings.Join(errs, ", "))
}
