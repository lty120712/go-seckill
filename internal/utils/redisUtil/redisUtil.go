package redisUtil

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

var ctx = context.Background()

func IsRateLimited(rdb *redis.Client, key string, limit int, window int) bool {

	count, err := rdb.Incr(ctx, key).Result()
	if err != nil {
		fmt.Println("Redis error:", err)
		return false
	}

	if count == 1 {
		rdb.Expire(ctx, key, time.Duration(window)*time.Second)
	}

	return count > int64(limit)
}
