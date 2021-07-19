package service

import "finserv/data"

type UserServiceImpl interface {
	CreateUser(data.User) (data.UserResponse, error)
	GetUser(string) (data.UserResponse, error)
	UpdateUser(string, data.User) (data.UserResponse, error)
	DeleteUser(string) error
}

type UserService struct {
	UserDB data.UsersDBImpl
}

func (u *UserService) CreateUser(usr data.User) (data.UserResponse, error) {
	return u.UserDB.Insert(usr)
}

func (u *UserService) GetUser(id string) (data.UserResponse, error) {
	return u.UserDB.FindById(id)
}

func (u *UserService) UpdateUser(id string, usr data.User) (data.UserResponse, error) {
	return u.UserDB.UpdateById(id, usr)
}

func (u *UserService) DeleteUser(id string) error {
	return u.UserDB.DeleteById(id)
}

func NewUsersService(userDb data.UsersDBImpl) UserServiceImpl {
	return &UserService{
		UserDB: userDb,
	}
}
