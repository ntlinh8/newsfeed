package main

import (
	"log"

	"newsfeed/pkg/handler/user_grpc"
	"newsfeed/pkg/repo"
	"newsfeed/pkg/service"
)

func main() {
	// TODO: read from env
	conf := user_grpc.UserGrpcConfig{Addr: ":44444"}
	mysqlConf := repo.MySQLConfig{
		Username:     "root",
		Password:     "123456",
		Addr:         "localhost:33334",
		DatabaseName: "newsfeed",
	}

	userRepo, err := repo.NewUserRepo(mysqlConf)
	if err != nil {
		log.Println("failed to create user repo", err)
		return
	}

	userSrv, err := service.NewUserService(userRepo)
	if err != nil {
		log.Println("failed to create user service", err)
		return
	}

	grpcWrapper, err := user_grpc.NewUserGrpc(conf, userSrv)
	if err != nil {
		log.Println("failed to create user grpc server", err)
		return
	}

	log.Println("starting user grpc server at", conf.Addr)
	err = grpcWrapper.Start()
	if err != nil {
		log.Println("failed to start user grpc server", err)
		return
	}
}

// create user: insert
// update/ delete
// get user: select
//
