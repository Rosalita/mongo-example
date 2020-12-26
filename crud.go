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

func createOne(c *mongo.Client) {
	collection := c.Database("employees").Collection("devs")

	// Create - Inserting Documents.
	// To insert a single document, use InsertOne.
	insertResult, err := collection.InsertOne(context.TODO(), Developer{"Rosie", 3, true})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

}

func createMany(c *mongo.Client) {
	collection := c.Database("employees").Collection("devs")

	// To insert multiple documents, use InsertMany.
	devs := []interface{}{
		Developer{"Alice", 5, true},
		Developer{"Bob", 1, false},
		Developer{"Cat", 10, true},
	}

	insertManyResult, err := collection.InsertMany(context.TODO(), devs)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted many documents: ", insertManyResult.InsertedIDs)
}

func updateOne(c *mongo.Client) {
	collection := c.Database("employees").Collection("devs")

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
		{"$inc", bson.D{ // $inc increments a value.
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
}

func findOne(c *mongo.Client) {
	collection := c.Database("employees").Collection("devs")

	// create a value into which the result can be decoded.
	var result Developer

	// create a filter which will be used for the search.
	filter := bson.D{{Key: "name", Value: "Rosie"}}

	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found a single document: %+v\n", result)
}

func findMany(c *mongo.Client) {

	collection := c.Database("employees").Collection("devs")

	// The find method requires options.
	findOptions := options.Find()
	findOptions.SetLimit(2)

	// Using an empty bson.D as a filter will match all documents in the collection.
	filter := bson.D{{}}

	// The find method returns a cursor. A cursor provides a stream of documents
	// that can be iterated over and decoded one at a time.
	cur, err := collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Create an array to store the decoded results.
	var results []*Developer

	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var dev Developer
		err := cur.Decode(&dev)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, &dev)
	}

	// Check the cursor for errors.
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	cur.Close(context.TODO())

	fmt.Printf("Found multiple documents (array of pointers): %+v\n", results)
}

func deleteOne(c *mongo.Client) {
	collection := c.Database("employees").Collection("devs")

	filter := bson.D{{
		Key:   "name",
		Value: "Cat",
	}}

	deleteResult, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Deleted %v document(s)\n", deleteResult.DeletedCount)
}

func deleteMany(c *mongo.Client) {
	collection := c.Database("employees").Collection("devs")

	filter := bson.D{{
		"name", bson.D{{
			"$in", bson.A{
				"Alice",
				"Bob",
			},
		}},
	}}

	deleteResult, err := collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Deleted %v document(s)\n", deleteResult.DeletedCount)
}
