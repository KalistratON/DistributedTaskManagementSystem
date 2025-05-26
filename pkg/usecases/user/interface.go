package usecases

import "dtms/pkg/domain"

type UserUseCases interface {
	// Create создаёт и обновлят историю пользователя
	Create(user *domain.User) (*domain.User, error)
	// Update обновляет пользователя и её историю
	Update(user *domain.User) (*domain.User, error)
	// Get Возвращает пользователя
	Get(id string) (*domain.User, error)
	// Get Возвращает все задачи по фильтру
	GetAll(filter map[string]interface{}) (*[]domain.User, error)
	// Delete удаляет пользователя и обновляет её историю
	Delete(id string) (*domain.User, error)
}
