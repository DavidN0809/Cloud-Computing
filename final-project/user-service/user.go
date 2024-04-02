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
	client, err = mongo.NewClient(options.Client().ApplyURI("mongodb://user-mongodb:27017"))
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

	// User endpoints
	mux.HandleFunc("/users/list", listUsers)
	mux.HandleFunc("/users/create", createUser)
	mux.HandleFunc("/users/get/", getUser)
	mux.HandleFunc("/users/update/", updateUser)
	mux.HandleFunc("/users/remove/", removeUser)

	// Start the server
	log.Println("User Service listening on port 8001...")
	log.Fatal(http.ListenAndServe(":8001", mux))
}

type User struct {
	ID       primitive.ObjectID `bson:"_id" json:"id"`
	Username string             `bson:"username" json:"username"`
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password" json:"password"`
}

func createUser(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user User
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	collection := client.Database("user").Collection("users")
	user.ID = primitive.NewObjectID()
	_, err = collection.InsertOne(context.TODO(), user)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func getUser(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := req.URL.Path[len("/users/get/"):]
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	collection := client.Database("user").Collection("users")
	filter := bson.M{"_id": objectID}

	var user User
	err = collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func updateUser(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := req.URL.Path[len("/users/update/"):]
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user User
	err = json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	collection := client.Database("user").Collection("users")
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{
		"username": user.Username,
		"email":    user.Email,
		"password": user.Password,
	}}

	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func removeUser(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := req.URL.Path[len("/users/remove/"):]
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	collection := client.Database("user-db").Collection("users")
	filter := bson.M{"_id": objectID}

	_, err = collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		http.Error(w, "Failed to remove user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func listUsers(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	collection := client.Database("user-db").Collection("users")
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, "Failed to list users", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var users []User
	err = cursor.All(context.Background(), &users)
	if err != nil {
		http.Error(w, "Failed to decode users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
