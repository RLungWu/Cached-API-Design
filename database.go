package main

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client   *mongo.Client
	ads      *mongo.Collection
	database = "adService"
	coll     = "ads"
)

func initialDatabaseConnection() {
	uri := os.Getenv("MONGODB_URI")
	if uri == ""{
		log.Fatalf("MONGODB_URI is not set")
	}
	
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil{
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil{
		log.Fatal(err)
	}

	ads = client.Database(database).Collection(coll)
	log.Println("Connected to MongoDB")

}
