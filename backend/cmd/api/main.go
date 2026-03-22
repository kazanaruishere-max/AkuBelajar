package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kazanaruishere-max/akubelajar/backend/config"
	"github.com/kazanaruishere-max/akubelajar/backend/pkg/cache"
	"github.com/kazanaruishere-max/akubelajar/backend/pkg/database"
	"github.com/kazanaruishere-max/akubelajar/backend/pkg/response"
	"github.com/kazanaruishere-max/akubelajar/backend/pkg/security"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Set Gin mode
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	ctx := context.Background()

	// Connect to PostgreSQL
	dbPool, err := database.NewPostgresPool(
		ctx,
		cfg.Database.DatabaseDSN(),
		cfg.Database.MaxConns,
		cfg.Database.MaxIdleTime,
	)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer dbPool.Close()
	log.Println("✅ Connected to PostgreSQL")

	// Connect to Redis
	redisClient, err := cache.NewRedisClient(
		ctx,
		cfg.Redis.URL,
		cfg.Redis.Host,
		cfg.Redis.Port,
		cfg.Redis.Password,
		cfg.Redis.DB,
	)
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
	defer redisClient.Close()
	log.Println("✅ Connected to Redis")

	// Initialize Paseto token maker
	tokenMaker, err := security.NewTokenMaker(cfg.Paseto.SymmetricKey)
	if err != nil {
		log.Fatalf("failed to create token maker: %v", err)
	}
	log.Println("✅ Paseto v4 token maker initialized")

	// Suppress unused variable warnings — will be used when handlers are registered
	_ = dbPool
	_ = redisClient
	_ = tokenMaker

	// Setup Gin router
	router := gin.Default()

	// CORS middleware
	router.Use(corsMiddleware(cfg.CORS.AllowOrigin))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		response.OK(c, gin.H{
			"status":  "ok",
			"version": "0.1.0",
			"time":    time.Now().UTC().Format(time.RFC3339),
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes (Sprint 1)
		_ = v1

		// TODO: Register domain handlers here
		// auth.RegisterRoutes(v1, authHandler)
		// academic.RegisterRoutes(v1, academicHandler)
	}

	// Start server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.App.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		log.Printf("🚀 Server starting on port %s (env: %s)", cfg.App.Port, cfg.App.Env)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("👋 Server exited gracefully")
}

// corsMiddleware handles CORS for cross-origin requests from the frontend.
func corsMiddleware(allowOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", allowOrigin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
