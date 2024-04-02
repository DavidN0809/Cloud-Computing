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
    "golang.org/x/crypto/bcrypt"
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
    mux.HandleFunc("/", handleUsers)
    mux.HandleFunc("/list", listUsers)
    mux.HandleFunc("/login", loginUser)

    // Start the server
    log.Println("User Service listening on port 8001...")
    log.Fatal(http.ListenAndServe(":8001", mux))
}

type User struct {
    ID       primitive.ObjectID `bson:"_id" json:"id"`
    Username string             `bson:"username" json:"username"`
    Email    string             `bson:"email" json:"email"`
    Password string             `bson:"password" json:"password"`
    Role     string             `bson:"role" json:"role"`
}

func handleUsers(w http.ResponseWriter, req *http.Request) {
    switch req.Method {
    case http.MethodPost:
        createUser(w, req)
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

func createUser(w http.ResponseWriter, req *http.Request) {
    var user User
    err := json.NewDecoder(req.Body).Decode(&user)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Hash the user password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "Failed to hash password", http.StatusInternalServerError)
        return
    }
    user.Password = string(hashedPassword)

    collection := client.Database("user").Collection("users")
    user.ID = primitive.NewObjectID()
    _, err = collection.InsertOne(context.TODO(), user)
    if err != nil {
        http.Error(w, "Failed to create user", http.StatusInternalServerError)
        return
    }

    // Remove the password from the response
    user.Password = ""

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

func loginUser(w http.ResponseWriter, req *http.Request) {
    var credentials struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    err := json.NewDecoder(req.Body).Decode(&credentials)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    collection := client.Database("user").Collection("users")
    filter := bson.M{"email": credentials.Email}

    var user User
    err = collection.FindOne(context.TODO(), filter).Decode(&user)
    if err != nil {
        log.Println("User not found:", err)
        http.Error(w, "Invalid email or password", http.StatusUnauthorized)
        return
    }

    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
    if err != nil {
        log.Println("Password comparison failed:", err)
        http.Error(w, "Invalid email or password", http.StatusUnauthorized)
        return
    }

    // Remove the password from the response
    user.Password = ""

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

func listUsers(w http.ResponseWriter, req *http.Request) {
    collection := client.Database("user").Collection("users")
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
