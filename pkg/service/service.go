package service

import (
	"fmt"

	"newsfeed/pkg/service/model"
)

type UserRepo interface {
	CreateUser(user *model.User) (*model.User, error)
}

type userService struct {
	userRepo UserRepo
}

func NewUserService(userRepo UserRepo) (*userService, error) {
	srv := &userService{
		userRepo: userRepo,
	}
	return srv, nil
}

func (usv *userService) CreateUser(user *model.User) (*model.User, error) {
	dbUser, err := usv.userRepo.CreateUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user in db: %s", err)
	}
	return dbUser, nil
}
