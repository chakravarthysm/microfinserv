package service

import "finserv/data"

type AuthService interface {
}

type AuthHandler struct {
	AuthDB data.AuthDB
}

func NewAuthService(auth data.AuthDB) AuthService {
	return &AccountHandler{
		AccountDB: auth,
	}
}
