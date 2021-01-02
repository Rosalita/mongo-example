package main

import (
	"context"
	"fmt"
	"log"

//	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
//	"go.mongodb.org/mongo-driver/mongo/options"
)

// Podcast represents the schema for the "Podcasts" collection
type Podcast struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Title  string             `bson:"title,omitempty"`
	Author string             `bson:"author,omitempty"`
	Tags   []string           `bson:"tags,omitempty"`
}

// Episode represents the schema for the "Episodes" collection
type Episode struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Podcast     primitive.ObjectID `bson:"podcast,omitempty"`
	Title       string             `bson:"title,omitempty"`
	Description string             `bson:"description,omitempty"`
	Duration    int32              `bson:"duration,omitempty"`
}

func aggregate(c *mongo.Client) {

	database := c.Database("quickstart")
	col := database.Collection("episodes")

	// seed some data into the database

	id1, err := primitive.ObjectIDFromHex("5e3b37e51c9d4400004117e6")
	id2, err := primitive.ObjectIDFromHex("5e3b37e51c9d4400004117e6")

	episodes := []interface{}{
		Episode{
			Podcast:     id1,
			Title:       "Episode #1",
			Description: "The first episode",
			Duration:    25,
		},
		Episode{
			Podcast:     id2,
			Title:       "Episode #2",
			Description: "The second episode",
			Duration:    30,
		},
	}

	insertManyResult, err := col.InsertMany(context.TODO(), episodes)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted many documents: ", insertManyResult.InsertedIDs)

}
