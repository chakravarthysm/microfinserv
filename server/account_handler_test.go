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

type MockAccountService struct {
}

func (a *MockAccountService) NewAccount(req data.NewAccountRequest) (*data.NewAccountResponse, error) {
	return &data.NewAccountResponse{
		AccountId: "26c84048-6321-4d5a-b932-33dde68af01d",
	}, nil
}

func (a *MockAccountService) MakeTransaction(req data.TransactionRequest) (*data.TransactionResponse, error) {
	return nil, nil
}

func mockAccountService(t *testing.T) {
	mockAccountService := MockAccountService{}
	ah := AccountHandler{&mockAccountService}
	router = mux.NewRouter()
	router.HandleFunc("/users/{user_id}/account", ah.NewAccount).Methods(http.MethodPost)
	router.HandleFunc("/users/{user_id}/account/{account_id}", ah.MakeTransaction).Methods(http.MethodPost)
}

var newAccountPayload = `{
    "user_id": "1c819246-d6dd-4916-b625-066b736cc4c6",
    "amount": 1000
}`

var newTransactionPaylod = `{
    "account_id": "26c84048-6321-4d5a-b932-33dde68af01d",
    "amount": 500,
    "transaction_type": "withdrawal"
}`

func Test_new_account(t *testing.T) {
	mockAccountService(t)

	request, err := http.NewRequest(http.MethodPost, "/users/1c819246-d6dd-4916-b625-066b736cc4c6/account", bytes.NewBuffer([]byte(newAccountPayload)))
	assert.NoError(t, err)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusCreated, recorder.Code, "Expected status code not recieved")
}

func Test_new_transaction(t *testing.T) {
	mockAccountService(t)

	request, err := http.NewRequest(http.MethodPost, "/users/1c819246-d6dd-4916-b625-066b736cc4c6/account/26c84048-6321-4d5a-b932-33dde68af01d", bytes.NewBuffer([]byte(newTransactionPaylod)))
	assert.NoError(t, err)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusOK, recorder.Code, "Expected status code not recieved")
}
