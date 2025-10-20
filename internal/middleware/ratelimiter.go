package middleware

import (
    "net/http"
    "sync"
    "time"
)

type RateLimiter struct {
    sync.Mutex
    Requests map[string][]time.Time
    Limit    int
    Window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
    return &RateLimiter{
        Requests: make(map[string][]time.Time),
        Limit:    limit,
        Window:   window,
    }
}

func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ip := r.RemoteAddr
        rl.Lock()
        now := time.Now()
        times := rl.Requests[ip]
        var newTimes []time.Time
        for _, t := range times {
            if now.Sub(t) < rl.Window {
                newTimes = append(newTimes, t)
            }
        }
        rl.Requests[ip] = newTimes
        if len(newTimes) >= rl.Limit {
            rl.Unlock()
            http.Error(w, "Too many requests", http.StatusTooManyRequests)
            return
        }
        rl.Requests[ip] = append(rl.Requests[ip], now)
        rl.Unlock()
        next.ServeHTTP(w, r)
    })
}
