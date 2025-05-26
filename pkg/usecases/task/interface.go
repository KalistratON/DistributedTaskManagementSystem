package usecases

import "dtms/pkg/domain"

type TaskUseCases interface {
	// Create создаёт и обновлят историю задачи
	Create(task *domain.Task) (*domain.Task, error)
	// Update обновляет задачу и её историю
	Update(task *domain.Task) (*domain.Task, error)
	// Get Возвращает задачу
	Get(id string) (*domain.Task, error)
	// Get Возвращает все задачи по фильтру
	GetAll(filter map[string]interface{}) (*[]domain.Task, error)
	// Delete удаляет задачу и обновляет её историю
	Delete(task *domain.Task) (*domain.Task, error)
}
