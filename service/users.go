package service

import "finserv/data"

type UserService interface {
	CreateUser(data.User) (data.UserResponse, error)
	GetUser(string) (data.User, error)
	UpdateUser(string, data.User) (data.User, error)
	DeleteUser(string) error
}

type UserHandler struct {
	UserDB data.UsersDBImpl
}

func (a *UserHandler) CreateUser(usr data.User) (data.UserResponse, error) {
	return a.UserDB.Insert(usr)
}

func (a *UserHandler) GetUser(id string) (data.User, error) {
	return a.UserDB.FindById(id)
}

func (a *UserHandler) UpdateUser(id string, usr data.User) (data.User, error) {
	return a.UserDB.UpdateById(id, usr)
}

func (a *UserHandler) DeleteUser(id string) error {
	return a.UserDB.DeleteById(id)
}

func NewUsersService(userDb data.UsersDBImpl) *UserHandler {
	return &UserHandler{
		UserDB: userDb,
	}
}
