package config

import (
	"os"
	"strconv"
	"time"
)

const (
	defaultPort            = "8080"
	defaultServiceName     = "go-crypto-product-service"
	defaultCoinbaseURL     = "https://api.coinbase.com"
	defaultRedisAddr       = "localhost:6379"
	defaultReadTimeout     = 5 * time.Second
	defaultWriteTimeout    = 10 * time.Second
	defaultIdleTimeout     = 30 * time.Second
	defaultUpstreamTimeout = 10 * time.Second
	defaultCacheTTL        = 60 * time.Second
)

type Config struct {
	Port            string
	ServiceName     string
	CoinbaseBaseURL string
	RedisAddr       string
	RedisPassword   string
	RedisDB         int
	CacheTTL        time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	UpstreamTimeout time.Duration
}

func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	coinbaseBaseURL := os.Getenv("COINBASE_BASE_URL")
	if coinbaseBaseURL == "" {
		coinbaseBaseURL = defaultCoinbaseURL
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = defaultRedisAddr
	}

	redisDB := 0
	if redisDBValue := os.Getenv("REDIS_DB"); redisDBValue != "" {
		if parsedRedisDB, err := strconv.Atoi(redisDBValue); err == nil {
			redisDB = parsedRedisDB
		}
	}

	cacheTTL := defaultCacheTTL
	if cacheTTLValue := os.Getenv("CACHE_TTL"); cacheTTLValue != "" {
		if parsedCacheTTL, err := time.ParseDuration(cacheTTLValue); err == nil {
			cacheTTL = parsedCacheTTL
		}
	}

	return Config{
		Port:            port,
		ServiceName:     defaultServiceName,
		CoinbaseBaseURL: coinbaseBaseURL,
		RedisAddr:       redisAddr,
		RedisPassword:   os.Getenv("REDIS_PASSWORD"),
		RedisDB:         redisDB,
		CacheTTL:        cacheTTL,
		ReadTimeout:     defaultReadTimeout,
		WriteTimeout:    defaultWriteTimeout,
		IdleTimeout:     defaultIdleTimeout,
		UpstreamTimeout: defaultUpstreamTimeout,
	}
}

func (c Config) Address() string {
	return ":" + c.Port
}
