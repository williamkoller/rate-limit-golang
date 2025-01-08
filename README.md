# How implement rate limiter in Go

### Structure

```bash
.
├── cmd
│   └── main.go
├── go.mod
├── go.sum
└── README.md
```

### Installation

```bash
go mod tidy
```

### Usage

```go
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
```

### Run

```bash

go run cmd/main.go
```

### Test

```bash
for i in {1..10}; do curl -i http://localhost:8080/; echo ""; done
```

### Output

```bash
HTTP/1.1 200 OK
Date: Wed, 08 Jan 2025 19:13:04 GMT
Content-Length: 25
Content-Type: text/plain; charset=utf-8

Requisição bem-sucedida
HTTP/1.1 200 OK
Date: Wed, 08 Jan 2025 19:13:04 GMT
Content-Length: 25
Content-Type: text/plain; charset=utf-8

Requisição bem-sucedida
HTTP/1.1 429 Too Many Requests
Content-Type: application/json
Date: Wed, 08 Jan 2025 19:13:04 GMT
Content-Length: 30

{"error":"Too many requests"}

HTTP/1.1 429 Too Many Requests
Content-Type: application/json
Date: Wed, 08 Jan 2025 19:13:04 GMT
Content-Length: 30

{"error":"Too many requests"}

HTTP/1.1 429 Too Many Requests
Content-Type: application/json
Date: Wed, 08 Jan 2025 19:13:04 GMT
Content-Length: 30

{"error":"Too many requests"}

HTTP/1.1 429 Too Many Requests
Content-Type: application/json
Date: Wed, 08 Jan 2025 19:13:04 GMT
Content-Length: 30

{"error":"Too many requests"}

HTTP/1.1 429 Too Many Requests
Content-Type: application/json
Date: Wed, 08 Jan 2025 19:13:04 GMT
Content-Length: 30

{"error":"Too many requests"}

HTTP/1.1 429 Too Many Requests
Content-Type: application/json
Date: Wed, 08 Jan 2025 19:13:04 GMT
Content-Length: 30

{"error":"Too many requests"}

HTTP/1.1 429 Too Many Requests
Content-Type: application/json
Date: Wed, 08 Jan 2025 19:13:04 GMT
Content-Length: 30

{"error":"Too many requests"}

HTTP/1.1 429 Too Many Requests
Content-Type: application/json
Date: Wed, 08 Jan 2025 19:13:04 GMT
Content-Length: 30

{"error":"Too many requests"}

```
