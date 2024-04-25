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
	client, err = mongo.NewClient(options.Client().ApplyURI("mongodb://task-mongodb:27017"))
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

	// Task endpoints
	mux.HandleFunc("/tasks/list", listTasks)
	mux.HandleFunc("/tasks/create", createTask)
	mux.HandleFunc("/tasks/get/", getTask)
	mux.HandleFunc("/tasks/update/", updateTask)
	mux.HandleFunc("/tasks/remove/", authMiddleware(adminMiddleware(removeTask)))
	mux.HandleFunc("/tasks/removeAllTasks", removeAllTasks)

	// Start the server
	log.Println("Task Service listening on port 8002...")
	log.Fatal(http.ListenAndServe(":8002", mux))
}

func ensureDatabaseAndCollection(client *mongo.Client) error {
	dbName := "taskmanagement"
	collectionName := "tasks"

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

type Task struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	AssignedTo  primitive.ObjectID `bson:"assigned_to" json:"assigned_to"`
	Status      string             `bson:"status" json:"status"`
	Hours       float64            `bson:"hours" json:"hours"`
	InvoiceID   primitive.ObjectID
	ParentTask  *primitive.ObjectID `bson:"parent_task,omitempty" json:"parent_task,omitempty"`
}
type Invoice struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`                    // Unique identifier for the invoice
	TaskID      primitive.ObjectID `bson:"task_id" json:"task_id"`           // Associated task ID
	UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`           // User ID of the person responsible for the task
	Description string             `bson:"description" json:"description"`   // Description or title of the invoice
	DateIssued  time.Time          `bson:"date_issued" json:"date_issued"`   // Date when the invoice was issued
	Hours       float64            `bson:"hours" json:"hours"`               // Total hours worked on the task
	HourlyRate  float64            `bson:"hourly_rate" json:"hourly_rate"`   // Hourly rate for the work
	Amount      float64            `bson:"total_amount" json:"total_amount"` // Total amount due
}

func createTask(w http.ResponseWriter, req *http.Request) {
	var task Task
	err := json.NewDecoder(req.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if task.ParentTask != nil {
		var parentTask Task
		err = client.Database("taskmanagement").Collection("tasks").FindOne(context.TODO(), bson.M{"_id": *task.ParentTask}).Decode(&parentTask)
		if err != nil {
			http.Error(w, "Parent task not found", http.StatusNotFound)
			return
		}
	}

	task.ID = primitive.NewObjectID()
	_, err = client.Database("taskmanagement").Collection("tasks").InsertOne(context.TODO(), task)
	if err != nil {
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func getTask(w http.ResponseWriter, req *http.Request) {
	taskID := req.URL.Path[len("/tasks/get/"):]
	objectID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var task Task
	err = client.Database("taskmanagement").Collection("tasks").FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&task)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	var subtasks []Task
	cursor, err := client.Database("taskmanagement").Collection("tasks").Find(context.TODO(), bson.M{"parent_task": objectID})
	if err == nil {
		defer cursor.Close(context.Background())
		cursor.All(context.Background(), &subtasks)
	}

	response := struct {
		Task     Task   `json:"task"`
		Subtasks []Task `json:"subtasks"`
	}{
		Task:     task,
		Subtasks: subtasks,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func updateTask(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	taskID := req.URL.Path[len("/tasks/update/"):]
	objectID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}
	err = json.NewDecoder(req.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Prepare update document
	updateDoc := bson.M{"$set": bson.M{}}
	for key, value := range updates {
		// Ensure only allowed fields are updated
		switch key {
		case "title", "description", "assigned_to", "status", "hours":
			updateDoc["$set"].(bson.M)[key] = value
		}
	}

	collection := client.Database("taskmanagement").Collection("tasks")
	// Fetch the current task to compare changes
	var currentTask Task
	err = collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&currentTask)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	// Handle InvoiceID creation if task status changes to 'done'
	if currentTask.Status != "done" && updates["status"] == "done" {
		// Generate a new ObjectID for InvoiceID if it's transitioning to 'done'
		invoiceID := primitive.NewObjectID()
		updateDoc["$set"].(bson.M)["InvoiceID"] = invoiceID
		log.Printf("Task updated to 'done'. New InvoiceID: %v generated", invoiceID)
	}

	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": objectID}, updateDoc)
	if err != nil {
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func removeTask(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	taskID := req.URL.Path[len("/tasks/remove/"):]
	objectID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	collection := client.Database("taskmanagement").Collection("tasks")
	filter := bson.M{"_id": objectID}

	_, err = collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		http.Error(w, "Failed to remove task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func listTasks(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	collection := client.Database("taskmanagement").Collection("tasks")
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, "Failed to list tasks", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var tasks []Task
	err = cursor.All(context.Background(), &tasks)
	if err != nil {
		http.Error(w, "Failed to decode tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func removeAllTasks(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	collection := client.Database("taskmanagement").Collection("tasks")

	_, err := collection.DeleteMany(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, "Failed to remove all tasks", http.StatusInternalServerError)
		return
	}
}

const hourlyRate = 100.0 // Adjust this value as necessary

func createInvoice(task Task) (primitive.ObjectID, error) {
	invoice := Invoice{
		TaskID:      task.ID,
		UserID:      task.AssignedTo, // assuming the assigned user is the one being billed
		Description: "Invoice for " + task.Title,
		DateIssued:  time.Now(),
		Hours:       task.Hours,
		Amount:      task.Hours * hourlyRate,
	}

	collection := client.Database("taskmanagement").Collection("invoices")
	result, err := collection.InsertOne(context.TODO(), invoice)
	if err != nil {
		log.Printf("Failed to create Invoice: %v", err)
		return primitive.NilObjectID, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}
func getInvoiceByTaskID(taskID primitive.ObjectID) (*Invoice, error) {
	var invoice Invoice
	collection := client.Database("taskmanagement").Collection("invoices")
	err := collection.FindOne(context.TODO(), bson.M{"task_id": taskID}).Decode(&invoice)
	if err != nil {
		return nil, err
	}
	return &invoice, nil
}
