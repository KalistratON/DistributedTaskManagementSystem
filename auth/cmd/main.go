package main

import (
	"context"
	"dtms/auth/internal"
	helper "dtms/helper/pkg"
	repository "dtms/internal/repository/auth"
	usecases "dtms/internal/usecases/auth"
	pb "dtms/specs/go/pkg"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	rdsCln, err := helper.ConnectRedis()
	if err != nil {
		log.Fatalf("task service down with fatal error related to redis connection = %v", err)
	}
	defer rdsCln.Close()

	repo, err := repository.NewAuthRepository(rdsCln, context.Background())
	if err != nil {
		log.Fatalf("can't create repo: %v", err)
	}

	uc, err := usecases.NewAuthUseCases(repo)
	if err != nil {
		log.Fatalf("can't create usecases: %v", err)
	}

	authServer := internal.NewAuthServer(uc)
	if authServer == nil {
		log.Printf("can't create server")
	}

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer listener.Close()

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, authServer)

	fmt.Println("Starting gRPC server on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
