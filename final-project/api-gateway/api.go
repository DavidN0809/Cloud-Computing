package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
//	"strings"
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
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
	    userServiceProxy.ServeHTTP(w,r)
	})

//        mux.HandleFunc("/users",func(w http.ResponseWriter, r *http.Request){
//             userServiceProxy.ServeHTTP(w,r)
//         r.URL.Path = strings.TrimPrefix(r.URL.Path, "/users")
 //        userServiceProxy.ServeHTTP(w,r)
//         })


//	mux.Handle("/users", http.StripPrefix("/users", userServiceProxy))
//	mux.Handle("/users/", http.StripPrefix("/users/", userServiceProxy))

	// Task Service
	taskServiceURL, _ := url.Parse("http://task-service:8002")
	taskServiceProxy := httputil.NewSingleHostReverseProxy(taskServiceURL)
	mux.Handle("/tasks", http.StripPrefix("/tasks", taskServiceProxy))
	mux.Handle("/tasks/", http.StripPrefix("/tasks/", taskServiceProxy))

	// Billing Service
	billingServiceURL, _ := url.Parse("http://billing-service:8003")
	billingServiceProxy := httputil.NewSingleHostReverseProxy(billingServiceURL)
	mux.Handle("/billings", http.StripPrefix("/billings", billingServiceProxy))
	mux.Handle("/billings/", http.StripPrefix("/billings/", billingServiceProxy))

	// Start the server
	log.Println("API Gateway listening on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", mux))
}
