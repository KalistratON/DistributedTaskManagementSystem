package usecases

import "dtms/pkg/domain"

type AuthUseCases interface {
	// SoftCreate create if not exist, otherwise return existed auth
	SoftCreate(auth *domain.Auth) (*domain.Auth, error)
	// HardCreate create or recreate cash for auth
	HardCreate(auth *domain.Auth) (*domain.Auth, error)
	// Get return cash if exist, otherwise return error
	Get(tokent string) (*domain.Auth, error)
	// Extend return and reset TTL for cash if exist, otherwise return error
	Extend(auth string) (*domain.Auth, error)
	// Delete cash if exist, otherwise return error
	Delete(token string) (*domain.Auth, error)
}
