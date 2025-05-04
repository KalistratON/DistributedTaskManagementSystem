package main

import (
	"dtms/gateway/internal/handler"
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

	r.HandleFunc("/user", handler.UserPost).Methods("POST")
	r.HandleFunc("/user/{user_id}", handler.UserPut).Methods("PUT")
	r.HandleFunc("/user/{user_id}", handler.UserGet).Methods("GET")
	r.HandleFunc("/user/{user_id}", handler.UserDelete).Methods("DELETE")

	r.HandleFunc("/task", handler.TaskPost).Methods("POST")
	r.HandleFunc("/task/{task_id}", handler.TaskPut).Methods("PUT")
	r.HandleFunc("/task/{task_id}", handler.TaskGet).Methods("GET")
	r.HandleFunc("/task/{task_id}", handler.TaskDelete).Methods("DELETE")

	log.Println("Starting server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
