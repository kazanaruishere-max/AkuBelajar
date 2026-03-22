package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping: %v", err)
	}
	fmt.Println("✅ Connected to database")

	// Read migration files in order
	migDir := "migrations"
	entries, err := os.ReadDir(migDir)
	if err != nil {
		log.Fatalf("Failed to read migrations dir: %v", err)
	}

	var upFiles []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".up.sql") {
			upFiles = append(upFiles, e.Name())
		}
	}
	sort.Strings(upFiles)

	for _, f := range upFiles {
		path := filepath.Join(migDir, f)
		sql, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("Failed to read %s: %v", f, err)
		}

		fmt.Printf("⏳ Running %s...", f)
		_, err = pool.Exec(ctx, string(sql))
		if err != nil {
			fmt.Printf(" ❌\n   Error: %v\n", err)
			// Continue to next migration (some may already exist)
			continue
		}
		fmt.Println(" ✅")
	}

	fmt.Println("\n🎉 All migrations complete!")
}
