// Example use of Go mongo-driver
package main

import (
	"context"
	"fmt"
	"log"
	"time"
	//"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongodbEndpoint = "mongodb://172.17.0.2:27017" // Find this from the Mongo container
)

// type Post struct {
// 	ID        primitive.ObjectID `bson:"_id"`
// 	Title     string             `bson:"title"`
// 	Body      string             `bson:"body"`
// 	Tags      []string           `bson:"tags"`
// 	Comments  uint64             `bson:"comments"`
// 	CreatedAt time.Time          `bson:"created_at"`
// 	UpdatedAt time.Time          `bson:"updated_at"`
// }

type dollars float32

type storeDB struct {
	ID        primitive.ObjectID `bson:"_id"`
	Item	  string 			 `bson:"item"`
	Price	  dollars			 `bson:"price"`
	Category  string			 `bson:"category"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}
func main() {
	// create a mongo client
	client, err := mongo.NewClient(
		options.Client().ApplyURI(mongodbEndpoint),
	)
	checkError(err)

	// Connect to mongo
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)

	// Disconnect
	defer client.Disconnect(ctx)

	col := client.Database("convenience_store").Collection("items")

	res, err := col.InsertOne(ctx, &storeDB{
		ID:			primitive.NewObjectID(),
		Item:		"Twix",
		Price: 		1.25,
		Category:   "Candy",
		CreatedAt: time.Now(),
	})

	
	
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
