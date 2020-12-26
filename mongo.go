package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Developer represents a developer.
type Developer struct {
	Name      string `bson:"name"`
	YearsExp  int    `bson:"yearsExperience"`
	LikesCats bool   `bson:"likesCats"`
}

func main() {
	// create a new mongo client
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Create one document
	createOne(client)

	// Create many documents
	createMany(client)

	// Update one document.
	updateOne(client)

	// Read a single document.
	findOne(client)

	// Read many documents.
	findMany(client)

	// Delete one document.
	deleteOne(client)

	// Delete many documents.
	deleteMany(client)

	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")

}
