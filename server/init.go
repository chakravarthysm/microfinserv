package server

import (
	"finserv/data"
	"finserv/service"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func Init() {
	r := mux.NewRouter()

	ah := AuthHandlers{service.NewAuthService(data.NewAuthDB())}
	uh := UserHandlers{service.NewUsersService(data.NewUsersDB())}

	// Auth routes
	r.HandleFunc("/login", ah.login).Methods(http.MethodGet)
	r.HandleFunc("/logout", ah.logout).Methods(http.MethodGet)

	// User routes
	r.HandleFunc("/account", uh.getUser).Methods(http.MethodGet).Name("GetUser")
	r.HandleFunc("/account", uh.createUser).Methods(http.MethodPost).Name("CreateUser")
	r.HandleFunc("/account", uh.getUser).Methods(http.MethodPut).Name("UpdateUser")
	r.HandleFunc("/account", uh.deleteUser).Methods(http.MethodDelete).Name("DeleteUser")
	am := AuthMiddleware{data.NewAuthDB()}
	r.Use(am.authorizationHandler())
	log.Fatal(http.ListenAndServe("localhost:3000", r))
}
