package server

import (
	"finserv/data"
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
			if currentRoute.GetName() == "CreateUser" || currentRoute.GetName() == "login" {
				next.ServeHTTP(w, r)
				return
			}

			if authHeader != "" {
				token := getTokenFromHeader(authHeader)

				isAuthorized := a.source.IsAuthorized(token)

				if isAuthorized {
					next.ServeHTTP(w, r)
					return
				} else {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
			} else {
				http.Error(w, "Missing Token", http.StatusUnauthorized)
				return
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
