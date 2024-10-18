package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"github.com/yeboahd24/rate-limiter/model"
	"github.com/yeboahd24/rate-limiter/util"
)

func RateLimiterMiddleware(ipLimiter *util.IPRateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Rate Limiting
			ip := r.RemoteAddr
			if host, _, err := net.SplitHostPort(ip); err == nil {
				ip = host
			}
			log.Printf("Incoming request from IP %s", ip)

			allowed, waitTime := ipLimiter.Allow(r.Context(), ip)
			if !allowed {
				log.Printf("Rate limit exceeded for IP: %s", r.RemoteAddr)

				// Set the Retry-After header
				w.Header().Set("Retry-After", fmt.Sprintf("%.0f", waitTime.Seconds()))

				// Set the Content-Type header
				w.Header().Set("Content-Type", "application/json")

				// Create a custom response message
				response := map[string]interface{}{
					"error":             "Too many requests",
					"retry_after":       int(waitTime.Seconds()),
					"retry_after_human": fmt.Sprintf("Please wait %d seconds before retrying.", int(waitTime.Seconds())),
				}

				// Convert the response to JSON
				jsonResponse, err := json.Marshal(response)
				if err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}

				// Write the response
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write(jsonResponse)
				return
			}
			log.Printf("Request allowed for IP: %s", r.RemoteAddr)

			// JWT Authentication
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing auth token", http.StatusUnauthorized)
				return
			}

			log.Println(authHeader)

			bearerToken := strings.Split(authHeader, " ")
			if len(bearerToken) != 2 {
				http.Error(w, "Invalid auth token", http.StatusUnauthorized)
				return
			}

			token, err := jwt.ParseWithClaims(bearerToken[1], &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(viper.GetString("JWT_SECRET")), nil
			})

			if err != nil {
				http.Error(w, "Invalid auth token, jwt", http.StatusUnauthorized)
				return
			}

			if claims, ok := token.Claims.(*model.Claims); ok && token.Valid {
				ctx := context.WithValue(r.Context(), "user", claims)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				http.Error(w, "Invalid auth token", http.StatusUnauthorized)
			}
		})
	}
}
