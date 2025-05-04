package helper

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type TaskStatusLog struct {
	TaskID string `json:"task_id"`
	Action string `json:"action"`
	UserID string `json:"user_id"`
}

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

func PushTaskStatusLog(db *sql.DB, statusLog TaskStatusLog) error {
	if db == nil {
		return fmt.Errorf("psql is nil")
	}

	query := `INSERT INTO tasks.status_log (task_id, action, user_id)
		VALUES ($1, $2, $3);
	`
	_, err := db.Exec(query, statusLog.TaskID, statusLog.Action, statusLog.UserID)
	if err != nil {
		log.Printf("Error while trying to exec query : %v", err)
		return err
	}
	return nil
}
