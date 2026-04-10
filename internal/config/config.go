package config

import (
	"os"
	"time"
)

const (
	defaultPort            = "8080"
	defaultServiceName     = "go-crypto-product-service"
	defaultCoinbaseURL     = "https://api.coinbase.com"
	defaultReadTimeout     = 5 * time.Second
	defaultWriteTimeout    = 10 * time.Second
	defaultIdleTimeout     = 30 * time.Second
	defaultUpstreamTimeout = 10 * time.Second
)

type Config struct {
	Port            string
	ServiceName     string
	CoinbaseBaseURL string
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

	return Config{
		Port:            port,
		ServiceName:     defaultServiceName,
		CoinbaseBaseURL: coinbaseBaseURL,
		ReadTimeout:     defaultReadTimeout,
		WriteTimeout:    defaultWriteTimeout,
		IdleTimeout:     defaultIdleTimeout,
		UpstreamTimeout: defaultUpstreamTimeout,
	}
}

func (c Config) Address() string {
	return ":" + c.Port
}
