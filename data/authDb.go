package data

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/lib/pq"
)

type AuthDBImpl interface {
	FindByUsername(username, password string) (*Login, error)
	IsAuthorized(token string) bool
}

type AuthDB struct {
	client *sql.DB
}

const HMAC_SAMPLE_SECRET = "hmacSampleSecret"
const ACCESS_TOKEN_DURATION = time.Hour

type AccessTokenClaims struct {
	UserId   int32   `json:"user_id"`
	Accounts []int32 `json:"accounts"`
	Username string  `json:"username"`
	jwt.StandardClaims
}

type AuthToken struct {
	token *jwt.Token
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type Login struct {
	UserId   sql.NullInt32   `db:"user_id"`
	Username string          `db:"username"`
	Accounts []sql.NullInt32 `db:"account_numbers"`
}

type UserMeta struct {
	UserId   string `db:"user_id"`
	Username string `db:"username"`
}

func (l Login) ClaimsForAccessToken() AccessTokenClaims {
	// accounts := []int32{}
	// for _, accnt := range l.Accounts {
	// 	accounts = append(accounts, accnt.Int32)
	// }
	return AccessTokenClaims{
		UserId: l.UserId.Int32,
		//Accounts: accounts,
		Username: l.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ACCESS_TOKEN_DURATION).Unix(),
		},
	}
}

func NewAuthToken(claims AccessTokenClaims) AuthToken {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return AuthToken{token: token}
}

func (t AuthToken) NewAccessToken() (string, error) {
	signedString, err := t.token.SignedString([]byte(HMAC_SAMPLE_SECRET))
	if err != nil {
		log.Println("Failed while signing access token: " + err.Error())
		return "", errors.New("cannot generate access token")
	}
	return signedString, nil
}

func (a *AuthDB) FindByUsername(username, password string) (*Login, error) {
	findUserQuery := "SELECT user_id, username FROM users where username=$1 AND password=$2"

	row := a.client.QueryRow(findUserQuery, username, password)
	var id int32
	var usrname string
	err := row.Scan(&id, &usrname)
	if err != nil {
		return &Login{}, err
	}

	return &Login{Username: usrname, UserId: sql.NullInt32{Int32: id, Valid: true}}, nil
}

func (a AuthDB) IsAuthorized(token string) bool {
	if jwtToken, err := jwtTokenFromString(token); err != nil {
		return false
	} else {
		if jwtToken.Valid {
			return true
		} else {
			return false
		}
	}
}

func jwtTokenFromString(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(HMAC_SAMPLE_SECRET), nil
	})
	if err != nil {
		log.Println("Error while parsing token: " + err.Error())
		return nil, err
	}
	return token, nil
}

func NewAuthDB(client *sql.DB) AuthDBImpl {
	return &AuthDB{client}
}
