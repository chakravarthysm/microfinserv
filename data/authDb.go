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
	IsAuthorized(token string, routeVars map[string]string) bool
}

type AuthDB struct {
	client *sql.DB
}

const HMACSecret = "hmac_secret"
const ACCESS_TOKEN_DURATION = time.Hour

type AccessTokenClaims struct {
	UserId   string   `json:"user_id"`
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
	UserId   sql.NullString `db:"user_id"`
	Username string         `db:"username"`
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
	signedString, err := t.token.SignedString([]byte(HMACSecret))
	if err != nil {
		log.Println("Failed while signing access token: " + err.Error())
		return "", errors.New("cannot generate access token")
	}
	return signedString, nil
}

func (a *AuthDB) FindByUsername(username, password string) (*Login, error) {
	findUserQuery := `SELECT u.user_id, u.username, array_agg(a.account_id) as account_numbers FROM users u 
	LEFT JOIN accounts a ON a.user_id = u.user_id 
	WHERE username = $1 and password = $2 and u.status = $3
	GROUP BY a.user_id, u.username, u.user_id`

	row := a.client.QueryRow(findUserQuery, username, password, 1)
	var id string
	var usrname string
	var accnts string
	err := row.Scan(&id, &usrname, &accnts)
	if err != nil {
		if err == sql.ErrNoRows {
			return &Login{}, errors.New("User not found")
		}
		return &Login{}, err
	}

	return &Login{
		Username: usrname,
		UserId:   sql.NullString{String: id, Valid: true},
		Accounts: sql.NullString{String: strings.TrimFunc(accnts, func(r rune) bool {
			return (r == '{' || r == '}')
		}), Valid: true},
	}, nil
}

func (a AuthDB) IsAuthorized(token string, routeVars map[string]string) bool {
	jwtToken, err := jwtTokenFromString(token)
	if err != nil {
		return false
	}

	if jwtToken.Valid {
		claims := jwtToken.Claims.(*AccessTokenClaims)
		return claims.VerifyCliams(routeVars)
	}

	return false
}

func (c AccessTokenClaims) VerifyCliams(routeVars map[string]string) bool {
	if c.UserId != routeVars["user_id"] {
		return false
	}

	if routeVars["account_id"] != "" {
		accountFound := false
		for _, a := range c.Accounts {
			if a == routeVars["account_id"] {
				accountFound = true
				break
			}
		}
		return accountFound
	}
	return true
}

func jwtTokenFromString(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(HMACSecret), nil
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
