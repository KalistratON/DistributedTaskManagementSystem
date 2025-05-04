package main

import (
	"context"
	helper "dtms/helper/pkg"
	repository "dtms/internal/repository/user"
	usecases "dtms/internal/usecases/user"
	pb "dtms/specs/go/pkg"
	"dtms/user/internal"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	mngCln, err := helper.ConnectMongo()
	if err != nil {
		log.Fatalf("task service down with fatal error related to mongo connection= %v", err)
	}
	defer mngCln.Disconnect(context.Background())

	repoUser, err := repository.NewUserRepository(mngCln.Database("user").Collection("list"), context.Background())
	if err != nil {
		log.Fatalf("can't create repo: %v", err)
	}

	uc, err := usecases.NewUserUseCases(repoUser)
	if err != nil {
		log.Fatalf("can't create usecases: %v", err)
	}

	userServer := internal.NewUserServer(uc)
	if userServer == nil {
		log.Printf("can't create server")
	}

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer listener.Close()

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, userServer)

	fmt.Println("Starting gRPC server on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
