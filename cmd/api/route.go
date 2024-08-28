package main

import (
	"database/sql"
	"log/slog"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/opplieam/bb-core-api/internal/middleware"
	"github.com/opplieam/bb-core-api/internal/store"
	"github.com/opplieam/bb-core-api/internal/utils"
	"github.com/opplieam/bb-core-api/internal/v1/auth"
	"github.com/opplieam/bb-core-api/internal/v1/probe"
	"github.com/opplieam/bb-core-api/internal/v1/product"
	"google.golang.org/grpc"
)

func setupRoutes(log *slog.Logger, db *sql.DB, grpcConn *grpc.ClientConn) *gin.Engine {
	var r *gin.Engine
	if utils.GetEnv("WEB_SERVICE_ENV", "dev") == "dev" {
		r = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
	}

	r.Use(gin.Recovery())
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:5174"}
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowHeaders("Authorization")
	r.Use(cors.New(corsConfig))

	r.Use(middleware.SLogger(log, []string{"/v1/liveness", "/v1/readiness"}))

	v1 := r.Group("/v1")

	probeH := probe.NewHandler(build, store.NewHealthCheckStore(db))
	v1.GET("/liveness", probeH.LivenessHandler)
	v1.GET("/readiness", probeH.ReadinessHandler)

	// TODO: Add some sort of csrf token for login button
	// TODO: Add refresh token endpoint
	userStore := store.NewUserStore(db)
	authH := auth.NewHandler(userStore)
	v1.GET("/auth/:provider", authH.ProviderHandler)
	v1.GET("/auth/:provider/callback", authH.CallbackHandler)
	v1.POST("/auth/token", authH.GetTokenHandler)
	v1.GET("/auth/:provider/logout", authH.LogoutHandler)

	// TODO: Add authorization middleware
	productH := product.NewHandler(grpcConn)
	v1.GET("/product", productH.GetAllProducts)
	return r
}
