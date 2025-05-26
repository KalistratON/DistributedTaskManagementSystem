package main

import (
	"dtms/cmd/gateway/handler"
	"dtms/cmd/gateway/middleware"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter().PathPrefix("/api/v1").Subrouter()

	r.HandleFunc("/auth/soft", handler.AuthSoftCreate).Methods("POST")
	r.HandleFunc("/auth/hard", handler.AuthHardCreate).Methods("POST")
	r.HandleFunc("/auth/{token}", handler.AuthGet).Methods("GET")
	r.HandleFunc("/auth/{token}", handler.AuthExtend).Methods("PUT")
	r.HandleFunc("/auth/{token}", handler.AuthDelete).Methods("DELETE")

	userRouter := r.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/", handler.UserPost).Methods("POST")
	userRouter.HandleFunc("/{user_id}", handler.UserPut).Methods("PUT")
	userRouter.HandleFunc("/{user_id}", handler.UserGet).Methods("GET")
	userRouter.HandleFunc("/{user_id}", handler.UserDelete).Methods("DELETE")

	taskRouter := r.PathPrefix("/task").Subrouter()
	taskRouter.HandleFunc("/", handler.TaskPost).Methods("POST")
	taskRouter.HandleFunc("/{task_id}", handler.TaskPut).Methods("PUT")
	taskRouter.HandleFunc("/{task_id}", handler.TaskGet).Methods("GET")
	taskRouter.HandleFunc("/{task_id}", handler.TaskDelete).Methods("DELETE")

	taskRouter.Use(middleware.Authenticate)

	log.Println("Starting server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
