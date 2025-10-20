package cache

import (
    "sync"
    "time"
)

// RateCache stores conversion rates with TTL expiration
type RateCache struct {
    sync.RWMutex
    Data      map[string]CachedRate // key -> cached rate data
    ExpirySec int                   // TTL in seconds
}

// CachedRate holds rate value and timestamp
type CachedRate struct {
    Rate      float64
    Timestamp time.Time
}

// NewRateCache creates a new cache with specified TTL
func NewRateCache(ttl int) *RateCache {
    return &RateCache{
        Data:      make(map[string]CachedRate),
        ExpirySec: ttl,
    }
}

// Get retrieves cached rate if not expired
func (c *RateCache) Get(key string) (float64, bool) {
    c.RLock()
    defer c.RUnlock()
    val, ok := c.Data[key]
    if !ok || time.Since(val.Timestamp) > time.Duration(c.ExpirySec)*time.Second {
        return 0, false
    }
    return val.Rate, true
}

// Set stores a rate with current timestamp
func (c *RateCache) Set(key string, rate float64) {
    c.Lock()
    defer c.Unlock()
    c.Data[key] = CachedRate{Rate: rate, Timestamp: time.Now()}
}
