package util

import (
	"log"
	"sync"
	"time"
)

type IPRateLimiter struct {
	ips        map[string]*TokenBucket
	mu         sync.Mutex
	max        int
	refillRate time.Duration
}

func NewIPRateLimiter(max int, refillRate time.Duration) *IPRateLimiter {
	return &IPRateLimiter{
		ips:        make(map[string]*TokenBucket),
		max:        max,
		refillRate: refillRate,
	}
}

func (i *IPRateLimiter) GetLimiter(ip string) *TokenBucket {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.ips[ip]
	if !exists {
		log.Printf("Creating new limiter for IP: %s", ip)
		limiter = NewTokenBucket(i.max, i.refillRate)
		i.ips[ip] = limiter
	}

	return limiter
}

func NewTokenBucket(maxTokens int, refillRate time.Duration) *TokenBucket {
	return &TokenBucket{
		tokens:     maxTokens, // Start with max tokens to allow initial burst
		maxTokens:  maxTokens,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

type TokenBucket struct {
	tokens     int
	maxTokens  int
	refillRate time.Duration
	lastRefill time.Time
	mu         sync.Mutex
}

func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefill)
	tokensToAdd := int(elapsed / tb.refillRate)

	log.Printf("Tokens before refill: %d, Tokens to add: %d", tb.tokens, tokensToAdd)

	if tokensToAdd > 0 {
		tb.tokens = min(tb.tokens+tokensToAdd, tb.maxTokens)
		tb.lastRefill = now
		log.Printf("Refilled tokens. New token count: %d", tb.tokens)
	}

	if tb.tokens > 0 {
		tb.tokens-- // Decrement the token count
		log.Printf("Request allowed. Remaining tokens after decrement: %d", tb.tokens)
		return true
	}

	log.Printf("Request denied. No tokens remaining. Next token in %v", tb.refillRate-elapsed%tb.refillRate)
	return false
}

func (tb *TokenBucket) GetWaitTime() time.Duration {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	elapsed := time.Since(tb.lastRefill)
	waitTime := tb.refillRate - (elapsed % tb.refillRate)
	return waitTime
}
