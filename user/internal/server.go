package internal

import (
	"context"
	"dtms/internal/domain"
	usecases "dtms/internal/usecases/user"
	pb "dtms/specs/go/pkg"
)

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
	result, err := s.uc.Create(convert(rqs))
	return revertConvert(result), err

}

func (s *userServer) UpdateUser(ctx context.Context, rqs *pb.UserMessage) (*pb.UserMessage, error) {
	result, err := s.uc.Update(convert(rqs))
	return revertConvert(result), err
}

func (s *userServer) GetUser(ctx context.Context, rqs *pb.UserMessage) (*pb.UserMessage, error) {
	result, err := s.uc.Get(rqs.Id)
	return revertConvert(result), err
}

func (s *userServer) DeleteUser(ctx context.Context, rqs *pb.UserMessage) (*pb.UserMessage, error) {
	result, err := s.uc.Delete(rqs.Id)
	return revertConvert(result), err
}
