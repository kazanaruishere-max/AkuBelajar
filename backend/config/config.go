package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application.
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Paseto   PasetoConfig
	Gemini   GeminiConfig
	Resend   ResendConfig
	Supabase SupabaseConfig
	CORS     CORSConfig
}

// AppConfig holds application-level configuration.
type AppConfig struct {
	Port        string
	Env         string // "development", "production"
	FrontendURL string
}

// DatabaseConfig holds database connection configuration.
type DatabaseConfig struct {
	URL         string
	Host        string
	Port        string
	User        string
	Password    string
	Name        string
	SSLMode     string
	MaxConns    int
	MaxIdleTime time.Duration
}

// RedisConfig holds Redis connection configuration.
type RedisConfig struct {
	URL      string
	Host     string
	Port     string
	Password string
	DB       int
}

// PasetoConfig holds Paseto token configuration.
type PasetoConfig struct {
	SymmetricKey       string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
}

// GeminiConfig holds Gemini AI configuration.
type GeminiConfig struct {
	APIKey string
	Model  string
}

// ResendConfig holds email service configuration.
type ResendConfig struct {
	APIKey    string
	FromEmail string
}

// SupabaseConfig holds Supabase storage configuration.
type SupabaseConfig struct {
	URL     string
	AnonKey string
	Bucket  string
}

// CORSConfig holds CORS configuration.
type CORSConfig struct {
	AllowOrigin string
}

// Load reads configuration from environment variables and returns a Config.
func Load() (*Config, error) {
	cfg := &Config{
		App: AppConfig{
			Port:        getEnvOrDefault("PORT", "8080"),
			Env:         getEnvOrDefault("APP_ENV", "development"),
			FrontendURL: getEnvOrDefault("FRONTEND_URL", "http://localhost:3000"),
		},
		Database: DatabaseConfig{
			URL:         os.Getenv("DATABASE_URL"),
			Host:        getEnvOrDefault("DB_HOST", "localhost"),
			Port:        getEnvOrDefault("DB_PORT", "5432"),
			User:        getEnvOrDefault("DB_USER", "akubelajar"),
			Password:    getEnvOrDefault("DB_PASSWORD", "localdev123"),
			Name:        getEnvOrDefault("DB_NAME", "akubelajar_dev"),
			SSLMode:     getEnvOrDefault("DB_SSL_MODE", "disable"),
			MaxConns:    getEnvInt("DB_MAX_CONNS", 25),
			MaxIdleTime: time.Duration(getEnvInt("DB_MAX_IDLE_TIME_MIN", 15)) * time.Minute,
		},
		Redis: RedisConfig{
			URL:      os.Getenv("REDIS_URL"),
			Host:     getEnvOrDefault("REDIS_HOST", "localhost"),
			Port:     getEnvOrDefault("REDIS_PORT", "6379"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		Paseto: PasetoConfig{
			SymmetricKey:       os.Getenv("PASETO_KEY"),
			AccessTokenExpiry:  time.Duration(getEnvInt("ACCESS_TOKEN_EXPIRY_MIN", 15)) * time.Minute,
			RefreshTokenExpiry: time.Duration(getEnvInt("REFRESH_TOKEN_EXPIRY_DAY", 7)) * 24 * time.Hour,
		},
		Gemini: GeminiConfig{
			APIKey: os.Getenv("GEMINI_API_KEY"),
			Model:  getEnvOrDefault("GEMINI_MODEL", "gemini-2.0-flash"),
		},
		Resend: ResendConfig{
			APIKey:    os.Getenv("RESEND_API_KEY"),
			FromEmail: getEnvOrDefault("RESEND_FROM_EMAIL", "noreply@akubelajar.id"),
		},
		Supabase: SupabaseConfig{
			URL:     os.Getenv("SUPABASE_URL"),
			AnonKey: os.Getenv("SUPABASE_ANON_KEY"),
			Bucket:  getEnvOrDefault("SUPABASE_BUCKET", "uploads"),
		},
		CORS: CORSConfig{
			AllowOrigin: getEnvOrDefault("CORS_ORIGIN", "http://localhost:3000"),
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// DatabaseDSN returns the database connection string.
func (c *DatabaseConfig) DatabaseDSN() string {
	if c.URL != "" {
		return c.URL
	}
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Name, c.SSLMode,
	)
}

// validate checks that required configuration values are present.
func (c *Config) validate() error {
	if c.Paseto.SymmetricKey == "" && c.App.Env == "production" {
		return fmt.Errorf("PASETO_KEY is required in production")
	}

	// Set a development key if not provided
	if c.Paseto.SymmetricKey == "" {
		c.Paseto.SymmetricKey = "dev-paseto-key-min-32-chars-long!"
	}

	return nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}
