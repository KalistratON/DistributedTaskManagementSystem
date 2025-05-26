package usecases

import (
	"dtms/pkg/domain"
	"dtms/pkg/errors"
	repository_user "dtms/pkg/repository/user"
	"log"
)

type userUseCases struct {
	userRepo repository_user.UserRepository
}

func NewUserUseCases(
	userRepo repository_user.UserRepository) (UserUseCases, error) {
	if userRepo == nil {
		log.Printf("repo is nil")
		return nil, errors.NilRepo{}
	}
	return &userUseCases{
		userRepo: userRepo,
	}, nil
}

var _ UserUseCases = &userUseCases{}

func (us *userUseCases) Create(user *domain.User) (*domain.User, error) {
	return us.userRepo.Create(user)
}

func (us *userUseCases) Update(user *domain.User) (*domain.User, error) {
	return us.userRepo.Update(user)
}

func (us *userUseCases) Get(id string) (*domain.User, error) {
	return us.userRepo.Get(id)
}

func (us *userUseCases) GetAll(filter map[string]interface{}) (*[]domain.User, error) {
	return us.userRepo.GetAll(filter)
}

func (us *userUseCases) Delete(id string) (*domain.User, error) {
	return us.userRepo.Delete(id)
}
