package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	ratelimiter "github.com/rate-limiter-as-service/limiter/rate_limiter"
)

// Accepts a limiter instance as dependency
func Checkhandler(limiter *ratelimiter.ClientLimiter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Check handler")
		var reqBody struct {
			Uid string `json:"uid"`
		}

		// Decode request body
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			slog.Error("Unable to parse request body", "error", err)
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		// Apply token bucket logic
		allowed, _, _ := limiter.PerClientLimiter(reqBody.Uid)

		// Send response
		if !allowed {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)

			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello"))

	}
}
