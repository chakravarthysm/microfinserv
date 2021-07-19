package data

const dateLayout = "2006-01-02 15:04:05"

type NewAccountRequest struct {
	UserId string  `json:"user_id"`
	Amount float64 `json:"amount"`
}

type NewAccountResponse struct {
	AccountId string `json:"account_id"`
}

type Account struct {
	AccountId string  `db:"account_id"`
	UserId    string  `db:"user_id"`
	CreatedOn string  `db:"created_on"`
	Balance   float64 `db:"balance"`
	Status    int     `db:"status"`
}
