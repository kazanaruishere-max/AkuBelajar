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
	"github.com/joho/godotenv"
	"github.com/kazanaruishere-max/akubelajar/backend/config"
	"github.com/kazanaruishere-max/akubelajar/backend/internal/academic"
	"github.com/kazanaruishere-max/akubelajar/backend/internal/assignment"
	"github.com/kazanaruishere-max/akubelajar/backend/internal/attendance"
	"github.com/kazanaruishere-max/akubelajar/backend/internal/auth"
	"github.com/kazanaruishere-max/akubelajar/backend/internal/grade"
	"github.com/kazanaruishere-max/akubelajar/backend/internal/middleware"
	"github.com/kazanaruishere-max/akubelajar/backend/internal/notification"
	"github.com/kazanaruishere-max/akubelajar/backend/internal/quiz"
	"github.com/kazanaruishere-max/akubelajar/backend/internal/upload"
	"github.com/kazanaruishere-max/akubelajar/backend/pkg/ai"
	"github.com/kazanaruishere-max/akubelajar/backend/pkg/cache"
	"github.com/kazanaruishere-max/akubelajar/backend/pkg/database"
	"github.com/kazanaruishere-max/akubelajar/backend/pkg/response"
	"github.com/kazanaruishere-max/akubelajar/backend/pkg/security"
	"github.com/kazanaruishere-max/akubelajar/backend/pkg/storage"
	"github.com/kazanaruishere-max/akubelajar/backend/pkg/validator"
)

func main() {
	// Load .env file (ignore error if not found — production uses real env vars)
	_ = godotenv.Load()

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

	// Connect to PostgreSQL (warn if unavailable — server starts in degraded mode)
	dbPool, err := database.NewPostgresPool(
		ctx,
		cfg.Database.DatabaseDSN(),
		cfg.Database.MaxConns,
		cfg.Database.MaxIdleTime,
	)
	dbConnected := err == nil
	if err != nil {
		log.Printf("⚠️  Database unavailable (server will start without DB): %v", err)
	} else {
		defer dbPool.Close()
		log.Println("✅ Connected to PostgreSQL")
	}

	// Connect to Redis (warn if unavailable)
	redisClient, err := cache.NewRedisClient(
		ctx,
		cfg.Redis.URL,
		cfg.Redis.Host,
		cfg.Redis.Port,
		cfg.Redis.Password,
		cfg.Redis.DB,
	)
	redisConnected := err == nil
	if err != nil {
		log.Printf("⚠️  Redis unavailable (server will start without cache): %v", err)
	} else {
		defer redisClient.Close()
		log.Println("✅ Connected to Redis")
	}

	// Initialize Paseto token maker
	tokenMaker, err := security.NewTokenMaker(cfg.Paseto.SymmetricKey)
	if err != nil {
		log.Fatalf("failed to create token maker: %v", err)
	}
	log.Println("✅ Paseto v4 token maker initialized")

	// Initialize validator
	v := validator.New()

	// Initialize modules (only if DB is available)
	var authHandler *auth.Handler
	var academicHandler *academic.Handler
	var assignmentHandler *assignment.Handler
	var quizHandler *quiz.Handler
	var attendanceHandler *attendance.Handler
	var gradeHandler *grade.Handler
	var notifHandler *notification.Handler
	var uploadHandler *upload.Handler
	if dbConnected {
		// Auth
		authRepo := auth.NewRepository(dbPool)
		authService := auth.NewService(authRepo, tokenMaker, cfg.Paseto.AccessTokenExpiry, cfg.Paseto.RefreshTokenExpiry)
		authHandler = auth.NewHandler(authService, v)
		log.Println("✅ Auth module initialized")

		// Academic
		academicRepo := academic.NewRepository(dbPool)
		academicHandler = academic.NewHandler(academicRepo, v)
		log.Println("✅ Academic module initialized")

		// Assignment
		assignmentRepo := assignment.NewRepository(dbPool)
		assignmentService := assignment.NewService(assignmentRepo)
		assignmentHandler = assignment.NewHandler(assignmentService, assignmentRepo, v)
		log.Println("✅ Assignment module initialized")

		// Quiz
		quizRepo := quiz.NewRepository(dbPool)
		quizService := quiz.NewService(quizRepo)
		geminiClient := ai.NewGeminiClient()
		var quizAISvc *quiz.AIService
		if geminiClient.IsAvailable() {
			quizAISvc = quiz.NewAIService(geminiClient, quizRepo)
			log.Println("✅ AI Quiz service initialized")
		}
		quizHandler = quiz.NewHandler(quizService, quizRepo, quizAISvc, v)
		log.Println("✅ Quiz module initialized")

		// Attendance
		attendanceRepo := attendance.NewRepository(dbPool)
		attendanceHandler = attendance.NewHandler(attendanceRepo, v)
		log.Println("✅ Attendance module initialized")

		// Grade
		gradeRepo := grade.NewRepository(dbPool)
		gradeHandler = grade.NewHandler(gradeRepo, v)
		log.Println("✅ Grade module initialized")

		// Notification
		notifRepo := notification.NewRepository(dbPool)
		notifHandler = notification.NewHandler(notifRepo, v)
		log.Println("✅ Notification module initialized")

		// Upload
		supaStore := storage.NewSupabaseStorage()
		uploadHandler = upload.NewHandler(supaStore)
		log.Println("✅ Upload module initialized")
	}

	// Setup Gin router
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.RequestLogger())
	router.Use(corsMiddleware(cfg.CORS.AllowOrigin))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		status := "ok"
		if !dbConnected || !redisConnected {
			status = "degraded"
		}
		response.OK(c, gin.H{
			"status":   status,
			"version":  "0.1.0",
			"time":     time.Now().UTC().Format(time.RFC3339),
			"database": dbConnected,
			"redis":    redisConnected,
		})
	})

	// Auth middleware instance
	authMW := middleware.AuthMiddleware(tokenMaker)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Rate limiter: 120 req/min for general API
		if redisConnected {
			v1.Use(middleware.RateLimiter(redisClient, 120, time.Minute, "api"))
		}

		// RBAC middleware
		adminMW := middleware.RequireRole("super_admin")
		teacherMW := middleware.RequireRole("super_admin", "teacher")
		studentMW := middleware.RequireRole("student", "class_leader")

		// Auth routes
		if authHandler != nil {
			auth.RegisterRoutes(v1, authHandler, authMW)
		}

		// Academic routes (admin only)
		if academicHandler != nil {
			academic.RegisterRoutes(v1, academicHandler, authMW, adminMW)
		}

		// Assignment routes (teacher + student)
		if assignmentHandler != nil {
			assignment.RegisterRoutes(v1, assignmentHandler, authMW, teacherMW, studentMW)
		}

		// Quiz routes (teacher + student)
		if quizHandler != nil {
			quiz.RegisterRoutes(v1, quizHandler, authMW, teacherMW, studentMW)
		}

		// Attendance routes (teacher + student)
		if attendanceHandler != nil {
			attendance.RegisterRoutes(v1, attendanceHandler, authMW, teacherMW, studentMW)
		}

		// Grade routes (teacher + student)
		if gradeHandler != nil {
			grade.RegisterRoutes(v1, gradeHandler, authMW, teacherMW, studentMW)
		}

		// Notification routes (all authenticated users)
		if notifHandler != nil {
			notification.RegisterRoutes(v1, notifHandler, authMW, teacherMW)
		}

		// Upload route (all authenticated users)
		if uploadHandler != nil {
			upload.RegisterRoutes(v1, uploadHandler, authMW)
		}
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
