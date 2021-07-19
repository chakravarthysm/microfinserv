package data

import (
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
)

type UsersDBImpl interface {
	Insert(User) (UserResponse, error)
	FindById(string) (UserResponse, error)
	UpdateById(string, User) (UserResponse, error)
	DeleteById(string) error
}

type User struct {
	UserId        string `json:"user_id" db:"user_id"`
	Username      string `json:"username" db:"username"`
	Password      string `json:"password" db:"password"`
	Name          string `json:"name" db:"name"`
	Location      string `json:"location" db:"location"`
	PAN           string `json:"pan" db:"pan"`
	Address       string `json:"address" db:"address"`
	ContactNumber int    `json:"contact_number" db:"contact_number"`
	Gender        string `json:"gender" db:"gender"`
	Nationality   string `json:"nationality" db:"nationality"`
}

type UserResponse struct {
	UserId        string `json:"user_id" db:"user_id"`
	Username      string `json:"username" db:"username"`
	Name          string `json:"name" db:"name"`
	Location      string `json:"location" db:"location"`
	PAN           string `json:"pan" db:"pan"`
	Address       string `json:"address" db:"address"`
	ContactNumber int    `json:"contact_number" db:"contact_number"`
	Gender        string `json:"gender" db:"gender"`
	Nationality   string `json:"nationality" db:"nationality"`
}

type UsersDB struct {
	client *sql.DB
}

func (a *UsersDB) Insert(usr User) (UserResponse, error) {
	insertUserQuery := `INSERT INTO users ("user_id", "username", "password", "name","location", "pan", "address", "contact_number", "gender", "nationality")
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING *`

	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(usr.Password), bcrypt.DefaultCost)
	// if err != nil {
	// 	return UserResponse{}, err
	// }

	row := a.client.QueryRow(insertUserQuery, uuid.NewV4(), usr.Username, usr.Password, usr.Name, usr.Location, usr.PAN, usr.Address, usr.ContactNumber, usr.Gender, usr.Nationality)
	var u UserResponse
	var password string
	var status int
	err := row.Scan(
		&u.UserId,
		&u.Username,
		&password,
		&u.Name,
		&u.Location,
		&u.PAN,
		&u.Address,
		&u.Gender,
		&u.Nationality,
		&u.ContactNumber,
		&status,
	)

	if err != nil {
		return UserResponse{}, err
	}

	return u, nil
}

func (a *UsersDB) FindById(id string) (UserResponse, error) {
	findUserQuery := "SELECT * FROM users WHERE user_id=$1 AND status=$2"
	var u UserResponse
	row := a.client.QueryRow(findUserQuery, id, 1)
	var password string
	var status int
	err := row.Scan(
		&u.UserId,
		&u.Username,
		&password,
		&u.Name,
		&u.Location,
		&u.PAN,
		&u.Address,
		&u.Gender,
		&u.Nationality,
		&u.ContactNumber,
		&status,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return UserResponse{}, errors.New("User not found")
		}
		return UserResponse{}, err
	}
	return u, nil
}

func (a *UsersDB) UpdateById(id string, usr User) (UserResponse, error) {
	var password string
	if usr.Password == "" {
		getPasswordQuery := `SELECT password 
	FROM users
	WHERE user_id = $1 AND status = $2`

		err := a.client.QueryRow(getPasswordQuery, id, 1).Scan(&password)
		if err != nil {
			return UserResponse{}, err
		}
	}

	updateUserQuery := `UPDATE users SET 
	username = $1,
	password = $2,
	name = $3,
	location = $4,
	pan = $5,
	address = $6,
	contact_number = $7,
	gender = $8,
	nationality = $9
	WHERE user_id = $10 AND status = $11
	RETURNING *`

	row := a.client.QueryRow(updateUserQuery, usr.Username, password, usr.Name, usr.Location, usr.PAN, usr.Address, usr.ContactNumber, usr.Gender, usr.Nationality, id, 1)

	var u UserResponse
	var status int
	err := row.Scan(
		&u.UserId,
		&u.Username,
		&password,
		&u.Name,
		&u.Location,
		&u.PAN,
		&u.Address,
		&u.Gender,
		&u.Nationality,
		&u.ContactNumber,
		&status,
	)

	if err != nil {
		return UserResponse{}, err
	}

	return u, nil
}

func (a *UsersDB) DeleteById(id string) error {
	deleteUsersAuthQuery := `UPDATE users SET
	status = 0
	WHERE user_id=$1`

	_, err := a.client.Query(deleteUsersAuthQuery, id)
	if err != nil {
		return err
	}

	deleteUsersQuery := `UPDATE accounts SET
	status = 0
	WHERE user_id=$1`

	_, err = a.client.Query(deleteUsersQuery, id)
	if err != nil {
		return err
	}
	return nil
}

func NewUsersDB(client *sql.DB) UsersDBImpl {
	return &UsersDB{client}
}
