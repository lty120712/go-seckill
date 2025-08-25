package db

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"go-chat/configs"
)

var Redis *redis.Client

func InitRedis() {
	redisConfig := configs.AppConfig.Redis
	Redis = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", redisConfig.Host, redisConfig.Port),
		Password: redisConfig.Password, // 没有密码，默认值
		DB:       redisConfig.DB,       // 默认DB 0
	})
}
