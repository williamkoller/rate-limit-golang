package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/go-chi/chi/v5"
	"golang.org/x/time/rate"
)

func main() {
    r := chi.NewRouter()

    r.Use(RateLimiter(rate.Limit(1), 2))

    r.Get("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Requisição bem-sucedida"))
    })

    fmt.Println("Server running on port 8080")
    http.ListenAndServe(":8080", r)
}

func RateLimiter(limit rate.Limit, burst int) func(next http.Handler) http.Handler {
    var mu sync.Mutex
    limiterMap := make(map[string]*rate.Limiter)

    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ip := strings.Split(r.RemoteAddr, ":")[0]

            mu.Lock()
            limiter, exists := limiterMap[ip]
            if !exists {
                limiter = rate.NewLimiter(limit, burst)
                limiterMap[ip] = limiter
            }
            mu.Unlock()

            if !limiter.Allow() {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusTooManyRequests)
                json.NewEncoder(w).Encode(map[string]string{"error": "Too many requests"})
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}