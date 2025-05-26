package repository

import "dtms/pkg/domain"

type UserRepository interface {
	Create(User *domain.User) (*domain.User, error)
	Update(User *domain.User) (*domain.User, error)
	Get(id string) (*domain.User, error)
	GetAll(filter map[string]interface{}) (*[]domain.User, error)
	Delete(id string) (*domain.User, error)
}
