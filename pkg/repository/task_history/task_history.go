package repository

import (
	"database/sql"
	"dtms/pkg/domain"
	"dtms/pkg/errors"
	"encoding/json"
	"fmt"
	"log"
)

type taskHisoryRepository struct {
	db *sql.DB
}

func NewTaskHisoryRepository(db *sql.DB) (TaskHistoryRepository, error) {
	if db == nil {
		return nil, errors.NilSqlDb{}
	}

	return &taskHisoryRepository{
		db: db,
	}, nil
}

var _ TaskHistoryRepository = &taskHisoryRepository{}

func (r *taskHisoryRepository) Create(taskHisory *domain.TaskHistory) error {
	if err := r.db.Ping(); err != nil {
		log.Printf("cat reach database")
		return errors.UnconnectSqlDb{}
	}

	query := `INSERT INTO tasks.history (task_id, user_id, action, task_data) 
		VALUES ($1, $2, $3, $4);`

	taskData, err := json.Marshal(taskHisory.Common)
	if err != nil {
		log.Printf("error while marshaling taskHistory.Common")
		return err
	}

	insRes, err := r.db.Exec(query, taskHisory.Common.Id, taskHisory.Common.AuthorId, taskHisory.Action, taskData)
	if err != nil {
		log.Printf("error while trying to create new status string: %v", err)
		return err
	}
	if i, err := insRes.RowsAffected(); i < 1 || err != nil {
		log.Printf("no error was insert")
		return fmt.Errorf("%v: %v", errors.NoStatusCreate{}, err)
	}
	return nil
}
