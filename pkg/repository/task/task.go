package repository

import (
	"context"
	"dtms/pkg/domain"
	"dtms/pkg/errors"
	"fmt"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type taskRepository struct {
	mdb *mongo.Collection
	ctx context.Context
}

func NewTaskRepostory(mdb *mongo.Collection, ctx context.Context) (TaskRepository, error) {
	if mdb == nil {
		return nil, errors.NilMdb{}
	}
	return &taskRepository{
		mdb: mdb,
		ctx: ctx,
	}, nil
}

func (r *taskRepository) fitTask(task *domain.Task) {
	task.Id = strings.TrimSpace(task.Id)
	task.AuthorId = strings.TrimSpace(task.AuthorId)
	task.Name = strings.TrimSpace(task.Name)
	task.Description = strings.TrimSpace(task.Description)
	task.Deadline = strings.TrimSpace(task.Deadline)
	task.Status = strings.TrimSpace(task.Status)
}

func (r *taskRepository) validateTask(task *domain.Task) error {
	if task == nil {
		return errors.NilTask{}
	}
	if strings.TrimSpace(task.Name) == "" {
		return errors.EmptyNameTask{}
	}
	if strings.TrimSpace(task.AuthorId) == "" {
		return errors.EmptyAuthorIdTask{}
	}
	if strings.TrimSpace(task.Status) == "" {
		return errors.EmptyStatusTask{}
	}

	return nil
}

var _ TaskRepository = &taskRepository{}

// Create создаёт новую задачу в коллекции.
// Если такая задача уже существует, то будет возвращена ошибка
func (r *taskRepository) Create(task *domain.Task) (*domain.Task, error) {
	if err := r.validateTask(task); err != nil {
		log.Printf("task is not vaild: %v", err)
		return nil, err
	}
	r.fitTask(task)

	insRes, err := r.mdb.InsertOne(r.ctx, *task)
	if err != nil {
		log.Printf("error during insert on task creation: %v", err)
		return nil, err
	}

	id := insRes.InsertedID.(primitive.ObjectID).Hex()
	result, err := r.Get(id)
	if err != nil {
		log.Printf("error while checking inserting in mdb: %v", err)
		return nil, err
	}

	return result, nil
}

// Update обновляет задачу согласно task.Id и данынм внутри структуры
func (r *taskRepository) Update(task *domain.Task) (*domain.Task, error) {
	if err := r.validateTask(task); err != nil {
		log.Printf("task is not vaild: %v", err)
		return nil, err
	}
	r.fitTask(task)

	objId, err := primitive.ObjectIDFromHex(task.Id)
	if err != nil {
		log.Printf("task has invalid id: %v", err)
		return nil, fmt.Errorf("%v: %v", errors.InvalidIdTask{}, err)
	}

	updRes, err := r.mdb.UpdateByID(r.ctx, objId, bson.M{
		"$set": bson.M{
			"author_id":   task.AuthorId,
			"name":        task.Name,
			"description": task.Description,
			"deadline":    task.Deadline,
			"status":      task.Status,
		},
	})
	if err != nil {
		log.Printf("error while updating on task update")
		return nil, err
	}

	if updRes.MatchedCount < 1 || updRes.ModifiedCount < 1 {
		log.Println("no data was updated")
		return nil, errors.EmptyUpdateResultTask{}
	}

	result, err := r.Get(task.Id)
	if err != nil {
		log.Printf("error while checking updating result")
		return nil, err
	}
	return result, nil
}

// Get возвращает задачу согласно task.Id
func (r *taskRepository) Get(id string) (*domain.Task, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("task has invalid id: %v", err)
		return nil, fmt.Errorf("%v: %v", errors.InvalidIdTask{}, err)
	}

	var result domain.Task
	err = r.mdb.FindOne(r.ctx, bson.M{
		"_id": objId,
	}).Decode(&result)
	if err != nil {
		log.Printf("error while getting data by id: %v", err)
		return nil, err
	}
	result.Id = id

	return &result, nil
}

// GetAll возвращает все задачи, которые соответствуют фильтру
func (r *taskRepository) GetAll(filter map[string]interface{}) (*[]domain.Task, error) {
	cur, err := r.mdb.Find(r.ctx, filter)
	if err != nil {
		log.Printf("error while trying get data with filter: %v", err)
		return nil, err
	}
	defer cur.Close(r.ctx)

	type mgdbBufferTask struct {
		Id primitive.ObjectID `bson:"_id"`

		domain.Task `bson:",inline"`
	}

	buffer := make([]mgdbBufferTask, 0)
	if err = cur.All(r.ctx, &buffer); err != nil {
		log.Printf("error while decoding data: %v", err)
		return nil, err
	}

	result := make([]domain.Task, 0, len(buffer))
	for _, v := range buffer {
		result = append(result, domain.Task{
			Id:          v.Id.Hex(),
			AuthorId:    v.AuthorId,
			Name:        v.Name,
			Description: v.Description,
			Deadline:    v.Deadline,
			Status:      v.Status,
		})
	}
	return &result, nil
}

// Delete удаляет задачу согласно task.Id
func (r *taskRepository) Delete(task *domain.Task) (*domain.Task, error) {
	objId, err := primitive.ObjectIDFromHex(task.Id)
	if err != nil {
		log.Printf("task has invalid id: %v", err)
		return nil, fmt.Errorf("%v: %v", errors.InvalidIdTask{}, err)
	}

	result, err := r.Get(task.Id)
	if err != nil {
		log.Printf("task can't be found: %v", err)
		return nil, err
	}

	delRes, err := r.mdb.DeleteOne(r.ctx, bson.M{"_id": objId})
	if err != nil {
		log.Printf("task can not be deleted: %v", err)
		return nil, err
	}

	if delRes.DeletedCount < 1 {
		log.Printf("no data was deleted")
		return nil, errors.EmptyDeleteResultTask{}
	}

	return result, nil
}
