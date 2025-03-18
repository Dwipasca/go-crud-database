package middleware

import (
	"go-crud-database/utils"
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	mu         sync.Mutex
	requests   map[string]int
	reset      map[string]time.Time
	rate       int           // max requests
	burst      int           // max burst
	resetAfter time.Duration // duration to reset the rate
}

func NewRateLimiter(rate int, burst int, resetAfter time.Duration) *RateLimiter {
	return &RateLimiter{
		requests:   make(map[string]int),
		reset:      make(map[string]time.Time),
		rate:       rate,
		burst:      burst,
		resetAfter: resetAfter,
	}
}

func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rl.mu.Lock()
		defer rl.mu.Unlock()

		ip := r.RemoteAddr // Get the IP address of the client
		now := time.Now()

		// Check if this IP has a record
		if _, exists := rl.requests[ip]; !exists {
			rl.requests[ip] = 0
			rl.reset[ip] = now.Add(rl.resetAfter)
		}

		// Check if the rate limit has expired
		if now.After(rl.reset[ip]) {
			rl.requests[ip] = 0
			rl.reset[ip] = now.Add(rl.resetAfter)
		}

		// Check if the request limit has been reached
		if rl.requests[ip] >= rl.rate+rl.burst {
			utils.WriteJson(w, http.StatusTooManyRequests, "error", nil, "Too Many Requests")
			return
		}

		// Increment the request count
		rl.requests[ip]++

		// Call the next middleware/handler
		next.ServeHTTP(w, r)
	})
}
