package internal

import (
	"context"
	pb "dtms/specs/go/pkg"

	"dtms/internal/domain"
	usecases "dtms/internal/usecases/task"
)

func convert(task *pb.TaskMessage) *domain.Task {
	if task == nil {
		return nil
	}

	return &domain.Task{
		Id:          task.Id,
		AuthorId:    task.AuthorId,
		Name:        task.Name,
		Description: task.Description,
		Deadline:    task.Deadline,
		Status:      task.Status,
	}
}

func revertConvert(task *domain.Task) *pb.TaskMessage {
	if task == nil {
		return nil
	}

	return &pb.TaskMessage{
		Id:          task.Id,
		AuthorId:    task.AuthorId,
		Name:        task.Name,
		Description: task.Description,
		Deadline:    task.Deadline,
		Status:      task.Status,
	}
}

type TaskServer struct {
	uc usecases.TaskUseCases

	pb.UnimplementedTaskServiceServer
}

var _ pb.TaskServiceServer = &TaskServer{}

func NewTaskServer(uc usecases.TaskUseCases) *TaskServer {
	return &TaskServer{
		uc: uc,
	}
}

func (s *TaskServer) CreateTask(ctx context.Context, rqs *pb.TaskMessage) (*pb.TaskMessage, error) {
	result, err := s.uc.Create(convert(rqs))
	//ToDo: event to kafka
	return revertConvert(result), err

}

func (s *TaskServer) UpdateTask(ctx context.Context, rqs *pb.TaskMessage) (*pb.TaskMessage, error) {
	result, err := s.uc.Update(convert(rqs))
	//ToDo: event to kafka
	return revertConvert(result), err
}

func (s *TaskServer) GetTask(ctx context.Context, rqs *pb.TaskMessage) (*pb.TaskMessage, error) {
	result, err := s.uc.Get(rqs.Id)
	return revertConvert(result), err
}

func (s *TaskServer) DeleteTask(ctx context.Context, rqs *pb.TaskMessage) (*pb.TaskMessage, error) {
	result, err := s.uc.Delete(convert(rqs))
	//ToDo: event to kafka
	return revertConvert(result), err
}
