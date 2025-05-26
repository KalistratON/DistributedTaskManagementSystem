package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"dtms/cmd/gateway/types"

	pb "dtms/specs/go/pkg"

	"google.golang.org/grpc/credentials/insecure"

	gorilla "github.com/gorilla/mux"
	"google.golang.org/grpc"
)

func TaskPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(types.TaskResponse{
			TaskId: "-1",
			Status: "error",
			ErrMsg: "method not allowed",
		})
		return
	}

	var rqs types.TaskMessage
	if err := json.NewDecoder(r.Body).Decode(&rqs); err != nil {
		log.Println("error while trying to decode request")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.TaskResponse{
			TaskId: "-1",
			Status: "error",
			ErrMsg: "error while trying to decode request",
		})
		return
	}
	defer r.Body.Close()

	taskUrl := os.Getenv("TASK_SERVICE_URL")
	if taskUrl == "" {
		log.Println("TASK_SERVICE_URL is empty")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.TaskResponse{
			TaskId: "-1",
			Status: "error",
			ErrMsg: "task is not allowed now",
		})
		return
	}

	conn, err := grpc.NewClient(taskUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to gRPC server (%s): %s", taskUrl, err.Error())
	}
	defer conn.Close()

	client := pb.NewTaskServiceClient(conn)
	reqeustGrpc := &pb.TaskMessage{
		Id:          "-1",
		AuthorId:    rqs.AuthorId,
		Name:        rqs.Name,
		Description: rqs.Description,
		Deadline:    rqs.Deadline,
		Status:      rqs.Status,
	}

	responseGrpc, err := client.CreateTask(context.Background(), reqeustGrpc)
	if responseGrpc == nil {
		log.Printf("error while grpc response: %v\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.TaskResponse{
			TaskId: "-1",
			Status: "error",
			ErrMsg: "task is not allowed now",
		})
		return
	} else if err != nil {
		log.Printf("error while grpc response: %v\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.TaskResponse{
			TaskId: "-1",
			Status: "error",
			ErrMsg: err.Error(),
		})
		return
	}

	if responseGrpc.Id != "-1" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(types.TaskMessage{
		Id:          responseGrpc.Id,
		AuthorId:    responseGrpc.AuthorId,
		Name:        responseGrpc.Name,
		Description: responseGrpc.Description,
		Deadline:    responseGrpc.Deadline,
		Status:      responseGrpc.Status,
	})
}

func TaskPut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(types.TaskResponse{
			TaskId: "-1",
			Status: "error",
			ErrMsg: "method not allowed",
		})
		return
	}

	var rqs types.TaskMessage
	if err := json.NewDecoder(r.Body).Decode(&rqs); err != nil {
		log.Println("error while trying to decode request")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.TaskResponse{
			TaskId: "-1",
			Status: "error",
			ErrMsg: "error while trying to decode request",
		})
		return
	}
	defer r.Body.Close()

	taskUrl := os.Getenv("TASK_SERVICE_URL")
	if taskUrl == "" {
		log.Println("TASK_SERVICE_URL is empty")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.TaskResponse{
			TaskId: "-1",
			Status: "error",
			ErrMsg: "task is not allowed now",
		})
		return
	}

	conn, err := grpc.NewClient(taskUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to gRPC server (%s): %s", taskUrl, err.Error())
	}
	defer conn.Close()

	client := pb.NewTaskServiceClient(conn)
	reqeustGrpc := &pb.TaskMessage{
		Id:          gorilla.Vars(r)["task_id"],
		AuthorId:    rqs.AuthorId,
		Name:        rqs.Name,
		Description: rqs.Description,
		Deadline:    rqs.Deadline,
		Status:      rqs.Status,
	}

	responseGrpc, err := client.UpdateTask(context.Background(), reqeustGrpc)
	if responseGrpc == nil {
		log.Printf("error while grpc response: %v\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.TaskResponse{
			TaskId: "-1",
			Status: "error",
			ErrMsg: "task is not allowed now",
		})
		return
	} else if err != nil {
		log.Printf("error while grpc response: %v\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.TaskResponse{
			TaskId: "-1",
			Status: "error",
			ErrMsg: err.Error(),
		})
		return
	}

	if responseGrpc.Id != "-1" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(types.TaskMessage{
		Id:          responseGrpc.Id,
		AuthorId:    responseGrpc.AuthorId,
		Name:        responseGrpc.Name,
		Description: responseGrpc.Description,
		Deadline:    responseGrpc.Deadline,
		Status:      responseGrpc.Status,
	})
}

func TaskGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(types.TaskResponse{
			TaskId: "-1",
			Status: "error",
			ErrMsg: "method not allowed",
		})
		return
	}

	taskUrl := os.Getenv("TASK_SERVICE_URL")
	if taskUrl == "" {
		log.Println("TASK_SERVICE_URL is empty")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.TaskResponse{
			TaskId: "-1",
			Status: "error",
			ErrMsg: "task is not allowed now",
		})
		return
	}

	var rqs types.TaskMessage
	if err := json.NewDecoder(r.Body).Decode(&rqs); err != nil {
		log.Println("error while trying to decode request")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.TaskResponse{
			TaskId: "-1",
			Status: "error",
			ErrMsg: "error while trying to decode request",
		})
		return
	}
	defer r.Body.Close()

	conn, err := grpc.NewClient(taskUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to gRPC server (%s): %s", taskUrl, err.Error())
	}
	defer conn.Close()

	client := pb.NewTaskServiceClient(conn)
	reqeustGrpc := &pb.TaskMessage{
		Id: gorilla.Vars(r)["task_id"],
	}

	responseGrpc, err := client.GetTask(context.Background(), reqeustGrpc)
	if responseGrpc == nil {
		log.Printf("error while grpc response: %v\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.TaskResponse{
			TaskId: "-1",
			Status: "error",
			ErrMsg: "task is not allowed now",
		})
		return
	} else if err != nil {
		log.Printf("error while grpc response: %v\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.TaskResponse{Status: "error", ErrMsg: err.Error()})
		return
	}

	if responseGrpc.Id != "-1" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(types.TaskMessage{
		Id:          responseGrpc.Id,
		AuthorId:    responseGrpc.AuthorId,
		Name:        responseGrpc.Name,
		Description: responseGrpc.Description,
		Deadline:    responseGrpc.Deadline,
		Status:      responseGrpc.Status,
	})
}

func TaskDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(types.TaskResponse{
			TaskId: "-1",
			Status: "error",
			ErrMsg: "method not allowed",
		})
		return
	}

	taskUrl := os.Getenv("TASK_SERVICE_URL")
	if taskUrl == "" {
		log.Println("TASK_SERVICE_URL is empty")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.TaskResponse{
			TaskId: "-1",
			Status: "error",
			ErrMsg: "task is not allowed now",
		})
		return
	}

	var rqs types.TaskMessage
	if err := json.NewDecoder(r.Body).Decode(&rqs); err != nil {
		log.Println("error while trying to decode request")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.TaskResponse{
			TaskId: "-1",
			Status: "error",
			ErrMsg: "error while trying to decode request",
		})
		return
	}
	defer r.Body.Close()

	conn, err := grpc.NewClient(taskUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to gRPC server (%s): %s", taskUrl, err.Error())
	}
	defer conn.Close()

	client := pb.NewTaskServiceClient(conn)
	reqeustGrpc := &pb.TaskMessage{
		Id:       gorilla.Vars(r)["task_id"],
		AuthorId: rqs.AuthorId,
	}

	responseGrpc, err := client.DeleteTask(context.Background(), reqeustGrpc)
	if responseGrpc == nil {
		log.Printf("error while grpc response: %v\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.TaskResponse{
			TaskId: "-1",
			Status: "error",
			ErrMsg: "task is not allowed now",
		})
		return
	} else if err != nil {
		log.Printf("error while grpc response: %v\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.TaskResponse{Status: "error", ErrMsg: err.Error()})
		return
	}

	if responseGrpc.Id != "-1" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(types.TaskMessage{
		Id:          responseGrpc.Id,
		AuthorId:    responseGrpc.AuthorId,
		Name:        responseGrpc.Name,
		Description: responseGrpc.Description,
		Deadline:    responseGrpc.Deadline,
		Status:      responseGrpc.Status,
	})
}
