package usecases

import (
	"dtms/pkg/domain"
	"dtms/pkg/errors"
	repository_auth "dtms/pkg/repository/auth"
	"log"
)

type authUseCases struct {
	authRepo repository_auth.AuthRepository
}

func NewAuthUseCases(authRepo repository_auth.AuthRepository) (AuthUseCases, error) {
	if authRepo == nil {
		log.Printf("repo is nil")
		return nil, errors.NilRepo{}
	}
	return &authUseCases{
		authRepo: authRepo,
	}, nil
}

var _ AuthUseCases = &authUseCases{}

func (r *authUseCases) SoftCreate(auth *domain.Auth) (*domain.Auth, error) {
	existedAuth, err := r.authRepo.GetToken(auth.Id)
	if err != nil {
		return r.authRepo.Create(auth)
	}

	return r.authRepo.Update(existedAuth.Token)
}

func (r *authUseCases) HardCreate(auth *domain.Auth) (*domain.Auth, error) {
	existedAuth, err := r.authRepo.GetToken(auth.Id)
	if err == nil {
		r.authRepo.Delete(existedAuth.Token)
	}
	return r.authRepo.Create(auth)
}

func (r *authUseCases) Get(token string) (*domain.Auth, error) {
	return r.authRepo.Get(token)
}

func (r *authUseCases) Extend(auth string) (*domain.Auth, error) {
	return r.authRepo.Update(auth)
}

func (r *authUseCases) Delete(token string) (*domain.Auth, error) {
	return r.authRepo.Delete(token)
}
