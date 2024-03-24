package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	clientOptions := options.Client().ApplyURI("mongodb://mongodb:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")

	// Call  to seed initial data into MongoDB.
        seedData(client)
	    
	mux := http.NewServeMux()
	db := &database{client: client}
	mux.HandleFunc("/list", db.list)
	mux.HandleFunc("/price", db.price)
	mux.HandleFunc("/create", db.create)
	mux.HandleFunc("/update", db.update)
	mux.HandleFunc("/delete", db.delete)

	log.Fatal(http.ListenAndServe(":8000", mux))
}

type database struct {
	client *mongo.Client
}

func seedData(client *mongo.Client) {
    collection := client.Database("blog").Collection("posts")

    // Define initial data
    initialData := []interface{}{
        Post{
            ID:        primitive.NewObjectID(),
            Title:     "Welcome",
            Body:      "Welcome to the blog!",
            Tags:      []string{"introduction"},
            Comments:  0,
            CreatedAt: time.Now(),
            UpdatedAt: time.Now(),
        },
        Post{
            ID:        primitive.NewObjectID(),
            Title:     "MongoDB",
            Body:      "MongoDB is a NoSQL database",
            Tags:      []string{"mongodb", "database"},
            Comments:  0,
            CreatedAt: time.Now(),
            UpdatedAt: time.Now(),
        },
    }

    // Upsert each item in the initial data
    for _, data := range initialData {
        filter := bson.M{"_id": data.(Post).ID}
        update := bson.M{"$set": data}
        opts := options.Update().SetUpsert(true)

        _, err := collection.UpdateOne(context.Background(), filter, update, opts)
        if err != nil {
            log.Printf("Failed to upsert initial data: %v", err)
        } else {
            log.Printf("Upserted initial data: %+v", data)
        }
    }
}

func (db *database) list(w http.ResponseWriter, req *http.Request) {
	collection := db.client.Database("myDB").Collection("inventory")
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		http.Error(w, "Failed to list items", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var item struct {
			Item  string  `bson:"item"`
			Price float64 `bson:"price"`
		}
		if err := cursor.Decode(&item); err != nil {
			http.Error(w, "Failed to decode item", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "%s: $%.2f\n", item.Item, item.Price)
	}
}

func (db *database) price(w http.ResponseWriter, req *http.Request) {
	itemQuery := req.URL.Query().Get("item")
	collection := db.client.Database("myDB").Collection("inventory")
	filter := bson.D{{Key: "item", Value: itemQuery}}

	var result struct {
		Item  string  `bson:"item"`
		Price float64 `bson:"price"`
	}

	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "%s: $%.2f\n", result.Item, result.Price)
}

func (db *database) create(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	priceStr := req.URL.Query().Get("price")
	price, err := strconv.ParseFloat(priceStr, 32)
	if err != nil {
		http.Error(w, "Invalid price format", http.StatusBadRequest)
		return
	}

	collection := db.client.Database("myDB").Collection("inventory")
	_, err = collection.InsertOne(context.TODO(), bson.D{
		{Key: "item", Value: item},
		{Key: "price", Value: price},
	})
	if err != nil {
		http.Error(w, "Failed to create item", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Created %s: $%.2f\n", item, price)
}

func (db *database) update(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	priceStr := req.URL.Query().Get("price")
	price, err := strconv.ParseFloat(priceStr, 32)
	if err != nil {
		http.Error(w, "Invalid price format", http.StatusBadRequest)
		return
	}

	collection := db.client.Database("myDB").Collection("inventory")
	filter := bson.D{{Key: "item", Value: item}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "price", Value: price}}}}

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

func (db *database) delete(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")

	collection := db.client.Database("myDB").Collection("inventory")
	filter := bson.D{{Key: "item", Value: item}}

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


