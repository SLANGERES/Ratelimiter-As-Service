package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/rate-limiter-as-service/api"
	ratelimiter "github.com/rate-limiter-as-service/limiter/rate_limiter"
	"github.com/redis/go-redis/v9"
)

func main() {
	fmt.Print("Welcome to rate limiter as a service")

	rds := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	limiter := ratelimiter.ClientLimiter{
		Redis: rds,
		Config: map[string]ratelimiter.RatelimiterConfig{
			"user123": {Limit: 10, Window: time.Minute},
		},
	}

	slog.Info("Redis connected sucessfully")

	router := http.NewServeMux()

	router.HandleFunc("GET /api/v1/", func(w http.ResponseWriter, r *http.Request) { slog.Info("home service") })

	router.HandleFunc("POST /api/v1/checks", api.Checkhandler(&limiter))

	err := http.ListenAndServe(
		"0.0.0.0:7823",
		router,
	)

	if err != nil {
		slog.Any("Server fail ..", err)
	}
	slog.Info("Server started sucessfully")

}
