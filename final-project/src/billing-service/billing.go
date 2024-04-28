package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func main() {
    // Create a new MongoDB client
    var err error
    client, err = mongo.NewClient(options.Client().ApplyURI("mongodb://billing-mongodb:27017"))
    if err != nil {
        log.Fatal(err)
    }

    // Connect to MongoDB
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    err = client.Connect(ctx)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Disconnect(ctx)

    // Check if the database and collection exist, create them if they don't
    err = ensureDatabaseAndCollection(client)
    if err != nil {
        log.Fatal(err)
    }

    // Create a new HTTP server
    mux := http.NewServeMux()

    // Billing endpoints
mux.Handle("/billings/list", authMiddleware(adminMiddleware(http.HandlerFunc(listBillings))))
mux.Handle("/billings/create", authMiddleware(taskServiceMiddleware(http.HandlerFunc(createBilling))))
mux.Handle("/billings/get/", authMiddleware(adminMiddleware(http.HandlerFunc(getBilling))))
mux.Handle("/billings/update/", authMiddleware(adminMiddleware(http.HandlerFunc(updateBilling))))
mux.Handle("/billings/remove/", authMiddleware(adminMiddleware(http.HandlerFunc(removeBilling))))
mux.Handle("/billings/removeAllBillings", http.HandlerFunc(removeAllBillings))

    // Start the server
    log.Println("Billing Service listening on port 8003...")
    log.Fatal(http.ListenAndServe(":8003", mux))
}
func ensureDatabaseAndCollection(client *mongo.Client) error {
    dbName := "billing"
    collectionName := "billings"

    // Check if the database exists
    databases, err := client.ListDatabaseNames(context.Background(), bson.M{})
    if err != nil {
        return err
    }

    dbExists := false
    for _, db := range databases {
        if db == dbName {
            dbExists = true
            break
        }
    }

    if !dbExists {
        // Create the database if it doesn't exist
        err = client.Database(dbName).CreateCollection(context.Background(), collectionName)
        if err != nil {
            return err
        }
        log.Printf("Created database '%s' and collection '%s'", dbName, collectionName)
    } else {
        // Check if the collection exists
        collections, err := client.Database(dbName).ListCollectionNames(context.Background(), bson.M{})
        if err != nil {
            return err
        }

        collectionExists := false
        for _, coll := range collections {
            if coll == collectionName {
                collectionExists = true
                break
            }
        }

        if !collectionExists {
            // Create the collection if it doesn't exist
            err = client.Database(dbName).CreateCollection(context.Background(), collectionName)
            if err != nil {
                return err
            }
            log.Printf("Created collection '%s' in database '%s'", collectionName, dbName)
        }
    }

    return nil
}


type Billing struct {
	ID     primitive.ObjectID `bson:"_id" json:"id"`
	UserID primitive.ObjectID `bson:"user_id" json:"user_id"`
	TaskID primitive.ObjectID `bson:"task_id" json:"task_id"`
	Hours  string             `bson:"hours" json:"hours"`
	Amount string             `bson:"amount" json:"amount"`
}

func createBilling(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var billing Billing
	err := json.NewDecoder(req.Body).Decode(&billing)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Parse hours and calculate the amount
	hours, err := strconv.ParseFloat(billing.Hours, 64)
	if err != nil {
		http.Error(w, "Invalid hours format", http.StatusBadRequest)
		return
	}
	hourlyRate := 100.0 // Adjust this value as necessary
	amount := hours * hourlyRate
	billing.Amount = strconv.FormatFloat(amount, 'f', 2, 64)

	collection := client.Database("billing").Collection("billings")
	billing.ID = primitive.NewObjectID()
	_, err = collection.InsertOne(context.TODO(), billing)
	if err != nil {
		http.Error(w, "Failed to create billing", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(billing)
}

func getBilling(w http.ResponseWriter, req *http.Request) {
    if req.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    billingID := req.URL.Path[len("/billings/get/"):]
    objectID, err := primitive.ObjectIDFromHex(billingID)
    if err != nil {
        http.Error(w, "Invalid billing ID", http.StatusBadRequest)
        return
    }

    collection := client.Database("billing").Collection("billings")
    filter := bson.M{"_id": objectID}

    var billing Billing
    err = collection.FindOne(context.TODO(), filter).Decode(&billing)
    if err != nil {
        http.Error(w, "Billing not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(billing)
}

func updateBilling(w http.ResponseWriter, req *http.Request) {
    if req.Method != http.MethodPut {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    billingID := req.URL.Path[len("/billings/update/"):]
    objectID, err := primitive.ObjectIDFromHex(billingID)
    if err != nil {
        http.Error(w, "Invalid billing ID", http.StatusBadRequest)
        return
    }

    var billing Billing
    err = json.NewDecoder(req.Body).Decode(&billing)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    collection := client.Database("billing").Collection("billings")
    filter := bson.M{"_id": objectID}
    update := bson.M{"$set": bson.M{
        "user_id": billing.UserID,
        "task_id": billing.TaskID,
        "hours":   billing.Hours,
        "amount":  billing.Amount,
    }}

    _, err = collection.UpdateOne(context.TODO(), filter, update)
    if err != nil {
        http.Error(w, "Failed to update billing", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

func removeBilling(w http.ResponseWriter, req *http.Request) {
    if req.Method != http.MethodDelete {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    billingID := req.URL.Path[len("/billings/remove/"):]
    objectID, err := primitive.ObjectIDFromHex(billingID)
    if err != nil {
        http.Error(w, "Invalid billing ID", http.StatusBadRequest)
        return
    }

    collection := client.Database("billing").Collection("billings")
    filter := bson.M{"_id": objectID}

    _, err = collection.DeleteOne(context.TODO(), filter)
    if err != nil {
        http.Error(w, "Failed to remove billing", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

func listBillings(w http.ResponseWriter, req *http.Request) {
    if req.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    collection := client.Database("billing").Collection("billings")
    cursor, err := collection.Find(context.TODO(), bson.M{})
    if err != nil {
        http.Error(w, "Failed to list billings", http.StatusInternalServerError)
        return
    }
    defer cursor.Close(context.Background())

    var billings []Billing
    err = cursor.All(context.Background(), &billings)
    if err != nil {
        http.Error(w, "Failed to decode billings", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(billings)
}


func removeAllBillings(w http.ResponseWriter, req *http.Request) {
    if req.Method != http.MethodDelete {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    collection := client.Database("billing").Collection("billings")

    _, err := collection.DeleteMany(context.TODO(), bson.M{})
    if err != nil {
        http.Error(w, "Failed to remove all billings", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
