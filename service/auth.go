package service

import (
	"finserv/data"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type AuthService interface {
	Login(LoginRequest) (*LoginResponse, *error)
	Verify(urlParams map[string]string) *error
}

type DefaultAuthService struct {
	source data.AuthDBImpl
}

func (s DefaultAuthService) Login(req LoginRequest) (*LoginResponse, error) {
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

	return &LoginResponse{AccessToken: accessToken}, nil

}

func NewAuthService(authDb data.AuthDBImpl) DefaultAuthService {
	return DefaultAuthService{authDb}
}
