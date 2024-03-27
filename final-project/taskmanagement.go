package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongodbEndpoint = "mongodb://mongodb-final:27017"
	dbName          = "taskmanagement"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username"`
	Email    string             `bson:"email"`
}

type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Title       string             `bson:"title"`
	Description string             `bson:"description"`
	AssignedTo  primitive.ObjectID `bson:"assignedTo"`
	Status      string             `bson:"status"`
	Hours       float64            `bson:"hours"`
}

type UserService struct {
	collection *mongo.Collection
}

func NewUserService(client *mongo.Client) *UserService {
	collection := client.Database(dbName).Collection("users")
	return &UserService{collection: collection}
}

func (s *UserService) CreateUser(ctx context.Context, user User) error {
	_, err := s.collection.InsertOne(ctx, user)
	return err
}

func (db *database) listUsers(w http.ResponseWriter, req *http.Request) {
    collection := db.client.Database(dbName).Collection("users")
    cursor, err := collection.Find(context.TODO(), bson.M{})
    if err != nil {
        http.Error(w, "Failed to list users", http.StatusInternalServerError)
        return
    }
    defer cursor.Close(context.Background())

    var users []User
    if err := cursor.All(context.TODO(), &users); err != nil {
        http.Error(w, "Failed to decode users", http.StatusInternalServerError)
        return
    }

    jsonData, err := json.Marshal(users)
    if err != nil {
        http.Error(w, "Failed to encode users as JSON", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(jsonData)
}

type TaskService struct {
	collection *mongo.Collection
}

func NewTaskService(client *mongo.Client) *TaskService {
	collection := client.Database(dbName).Collection("tasks")
	return &TaskService{collection: collection}
}

func (s *TaskService) CreateTask(ctx context.Context, task Task) error {
	_, err := s.collection.InsertOne(ctx, task)
	return err
}

func (s *TaskService) ListTasks(ctx context.Context) ([]Task, error) {
	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []Task
	for cursor.Next(ctx) {
		var task Task
		if err := cursor.Decode(&task); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, taskID primitive.ObjectID, update bson.M) error {
	_, err := s.collection.UpdateOne(ctx, bson.M{"_id": taskID}, update)
	return err
}

type BillingService struct {
	collection *mongo.Collection
}

func NewBillingService(client *mongo.Client) *BillingService {
	collection := client.Database(dbName).Collection("billing")
	return &BillingService{collection: collection}
}

func (s *BillingService) RecordHours(ctx context.Context, taskID primitive.ObjectID, hours float64) error {
	_, err := s.collection.InsertOne(ctx, bson.M{"taskId": taskID, "hours": hours})
	return err
}

func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongodbEndpoint))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	userService := NewUserService(client)
	taskService := NewTaskService(client)
	billingService := NewBillingService(client)

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var user User
			if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if err := userService.CreateUser(ctx, user); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var task Task
			if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if err := taskService.CreateTask(ctx, task); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
		} else if r.Method == http.MethodGet {
			tasks, err := taskService.ListTasks(ctx)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if err := json.NewEncoder(w).Encode(tasks); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/tasks/update", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			var update struct {
				TaskID primitive.ObjectID `json:"taskId"`
				Status string             `json:"status"`
			}
			if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if err := taskService.UpdateTask(ctx, update.TaskID, bson.M{"$set": bson.M{"status": update.Status}}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/billing", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var record struct {
				TaskID primitive.ObjectID `json:"taskId"`
				Hours  float64            `json:"hours"`
			}
			if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if err := billingService.RecordHours(ctx, record.TaskID, record.Hours); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server started on port 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
