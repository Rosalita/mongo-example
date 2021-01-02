package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	//seedData(col)

	podcastID, _ := primitive.ObjectIDFromHex("5e3b37e51c9d4400004117e6")

	// Matching stage that matches all documents that have a specific podcastID
	matchStage := bson.D{{
		Key: "$match",
		Value: bson.D{{
			Key:   "podcast",
			Value: podcastID,
		}},
	}}

	// Grouping stage groups the data by podcastID, which means there's only one result in the group.
	// Then it sums all the episode durations for each item in the group.
	groupStage := bson.D{{
		Key: "$group",
		Value: bson.D{
			{
				Key:   "_id",
				Value: "$podcast",
			},
			{
				Key: "total",
				Value: bson.D{
					{
						Key:   "$sum",
						Value: "$duration",
					},
				},
			},
		},
	}}

	showInfoCursor, err := col.Aggregate(context.TODO(), mongo.Pipeline{matchStage, groupStage})
	if err != nil {
		log.Fatal(err)
	}

	var showsWithInfo []bson.M
	if err = showInfoCursor.All(context.TODO(), &showsWithInfo); err != nil {
		log.Fatal(err)
	}
	fmt.Println(showsWithInfo)
}

func seedData(c *mongo.Collection) {
	podcastID, _ := primitive.ObjectIDFromHex("5e3b37e51c9d4400004117e6")

	episodes := []interface{}{
		Episode{
			Podcast:     podcastID,
			Title:       "Episode #1",
			Description: "The first episode",
			Duration:    25,
		},
		Episode{
			Podcast:     podcastID,
			Title:       "Episode #2",
			Description: "The second episode",
			Duration:    30,
		},
	}

	insertManyResult, err := c.InsertMany(context.TODO(), episodes)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted many documents: ", insertManyResult.InsertedIDs)
}
