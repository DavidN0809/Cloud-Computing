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
    mux.HandleFunc("/tasks/remove/", removeTask)
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
}

func createTask(w http.ResponseWriter, req *http.Request) {
    if req.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

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

func getTask(w http.ResponseWriter, req *http.Request) {
    if req.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    taskID := req.URL.Path[len("/tasks/get/"):]
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
