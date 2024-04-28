package main

import (
    "context"
    "fmt"
    "net/http"
    "strings"

    "github.com/dgrijalva/jwt-go"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8003")
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
        role := req.Context().Value("role")
        if role != "admin" {
            http.Error(w, "Unauthorized", http.StatusForbidden)
            return
        }
        next(w, req)
    }
}

func isAdmin(req *http.Request) bool {
    role := req.Context().Value("role")
    return role == "admin"
}
