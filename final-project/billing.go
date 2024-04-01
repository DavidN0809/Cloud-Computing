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

	// Create a new HTTP server
	mux := http.NewServeMux()

	// Billing endpoints
	mux.HandleFunc("/billings", handleBillings)
	mux.HandleFunc("/billings/", handleBilling)

	// Start the server
	log.Println("Billing Service listening on port 8003...")
	log.Fatal(http.ListenAndServe(":8003", mux))
}

type Billing struct {
	ID     primitive.ObjectID `bson:"_id" json:"id"`
	UserID primitive.ObjectID `bson:"user_id" json:"user_id"`
	TaskID primitive.ObjectID `bson:"task_id" json:"task_id"`
	Hours  float64            `bson:"hours" json:"hours"`
	Amount float64            `bson:"amount" json:"amount"`
}

func handleBillings(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		listBillings(w, req)
	case http.MethodPost:
		createBilling(w, req)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleBilling(w http.ResponseWriter, req *http.Request) {
	billingID := req.URL.Path[len("/billings/"):]

	switch req.Method {
	case http.MethodGet:
		getBilling(w, req, billingID)
	case http.MethodPut:
		updateBilling(w, req, billingID)
	case http.MethodDelete:
		removeBilling(w, req, billingID)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func createBilling(w http.ResponseWriter, req *http.Request) {
	var billing Billing
	err := json.NewDecoder(req.Body).Decode(&billing)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

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

func getBilling(w http.ResponseWriter, req *http.Request, billingID string) {
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

func updateBilling(w http.ResponseWriter, req *http.Request, billingID string) {
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

func removeBilling(w http.ResponseWriter, req *http.Request, billingID string) {
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
