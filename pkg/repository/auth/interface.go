package repository

import "dtms/pkg/domain"

type AuthRepository interface {
	Create(auth *domain.Auth) (*domain.Auth, error)
	Get(token string) (*domain.Auth, error)
	GetToken(id string) (*domain.Auth, error)
	Update(token string) (*domain.Auth, error)
	Delete(token string) (*domain.Auth, error)
}
