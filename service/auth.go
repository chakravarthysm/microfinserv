package service

import (
	"errors"
	"finserv/common"
	"finserv/data"
)

type AuthService struct {
	source data.AuthDBImpl
}

func (s AuthService) Login(req data.LoginRequest) (*data.LoginResponse, error) {
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

func (s AuthService) Logout(jwt string) error {
	var err error = errors.New("Error occured when loggin out")
	redisClient, err := common.NewRedisClient()
	if err != nil {
		return err
	}

	if err = redisClient.AddToBlacklist(jwt); err != nil {
		return err
	}

	return nil
}

func NewAuthService(authDb data.AuthDBImpl) AuthService {
	return AuthService{authDb}
}
