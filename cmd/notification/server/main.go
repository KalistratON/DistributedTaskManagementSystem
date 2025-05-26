package main

import (
	"context"
	types "dtms/cmd/notification/types"
	helper "dtms/pkg/database"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	mngCln, err := helper.ConnectMongo()
	if err != nil {
		log.Fatalf("task service down with fatal error related to mongo connection= %v", err)
	}
	defer mngCln.Disconnect(context.Background())

	service := NewNotificationService(mngCln)
	service.Run()
}

func parseMessageToTaskInfo(msg *kafka.Message) (types.NotificationTaskInfo, error) {
	var taskInfo types.NotificationTaskInfo
	if err := json.Unmarshal(msg.Value, &taskInfo); err != nil {
		log.Printf("error while trying parse message from kafka %v : %s", err, msg.Value)
		return types.NotificationTaskInfo{}, err
	}

	return taskInfo, nil
}

func getTaskLogInfoFromTaskInfo(clnTask *mongo.Collection, clnAuthor *mongo.Collection, taskInfo *types.NotificationTaskInfo) (types.NotificationTaskLogInfo, error) {
	var taskInfoLog types.NotificationTaskLogInfo

	objId, _ := primitive.ObjectIDFromHex(taskInfo.TaskId)
	filter := bson.M{"_id": objId}
	var result bson.M
	err := clnTask.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		log.Printf("error while trying get mongo-task: %v", err)
		return taskInfoLog, err
	}

	taskInfoLog.Name = result["name"].(string)
	taskInfoLog.Description = result["description"].(string)
	taskInfoLog.Deadline = result["deadline"].(string)

	taskInfoLog.Status = taskInfo.Status

	objId, _ = primitive.ObjectIDFromHex(taskInfo.AuthorId)
	filter = bson.M{"_id": objId}

	var resultAuthor bson.M
	err = clnAuthor.FindOne(context.Background(), filter).Decode(&resultAuthor)
	if err != nil {
		log.Printf("error while trying get mongo-task: %v", err)
		return taskInfoLog, err
	}
	if err != nil {
		log.Printf("error while trying get mongo-task: %v", err)
		return taskInfoLog, err
	}

	login, ok := resultAuthor["login"]
	if !ok {
		err = fmt.Errorf("login can't be found")
		log.Printf("error while trying get login from mongo-user: %v", err)
	}

	taskInfoLog.Author = login.(string)
	return taskInfoLog, nil
}

func pushToLog(taskInfo *types.NotificationTaskLogInfo) {
	log.Println("-------------")
	log.Println(taskInfo)
	log.Println("-------------")
}

type NotificationService struct {
	clnTask   *mongo.Collection
	clnAuthor *mongo.Collection
}

func NewNotificationService(client *mongo.Client) *NotificationService {
	return &NotificationService{
		clnTask:   client.Database("task").Collection("users"),
		clnAuthor: client.Database("auth").Collection("users"),
	}
}

func (service *NotificationService) Run() {
	topic := os.Getenv("KAFKA_TOPIC")
	if topic == "" {
		log.Printf("KAFKA_TOPIC is empty")
		return
	}

	kafkaCh, closeFunc, err := helper.Consume(topic)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	defer closeFunc()

	log.Println("ready to parse messages")
	for msg := range kafkaCh {
		log.Println("strart parsing ...")
		taskIdInfo, err := parseMessageToTaskInfo(&msg)
		if err != nil {
			log.Println("error while trying parse msg from kafka: %v", err)
		}

		taskInfo, err := getTaskLogInfoFromTaskInfo(service.clnTask, service.clnAuthor, &taskIdInfo)
		if err != nil {
			log.Println("error while trying translate msg from kafka: %v", err)
		}

		pushToLog(&taskInfo)
		log.Println("end parsing...")
	}
}
