package server

import (
	"finserv/data"
	"finserv/service"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func InitRoutes() {
	r := mux.NewRouter()

	dbClient, err := data.GetDbClient()
	if err != nil {
		log.Fatal(err)
	}
	ah := AuthHandler{service.NewAuthService(data.NewAuthDB(dbClient))}
	ach := AccountHandler{service.NewAccountService(data.NewAccountsDB(dbClient))}
	uh := UserHandler{service.NewUsersService(data.NewUsersDB(dbClient))}

	// Auth routes
	r.HandleFunc("/login", ah.login).Methods(http.MethodPost).Name("Login")
	r.HandleFunc("/logout", ah.logout).Methods(http.MethodPost).Name("Logout")

	// User routes
	r.HandleFunc("/users/{user_id}", uh.getUser).Methods(http.MethodGet).Name("GetUser")
	r.HandleFunc("/users", uh.createUser).Methods(http.MethodPost).Name("CreateUser")
	r.HandleFunc("/users/{user_id}", uh.updateUser).Methods(http.MethodPut).Name("UpdateUser")
	r.HandleFunc("/users/{user_id}", uh.deleteUser).Methods(http.MethodDelete).Name("DeleteUser")

	// Account routes
	r.HandleFunc("/users/{user_id}/account", ach.NewAccount).Methods(http.MethodPost).Name("CreateAccount")
	r.HandleFunc("/users/{user_id}/account/{account_id}", ach.MakeTransaction).Methods(http.MethodPost).Name("Transaction")

	am := AuthMiddleware{data.NewAuthDB(dbClient)}
	r.Use(am.authorizationHandler())
	log.Println("Starting server on port 3000")
	log.Fatal(http.ListenAndServe("localhost:3000", r))
}
