package main

import (
	"context"
	helper "dtms/helper/pkg"
	repository "dtms/internal/repository/task"
	repository_task_history "dtms/internal/repository/task_history"
	usecases "dtms/internal/usecases/task"
	pb "dtms/specs/go/pkg"
	"dtms/task/internal"
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

	taskServer := internal.NewTaskServer(uc)
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
