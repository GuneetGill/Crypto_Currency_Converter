package cache

import (
    "sync"
    "time"
)

type RateCache struct {
    sync.RWMutex
    Data      map[string]CachedRate
    ExpirySec int
}

type CachedRate struct {
    Rate      float64
    Timestamp time.Time
}

func NewRateCache(ttl int) *RateCache {
    return &RateCache{
        Data:      make(map[string]CachedRate),
        ExpirySec: ttl,
    }
}

func (c *RateCache) Get(key string) (float64, bool) {
    c.RLock()
    defer c.RUnlock()
    val, ok := c.Data[key]
    if !ok || time.Since(val.Timestamp) > time.Duration(c.ExpirySec)*time.Second {
        return 0, false
    }
    return val.Rate, true
}

func (c *RateCache) Set(key string, rate float64) {
    c.Lock()
    defer c.Unlock()
    c.Data[key] = CachedRate{Rate: rate, Timestamp: time.Now()}
}
