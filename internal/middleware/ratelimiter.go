package middleware

import (
    "net/http"
    "sync"
    "time"
)

// RateLimiter implements sliding window rate limiting per IP address
type RateLimiter struct {
    sync.Mutex
    Requests map[string][]time.Time // IP -> request timestamps
    Limit    int                   // Max requests per window
    Window   time.Duration         // Time window duration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
    return &RateLimiter{
        Requests: make(map[string][]time.Time),
        Limit:    limit,
        Window:   window,
    }
}

// Limit returns HTTP middleware that enforces rate limiting
func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ip := r.RemoteAddr
        rl.Lock()
        now := time.Now()
        
        // Clean up old requests outside the window
        times := rl.Requests[ip]
        var newTimes []time.Time
        for _, t := range times {
            if now.Sub(t) < rl.Window {
                newTimes = append(newTimes, t)
            }
        }
        rl.Requests[ip] = newTimes
        
        // Check rate limit
        if len(newTimes) >= rl.Limit {
            rl.Unlock()
            http.Error(w, "Too many requests", http.StatusTooManyRequests)
            return
        }
        
        // Record this request
        rl.Requests[ip] = append(rl.Requests[ip], now)
        rl.Unlock()
        next.ServeHTTP(w, r)
    })
}
