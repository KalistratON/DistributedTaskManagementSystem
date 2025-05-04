package repository

import "dtms/internal/domain"

type TaskHistoryRepository interface {
	Create(taskHisory *domain.TaskHistory) error
}
