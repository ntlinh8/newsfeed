package user_grpc

import (
	"context"
	"fmt"
	"log"

	"newsfeed/pkg/service/model"
	gen "newsfeed/proto/gen/proto"
)

type UserService interface {
	CreateUser(user *model.User) (*model.User, error)
}

type grpcHandler struct {
	gen.UnimplementedUserServiceServer

	userSrv UserService
}

func (s *grpcHandler) CreateUser(ctx context.Context, req *gen.CreateUserRequest) (*gen.CreateUserResponse, error) {
	log.Println("[DEBUG] receipt req", req)

	createdUser, err := s.userSrv.CreateUser(&model.User{
		Username:  req.Username,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		DOB:       int(req.Dob),
		Email:     req.Email,
	})
	if err != nil {
		log.Println("[DEBUG] failed to create user", err)
		return nil, fmt.Errorf("failed to create user: %s", err)
	}

	return &gen.CreateUserResponse{
		ErrCode: "OK",
		Message: "OKOKOKOKO",
		User: &gen.User{
			Id:        int64(createdUser.UserId),
			Username:  createdUser.Username,
			FirstName: createdUser.FirstName,
			LastName:  createdUser.LastName,
			Dob:       int32(createdUser.DOB),
			Email:     createdUser.Email,
		},
	}, nil
}
