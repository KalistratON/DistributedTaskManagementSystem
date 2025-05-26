package repository

import "dtms/pkg/domain"

type TaskHistoryRepository interface {
	Create(taskHisory *domain.TaskHistory) error
}
