package domain

type AccountsRepository interface {
	GetAccount() Account
}

type Account struct {
	Name          string
	Location      string
	PAN           string
	Address       string
	ContactNumber int
	Gender        string
	Nationality   string
}

type AccountStub struct {
}

func (a *AccountStub) GetAccount() Account {
	account := Account{"Chakra", "Bangalore", "AFGYT43519K", "#256, 2nd Cross, Some Extension", 8867508500, "Male", "Indian"}
	return account
}

func NewAccountRepository() AccountsRepository {
	return &AccountStub{}
}
