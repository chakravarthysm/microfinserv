package server

import (
	"finserv/data"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type AuthMiddleware struct {
	source data.AuthDBImpl
}

func (a AuthMiddleware) authorizationHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			currentRoute := mux.CurrentRoute(r)
			fmt.Println("name", currentRoute.GetName())
			if currentRoute.GetName() == "CreateUser" {
				next.ServeHTTP(w, r)
			}

			if authHeader != "" {
				token := getTokenFromHeader(authHeader)

				isAuthorized := a.source.IsAuthorized(token)

				if isAuthorized {
					next.ServeHTTP(w, r)
				} else {
					w.WriteHeader(http.StatusForbidden)
					fmt.Fprint(w, []byte("Unauthorized"))
				}
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprint(w, []byte("Missing token"))
			}
		})
	}
}

func getTokenFromHeader(header string) string {
	splitToken := strings.Split(header, "Bearer")
	if len(splitToken) == 2 {
		return strings.TrimSpace(splitToken[1])
	}
	return ""
}
