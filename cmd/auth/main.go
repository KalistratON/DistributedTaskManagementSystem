package main

import (
	"context"
	helper "dtms/pkg/database"
	"dtms/pkg/domain"
	repository "dtms/pkg/repository/auth"
	usecases "dtms/pkg/usecases/auth"
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

	authServer := NewAuthServer(uc)
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

func convert(auth *pb.AuthMessage) *domain.Auth {
	if auth == nil {
		return nil
	}

	return &domain.Auth{
		Id:    auth.Id,
		Token: auth.Token,
	}
}

func revertConvert(auth *domain.Auth) *pb.AuthMessage {
	if auth == nil {
		return nil
	}

	return &pb.AuthMessage{
		Id:    auth.Id,
		Token: auth.Token,
	}
}

type AuthServer struct {
	uc usecases.AuthUseCases

	pb.UnimplementedAuthServiceServer
}

var _ pb.AuthServiceServer = &AuthServer{}

func NewAuthServer(uc usecases.AuthUseCases) *AuthServer {
	return &AuthServer{
		uc: uc,
	}
}

func (s *AuthServer) SoftCreate(ctx context.Context, rqs *pb.AuthMessage) (*pb.AuthMessage, error) {
	input := make(chan domain.Auth)
	errChannel := make(chan error)
	go func() {
		result, err := s.uc.SoftCreate(convert(rqs))
		if err != nil {
			errChannel <- err
		} else {
			input <- *result
		}
	}()
	defer close(input)
	defer close(errChannel)

	select {
	case result := <-input:
		return revertConvert(&result), nil
	case err := <-errChannel:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (s *AuthServer) HardCreate(ctx context.Context, rqs *pb.AuthMessage) (*pb.AuthMessage, error) {
	input := make(chan domain.Auth)
	errChannel := make(chan error)
	go func() {
		result, err := s.uc.HardCreate(convert(rqs))
		if err != nil {
			errChannel <- err
		} else {
			input <- *result
		}
	}()
	defer close(input)
	defer close(errChannel)

	select {
	case result := <-input:
		return revertConvert(&result), nil
	case err := <-errChannel:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (s *AuthServer) Get(ctx context.Context, rqs *pb.AuthMessage) (*pb.AuthMessage, error) {
	input := make(chan domain.Auth)
	errChannel := make(chan error)
	go func() {
		result, err := s.uc.Get(rqs.Token)
		if err != nil {
			errChannel <- err
		} else {
			input <- *result
		}
	}()
	defer close(input)
	defer close(errChannel)

	select {
	case result := <-input:
		return revertConvert(&result), nil
	case err := <-errChannel:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (s *AuthServer) Extend(ctx context.Context, rqs *pb.AuthMessage) (*pb.AuthMessage, error) {
	input := make(chan domain.Auth)
	errChannel := make(chan error)
	go func() {
		result, err := s.uc.Get(rqs.Token)
		if err != nil {
			errChannel <- err
		} else {
			input <- *result
		}
	}()
	defer close(input)
	defer close(errChannel)

	select {
	case result := <-input:
		return revertConvert(&result), nil
	case err := <-errChannel:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (s *AuthServer) Delete(ctx context.Context, rqs *pb.AuthMessage) (*pb.AuthMessage, error) {
	input := make(chan domain.Auth)
	errChannel := make(chan error)
	go func() {
		result, err := s.uc.Get(rqs.Token)
		if err != nil {
			errChannel <- err
		} else {
			input <- *result
		}
	}()
	defer close(input)
	defer close(errChannel)

	select {
	case result := <-input:
		return revertConvert(&result), nil
	case err := <-errChannel:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
