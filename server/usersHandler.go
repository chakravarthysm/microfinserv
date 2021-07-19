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

type UserHandler struct {
	service service.UserServiceImpl
}

func (uh *UserHandler) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["user_id"]

	user, err := uh.service.GetUser(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (uh *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
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

	userCreated, err := uh.service.CreateUser(usr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userCreated)
}

func (uh *UserHandler) updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["user_id"]

	var usr data.User
	err := json.NewDecoder(r.Body).Decode(&usr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	existingUserData, err := uh.service.GetUser(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}

	prepareUpdatePayload(&usr, existingUserData)

	user, err := uh.service.UpdateUser(userId, usr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (uh *UserHandler) deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["user_id"]

	err := uh.service.DeleteUser(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
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

	if len(errs) > 0 {
		return fmt.Errorf("%s is/are required", strings.Join(errs, ", "))
	}

	return nil
}

func prepareUpdatePayload(usr *data.User, existingUserData data.UserResponse) {
	var user data.User
	if usr.Username == user.Username {
		usr.Username = existingUserData.Username
	}
	if usr.Name == user.Name {
		usr.Name = existingUserData.Name
	}
	if usr.Location == user.Location {
		usr.Location = existingUserData.Location
	}
	if usr.PAN == user.PAN {
		usr.PAN = existingUserData.PAN
	}
	if usr.Address == user.Address {
		usr.Address = existingUserData.Address
	}
	if usr.ContactNumber == user.ContactNumber {
		usr.ContactNumber = existingUserData.ContactNumber
	}
	if usr.Gender == user.Gender {
		usr.Gender = existingUserData.Gender
	}
	if usr.Nationality == user.Nationality {
		usr.Nationality = existingUserData.Nationality
	}
}
