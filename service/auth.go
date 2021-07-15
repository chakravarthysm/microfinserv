package service

import (
	"finserv/data"
)

type AuthService interface {
	Login(data.LoginRequest) (*data.LoginResponse, *error)
	Verify(urlParams map[string]string) *error
}

type DefaultAuthService struct {
	source data.AuthDBImpl
}

func (s DefaultAuthService) Login(req data.LoginRequest) (*data.LoginResponse, error) {
	var err error
	var login *data.Login

	if login, err = s.source.FindByUsername(req.Username, req.Password); err != nil {
		return nil, err
	}

	claims := login.ClaimsForAccessToken()
	authToken := data.NewAuthToken(claims)

	var accessToken string
	if accessToken, err = authToken.NewAccessToken(); err != nil {
		return nil, err
	}

	return &data.LoginResponse{AccessToken: accessToken}, nil

}

func NewAuthService(authDb data.AuthDBImpl) DefaultAuthService {
	return DefaultAuthService{authDb}
}
