package main

import (
    "encoding/json"
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
    "path"
)

func main() {
    // Create a new HTTP server
    mux := http.NewServeMux()

    // User Service
    userServiceURL, _ := url.Parse("http://user-service:8001")
    userServiceProxy := httputil.NewSingleHostReverseProxy(userServiceURL)
    userServiceProxy.Director = func(req *http.Request) {
        req.URL.Scheme = userServiceURL.Scheme
        req.URL.Host = userServiceURL.Host
        req.URL.Path = path.Join(userServiceURL.Path, req.URL.Path)
    }
    mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
        userServiceProxy.ServeHTTP(w, r)
    })

    // Task Service
    taskServiceURL, _ := url.Parse("http://task-service:8002")
    taskServiceProxy := httputil.NewSingleHostReverseProxy(taskServiceURL)
    mux.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
        taskServiceProxy.ServeHTTP(w, r)
    })

    // Billing Service
    billingServiceURL, _ := url.Parse("http://billing-service:8003")
    billingServiceProxy := httputil.NewSingleHostReverseProxy(billingServiceURL)
    mux.HandleFunc("/billings/", func(w http.ResponseWriter, r *http.Request) {
        billingServiceProxy.ServeHTTP(w, r)
    })

    // User Types
    mux.HandleFunc("/auth/login", handleLogin)
    mux.HandleFunc("/auth/register", handleRegister)

    // Start the server
    log.Println("API Gateway listening on port 8000...")
    log.Fatal(http.ListenAndServe(":8000", mux))
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
    // Forward the request to the user service
    userServiceURL, _ := url.Parse("http://user-service:8001")
    userServiceProxy := httputil.NewSingleHostReverseProxy(userServiceURL)
    r.URL.Path = "/login"
    userServiceProxy.ServeHTTP(w, r)
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
    var user struct {
        Username string `json:"username"`
        Email    string `json:"email"`
        Password string `json:"password"`
        Role     string `json:"role"`
    }
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Validate user role
    if user.Role != "admin" && user.Role != "regular" {
        http.Error(w, "Invalid user role", http.StatusBadRequest)
        return
    }

    // Forward the request to the user service
    userServiceURL, _ := url.Parse("http://user-service:8001")
    userServiceProxy := httputil.NewSingleHostReverseProxy(userServiceURL)
    userServiceProxy.ServeHTTP(w, r)
}
