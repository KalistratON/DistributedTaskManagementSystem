package helper

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongo() (*mongo.Client, error) {
	connectioStringMongo := os.Getenv("MONGODB_CONNECTION_STRING")
	if connectioStringMongo == "" {
		return nil, fmt.Errorf("MONGODB_CONNECTION_STRING env is not setted")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectioStringMongo))
	if err != nil {
		return nil, err
	}
	return client, nil
}

func DisconnectMongo(client *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if client != nil {
		client.Disconnect(ctx)
	}
}

func InserToMongo(client *mongo.Client, database string, collection string, document interface{}) error {
	db := client.Database(database)
	if db == nil {
		return fmt.Errorf("can't create auth-database in mongo")
	}

	cln := db.Collection(collection)
	if cln == nil {
		return fmt.Errorf("can't createusers-collection in mongo")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := cln.InsertOne(ctx, document)
	if err != nil {
		return err
	}
	return nil
}

func FindMongo(client *mongo.Client, database string, collection string, filters map[string]string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	db := client.Database(database)
	if db == nil {
		return result, fmt.Errorf("can't create auth-database in mongo")
	}

	cln := db.Collection(collection)
	if cln == nil {
		return result, fmt.Errorf("can't createusers-collection in mongo")
	}

	bsonMap := make(bson.M)
	for key, val := range filters {
		bsonMap[key] = val
	}

	ctx := context.Background()

	err := cln.FindOne(ctx, bsonMap).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func FindInCollectionMongo(cln *mongo.Collection, filters map[string]string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	bsonMap := make(bson.M)
	for key, val := range filters {
		bsonMap[key] = val
	}

	ctx := context.Background()

	err := cln.FindOne(ctx, bsonMap).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}
