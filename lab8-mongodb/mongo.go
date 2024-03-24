// Example use of Go mongo-driver
package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"net/http"
	"strconv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongodbEndpoint = "mongodb://172.17.0.2:27017" // Find this from the Mongo container
)

type database struct{col *mongo.Collection}

type dollars float32
func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

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

	db := client.Database("convenience_store")
	collection := db.Collection("items")

	// Handlers
	

	//Additional Handlers
	

	res, err := collection.InsertOne(ctx, &storeDB{
		ID:			primitive.NewObjectID(),
		Item:		"Twix",
		Price: 		1.25,
		Category:   "Candy",
		CreatedAt: time.Now(),
	})
	fmt.Printf("inserted id: %s\n", res.InsertedID.(primitive.ObjectID).Hex())
	checkError(err)

	databaseInstance := database{col: collection}
	mux := http.NewServeMux()

	mux.HandleFunc("/list", databaseInstance.list)
	mux.HandleFunc("/price", databaseInstance.price)
	mux.HandleFunc("/update", databaseInstance.update)
	mux.HandleFunc("/create", databaseInstance.create)
	mux.HandleFunc("/read", databaseInstance.read)
	mux.HandleFunc("/delete", databaseInstance.delete)
	log.Fatal(http.ListenAndServe(":8000", mux))
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (db *database) list(w http.ResponseWriter, req *http.Request) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur, err := db.col.Find(ctx, bson.M{})
	checkError(err)

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var result storeDB
		err := cur.Decode(&result)
		checkError(err)

		fmt.Fprintf(w, "========================================================\n")
		fmt.Fprintf(w, "Item: %s\nCategory: %s\nPrice: %s\n", result.Item, result.Category, result.Price)
		fmt.Fprintf(w, "========================================================\n")
	}

	checkError(cur.Err())
}

func (db *database) price(w http.ResponseWriter, req *http.Request){
	
	itemName := req.URL.Query().Get("item")

	if itemName == "" {
		http.Error(w, "Missing item name parameter", http.StatusBadRequest)
        return
	}

	filter := bson.M{"item": itemName}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur := db.col.FindOne(ctx, filter)

	var result storeDB
	err := cur.Decode(&result)
	checkError(err)

	fmt.Fprintf(w, "========================================================\n")
	fmt.Fprintf(w, "Item: %s\nPrice: %s\n", result.Item, result.Price)
	fmt.Fprintf(w, "========================================================\n")
}

func (db *database) update(w http.ResponseWriter, req *http.Request){

	itemName := req.URL.Query().Get("item")
	updatedPrice := req.URL.Query().Get("price")

	price, err := strconv.ParseFloat(updatedPrice, 32)

	if err != nil {
		fmt.Println("Error encountered: ", err)
		fmt.Fprintf(w, "'%s' is not a valid price, please try again.\n", updatedPrice)
		return
	}

	if itemName == "" {
		http.Error(w, "Missing item name parameter", http.StatusBadRequest)
        return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"item": itemName}
	update := bson.M{"$set": bson.M{"price": price, "updated_at": time.Now()}}

	_, err = db.col.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println("Error updating document:", err)
		http.Error(w, "Failed to update price", http.StatusInternalServerError)
		return
	}


	fmt.Fprintf(w, "Price updated successfully for item: %s\n", itemName)
}

func (db *database) create(w http.ResponseWriter, req *http.Request){

	// Grab the item name, price, and category
	itemName := req.URL.Query().Get("item")
	PRICE := req.URL.Query().Get("price")
	itemCategory := req.URL.Query().Get("category")

	// Check for errors
	price, err := strconv.ParseFloat(PRICE, 32)
	if err != nil {
		fmt.Println("Error encountered: ", err)
		fmt.Fprintf(w, "'%s' is not a valid price, please try again.\n", PRICE)
		return
	}

	if itemName == "" {
		http.Error(w, "Missing item name parameter", http.StatusBadRequest)
        return
	}
	
	if itemCategory == "" {
		http.Error(w, "Missing category parameter", http.StatusBadRequest)
        return
	}

	// Start the context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Insert the information grabbed from the request
	doc := bson.M{
		"_id":			primitive.NewObjectID(),
		"item":     	itemName,
		"price":    	price,
		"category": 	itemCategory,
		"created_at":	time.Now(),	
	}

	// Inbsert the information into the database
	_, err = db.col.InsertOne(ctx, doc)

	// Check for errors
	if err != nil {
		fmt.Println("Error creating document:", err)
		http.Error(w, "Failed to create document", http.StatusInternalServerError)
		return
	}

	//Print out the item that was created
	fmt.Fprintf(w, "Item %s created successfully\n", itemName)
}

func (db *database) read(w http.ResponseWriter, req *http.Request) {
	db.list(w, req)
}

func (db *database) delete(w http.ResponseWriter, req *http.Request) {

	itemName := req.URL.Query().Get("item")

	if itemName == "" {
		http.Error(w, "Missing item name parameter", http.StatusBadRequest)
        return
	}

	filter := bson.M{"item": itemName}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := db.col.DeleteOne(ctx, filter)

	if err != nil {
		fmt.Println("Error deleting item:", err)
		http.Error(w, "Failed to delete item", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Item %s deleted successfully\n", itemName)

}