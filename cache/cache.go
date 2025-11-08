package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/go-redis/redis/v8"
	t "github.com/instaUpload/user-service/types"
)

type Cacher interface {
	Close(context.Context)
	store(key string, value any) error
	retrieve(key string) (interface{}, error)
	StoreUser(user t.User) error
	RetrieveUser(userID int64) (t.User, error)
}

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(ctx context.Context) *RedisCache {
	// Initialize and return a new RedisCache instance
	config := t.NewCacheConfig()
	client := redis.NewClient(&redis.Options{
		Addr:     config.Address(),
		Password: config.Password,
		DB:       config.DB,
	})

	cache := &RedisCache{
		client: client,
	}

	go cache.Close(ctx)
	res := client.Ping(ctx)
	if res.Err() != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", res.Err()))
	}
	slog.Info("Connected to Redis cache successfully", slog.String("ping:", res.Val()))
	return cache
}

func (r *RedisCache) Close(ctx context.Context) {
	<-ctx.Done()
	r.client.Close()
	fmt.Println("Redis cache connection closed")
}

func (r *RedisCache) store(key string, value any) error {
	err := r.client.Set(context.Background(), key, value, 0).Err()
	return err
}

func (r *RedisCache) retrieve(key string) (interface{}, error) {
	val, err := r.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("key %s does not exist", key)
	} else if err != nil {
		return nil, err
	}
	return val, nil
}

func (r *RedisCache) StoreUser(user t.User) error {
	key := fmt.Sprintf("user:%d", user.ID)
	value, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return r.store(key, string(value))
}

func (r *RedisCache) RetrieveUser(userID int64) (t.User, error) {
	key := fmt.Sprintf("user:%d", userID)
	val, err := r.retrieve(key)
	if err != nil {
		return t.User{}, err
	}
	var user t.User
	err = json.Unmarshal([]byte(val.(string)), &user)
	if err != nil {
		return t.User{}, err
	}
	return user, nil
}
