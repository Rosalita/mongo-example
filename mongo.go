package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// create a new mongo client
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// // CRUD operations
	// // Create one document
	// createOne(client)

	// // Create many documents
	// createMany(client)

	// // Update one document.
	// updateOne(client)

	// // Read a single document.
	// findOne(client)

	// // Read many documents.
	// find(client)

	// // Delete one document.
	// deleteOne(client)

	// // Delete many documents.
	// deleteMany(client)

	// Aggregation Pipeline
	// aggregate(client)

	// Graph Tree
	graphTree(client)

	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")

}
