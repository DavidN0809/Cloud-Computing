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
        "github.com/dgrijalva/jwt-go"
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

	// Check if the database and collection exist, create them if they don't
	err = ensureDatabaseAndCollection(client)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new HTTP server
	mux := http.NewServeMux()

	// User endpoints
	mux.Handle("/users/list", authMiddleware(adminMiddleware(http.HandlerFunc(listUsers))))
	mux.Handle("/users/create", http.HandlerFunc(createUser))
mux.Handle("/users/get/", authMiddleware(adminMiddleware(http.HandlerFunc(getUser))))
mux.Handle("/users/update/", authMiddleware(adminMiddleware(http.HandlerFunc(updateUser))))
mux.Handle("/users/remove/", authMiddleware(adminMiddleware(http.HandlerFunc(removeUser))))
mux.Handle("/users/delete-all", http.HandlerFunc(deleteAllUsers))
mux.Handle("/users/login", http.HandlerFunc(loginUser))

	// Start the server
	log.Println("User Service listening on port 8001...")
	log.Fatal(http.ListenAndServe(":8001", mux))
}

func ensureDatabaseAndCollection(client *mongo.Client) error {
	dbName := "user"
	collectionName := "users"

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

type User struct {
	ID       primitive.ObjectID `bson:"_id" json:"id"`
	Username string             `bson:"username" json:"username"`
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password" json:"password"`
        Role     string             `bson:"role" json:"role"`
}

func createUser(w http.ResponseWriter, req *http.Request) {
    log.Println("Starting to create user")

    var user User
    err := json.NewDecoder(req.Body).Decode(&user)
    if err != nil {
        log.Printf("Invalid request body: %v", err)
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Set default role to "regular" if not specified
    if user.Role == "" {
        user.Role = "regular"
    }

    log.Printf("Attempting to insert user: %+v", user)

    collection := client.Database("user").Collection("users")

    // Check if a user with the same username already exists
    usernameFilter := bson.M{"username": user.Username}
    existingUserByUsername := &User{}
    err = collection.FindOne(context.TODO(), usernameFilter).Decode(existingUserByUsername)
    if err == nil {
        log.Printf("User with the same username already exists: %+v", existingUserByUsername)
        http.Error(w, "User with the same username already exists", http.StatusConflict)
        return
    }

    // Check if a user with the same email already exists
    emailFilter := bson.M{"email": user.Email}
    existingUserByEmail := &User{}
    err = collection.FindOne(context.TODO(), emailFilter).Decode(existingUserByEmail)
    if err == nil {
        log.Printf("User with the same email already exists: %+v", existingUserByEmail)
        http.Error(w, "User with the same email already exists", http.StatusConflict)
        return
    }

    user.ID = primitive.NewObjectID()
    _, err = collection.InsertOne(context.TODO(), user)
    if err != nil {
        log.Printf("Failed to create user: %v", err)
        http.Error(w, "Failed to create user", http.StatusInternalServerError)
        return
    }

    log.Printf("User created successfully: %+v", user)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}


func loginUser(w http.ResponseWriter, req *http.Request) {
    log.Println("Received request to login user")

    if req.Method != http.MethodPost {
        log.Println("Invalid request method for user login")
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var credentials struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }

    err := json.NewDecoder(req.Body).Decode(&credentials)
    if err != nil {
        log.Println("Failed to decode request body:", err)
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    log.Printf("Login attempt for username: %s", credentials.Username)

    collection := client.Database("user").Collection("users")
    filter := bson.M{"username": credentials.Username, "password": credentials.Password}

    var user User
    err = collection.FindOne(context.TODO(), filter).Decode(&user)
    if err != nil {
        log.Println("Invalid username or password")
        http.Error(w, "Invalid username or password", http.StatusUnauthorized)
        return
    }

    log.Printf("User logged in successfully: %+v", user)


    // Generate JWT token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "userID": user.ID.Hex(),
        "role":   user.Role,
        "exp":    time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
    })

    // Sign the token with a secret key
    secretKey := []byte("your-secret-key")
    tokenString, err := token.SignedString(secretKey)
    if err != nil {
        log.Println("Failed to generate JWT token:", err)
        http.Error(w, "Failed to generate token", http.StatusInternalServerError)
        return
    }

    // Send the token in the response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "token": tokenString,
    })

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK) // Explicitly set the 200 OK status
    json.NewEncoder(w).Encode(user)
}

func getUser(w http.ResponseWriter, req *http.Request) {
	log.Println("Received request to get user")

	if req.Method != http.MethodGet {
		log.Println("Invalid request method")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := req.URL.Path[len("/users/get/"):]
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Printf("Invalid user ID: %v", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	log.Printf("Getting user with ID: %s", userID)

	collection := client.Database("user").Collection("users")
	filter := bson.M{"_id": objectID}

	var user User
	err = collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		log.Printf("User not found: %v", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	log.Printf("User found: %+v", user)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(user)
}

func updateUser(w http.ResponseWriter, req *http.Request) {
	log.Println("Received request to update user")

	if req.Method != http.MethodPut {
		log.Println("Invalid request method")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := req.URL.Path[len("/users/update/"):]
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Printf("Invalid user ID: %v", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	log.Printf("Updating user with ID: %s", userID)

	var user User
	err = json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		log.Printf("Invalid request body: %v", err)
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
		log.Printf("Failed to update user: %v", err)
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	log.Printf("User updated successfully: %+v", user)
	w.WriteHeader(http.StatusNoContent)
}

func removeUser(w http.ResponseWriter, req *http.Request) {
	log.Println("Received request to remove user")

	if req.Method != http.MethodDelete {
		log.Println("Invalid request method")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := req.URL.Path[len("/users/remove/"):]
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Printf("Invalid user ID: %v", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	log.Printf("Removing user with ID: %s", userID)

	collection := client.Database("user").Collection("users")
	filter := bson.M{"_id": objectID}

	_, err = collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Printf("Failed to remove user: %v", err)
		http.Error(w, "Failed to remove user", http.StatusInternalServerError)
		return
	}

	log.Printf("User removed successfully: %s", userID)
	w.WriteHeader(http.StatusNoContent)
}

func listUsers(w http.ResponseWriter, req *http.Request) {
	log.Println("Received request to list users")

	if req.Method != http.MethodGet {
		log.Println("Invalid request method")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	collection := client.Database("user").Collection("users")
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Printf("Failed to list users: %v", err)
		http.Error(w, "Failed to list users", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var users []User
	err = cursor.All(context.Background(), &users)
	if err != nil {
		log.Printf("Failed to decode users: %v", err)
		http.Error(w, "Failed to decode users", http.StatusInternalServerError)
		return
	}

	log.Printf("Users listed successfully: %+v", users)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func deleteAllUsers(w http.ResponseWriter, req *http.Request) {
	log.Println("Received request to delete all users")

	if req.Method != http.MethodDelete {
		log.Println("Invalid request method")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	collection := client.Database("user").Collection("users")
	_, err := collection.DeleteMany(context.TODO(), bson.M{})
	if err != nil {
		log.Printf("Failed to delete users: %v", err)
		http.Error(w, "Failed to delete users", http.StatusInternalServerError)
		return
	}

	log.Println("All users deleted successfully")
	w.WriteHeader(http.StatusNoContent)
}
