package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	collection := client.Database("employees").Collection("devs")

	// Create - Inserting Documents.
	// To insert a single document, use InsertOne.
	insertResult, err := collection.InsertOne(context.TODO(), Developer{"Rosie", 3, true})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	// To insert multiple documents, use InsertMany.
	devs := []interface{}{
		Developer{"Alice", 5, true},
		Developer{"Bob", 1, false},
		Developer{"Cat", 10, true},
	}

	insertManyResult, err := collection.InsertMany(context.TODO(), devs)
	fmt.Println("Inserted many documents: ", insertManyResult.InsertedIDs)

	// Updating documents.
	// To update single document, use updateOne.
	// updateOne requires a filter document to match documents in the database.
	// updateOne also requires an update document to describe the update.
	// The filter and update documents can be built using the bson.D type.

	// This is the simplest way to make a filter with bson.D
	filter := bson.D{{"name", "Rosie"}}

	// however this shows a warning in IDE to safeguard from changes in thirdparty code
	// causing silent breakages. e.g. if the order of the key and value swapped around,
	// things would break silently.

	// The full way to define a filter with a bson.D is like
	filter = bson.D{primitive.E{Key: "name", Value: "Rosie"}}

	// primitive.E represents a bson element for a D.
	// The type can be omitted however to safeguard against thirdparty changes
	// I like naming the Key and Value fields.
	filter = bson.D{{Key: "name", Value: "Rosie"}}


	// It's the same with creating updates from bson.D.
	update := bson.D{
		{"$inc", bson.D{ // $inc increments a value
			{"yearsExperience", 1},
		}},
	}

	// The IDE will complain unless the struct fields are explicitly named.
	update = bson.D{
		{Key: "$inc", Value: bson.D{
			{Key: "yearsExperience", Value: 1},
		}},
	}


	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")

}
