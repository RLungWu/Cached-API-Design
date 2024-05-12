package db

import (
	"context"
	"errors"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client   *mongo.Client
	ads      *mongo.Collection
	database = "adservice"
	coll     = "ads"
)

func InitDatabaseConnection() error{
	uri := os.Getenv("MONGO_URI")
	if uri == ""{
		return errors.New("MONGO_URI is not set")
	}

	var err error
	clientOptions := options.Client().ApplyURI(uri)
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil{
		return err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil{
		return err
	}

	ads = client.Database(database).Collection(coll)
	return nil
}