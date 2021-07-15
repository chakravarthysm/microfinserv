package data

import (
	"database/sql"

	"log"

	_ "github.com/lib/pq"
)

type UsersDBImpl interface {
	Insert(User) (UserResponse, error)
	FindById(string) (User, error)
	UpdateById(string, User) (User, error)
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
	insertUserQuery := `INSERT INTO users ("username", "password", "name","location", "pan", "address", "contact_number", "gender", "nationality")
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`

	rows, err := a.client.Query(insertUserQuery, usr.Username, usr.Password, usr.Name, usr.Location, usr.PAN, usr.Address, usr.ContactNumber, usr.Gender, usr.Nationality)
	if err != nil {
		return UserResponse{}, err
	}
	var u UserResponse
	for rows.Next() {
		err := rows.Scan(&u)
		if err != nil {
			return UserResponse{}, err
		}
	}

	return u, nil
}

func (a *UsersDB) FindById(id string) (User, error) {
	findUserQuery := "SELECT * FROM users WHERE user_id=$1 AND status=$2"

	rows, err := a.client.Query(findUserQuery, id, 1)
	if err != nil {
		return User{}, err
	}

	var u User
	for rows.Next() {
		err := rows.Scan(&u)
		if err != nil {
			return User{}, err
		}
	}

	return u, nil
}

func (a *UsersDB) UpdateById(id string, usr User) (User, error) {
	findUserQuery := `INSERT INTO users ("username", "password","name", "location", "pan", "address", "contact_number", "gender", "nationality")
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	WHERE user_id=id 
	ON CONFLICT ("name", "location", "pan", "address", "contact_number", "gender", "nationality")
	DO NOTHING`

	rows, err := a.client.Query(findUserQuery, usr.Username, usr.Password, usr.Name, usr.Location, usr.PAN, usr.Address, usr.ContactNumber, usr.Gender, usr.Nationality)
	if err != nil {
		return User{}, err
	}

	var u User
	for rows.Next() {
		err := rows.Scan(&u)
		if err != nil {
			return User{}, err
		}
	}

	return u, nil
}

func (a *UsersDB) DeleteById(id string) error {
	deleteUsersAuthQuery := `INSERT INTO users (status)
	VALUES (0)
	WHERE user_id=id`

	_, err := a.client.Query(deleteUsersAuthQuery)
	if err != nil {
		return err
	}

	deleteUsersQuery := `INSERT INTO accounts (status)
	VALUES (0)
	WHERE user_id=id`

	_, err = a.client.Query(deleteUsersQuery)
	if err != nil {
		return err
	}
	return nil
}

func NewUsersDB() UsersDBImpl {
	connStr := "user=postgres dbname=finserv password=postgres host=localhost sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	return &UsersDB{db}
}
