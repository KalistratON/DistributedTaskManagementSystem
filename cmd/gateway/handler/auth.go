package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"dtms/cmd/gateway/types"

	pb "dtms/specs/go/pkg"

	gorilla "github.com/gorilla/mux"
	"google.golang.org/grpc/credentials/insecure"

	"google.golang.org/grpc"
)

func Authenticate(r *http.Request) bool {
	authUrl := os.Getenv("AUTH_SERVICE_URL")
	if authUrl == "" {
		log.Println("AUTH_SERVICE_URL is empty")
		return false
	}

	conn, err := grpc.NewClient(authUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to gRPC server (%s): %s", authUrl, err.Error())
	}
	defer conn.Close()

	client := pb.NewAuthServiceClient(conn)
	token := r.Header.Get("Authorization")
	if token == "" || !strings.HasPrefix(token, "Bearer ") {
		log.Println("Authorization is empty")
		return false
	}
	token = strings.TrimPrefix(token, "Bearer ")

	reqeustGrpc := &pb.AuthMessage{
		Id:    "-1",
		Token: token,
	}

	responseGrpc, err := client.Get(context.Background(), reqeustGrpc)
	if responseGrpc == nil || err != nil {
		log.Printf("error while grpc response: %v\n", err.Error())
		return false
	} else if responseGrpc.Id == "-1" {
		log.Printf("error while grpc response: responseGrpc.Id = -1: %v\n", err.Error())
		return false
	}

	return true
}

func AuthSoftCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(types.AuthMessage{
			Id:    "-1",
			Token: "-1",
		})
		return
	}

	var rqs types.AuthMessage
	if err := json.NewDecoder(r.Body).Decode(&rqs); err != nil {
		log.Println("error while trying to decode request")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.AuthMessage{
			Id:    "-1",
			Token: "-1",
		})
		return
	}
	defer r.Body.Close()

	authUrl := os.Getenv("AUTH_SERVICE_URL")
	if authUrl == "" {
		log.Println("AUTH_SERVICE_URL is empty")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.AuthMessage{
			Id:    "-1",
			Token: "-1",
		})
		return
	}

	conn, err := grpc.NewClient(authUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to gRPC server (%s): %s", authUrl, err.Error())
	}
	defer conn.Close()

	client := pb.NewAuthServiceClient(conn)
	reqeustGrpc := &pb.AuthMessage{
		Id:    rqs.Id,
		Token: rqs.Token,
	}

	responseGrpc, err := client.SoftCreate(context.Background(), reqeustGrpc)
	if responseGrpc == nil {
		log.Printf("error while grpc response: %v\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.AuthMessage{
			Id:    "-1",
			Token: "-1",
		})
		return
	} else if err != nil {
		log.Printf("error while grpc response: %v\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.AuthMessage{
			Id:    "-1",
			Token: "-1",
		})
		return
	}

	if responseGrpc.Token != "" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(types.AuthMessage{
		Id:    responseGrpc.Id,
		Token: responseGrpc.Token,
	})
}

func AuthHardCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(types.AuthMessage{
			Id:    "-1",
			Token: "-1",
		})
		return
	}

	var rqs types.AuthMessage
	if err := json.NewDecoder(r.Body).Decode(&rqs); err != nil {
		log.Println("error while trying to decode request")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.AuthMessage{
			Id:    "-1",
			Token: "-1",
		})
		return
	}
	defer r.Body.Close()

	authUrl := os.Getenv("AUTH_SERVICE_URL")
	if authUrl == "" {
		log.Println("AUTH_SERVICE_URL is empty")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.AuthMessage{
			Id:    "-1",
			Token: "-1",
		})
		return
	}

	conn, err := grpc.NewClient(authUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to gRPC server (%s): %s", authUrl, err.Error())
	}
	defer conn.Close()

	client := pb.NewAuthServiceClient(conn)
	reqeustGrpc := &pb.AuthMessage{
		Id:    rqs.Id,
		Token: rqs.Token,
	}

	responseGrpc, err := client.HardCreate(context.Background(), reqeustGrpc)
	if responseGrpc == nil {
		log.Printf("error while grpc response: %v\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.AuthMessage{
			Id:    "-1",
			Token: "-1",
		})
		return
	} else if err != nil {
		log.Printf("error while grpc response: %v\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.AuthMessage{
			Id:    "-1",
			Token: "-1",
		})
		return
	}

	if responseGrpc.Token != "" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(types.AuthMessage{
		Id:    responseGrpc.Id,
		Token: responseGrpc.Token,
	})
}

func AuthGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(types.AuthMessage{
			Id:    "-1",
			Token: "-1",
		})
		return
	}

	authUrl := os.Getenv("AUTH_SERVICE_URL")
	if authUrl == "" {
		log.Println("AUTH_SERVICE_URL is empty")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.AuthMessage{
			Id:    "-1",
			Token: "-1",
		})
		return
	}

	conn, err := grpc.NewClient(authUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to gRPC server (%s): %s", authUrl, err.Error())
	}
	defer conn.Close()

	client := pb.NewAuthServiceClient(conn)
	reqeustGrpc := &pb.AuthMessage{
		Id:    "-1",
		Token: gorilla.Vars(r)["token"],
	}

	responseGrpc, err := client.Get(context.Background(), reqeustGrpc)
	if responseGrpc == nil {
		log.Printf("error while grpc response: %v\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.AuthMessage{
			Id:    "-1",
			Token: "-1",
		})
		return
	} else if err != nil {
		log.Printf("error while grpc response: %v\n", err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(types.AuthMessage{
			Id:    "-1",
			Token: "-1",
		})
		return
	}

	if responseGrpc.Id != "-1" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(types.AuthMessage{
		Id:    responseGrpc.Id,
		Token: responseGrpc.Token,
	})
}

func AuthExtend(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(types.AuthMessage{
			Id:    "-1",
			Token: "-1",
		})
		return
	}

	authUrl := os.Getenv("AUTH_SERVICE_URL")
	if authUrl == "" {
		log.Println("AUTH_SERVICE_URL is empty")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.AuthMessage{
			Id:    "-1",
			Token: "-1",
		})
		return
	}

	conn, err := grpc.NewClient(authUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to gRPC server (%s): %s", authUrl, err.Error())
	}
	defer conn.Close()

	client := pb.NewAuthServiceClient(conn)
	reqeustGrpc := &pb.AuthMessage{
		Id:    "-1",
		Token: gorilla.Vars(r)["token"],
	}

	responseGrpc, err := client.Extend(context.Background(), reqeustGrpc)
	if responseGrpc == nil {
		log.Printf("error while grpc response: %v\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.AuthMessage{
			Id:    "-1",
			Token: "-1",
		})
		return
	} else if err != nil {
		log.Printf("error while grpc response: %v\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.AuthMessage{
			Id:    "-1",
			Token: "-1",
		})
		return
	}

	if responseGrpc.Id != "-1" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(types.AuthMessage{
		Id:    responseGrpc.Id,
		Token: responseGrpc.Token,
	})
}

func AuthDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(types.AuthMessage{
			Id:    "-1",
			Token: "-1",
		})
		return
	}

	authUrl := os.Getenv("AUTH_SERVICE_URL")
	if authUrl == "" {
		log.Println("AUTH_SERVICE_URL is empty")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.AuthMessage{
			Id:    "-1",
			Token: "-1",
		})
		return
	}

	conn, err := grpc.NewClient(authUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to gRPC server (%s): %s", authUrl, err.Error())
	}
	defer conn.Close()

	client := pb.NewAuthServiceClient(conn)
	reqeustGrpc := &pb.AuthMessage{
		Id:    "-1",
		Token: gorilla.Vars(r)["token"],
	}

	responseGrpc, err := client.Delete(context.Background(), reqeustGrpc)
	if responseGrpc == nil {
		log.Printf("error while grpc response: %v\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.AuthMessage{
			Id:    "-1",
			Token: "-1",
		})
		return
	} else if err != nil {
		log.Printf("error while grpc response: %v\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.AuthMessage{
			Id:    "-1",
			Token: "-1",
		})
		return
	}

	if responseGrpc.Id != "-1" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(types.AuthMessage{
		Id:    responseGrpc.Id,
		Token: responseGrpc.Token,
	})
}
