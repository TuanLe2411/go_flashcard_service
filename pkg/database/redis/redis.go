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

func (r *RedisDatabase) Set(key string, value string, expiredTimeInSec int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return r.redis.Set(ctx, key, value, time.Duration(expiredTimeInSec)*time.Second).Err()
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

func (r *RedisDatabase) HSet(key, field string, value any) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return r.redis.HSet(ctx, key, field, value).Err()
}

func (r *RedisDatabase) HGet(key, field string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return r.redis.HGet(ctx, key, field).Result()
}

func (r *RedisDatabase) HGetAll(key string) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return r.redis.HGetAll(ctx, key).Result()
}

func (r *RedisDatabase) HMSet(key string, fields map[string]any) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return r.redis.HMSet(ctx, key, fields).Err()
}

func (r *RedisDatabase) HDel(key string, fields ...string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return r.redis.HDel(ctx, key, fields...).Err()
}

func (r *RedisDatabase) HSetWithExpiry(key, field string, value any, expiredTimeInSec int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pipe := r.redis.Pipeline()
	pipe.HSet(ctx, key, field, value)
	pipe.Expire(ctx, key, time.Duration(expiredTimeInSec)*time.Second)
	_, err := pipe.Exec(ctx)
	return err
}

func (r *RedisDatabase) HMSetWithExpiry(key string, fields map[string]any, expiredTimeInSec int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pipe := r.redis.Pipeline()
	pipe.HMSet(ctx, key, fields)
	pipe.Expire(ctx, key, time.Duration(expiredTimeInSec)*time.Second)
	_, err := pipe.Exec(ctx)
	return err
}
