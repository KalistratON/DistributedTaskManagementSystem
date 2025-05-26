package main

import (
	"context"
	helper "dtms/pkg/database"
	"dtms/pkg/domain"
	repository "dtms/pkg/repository/task"
	repository_task_history "dtms/pkg/repository/task_history"
	usecases "dtms/pkg/usecases/task"
	pb "dtms/specs/go/pkg"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	mngCln, err := helper.ConnectMongo()
	if err != nil {
		log.Fatalf("task service down with fatal error related to mongo connection= %v", err)
	}
	defer mngCln.Disconnect(context.Background())

	rdsCln, err := helper.ConnectRedis()
	if err != nil {
		log.Fatalf("task service down with fatal error related to redis connection = %v", err)
	}
	defer rdsCln.Close()

	pdb, err := helper.ConnectPsql()
	if err != nil {
		log.Fatalf("task service down with fatal error related to postgres connection = %v", err)
	}
	defer pdb.Close()

	repoTask, err := repository.NewTaskRepostory(mngCln.Database("task").Collection("list"), context.Background())
	if err != nil {
		log.Fatalf("can't create repo: %v", err)
	}

	repoTaskHisory, err := repository_task_history.NewTaskHisoryRepository(pdb)
	if err != nil {
		log.Fatalf("can't create repo: %v", err)
	}

	uc, err := usecases.NewTaskUseCases(repoTask, repoTaskHisory)
	if err != nil {
		log.Fatalf("can't create usecases: %v", err)
	}

	taskServer := NewTaskServer(uc)
	if taskServer == nil {
		log.Printf("can't create server")
	}

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer listener.Close()

	grpcServer := grpc.NewServer()
	pb.RegisterTaskServiceServer(grpcServer, taskServer)

	fmt.Println("Starting gRPC server on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

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
	input := make(chan domain.Task)
	errChannel := make(chan error)
	go func() {
		result, err := s.uc.Create(convert(rqs))
		if err != nil {
			errChannel <- err
		} else {
			input <- *result
		}
	}()
	defer close(input)
	defer close(errChannel)

	select {
	case result := <-input:
		return revertConvert(&result), nil
	case err := <-errChannel:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (s *TaskServer) UpdateTask(ctx context.Context, rqs *pb.TaskMessage) (*pb.TaskMessage, error) {
	input := make(chan domain.Task)
	errChannel := make(chan error)
	go func() {
		result, err := s.uc.Update(convert(rqs))
		if err != nil {
			errChannel <- err
		} else {
			input <- *result
		}
	}()
	defer close(input)
	defer close(errChannel)

	select {
	case result := <-input:
		return revertConvert(&result), nil
	case err := <-errChannel:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (s *TaskServer) GetTask(ctx context.Context, rqs *pb.TaskMessage) (*pb.TaskMessage, error) {
	input := make(chan domain.Task)
	errChannel := make(chan error)
	go func() {
		result, err := s.uc.Get(rqs.Id)
		if err != nil {
			errChannel <- err
		} else {
			input <- *result
		}
	}()
	defer close(input)
	defer close(errChannel)

	select {
	case result := <-input:
		return revertConvert(&result), nil
	case err := <-errChannel:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (s *TaskServer) DeleteTask(ctx context.Context, rqs *pb.TaskMessage) (*pb.TaskMessage, error) {
	input := make(chan domain.Task)
	errChannel := make(chan error)
	go func() {
		result, err := s.uc.Get(rqs.Id)
		if err != nil {
			errChannel <- err
		} else {
			input <- *result
		}
	}()
	defer close(input)
	defer close(errChannel)

	select {
	case result := <-input:
		return revertConvert(&result), nil
	case err := <-errChannel:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
