package repository

import "dtms/pkg/domain"

type TaskRepository interface {
	Create(task *domain.Task) (*domain.Task, error)
	Update(task *domain.Task) (*domain.Task, error)
	Get(id string) (*domain.Task, error)
	GetAll(filter map[string]interface{}) (*[]domain.Task, error)
	Delete(task *domain.Task) (*domain.Task, error)
}
