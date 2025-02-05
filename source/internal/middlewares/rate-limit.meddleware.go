package middlewares

import (
	"context"
	"fmt"
	"net/http"

	"time"

	errors "github.com/Lucasvmarangoni/logella/err"
	redisrate "github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
)

func SetupRedisRateLimiter(addr, password string) *redisrate.Limiter {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})	
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		errors.PanicErr(err, "SetupRedisRateLimiter")
	}
	return redisrate.NewLimiter(client)
}

type UserRateLimitMiddleware struct {
	RedisLimiter *redisrate.Limiter
}

func NewUserRateLimit(addr, password string) *UserRateLimitMiddleware {
	return &UserRateLimitMiddleware{SetupRedisRateLimiter(addr, password)}
}

func (u *UserRateLimitMiddleware) Handler() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := getSub(w, r)
			key := fmt.Sprintf("user:%s", id)
			ctx := r.Context()
			res, _ := u.RedisLimiter.Allow(ctx, key, redisrate.Limit{
				Rate:   10,
				Burst:  20,
				Period: 60 * time.Second,
			})
			if res.Allowed <= 0 {
				http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
