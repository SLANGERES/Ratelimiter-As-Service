package ratelimiter

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RatelimiterConfig struct {
	Limit  int           // max requests
	Window time.Duration // e.g. 1 * time.Minute
}

type ClientLimiter struct {
	Redis  *redis.Client
	Config map[string]RatelimiterConfig
}

// Fixed Window Rate Limiting
func (C *ClientLimiter) PerClientLimiter(clientID string) (bool, int, int) {
	ctx := context.Background()

	// Get client config or fallback
	cnfg, ok := C.Config[clientID]
	if !ok {
		cnfg = RatelimiterConfig{
			Limit:  5,
			Window: time.Minute,
		}
	}

	key := fmt.Sprintf("fwrl:%s", clientID)

	// Try to decrement the token count
	remaining, err := C.Redis.Decr(ctx, key).Result()
	fmt.Println(clientID, "+", remaining)
	if err != nil {
		// Key doesn't exist — first request
		C.Redis.Set(ctx, key, cnfg.Limit-1, cnfg.Window)
		return true, cnfg.Limit - 1, cnfg.Limit
	}

	// Check if still within limit
	if remaining >= 0 {
		return true, int(remaining), cnfg.Limit
	}

	// Reached or exceeded limit
	return false, int(remaining), cnfg.Limit
}
