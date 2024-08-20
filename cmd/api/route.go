package main

import (
	"database/sql"
	"log/slog"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/opplieam/bb-core-api/internal/middleware"
	"github.com/opplieam/bb-core-api/internal/store"
	"github.com/opplieam/bb-core-api/internal/utils"
	"github.com/opplieam/bb-core-api/internal/v1/probe"
)

func setupRoutes(log *slog.Logger, db *sql.DB) *gin.Engine {
	var r *gin.Engine
	if utils.GetEnv("WEB_SERVICE_ENV", "dev") == "dev" {
		r = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
	}

	r.Use(gin.Recovery())
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:5173"}
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowHeaders("Authorization")
	r.Use(cors.New(corsConfig))

	r.Use(middleware.SLogger(log, []string{"/v1/liveness", "/v1/readiness"}))

	v1 := r.Group("/v1")

	probeH := probe.NewHandler(build, store.NewHealthCheckStore(db))
	v1.GET("/liveness", probeH.LivenessHandler)
	v1.GET("/readiness", probeH.ReadinessHandler)
	return r
}
