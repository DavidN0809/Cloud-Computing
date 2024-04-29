package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
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
mux.Handle("/billings/create", authMiddleware(http.HandlerFunc(createBilling)))
mux.Handle("/billings/get/", authMiddleware(adminMiddleware(http.HandlerFunc(getBilling))))
mux.Handle("/billings/update/", authMiddleware(adminMiddleware(http.HandlerFunc(updateBilling))))
mux.Handle("/billings/remove/", authMiddleware(adminMiddleware(http.HandlerFunc(removeBilling))))
mux.Handle("/billings/removeAllBillings", http.HandlerFunc(removeAllBillings))
mux.Handle("/billings/createForTaskService", http.HandlerFunc(createBilling))
mux.Handle("/billings/listByUserID", http.HandlerFunc(listBillingsUserID))


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
	Hours  float64             `bson:"hours" json:"hours"`
        HourlyRate *float64           `bson:"hourly_rate,omitempty" json:"hourly_rate,omitempty"`
	Amount float64             `bson:"amount" json:"amount"`
}

func createBilling(w http.ResponseWriter, req *http.Request) {
    log.Println("Received request to create billing")  // Log the start of the operation

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

    defaultRate := 100.0
    if billing.HourlyRate == nil {
        billing.HourlyRate = &defaultRate
    }
    billing.Amount = billing.Hours * *billing.HourlyRate

    collection := client.Database("billing").Collection("billings")
    billing.ID = primitive.NewObjectID()
    _, err = collection.InsertOne(context.TODO(), billing)
    if err != nil {
        http.Error(w, "Failed to create billing", http.StatusInternalServerError)
        return
    }
    log.Printf("Billing created successfully: %+v", billing)  // Confirm successful creation
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(billing)
}

func getBilling(w http.ResponseWriter, req *http.Request) {
    if req.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    billingID := req.URL.Path[len("/billings/get/"):]
    log.Printf("Received request to get billing with ID: %s", billingID)  // Log the billing ID being queried
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


    log.Printf("Billing retrieved successfully: %+v", billing)  // Confirm successful retrieval
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

    var input struct {
        UserID     *primitive.ObjectID `json:"user_id"`
        TaskID     *primitive.ObjectID `json:"task_id"`
        Hours      *float64            `json:"hours"`
        HourlyRate *float64            `json:"hourly_rate"`
        Amount     *float64            `json:"amount"`
    }
    err = json.NewDecoder(req.Body).Decode(&input)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    collection := client.Database("billing").Collection("billings")
    filter := bson.M{"_id": objectID}

    // Fetch the current data to handle calculations properly
    var current Billing
    err = collection.FindOne(context.TODO(), filter).Decode(&current)
    if err != nil {
        http.Error(w, "Billing not found", http.StatusNotFound)
        return
    }

    update := bson.M{}
    if input.UserID != nil {
        update["user_id"] = *input.UserID
    }
    if input.TaskID != nil {
        update["task_id"] = *input.TaskID
    }
    if input.Hours != nil {
        update["hours"] = *input.Hours
    }
    if input.HourlyRate != nil {
        update["hourly_rate"] = *input.HourlyRate
    }

    // Determine amount calculation logic
    finalHours := current.Hours
    if input.Hours != nil {
        finalHours = *input.Hours
    }

    finalRate := current.HourlyRate
    if input.HourlyRate != nil {
        finalRate = input.HourlyRate
    } else if finalRate == nil {
        defaultRate := 100.0 // Default rate if no rate is recorded or provided
        finalRate = &defaultRate
    }

    if input.HourlyRate != nil || input.Hours != nil {
        update["amount"] = *finalRate * finalHours
    }
    if input.Amount != nil {
        update["amount"] = *input.Amount // Override calculated amount if direct amount is provided
    }

    _, err = collection.UpdateOne(context.TODO(), filter, bson.M{"$set": update})
    if err != nil {
        http.Error(w, "Failed to update billing", http.StatusInternalServerError)
        return
    }

    log.Println("Billing updated successfully")
    w.WriteHeader(http.StatusNoContent)
}

func removeBilling(w http.ResponseWriter, req *http.Request) {
    log.Println("Received request to remove all billings")  // Log the start of the operation

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
    log.Println("Billing removed successfully")  // Confirm successful update
    w.WriteHeader(http.StatusNoContent)
}

func listBillings(w http.ResponseWriter, req *http.Request) {
   log.Println("Recieved request to list billing") 

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

    log.Println("Billings listed successfully")  // Confirm successful operation
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(billings)
}


func removeAllBillings(w http.ResponseWriter, req *http.Request) {
    log.Println("Received request to remove all billings")  // Log the start of the operation
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

    log.Printf("All billings removed successfully, count: %d")  // Log the count of billings removed
    w.WriteHeader(http.StatusNoContent)
}
func listBillingsUserID(w http.ResponseWriter, req *http.Request) {
    log.Println("Received request to list billings")
    if req.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var filter bson.M = bson.M{}
    if userIDParam := req.URL.Query().Get("user_id"); userIDParam != "" {
        userID, err := primitive.ObjectIDFromHex(userIDParam)
        if err != nil {
            http.Error(w, "Invalid user ID", http.StatusBadRequest)
            return
        }
        filter["user_id"] = userID
    }

    collection := client.Database("billing").Collection("billings")
    cursor, err := collection.Find(context.TODO(), filter)
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

    log.Println("Billings listed successfully")
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(billings)
}
