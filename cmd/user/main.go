package main

import (
	"context"
	helper "dtms/pkg/database"
	"dtms/pkg/domain"
	repository "dtms/pkg/repository/user"
	usecases "dtms/pkg/usecases/user"
	pb "dtms/specs/go/pkg"
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

	userServer := NewUserServer(uc)
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

func convert(user *pb.UserMessage) *domain.User {
	if user == nil {
		return nil
	}

	return &domain.User{
		Id:       user.Id,
		Login:    user.Login,
		Email:    user.Email,
		Password: user.Password,
	}
}

func revertConvert(user *domain.User) *pb.UserMessage {
	if user == nil {
		return nil
	}

	return &pb.UserMessage{
		Id:       user.Id,
		Login:    user.Login,
		Email:    user.Email,
		Password: user.Password,
	}
}

type userServer struct {
	uc usecases.UserUseCases

	pb.UnimplementedUserServiceServer
}

var _ pb.UserServiceServer = &userServer{}

func NewUserServer(uc usecases.UserUseCases) *userServer {
	return &userServer{
		uc: uc,
	}
}

func (s *userServer) CreateUser(ctx context.Context, rqs *pb.UserMessage) (*pb.UserMessage, error) {
	input := make(chan domain.User)
	errChannel := make(chan error)
	go func() {
		result, err := s.uc.Create(convert(rqs))
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

func (s *userServer) UpdateUser(ctx context.Context, rqs *pb.UserMessage) (*pb.UserMessage, error) {
	input := make(chan domain.User)
	errChannel := make(chan error)
	go func() {
		result, err := s.uc.Update(convert(rqs))
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

func (s *userServer) GetUser(ctx context.Context, rqs *pb.UserMessage) (*pb.UserMessage, error) {
	input := make(chan domain.User)
	errChannel := make(chan error)
	go func() {
		result, err := s.uc.Get(rqs.Id)
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

func (s *userServer) DeleteUser(ctx context.Context, rqs *pb.UserMessage) (*pb.UserMessage, error) {
	input := make(chan domain.User)
	errChannel := make(chan error)
	go func() {
		result, err := s.uc.Get(rqs.Id)
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
