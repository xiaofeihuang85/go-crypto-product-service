package config

import (
	"os"
	"time"
)

const (
	defaultPort         = "8080"
	defaultServiceName  = "go-crypto-product-service"
	defaultReadTimeout  = 5 * time.Second
	defaultWriteTimeout = 10 * time.Second
	defaultIdleTimeout  = 30 * time.Second
)

type Config struct {
	Port         string
	ServiceName  string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	return Config{
		Port:         port,
		ServiceName:  defaultServiceName,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		IdleTimeout:  defaultIdleTimeout,
	}
}

func (c Config) Address() string {
	return ":" + c.Port
}
