package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongo() (*mongo.Client, error) {
	connectioStringMongo := os.Getenv("MONGODB_CONNECTION_STRING")
	if connectioStringMongo == "" {
		return nil, fmt.Errorf("MONGODB_CONNECTION_STRING env is not setted")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectioStringMongo))
	if err != nil {
		return nil, err
	}
	return client, nil
}
