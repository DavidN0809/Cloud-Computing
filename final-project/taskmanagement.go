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

var client *mongo.Client

func main() {
    // Create a new MongoDB client
    var err error
    client, err = mongo.NewClient(options.Client().ApplyURI("mongodb://mongodb-final:27017"))
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
    mux.HandleFunc("/users/create", createUser)
    mux.HandleFunc("/users/get", getUser)

    // Task endpoints
    mux.HandleFunc("/tasks/create", createTask)
    mux.HandleFunc("/tasks/list", listTasks)
    mux.HandleFunc("/tasks/update", updateTask)

    // Billing endpoints
    mux.HandleFunc("/billings/create", createBilling)
    mux.HandleFunc("/billings/get", getBilling)

    // Start the server
    log.Println("Server listening on port 8000...")
    log.Fatal(http.ListenAndServe(":8000", mux))
}

type User struct {
    ID       primitive.ObjectID `bson:"_id"`
    Username string             `bson:"username"`
    Email    string             `bson:"email"`
    Password string             `bson:"password"`
}

func createUser(w http.ResponseWriter, req *http.Request) {
    username := req.URL.Query().Get("username")
    email := req.URL.Query().Get("email")
    password := req.URL.Query().Get("password")

    collection := client.Database("taskmanagement").Collection("users")
    newUser := User{
        ID:       primitive.NewObjectID(),
        Username: username,
        Email:    email,
        Password: password,
    }
    _, err := collection.InsertOne(context.TODO(), newUser)
    if err != nil {
        http.Error(w, "Failed to create user", http.StatusInternalServerError)
        return
    }

    fmt.Fprintf(w, "Created user: %s\n", username)
}

func getUser(w http.ResponseWriter, req *http.Request) {
    username := req.URL.Query().Get("username")

    collection := client.Database("taskmanagement").Collection("users")
    filter := bson.M{"username": username}

    var user User
    err := collection.FindOne(context.TODO(), filter).Decode(&user)
    if err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    fmt.Fprintf(w, "User: %s, Email: %s\n", user.Username, user.Email)
}

type Task struct {
    ID          primitive.ObjectID `bson:"_id"`
    Title       string             `bson:"title"`
    Description string             `bson:"description"`
    AssignedTo  primitive.ObjectID `bson:"assigned_to"`
    Status      string             `bson:"status"`
    Hours       float64            `bson:"hours"`
}

func createTask(w http.ResponseWriter, req *http.Request) {
    title := req.URL.Query().Get("title")
    description := req.URL.Query().Get("description")
    assignedTo := req.URL.Query().Get("assigned_to")
    status := req.URL.Query().Get("status")
    hours, _ := strconv.ParseFloat(req.URL.Query().Get("hours"), 64)

    assignedToID, err := primitive.ObjectIDFromHex(assignedTo)
    if err != nil {
        http.Error(w, "Invalid assigned_to ID", http.StatusBadRequest)
        return
    }

    collection := client.Database("taskmanagement").Collection("tasks")
    newTask := Task{
        ID:          primitive.NewObjectID(),
        Title:       title,
        Description: description,
        AssignedTo:  assignedToID,
        Status:      status,
        Hours:       hours,
    }
    _, err = collection.InsertOne(context.TODO(), newTask)
    if err != nil {
        http.Error(w, "Failed to create task", http.StatusInternalServerError)
        return
    }

    fmt.Fprintf(w, "Created task: %s\n", title)
}

func listTasks(w http.ResponseWriter, req *http.Request) {
    collection := client.Database("taskmanagement").Collection("tasks")
    cursor, err := collection.Find(context.TODO(), bson.M{})
    if err != nil {
        http.Error(w, "Failed to list tasks", http.StatusInternalServerError)
        return
    }
    defer cursor.Close(context.Background())

    for cursor.Next(context.Background()) {
        var task Task
        if err := cursor.Decode(&task); err != nil {
            http.Error(w, "Failed to decode task", http.StatusInternalServerError)
            return
        }
        fmt.Fprintf(w, "Task: %s, Assigned To: %s, Status: %s, Hours: %.2f\n",
            task.Title, task.AssignedTo.Hex(), task.Status, task.Hours)
    }
}

func updateTask(w http.ResponseWriter, req *http.Request) {
    taskID := req.URL.Query().Get("task_id")
    status := req.URL.Query().Get("status")
    hours, _ := strconv.ParseFloat(req.URL.Query().Get("hours"), 64)

    taskObjectID, err := primitive.ObjectIDFromHex(taskID)
    if err != nil {
        http.Error(w, "Invalid task ID", http.StatusBadRequest)
        return
    }

    collection := client.Database("taskmanagement").Collection("tasks")
    filter := bson.M{"_id": taskObjectID}
    update := bson.M{"$set": bson.M{"status": status, "hours": hours}}

    result, err := collection.UpdateOne(context.TODO(), filter, update)
    if err != nil {
        http.Error(w, "Failed to update task", http.StatusInternalServerError)
        return
    }
    if result.ModifiedCount == 0 {
        http.Error(w, "Task not found", http.StatusNotFound)
        return
    }

    fmt.Fprintf(w, "Updated task: %s\n", taskID)
}

type Billing struct {
    ID     primitive.ObjectID `bson:"_id"`
    UserID primitive.ObjectID `bson:"user_id"`
    TaskID primitive.ObjectID `bson:"task_id"`
    Hours  float64            `bson:"hours"`
    Amount float64            `bson:"amount"`
}

func createBilling(w http.ResponseWriter, req *http.Request) {
    userID := req.URL.Query().Get("user_id")
    taskID := req.URL.Query().Get("task_id")
    hours, _ := strconv.ParseFloat(req.URL.Query().Get("hours"), 64)
    amount, _ := strconv.ParseFloat(req.URL.Query().Get("amount"), 64)

    userObjectID, err := primitive.ObjectIDFromHex(userID)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    taskObjectID, err := primitive.ObjectIDFromHex(taskID)
    if err != nil {
        http.Error(w, "Invalid task ID", http.StatusBadRequest)
        return
    }

    collection := client.Database("taskmanagement").Collection("billings")
    newBilling := Billing{
        ID:     primitive.NewObjectID(),
        UserID: userObjectID,
        TaskID: taskObjectID,
        Hours:  hours,
        Amount: amount,
    }
    _, err = collection.InsertOne(context.TODO(), newBilling)
    if err != nil {
        http.Error(w, "Failed to create billing", http.StatusInternalServerError)
        return
    }

    fmt.Fprintf(w, "Created billing for user: %s, task: %s\n", userID, taskID)
}

func getBilling(w http.ResponseWriter, req *http.Request) {
    userID := req.URL.Query().Get("user_id")

    userObjectID, err := primitive.ObjectIDFromHex(userID)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    collection := client.Database("taskmanagement").Collection("billings")
    filter := bson.M{"user_id": userObjectID}

    cursor, err := collection.Find(context.TODO(), filter)
    if err != nil {
        http.Error(w, "Failed to retrieve billing", http.StatusInternalServerError)
        return
    }
    defer cursor.Close(context.Background())

    for cursor.Next(context.Background()) {
        var billing Billing
        if err := cursor.Decode(&billing); err != nil {
            http.Error(w, "Failed to decode billing", http.StatusInternalServerError)
            return
        }
        fmt.Fprintf(w, "Billing: User ID: %s, Task ID: %s, Hours: %.2f, Amount: %.2f\n",
            billing.UserID.Hex(), billing.TaskID.Hex(), billing.Hours, billing.Amount)
    }
}
