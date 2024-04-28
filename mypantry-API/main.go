package main

import (
    "context"
    "log"
    "os"
  //  "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "github.com/hburnt/mypantry-API/recipeapi"
    "github.com/hburnt/MyPantry-API/recipeinfoapi"
)

const mongoURI = "mongodb://127.17.0.2:27017"

func main() {
    // Set up MongoDB connection
    clientOptions := options.Client().ApplyURI(mongoURI)
    client, err := mongo.Connect(context.Background(), clientOptions)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Disconnect(context.Background())

    // Access a MongoDB collection
    collection := client.Database("MyPantryDB").Collection("recipes")

    // Get API key from environment variable
	  apiKey := os.Getenv("SPOONACULAR_API_KEY")

    // Create a client with API key
    apiClient := recipeapi.NewClient(apiKey)

    // Query input
    query := "pasta"
    recipeID := 654959 
    // Make API call to get recipe
    recipe, err := apiClient.GetRecipe(query)
    if err != nil {
        log.println("error:", err)
        return
    }

    // Insert recipe into the database
    _, err = collection.InsertOne(context.Background(), recipe)
    if err != nil {
        log.Fatal(err)
    }

    collection := client.Database("MyPantryDB").Collection("recipe_info")
    recipe_info, err := apiClient.GetRecipeInfo(recipeID)
    if err != nil {
        log.println("error:", err)
        return
    }
    _, err = collection.InsertOne(context.Background(), recipe_info)
    if err != nil {
      log.Fatal(err)
     }
    log.Println("Recipe inserted successfully.")
}
