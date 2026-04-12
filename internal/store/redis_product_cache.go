package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/xiaofeihuang85/go-crypto-product-service/internal/model"
)

var ErrCacheMiss = errors.New("product cache miss")

type RedisProductCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisProductCache(addr, password string, db int, ttl time.Duration) *RedisProductCache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisProductCache{
		client: client,
		ttl:    ttl,
	}
}

func (c *RedisProductCache) GetProduct(ctx context.Context, productID string) (model.ProductResponse, error) {
	value, err := c.client.Get(ctx, productKey(productID)).Result()
	if errors.Is(err, redis.Nil) {
		return model.ProductResponse{}, ErrCacheMiss
	}
	if err != nil {
		return model.ProductResponse{}, fmt.Errorf("get product from redis: %w", err)
	}

	var product model.ProductResponse
	if err := json.Unmarshal([]byte(value), &product); err != nil {
		return model.ProductResponse{}, fmt.Errorf("decode cached product: %w", err)
	}

	return product, nil
}

func (c *RedisProductCache) SetProduct(ctx context.Context, product model.ProductResponse) error {
	payload, err := json.Marshal(product)
	if err != nil {
		return fmt.Errorf("encode product for redis: %w", err)
	}

	if err := c.client.Set(ctx, productKey(product.ProductID), payload, c.ttl).Err(); err != nil {
		return fmt.Errorf("set product in redis: %w", err)
	}

	return nil
}

func productKey(productID string) string {
	return "product:" + productID
}
