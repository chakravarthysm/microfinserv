package data

import (
	"database/sql"
	"errors"
	"log"
	"strings"
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
	UserId   string   `json:"customer_id"`
	Accounts []string `json:"accounts"`
	Username string   `json:"username"`
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
	Username string         `db:"username"`
	UserId   sql.NullString `db:"user_id"`
	Accounts sql.NullString `db:"account_numbers"`
}

func (l Login) ClaimsForAccessToken() AccessTokenClaims {
	accounts := strings.Split(l.Accounts.String, ",")
	return AccessTokenClaims{
		UserId:   l.UserId.String,
		Accounts: accounts,
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
	findAccountQuery := "SELECT username, user_id, account_numbers FROM USERS where userame=$1 AND password=$2"

	rows, err := a.client.Query(findAccountQuery)
	if err != nil {
		return &Login{}, err
	}

	var lgn Login
	for rows.Next() {
		err := rows.Scan(&lgn)
		if err != nil {
			return &Login{}, err
		}
	}

	return &lgn, nil
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
	token, err := jwt.ParseWithClaims(tokenString, AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(HMAC_SAMPLE_SECRET), nil
	})
	if err != nil {
		log.Println("Error while parsing token: " + err.Error())
		return nil, err
	}
	return token, nil
}

func NewAuthDB() AuthDBImpl {
	connStr := "user=postgres dbname=finserv password=root host=localhost sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	return &AuthDB{db}
}
