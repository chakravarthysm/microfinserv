package server

import (
	"bytes"
	"finserv/data"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

var router *mux.Router

type MockUserService struct {
}

func (a *MockUserService) CreateUser(usr data.User) (data.UserResponse, error) {
	return data.UserResponse{
		Username:      usr.Username,
		Name:          usr.Name,
		Location:      usr.Location,
		PAN:           usr.PAN,
		Address:       usr.Address,
		ContactNumber: usr.ContactNumber,
		Gender:        usr.Gender,
		Nationality:   usr.Nationality,
	}, nil
}

func (a *MockUserService) GetUser(id string) (data.UserResponse, error) {
	return data.UserResponse{}, nil
}

func (a *MockUserService) UpdateUser(id string, usr data.User) (data.UserResponse, error) {
	return data.UserResponse{}, nil
}

func (a *MockUserService) DeleteUser(id string) error {
	return nil
}

func mockUserService(t *testing.T) {
	mockUserService := MockUserService{}
	uh := UserHandler{&mockUserService}
	router = mux.NewRouter()
	router.HandleFunc("/users", uh.createUser).Methods(http.MethodPost)
	router.HandleFunc("/users/{user_id}", uh.deleteUser).Methods(http.MethodDelete)
}

var createUserPaylod = `{
	"username": "chakrasm",
	"password": "pass@456",
	"name":"Mithun",
	"location":"DVG",
	"pan": "ATIJJ1962I",
	"address":"#9, 7th Cross, Vijayanagar",
	"contact_number": 9035968097,
	"gender": "Male",
	"nationality": "Indian"
}`

func Test_create_user(t *testing.T) {
	mockUserService(t)

	request, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer([]byte(createUserPaylod)))
	assert.NoError(t, err)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusCreated, recorder.Code, "Expected status code not recieved")
}

func Test_delete_user(t *testing.T) {
	mockUserService(t)

	request, err := http.NewRequest(http.MethodDelete, "/users/dummy", nil)
	assert.NoError(t, err)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusNoContent, recorder.Code, "Expected status code not recieved")
}
