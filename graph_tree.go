package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

// Graph data can be modelled in Mongo as a tree using the Parent References pattern.
// With this pattern, each node in the tree stores a reference to it's parent.
// See https://docs.mongodb.com/manual/tutorial/model-tree-structures-with-parent-references/

// This example creates a tree structure to model the data at
// https://en.wikipedia.org/wiki/Family_tree_of_the_Greek_gods

// Deity is a node in the tree representing a god or a goddess
type Deity struct {
	ID      string `bson:"_id"`
	Name    string `bson:"name"`
	KnownAs string `bson:"knownAs"`
	Parent  string `bson:"parent"`
}

func graphTree(c *mongo.Client) {

	collection := c.Database("trees").Collection("godsAndGoddesses")

	fmt.Println(collection)

	// seed some data into the database
	// seedTreeData(collection)

	// Aggregate Query with Graph Lookup
	// Get all the ancestors of the descendant known as Death

	// $graphLookup is a recursive search on a collection
	// options are used to restrict the search

	// Start with a filter which will capture the descendant known as Death
	filter := bson.M{
		"knownAs": "Death",
	}

	// Create an aggregate pipeline
	pipeline := bson.A{
		// first step in the pipeline is to match on the filter
		bson.M{"$match": filter},
		// This match will return some documents.
		// The second step is to perform a graph lookup for each matched document.
		bson.M{"$graphLookup": bson.M{
			// from the collection godsAndGoddesses - note this collection cannot be sharded.
			"from": "godsAndGoddesses",
			// startWith is the starting value, in this example it's the value stored in field parent
			// for the document(s) identified in the first $match step.
			// $graphLookup then begins recursively matching
			// it finds all the document(s) where the starting value is present in the connectToField
			// In this case this is where the field _id matches the value $parent
			// for this next round of matches it looks at the value in the connectFromField
			// if it matches the value in the connectToField it adds the matching document
			// to an result array, named by the 'as' parameter.
			// $graphLookup continues matching and building up the result array.
			// It stops when no more matches are found.
			"startWith":        "$parent",
			"connectToField":   "_id",
			"connectFromField": "parent",
			"as":               "ancestors",
		}},
	}
	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		log.Fatal(err)
	}

	var result []bson.M
	if err := cursor.All(context.TODO(), &result); err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
}

func seedTreeData(c *mongo.Collection) {

	// Create the root node of the tree.
	// The root node has no parent.
	rootID := newUUID()
	if _, err := c.InsertOne(context.TODO(), Deity{
		ID:      rootID,
		Name:    "Chaos",
		KnownAs: "The Void",
	}); err != nil {
		log.Fatal(err)
	}

	// generate some IDs to seed some children which have children
	gaiaID := newUUID()
	nyxID := newUUID()

	childrenOfChaos := []interface{}{
		Deity{newUUID(), "Tartarus", "The Abyss", rootID},
		Deity{gaiaID, "Gaia", "The Earth", rootID},
		Deity{newUUID(), "Eros", "Desire", rootID},
		Deity{newUUID(), "Erebus", "Darkness", rootID},
		Deity{nyxID, "Nyx", "The Night", rootID},
	}

	if _, err := c.InsertMany(context.TODO(), childrenOfChaos); err != nil {
		log.Fatal(err)
	}

	childrenOfGaia := []interface{}{
		Deity{newUUID(), "Typhon", "The Storms", gaiaID},
		Deity{newUUID(), "Uranus", "The Sky", gaiaID},
		Deity{newUUID(), "Ourea", "Mountains", gaiaID},
		Deity{newUUID(), "Pontus", "The Sea", gaiaID},
	}

	if _, err := c.InsertMany(context.TODO(), childrenOfGaia); err != nil {
		log.Fatal(err)
	}

	childrenOfNyx := []interface{}{
		Deity{newUUID(), "Moros", "Doom", nyxID},
		Deity{newUUID(), "Oneiroi", "Dreams", nyxID},
		Deity{newUUID(), "Nemesis", "Retribution", nyxID},
		Deity{newUUID(), "Momus", "Blame", nyxID},
		Deity{newUUID(), "Philotes", "Affection", nyxID},
		Deity{newUUID(), "Geras", "Aging", nyxID},
		Deity{newUUID(), "Thanatos", "Death", nyxID},
		Deity{newUUID(), "Hypnos", "Sleep", nyxID},
		Deity{newUUID(), "Eris", "Strife", nyxID},
		Deity{newUUID(), "Apate", "Deceit", nyxID},
		Deity{newUUID(), "Oizys", "Distress", nyxID},
	}

	if _, err := c.InsertMany(context.TODO(), childrenOfNyx); err != nil {
		log.Fatal(err)
	}
}

func newUUID() string {
	return uuid.New().String()
}
