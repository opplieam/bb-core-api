package main

import (
	"context"
	"database/sql"
	"log/slog"
	"time"
)

func setupDB(cfg *Config, log *slog.Logger) (*sql.DB, error) {
	log.Info("database start up", "driver", cfg.DB.Driver, "DSN", cfg.DB.DSN)
	db, err := sql.Open(cfg.DB.Driver, cfg.DB.DSN)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	db.SetMaxIdleConns(cfg.DB.MaxIdleConns)
	duration, err := time.ParseDuration(cfg.DB.MaxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	return db, nil
}
