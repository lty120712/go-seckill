package tests

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"testing"
	"time"
)

func TestRedis(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "123456",
		DB:       0,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client.Set(ctx, "test", "test", 10*time.Second)
	result, err := client.Get(ctx, "test").Result()
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			logrus.Error("timeout")
		}
		logrus.Error(err)
	}
	println(result)
}
