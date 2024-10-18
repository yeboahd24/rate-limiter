package router

import (
	"github.com/yeboahd24/rate-limiter/handler"
	"github.com/yeboahd24/rate-limiter/middleware"
	"github.com/yeboahd24/rate-limiter/util"
	"net/http"
	"time"
)

func SetupAuthRoutes(mux *http.ServeMux, authHandler *handler.AuthHandler) {
	// Create a single IPRateLimiter for all routes
	ipLimiter := util.NewIPRateLimiter(5, time.Minute)

	// Apply rate limiting to all routes
	applyMiddleware := func(h http.HandlerFunc) http.Handler {
		return middleware.RateLimiterMiddleware(ipLimiter)(h)
	}

	// Public routes with rate limiting
	mux.Handle("/register", http.HandlerFunc(authHandler.RegisterHandler))
	mux.Handle("/login", http.HandlerFunc(authHandler.LoginHandler))

	// Protected route with rate limiting and authentication
	mux.Handle("/protected", applyMiddleware((authHandler.ProfileHandler)))
}
