package main

import (
	"context"
	"fmt"
	"log"
	//	"go.mongodb.org/mongo-driver/bson"
	//	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
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
