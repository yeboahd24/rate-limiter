package util

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type IPRateLimiter struct {
	redisClient *redis.Client
	max         int
	refillRate  time.Duration
}

func NewIPRateLimiter(redisAddr string, max int, refillRate time.Duration) *IPRateLimiter {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &IPRateLimiter{
		redisClient: redisClient,
		max:         max,
		refillRate:  refillRate,
	}
}

func (i *IPRateLimiter) Allow(ctx context.Context, ip string) (bool, time.Duration) {
	key := "ratelimit:" + ip
	count, err := i.redisClient.Incr(ctx, key).Result()
	if err != nil {
		log.Printf("Error incrementing rate limit for IP %s: %v", ip, err)
		return false, 0
	}

	if count > int64(i.max) {
		ttl, err := i.redisClient.TTL(ctx, key).Result()
		if err != nil {
			log.Printf("Error getting TTL for IP %s: %v", ip, err)
			return false, 0
		}

		log.Printf("Request denied for IP %s. Retry in %v", ip, ttl)
		return false, ttl
	}

	// Set the expiration time for the key
	_, err = i.redisClient.Expire(ctx, key, i.refillRate).Result()
	if err != nil {
		log.Printf("Error setting expiration for IP %s: %v", ip, err)
		return false, 0
	}

	log.Printf("Request allowed for IP %s", ip)
	return true, i.refillRate
}

func (i *IPRateLimiter) ResetLimit(ctx context.Context, ip string) error {
	key := "ratelimit:" + ip
	_, err := i.redisClient.Del(ctx, key).Result()
	if err != nil {
		log.Printf("Error resetting rate limit for IP %s: %v", ip, err)
		return err
	}
	return nil
}
