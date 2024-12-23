package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

var build = "dev"

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger = logger.With("service", "bb-core-api", "build", build)

	if err := run(logger); err != nil {
		logger.Error("run server", "error", err)
	}
}

func run(log *slog.Logger) error {
	log.Info("start up", "GOMAXPROCS", runtime.GOMAXPROCS(0))
	// Setup config
	cfg := NewConfig()

	// Setup Tracer
	err := initTracerProvider()
	if err != nil {
		return err
	}

	// Setup database
	db, err := setupDB(cfg, log)
	if err != nil {
		return err
	}
	defer db.Close()

	// setup OAuth provider
	setupProvider()

	// setup grpc
	conn, err := setupGRPC(cfg)
	defer conn.Close()

	// Setup routes
	r := setupRoutes(log, db, conn)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         cfg.Web.Addr,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		Handler:      r.Handler(),
	}

	serverErrors := make(chan error, 1)
	go func() {
		defer close(serverErrors)
		log.Info("start up", "status", "api router started", "address", srv.Addr)
		serverErrors <- srv.ListenAndServe()
	}()

	// Shutdown
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdown:
		log.Info("shutdown", "status", "shutdown started", "signal", sig)
		defer log.Info("shutdown", "status", "shutdown complete", "signal", sig)
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			_ = srv.Close()
			return fmt.Errorf("cannot not shutdown gratefuly: %w", err)
		}
	}

	return nil
}
