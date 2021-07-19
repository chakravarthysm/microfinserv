package server

import (
	"finserv/common"
	"finserv/data"
	"net/http"

	"github.com/gorilla/mux"
)

type AuthMiddleware struct {
	source data.AuthDBImpl
}

var excludedRoutes = map[string]struct{}{"CreateUser": struct{}{}, "Login": struct{}{}, "Logout": struct{}{}}

func (a AuthMiddleware) authorizationHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			currentRoute := mux.CurrentRoute(r)
			currentRouteVars := mux.Vars(r)

			if _, ok := excludedRoutes[currentRoute.GetName()]; ok {
				next.ServeHTTP(w, r)
				return
			}

			if authHeader != "" {
				token := common.GetTokenFromHeader(authHeader)
				redisClient, err := common.NewRedisClient()
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				if err := redisClient.IsBlacklisted(token); err != nil {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}

				isAuthorized := a.source.IsAuthorized(token, currentRouteVars)
				if isAuthorized {
					next.ServeHTTP(w, r)
					return
				}

				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return

			}
			http.Error(w, "Missing Token", http.StatusUnauthorized)
		})
	}
}
