package main

import (
    "context"
    "fmt"
    "net/http"
    "strings"

    "github.com/dgrijalva/jwt-go"
)

//func corsMiddleware(next http.Handler) http.Handler {
//    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//        // Allow all origins for testing purposes
//        origin := r.Header.Get("Origin")

        // Check if the CORS headers are already set
//        if w.Header().Get("Access-Control-Allow-Origin") == "" {
//            w.Header().Set("Access-Control-Allow-Origin", origin)
//        }
//        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
//        w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
//        w.Header().Set("Access-Control-Allow-Credentials", "true")

        // Handle preflight requests
//        if r.Method == http.MethodOptions {
//            w.WriteHeader(http.StatusOK)
//            return
//        }

//        next.ServeHTTP(w, r)
//    })
//}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        tokenString := req.Header.Get("Authorization")
        if tokenString == "" {
            http.Error(w, "Missing token", http.StatusUnauthorized)
            return
        }

        // Remove the "Bearer " prefix from the token string
        tokenString = strings.TrimPrefix(tokenString, "Bearer ")

        // Parse and validate the JWT token
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            secretKey := []byte("your-secret-key")
            return secretKey, nil
        })

        if err != nil {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
            userID := claims["userID"].(string)
            role := claims["role"].(string)

            // Set the user ID and role in the request context
            ctx := context.WithValue(req.Context(), "userID", userID)
            ctx = context.WithValue(ctx, "role", role)
            req = req.WithContext(ctx)

            next(w, req)
        } else {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
        }
    }
}


func adminMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        tokenString := req.Header.Get("Authorization")[7:] // Assuming 'Bearer ' prefix
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            // Ensure the token algorithm is what you expect:
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            return []byte("your-secret-key"), nil
        })

        if err != nil {
            // If there's an error parsing the token, return an unauthorized error.
            http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
            return
        }

        if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
            if role, ok := claims["role"].(string); ok && role == "admin" {
                next(w, req)
                return
            }
        }
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
    }
}



func isAdmin(req *http.Request) bool {
    role := req.Context().Value("role")
    return role == "admin"
}
