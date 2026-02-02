package database

import (
	"context"

	"github.com/gofiber/fiber/v2/log"
	"github.com/redis/go-redis/v9"
	"github.com/shuTwT/hoshikuzu/pkg/config"
)

func NewRedisClient(ctx context.Context) (*redis.Client, error) {
	if !config.GetBool(config.Redis_Enable) {
		return nil, nil
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.GetString(config.REDIS_ADDR),
		Password: config.GetString(config.REDIS_PASSWORD),
		DB:       config.GetInt(config.REDIS_DB),
	})
	// 测试连接
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Error("Failed to connect to Redis", "error", err.Error())
		err = rdb.Close()
		if err != nil {
			log.Error("Failed to close Redis connection", "error", err.Error())
		}
		return nil, err
	}
	return rdb, nil
}
