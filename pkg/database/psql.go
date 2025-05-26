package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func ConnectPsql() (*sql.DB, error) {
	psqlUrl := os.Getenv("POSTGRES_CONNECTION_STRING")

	if psqlUrl == "" {
		return nil, fmt.Errorf("POSTGRES_CONNECTION_STRING does'n setted")
	}

	db, err := sql.Open("postgres", psqlUrl)
	if err != nil {
		log.Printf("Error while trying connect to : %v", err)
		return nil, err
	}
	return db, nil
}
