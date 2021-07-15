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

	auh := AuthHandlers{service.NewAuthService(data.NewAuthDB())}
	ach := AccountHandlers{service.NewAccountService(data.NewAccountsDB())}

	// Auth routes
	r.HandleFunc("/login", auh.login).Methods(http.MethodGet)
	r.HandleFunc("/logout", auh.logout).Methods(http.MethodGet)

	// Account routes
	r.HandleFunc("/account", ach.getAccount).Methods(http.MethodGet)
	r.HandleFunc("/account", ach.createAccount).Methods(http.MethodPost)
	r.HandleFunc("/account", ach.getAccount).Methods(http.MethodPut)
	r.HandleFunc("/account", ach.deleteAccount).Methods(http.MethodDelete)
	am := AuthMiddleware{data.NewAuthDB()}
	r.Use(am.authorizationHandler())
	log.Fatal(http.ListenAndServe("localhost:3000", r))
}
