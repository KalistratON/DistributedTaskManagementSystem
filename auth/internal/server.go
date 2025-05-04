package internal

import (
	"context"
	pb "dtms/specs/go/pkg"

	"dtms/internal/domain"
	usecases "dtms/internal/usecases/auth"
)

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
	result, err := s.uc.SoftCreate(convert(rqs))
	return revertConvert(result), err
}

func (s *AuthServer) HardCreate(ctx context.Context, rqs *pb.AuthMessage) (*pb.AuthMessage, error) {
	result, err := s.uc.HardCreate(convert(rqs))
	return revertConvert(result), err
}

func (s *AuthServer) Get(ctx context.Context, rqs *pb.AuthMessage) (*pb.AuthMessage, error) {
	result, err := s.uc.Get(rqs.Token)
	return revertConvert(result), err
}

func (s *AuthServer) Extend(ctx context.Context, rqs *pb.AuthMessage) (*pb.AuthMessage, error) {
	result, err := s.uc.Extend(rqs.Token)
	return revertConvert(result), err
}

func (s *AuthServer) Delete(ctx context.Context, rqs *pb.AuthMessage) (*pb.AuthMessage, error) {
	result, err := s.uc.Delete(rqs.Token)
	return revertConvert(result), err
}
