package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongodbEndpoint = "mongodb://mongodb:27017"
)

type Item struct {
	ID    primitive.ObjectID `bson:"_id"`
	Name  string             `bson:"item"`
	Price float64            `bson:"price"`
}

func main() {
	// Create a MongoDB client
	client, err := mongo.NewClient(options.Client().ApplyURI(mongodbEndpoint))
	checkError(err)

	// Connect to MongoDB
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	checkError(err)

	// Disconnect after the main function completes
	defer client.Disconnect(ctx)

	// Call seedData to seed initial data into MongoDB
	seedData(client)

	// Create a new ServeMux and database instance
	mux := http.NewServeMux()
	db := &database{client: client}

	// Register handlers for different routes
	mux.HandleFunc("/list", db.list)
	mux.HandleFunc("/price", db.price)
	mux.HandleFunc("/create", db.create)
	mux.HandleFunc("/update", db.update)
	mux.HandleFunc("/delete", db.delete)

	// Start the server
	log.Fatal(http.ListenAndServe(":8000", mux))
}

// database is a struct that holds a reference to the MongoDB client
type database struct {
	client *mongo.Client
}

func seedData(client *mongo.Client) {
    collection := client.Database("myDB").Collection("inventory")

    // Define initial data
    initialData := []interface{}{
        Item{Name: "shoes", Price: 50},
        Item{Name: "socks", Price: 5},
    }

    // Insert each item in the database if it doesn't exist
    for _, data := range initialData {
        filter := bson.M{"item": data.(Item).Name}
        update := bson.M{
            "$setOnInsert": bson.M{
                "_id":   primitive.NewObjectID(),
                "item":  data.(Item).Name,
                "price": data.(Item).Price,
            },
        }
        opts := options.FindOneAndUpdate().SetUpsert(true)

        var result Item
        err := collection.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&result)
        if err != nil {
            if err == mongo.ErrNoDocuments {
                log.Printf("Inserted initial data: %v", data)
            } else {
                log.Printf("Failed to insert initial data: %v", err)
            }
        } else {
            log.Printf("Initial data already exists: %v", data)
        }
    }
}

// list handles the "/list" route and lists all items in the inventory
func (db *database) list(w http.ResponseWriter, req *http.Request) {
	collection := db.client.Database("myDB").Collection("inventory")
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, "Failed to list items", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var item Item
		if err := cursor.Decode(&item); err != nil {
			http.Error(w, "Failed to decode item", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "%s: $%.2f\n", item.Name, item.Price)
	}
}

// price handles the "/price" route and retrieves the price of a specific item
func (db *database) price(w http.ResponseWriter, req *http.Request) {
	itemQuery := req.URL.Query().Get("item")
	collection := db.client.Database("myDB").Collection("inventory")
	filter := bson.M{"item": itemQuery}

	var result Item
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "%s: $%.2f\n", result.Name, result.Price)
}

// create handles the "/create" route and creates a new item in the inventory
func (db *database) create(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	priceStr := req.URL.Query().Get("price")
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		http.Error(w, "Invalid price format", http.StatusBadRequest)
		return
	}

	collection := db.client.Database("myDB").Collection("inventory")
	newItem := Item{
		ID:    primitive.NewObjectID(),
		Name:  item,
		Price: price,
	}
	_, err = collection.InsertOne(context.TODO(), newItem)
	if err != nil {
		http.Error(w, "Failed to create item", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Created %s: $%.2f\n", item, price)
}

// update handles the "/update" route and updates the price of an existing item
func (db *database) update(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	priceStr := req.URL.Query().Get("price")
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		http.Error(w, "Invalid price format", http.StatusBadRequest)
		return
	}

	collection := db.client.Database("myDB").Collection("inventory")
	filter := bson.M{"item": item}
	update := bson.M{"$set": bson.M{"price": price}}

	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		http.Error(w, "Failed to update item", http.StatusInternalServerError)
		return
	}
	if result.ModifiedCount == 0 {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Updated %s: $%.2f\n", item, price)
}

// delete handles the "/delete" route and deletes an item from the inventory
func (db *database) delete(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")

	collection := db.client.Database("myDB").Collection("inventory")
	filter := bson.M{"item": item}

	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		http.Error(w, "Failed to delete item", http.StatusInternalServerError)
		return
	}
	if result.DeletedCount == 0 {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Deleted %s\n", item)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
