package repository

import "dtms/internal/domain"

type AuthRepository interface {
	Create(auth *domain.Auth) (*domain.Auth, error)
	Get(token string) (*domain.Auth, error)
	Update(token string) (*domain.Auth, error)
	Delete(token string) (*domain.Auth, error)
}
