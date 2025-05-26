package usecases

import (
	"dtms/pkg/domain"
	"dtms/pkg/errors"
	repository_task "dtms/pkg/repository/task"
	repository_task_history "dtms/pkg/repository/task_history"
	"log"
)

type taskUseCases struct {
	taskRepo       repository_task.TaskRepository
	taskStatusRepo repository_task_history.TaskHistoryRepository
}

func NewTaskUseCases(
	taskRepo repository_task.TaskRepository, taskStatusRepo repository_task_history.TaskHistoryRepository) (TaskUseCases, error) {
	if taskRepo == nil || taskStatusRepo == nil {
		log.Printf("repo is nil")
		return nil, errors.NilRepo{}
	}
	return &taskUseCases{
		taskRepo:       taskRepo,
		taskStatusRepo: taskStatusRepo,
	}, nil
}

var _ TaskUseCases = &taskUseCases{}

func (us *taskUseCases) Create(task *domain.Task) (*domain.Task, error) {
	result, err := us.taskRepo.Create(task)
	if err != nil {
		return nil, err
	}

	taskHistory := domain.TaskHistory{
		Common: *task,
		Action: string(domain.CREATE),
	}
	if err = us.taskStatusRepo.Create(&taskHistory); err != nil {
		return nil, err
	}
	return result, nil
}

func (us *taskUseCases) Update(task *domain.Task) (*domain.Task, error) {
	result, err := us.taskRepo.Update(task)
	if err != nil {
		return nil, err
	}

	taskHistory := domain.TaskHistory{
		Common: *task,
		Action: string(domain.UPDATE),
	}
	if err = us.taskStatusRepo.Create(&taskHistory); err != nil {
		return nil, err
	}
	return result, nil
}

func (us *taskUseCases) Get(id string) (*domain.Task, error) {
	return us.taskRepo.Get(id)
}

func (us *taskUseCases) GetAll(filter map[string]interface{}) (*[]domain.Task, error) {
	return us.taskRepo.GetAll(filter)
}

func (us *taskUseCases) Delete(task *domain.Task) (*domain.Task, error) {
	result, err := us.taskRepo.Delete(task)
	if err != nil {
		return nil, err
	}

	taskHistory := domain.TaskHistory{
		Common: *task,
		Action: string(domain.DELETE),
	}
	if err = us.taskStatusRepo.Create(&taskHistory); err != nil {
		return nil, err
	}
	return result, nil
}
