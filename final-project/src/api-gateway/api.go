package main

import (
    "encoding/json"
    "bytes"
    "io"
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
)

func main() {
    // Create a new HTTP server
    mux := http.NewServeMux()

    // User Service
    userServiceURL, _ := url.Parse("http://user-service:8001")
    userServiceProxy := httputil.NewSingleHostReverseProxy(userServiceURL)
    mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
        userServiceProxy.ServeHTTP(w, r)
    })

    // Task Service
    mux.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
        forwardRequest(w, r, "http://task-service:8002")
    })

    // Billing Service
    mux.HandleFunc("/billings/", func(w http.ResponseWriter, r *http.Request) {
        forwardRequest(w, r, "http://billing-service:8003")
    })

    // User Types
    mux.HandleFunc("/auth/login", handleLogin)
    mux.HandleFunc("/auth/register", handleRegister)

    // Start the server
    log.Println("API Gateway listening on port 8000...")
    log.Fatal(http.ListenAndServe(":8000", mux))
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
    var user struct {
        Username string `json:"username"`
        Email    string `json:"email"`
        Password string `json:"password"`
        Role     string `json:"role"`
    }

    // Log the raw request body
    body, err := io.ReadAll(r.Body)
    if err != nil {
        log.Println("Failed to read request body:", err)
        http.Error(w, "Failed to read request body", http.StatusInternalServerError)
        return
    }
    log.Printf("Request body: %s", string(body))

    err = json.Unmarshal(body, &user)
    if err != nil {
        log.Println("Failed to parse request body:", err)
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Validate user role
    if user.Role != "admin" && user.Role != "regular" {
        log.Println("Invalid user role")
        http.Error(w, "Invalid user role", http.StatusBadRequest)
        return
    }

    // Forward the request to the user service
    userServiceURL, _ := url.Parse("http://user-service:8001")
    userServiceProxy := httputil.NewSingleHostReverseProxy(userServiceURL)
    r.URL.Path = "/users/create"

    // Forward the request body
    r.Body = io.NopCloser(bytes.NewBuffer(body))

    userServiceProxy.ServeHTTP(w, r)
}


func handleLogin(w http.ResponseWriter, r *http.Request) {
    log.Println("Received request to login user")

    // Forward the request to the user service
    userServiceURL, _ := url.Parse("http://user-service:8001")
    userServiceProxy := httputil.NewSingleHostReverseProxy(userServiceURL)
    r.URL.Path = "/users/login"
    userServiceProxy.ServeHTTP(w, r)
}


func forwardRequest(w http.ResponseWriter, r *http.Request, serviceURL string) {
    // Forward the JWT token to the downstream service
    token := r.Header.Get("Authorization")
    if token != "" {
        r.Header.Set("Authorization", token)
    }

    url, _ := url.Parse(serviceURL)
    proxy := httputil.NewSingleHostReverseProxy(url)
    proxy.ServeHTTP(w, r)
}
