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

func UserPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(types.UserResponse{
			Id:     "-1",
			Status: "error",
			ErrMsg: "method not allowed",
		})
		return
	}

	var rqs types.UserMessage
	if err := json.NewDecoder(r.Body).Decode(&rqs); err != nil {
		log.Println("error while trying to decode request")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.UserResponse{
			Id:     "-1",
			Status: "error",
			ErrMsg: "error while trying to decode request",
		})
		return
	}
	defer r.Body.Close()

	userUrl := os.Getenv("USER_SERVICE_URL")
	if userUrl == "" {
		log.Println("USER_SERVICE_URL is empty")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.UserResponse{
			Id:     "-1",
			Status: "error",
			ErrMsg: "user is not allowed now",
		})
		return
	}

	conn, err := grpc.NewClient(userUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to gRPC server (%s): %s", userUrl, err.Error())
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)
	reqeustGrpc := &pb.UserMessage{
		Id:       "-1",
		Login:    rqs.Login,
		Email:    rqs.Email,
		Password: rqs.Password,
	}

	responseGrpc, err := client.CreateUser(context.Background(), reqeustGrpc)
	if err != nil {
		log.Printf("error while grpc response: %v\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.UserResponse{
			Id:     "-1",
			Status: "error",
			ErrMsg: err.Error(),
		})
		return
	} else if responseGrpc == nil {
		log.Printf("error while grpc response: %v\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.UserResponse{
			Id:     "-1",
			Status: "error",
			ErrMsg: "user is not allowed now",
		})
		return
	}

	if responseGrpc.Id != "-1" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(types.UserMessage{
		Id:       responseGrpc.Id,
		Login:    responseGrpc.Login,
		Email:    responseGrpc.Email,
		Password: responseGrpc.Password,
	})
}

func UserPut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(types.UserResponse{
			Id:     "-1",
			Status: "error",
			ErrMsg: "method not allowed",
		})
		return
	}

	var rqs types.UserMessage
	if err := json.NewDecoder(r.Body).Decode(&rqs); err != nil {
		log.Println("error while trying to decode request")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.UserResponse{
			Id:     "-1",
			Status: "error",
			ErrMsg: "error while trying to decode request",
		})
		return
	}
	defer r.Body.Close()

	userUrl := os.Getenv("USER_SERVICE_URL")
	if userUrl == "" {
		log.Println("USER_SERVICE_URL is empty")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.UserResponse{
			Id:     "-1",
			Status: "error",
			ErrMsg: "user is not allowed now",
		})
		return
	}

	conn, err := grpc.NewClient(userUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to gRPC server (%s): %s", userUrl, err.Error())
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)
	reqeustGrpc := &pb.UserMessage{
		Id:       gorilla.Vars(r)["user_id"],
		Login:    rqs.Login,
		Email:    rqs.Email,
		Password: rqs.Password,
	}

	responseGrpc, err := client.UpdateUser(context.Background(), reqeustGrpc)
	if responseGrpc == nil {
		log.Printf("error while grpc response: %v\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.UserResponse{
			Id:     "-1",
			Status: "error",
			ErrMsg: "user is not allowed now",
		})
		return
	} else if err != nil {
		log.Printf("error while grpc response: %v\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.UserResponse{
			Id:     "-1",
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

	json.NewEncoder(w).Encode(types.UserMessage{
		Id:       responseGrpc.Id,
		Login:    responseGrpc.Login,
		Email:    responseGrpc.Email,
		Password: responseGrpc.Password,
	})
}

func UserGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(types.UserResponse{
			Id:     "-1",
			Status: "error",
			ErrMsg: "method not allowed",
		})
		return
	}

	userUrl := os.Getenv("USER_SERVICE_URL")
	if userUrl == "" {
		log.Println("USER_SERVICE_URL is empty")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.UserResponse{
			Id:     "-1",
			Status: "error",
			ErrMsg: "user is not allowed now",
		})
		return
	}

	conn, err := grpc.NewClient(userUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to gRPC server (%s): %s", userUrl, err.Error())
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)
	reqeustGrpc := &pb.UserMessage{
		Id: gorilla.Vars(r)["user_id"],
	}

	responseGrpc, err := client.GetUser(context.Background(), reqeustGrpc)
	if responseGrpc == nil {
		log.Printf("error while grpc response: %v\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.UserResponse{
			Id:     "-1",
			Status: "error",
			ErrMsg: "user is not allowed now",
		})
		return
	} else if err != nil {
		log.Printf("error while grpc response: %v\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.UserResponse{
			Id:     "-1",
			Status: "error",
			ErrMsg: err.Error()})
		return
	}

	if responseGrpc.Id != "-1" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(types.UserMessage{
		Id:       responseGrpc.Id,
		Login:    responseGrpc.Login,
		Email:    responseGrpc.Email,
		Password: responseGrpc.Password,
	})
}

func UserDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(types.UserResponse{
			Id:     "-1",
			Status: "error",
			ErrMsg: "method not allowed",
		})
		return
	}

	userUrl := os.Getenv("USER_SERVICE_URL")
	if userUrl == "" {
		log.Println("USER_SERVICE_URL is empty")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.UserResponse{
			Id:     "-1",
			Status: "error",
			ErrMsg: "user is not allowed now",
		})
		return
	}

	conn, err := grpc.NewClient(userUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to gRPC server (%s): %s", userUrl, err.Error())
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)
	reqeustGrpc := &pb.UserMessage{
		Id: gorilla.Vars(r)["user_id"],
	}

	responseGrpc, err := client.DeleteUser(context.Background(), reqeustGrpc)
	if responseGrpc == nil {
		log.Printf("error while grpc response: %v\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.UserResponse{
			Id:     "-1",
			Status: "error",
			ErrMsg: "user is not allowed now",
		})
		return
	} else if err != nil {
		log.Printf("error while grpc response: %v\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.UserResponse{
			Id:     "-1",
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

	json.NewEncoder(w).Encode(types.UserMessage{
		Id:       responseGrpc.Id,
		Login:    responseGrpc.Login,
		Email:    responseGrpc.Email,
		Password: responseGrpc.Password,
	})
}
