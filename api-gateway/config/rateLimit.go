package config

import (
	"api-gateway/webResponse"
	"github.com/labstack/echo/v4"
	"net/http"
	"sync"
	"time"
)

type RateLimitConfig struct {
	RequestTimeout time.Duration
	RateLimit      int
	Enable         bool
}

type NewRateLimiterStruct struct {
	Config     *RateLimitConfig
	Requests   map[string]int
	LastAccess map[string]time.Time
	mu         sync.Mutex
}

// LoadRateLimitConfig loads rate limit configuration
func LoadRateLimitConfig() *RateLimitConfig {
	timeout := 10 * time.Second
	rateLimit := 100

	return &RateLimitConfig{
		RequestTimeout: timeout,
		RateLimit:      rateLimit,
		Enable:         true,
	}
}
func NewRateLimiter(config *RateLimitConfig) *NewRateLimiterStruct {
	return &NewRateLimiterStruct{
		Config:     config,
		Requests:   make(map[string]int),
		LastAccess: make(map[string]time.Time),
	}
}

func (l *NewRateLimiterStruct) Allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()

	if last, found := l.LastAccess[ip]; found && now.Sub(last) > time.Minute {
		l.Requests[ip] = 0
		l.LastAccess[ip] = now
	}

	l.Requests[ip]++
	if l.Requests[ip] > l.Config.RateLimit {
		return false
	}

	l.LastAccess[ip] = now
	return true
}

// CheckRateLimit checks if the rate limit is enabled and if the request is allowed
func CheckRateLimit(c echo.Context) error {
	rateLimit := LoadRateLimitConfig()
	limiter := NewRateLimiter(rateLimit)

	ip := c.RealIP()
	if rateLimit.Enable && !limiter.Allow(ip) {
		return webResponse.ResponseJson(c, http.StatusTooManyRequests, nil, "Too many requests")
	}
	
	return nil
}
