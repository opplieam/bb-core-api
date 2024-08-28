package main

import (
	"strconv"
	"time"

	"github.com/opplieam/bb-core-api/internal/utils"
)

type Config struct {
	Web     WebConfig
	DB      DBConfig
	Service ServiceConfig
}

type WebConfig struct {
	Addr            string
	WriteTimeout    time.Duration
	ReadTimeout     time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

type DBConfig struct {
	Driver       string
	DSN          string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

type ServiceConfig struct {
	ProductAddr string
}

func NewConfig() *Config {
	writeTimeout, _ := strconv.Atoi(utils.GetEnv("WEB_WRITE_TIMEOUT", "10"))
	readTimeout, _ := strconv.Atoi(utils.GetEnv("WEB_READ_TIMEOUT", "5"))
	idleTimeout, _ := strconv.Atoi(utils.GetEnv("WEB_IDLE_TIMEOUT", "120"))
	shutDownTimeout, _ := strconv.Atoi(utils.GetEnv("WEB_SHUTDOWN_TIMEOUT", "20"))

	dbMaxOpenConns, _ := strconv.Atoi(utils.GetEnv("DB_MAX_OPEN_CONNS", "25"))
	dbMaxIdleConns, _ := strconv.Atoi(utils.GetEnv("DB_MAX_IDLE_CONNS", "25"))

	return &Config{
		Web: WebConfig{
			Addr:            utils.GetEnv("WEB_ADDR", ":3030"),
			WriteTimeout:    time.Duration(writeTimeout) * time.Second,
			ReadTimeout:     time.Duration(readTimeout) * time.Second,
			IdleTimeout:     time.Duration(idleTimeout) * time.Second,
			ShutdownTimeout: time.Duration(shutDownTimeout) * time.Second,
		},
		DB: DBConfig{
			Driver:       utils.GetEnv("DB_DRIVER", "postgres"),
			DSN:          utils.GetEnv("DB_DSN", "Postgres DSN"),
			MaxOpenConns: dbMaxOpenConns,
			MaxIdleConns: dbMaxIdleConns,
			MaxIdleTime:  utils.GetEnv("DB_MAX_IDLE_TIME", "15m"),
		},
		Service: ServiceConfig{
			ProductAddr: utils.GetEnv("PRODUCT_SERVICE_ADDR", "localhost:3031"),
		},
	}
}
