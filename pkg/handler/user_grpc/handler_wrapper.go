package user_grpc

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	gen "newsfeed/proto/gen/proto"
)

type userGrpcHandlerWrapper struct {
	conf UserGrpcConfig

	server *grpc.Server // grpc server (listen req, connections, ...)

	handler *grpcHandler // handle biz logic API
}

type UserGrpcConfig struct {
	Addr string
}

func NewUserGrpc(conf UserGrpcConfig, userSrv UserService) (*userGrpcHandlerWrapper, error) {
	s := grpc.NewServer()
	userHandler := &grpcHandler{
		userSrv: userSrv,
	}

	gen.RegisterUserServiceServer(s, userHandler)

	h := &userGrpcHandlerWrapper{
		conf:    conf,
		server:  s,
		handler: userHandler,
	}
	return h, nil
}

func (h *userGrpcHandlerWrapper) Start() error {
	// create listener: addr
	lis, err := net.Listen("tcp", fmt.Sprintf("%s", h.conf.Addr))
	if err != nil {
		log.Println("failed to listen: %v", err)
		return fmt.Errorf("failed to listen user GRPC server on %s: %s", h.conf.Addr, err)
	}

	log.Println("user grpc server listening at ", lis.Addr())
	if err := h.server.Serve(lis); err != nil { // block
		log.Fatalf("failed to serve: %v", err)
	}

	return nil
}
