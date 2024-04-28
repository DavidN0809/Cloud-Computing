package main

import (
    "bytes"
    "encoding/json"
    "io"
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
)

func main() {
    mux := http.NewServeMux()

    // Wrap handlers with CORS middleware using mux.Handle
    mux.Handle("/users/", corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        userServiceURL, _ := url.Parse("http://user-service:8001")
        userServiceProxy := httputil.NewSingleHostReverseProxy(userServiceURL)
        userServiceProxy.ServeHTTP(w, r)
    })))

    mux.Handle("/tasks/", corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        forwardRequest(w, r, "http://task-service:8002")
    })))

    mux.Handle("/billings/", corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        forwardRequest(w, r, "http://billing-service:8003")
    })))

    // Note the change to mux.Handle here as well
    mux.Handle("/auth/login", corsMiddleware(http.HandlerFunc(handleLogin)))
    mux.Handle("/auth/register", corsMiddleware(http.HandlerFunc(handleRegister)))

    log.Println("API Gateway listening on port 8000...")
    log.Fatal(http.ListenAndServe(":8000", mux))
}

func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Allow all origins for testing purposes
        origin := r.Header.Get("Origin")

        // Check if the CORS headers are already set
        if w.Header().Get("Access-Control-Allow-Origin") == "" {
            w.Header().Set("Access-Control-Allow-Origin", origin)
        }
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
        w.Header().Set("Access-Control-Allow-Credentials", "true")

        // Handle preflight requests
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
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


    log.Printf("Before parsing JSON, body is: %s", string(body))
    err = json.Unmarshal(body, &user)
    log.Printf("Parsed user: %+v", user)
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
