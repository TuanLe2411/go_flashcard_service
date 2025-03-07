package redis

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisDatabase struct {
	redis *redis.Client
}

func NewRedisClient() *RedisDatabase {
	return &RedisDatabase{}
}

func (r *RedisDatabase) Connect() error {
	db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	})

	r.redis = rdb
	return r.Ping()
}

func (r *RedisDatabase) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	ping := r.redis.Ping(ctx)
	if ping.Err() != nil {
		return ping.Err()
	}
	return nil
}

func (r *RedisDatabase) Close() error {
	return r.redis.Close()
}

func (r *RedisDatabase) Set(key string, value string, expiredTimeInMs int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return r.redis.Set(ctx, key, value, time.Duration(expiredTimeInMs)*time.Millisecond).Err()
}

func (r *RedisDatabase) Get(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return r.redis.Get(ctx, key).Result()
}

func (r *RedisDatabase) Del(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return r.redis.Del(ctx, key).Err()
}
