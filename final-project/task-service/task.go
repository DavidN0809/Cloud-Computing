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

    // Create a new HTTP server
    mux := http.NewServeMux()

    // Task endpoints
    mux.HandleFunc("/", handleTasks)
    mux.HandleFunc("/", handleTask)

    // Start the server
    log.Println("Task Service listening on port 8002...")
    log.Fatal(http.ListenAndServe(":8002", mux))
}

type Task struct {
    ID          primitive.ObjectID `bson:"_id" json:"id"`
    Title       string             `bson:"title" json:"title"`
    Description string             `bson:"description" json:"description"`
    AssignedTo  primitive.ObjectID `bson:"assigned_to" json:"assigned_to"`
    Status      string             `bson:"status" json:"status"`
    Hours       float64            `bson:"hours" json:"hours"`
}

func handleTasks(w http.ResponseWriter, req *http.Request) {
    switch req.Method {
    case http.MethodGet:
        listTasks(w, req)
    case http.MethodPost:
        createTask(w, req)
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

func handleTask(w http.ResponseWriter, req *http.Request) {
    taskID := req.URL.Path[len("/"):]

    switch req.Method {
    case http.MethodGet:
        getTask(w, req, taskID)
    case http.MethodPut:
        updateTask(w, req, taskID)
    case http.MethodDelete:
        removeTask(w, req, taskID)
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

func createTask(w http.ResponseWriter, req *http.Request) {
    var task Task
    err := json.NewDecoder(req.Body).Decode(&task)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    collection := client.Database("taskmanagement").Collection("tasks")
    task.ID = primitive.NewObjectID()
    _, err = collection.InsertOne(context.TODO(), task)
    if err != nil {
        http.Error(w, "Failed to create task", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(task)
}

func getTask(w http.ResponseWriter, req *http.Request, taskID string) {
    objectID, err := primitive.ObjectIDFromHex(taskID)
    if err != nil {
        http.Error(w, "Invalid task ID", http.StatusBadRequest)
        return
    }

    collection := client.Database("taskmanagement").Collection("tasks")
    filter := bson.M{"_id": objectID}

    var task Task
    err = collection.FindOne(context.TODO(), filter).Decode(&task)
    if err != nil {
        http.Error(w, "Task not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(task)
}

func updateTask(w http.ResponseWriter, req *http.Request, taskID string) {
    objectID, err := primitive.ObjectIDFromHex(taskID)
    if err != nil {
        http.Error(w, "Invalid task ID", http.StatusBadRequest)
        return
    }

    var task Task
    err = json.NewDecoder(req.Body).Decode(&task)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    collection := client.Database("taskmanagement").Collection("tasks")
    filter := bson.M{"_id": objectID}
    update := bson.M{"$set": bson.M{
        "title":       task.Title,
        "description": task.Description,
        "assigned_to": task.AssignedTo,
        "status":      task.Status,
        "hours":       task.Hours,
    }}

    _, err = collection.UpdateOne(context.TODO(), filter, update)
    if err != nil {
        http.Error(w, "Failed to update task", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

func removeTask(w http.ResponseWriter, req *http.Request, taskID string) {
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
