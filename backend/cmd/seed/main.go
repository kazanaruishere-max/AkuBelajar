package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/kazanaruishere-max/akubelajar/backend/pkg/security"
)

type user struct {
	Email    string
	Password string
	Role     string
}

func main() {
	_ = godotenv.Load()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer pool.Close()
	fmt.Println("✅ Connected to database")

	schoolID := "a0000000-0000-0000-0000-000000000001"

	// 1. Seed school (skip if exists)
	_, err = pool.Exec(ctx, `
		INSERT INTO schools (id, name, code, address, config)
		VALUES ($1, 'SMP Nusantara Demo', 'SMND01', 'Jl. Pendidikan No. 1, Jakarta',
		 '{"grading_weights": {"assignment": 60, "quiz": 40}, "kkm_default": 70}')
		ON CONFLICT (id) DO NOTHING
	`, schoolID)
	if err != nil {
		log.Fatalf("Failed to seed school: %v", err)
	}
	fmt.Println("✅ School: SMP Nusantara Demo")

	// 2. Seed users with Argon2id hashed passwords
	users := []user{
		{Email: "admin@akubelajar.id", Password: "Admin@123!", Role: "super_admin"},
		{Email: "guru@akubelajar.id", Password: "Guru@123!", Role: "teacher"},
		{Email: "guru2@akubelajar.id", Password: "Guru@123!", Role: "teacher"},
		{Email: "siswa@akubelajar.id", Password: "Siswa@123!", Role: "student"},
		{Email: "siswa2@akubelajar.id", Password: "Siswa@123!", Role: "student"},
		{Email: "siswa3@akubelajar.id", Password: "Siswa@123!", Role: "student"},
		{Email: "siswa4@akubelajar.id", Password: "Siswa@123!", Role: "student"},
		{Email: "siswa5@akubelajar.id", Password: "Siswa@123!", Role: "student"},
	}

	for _, u := range users {
		hash, err := security.HashPassword(u.Password)
		if err != nil {
			log.Fatalf("Failed to hash password for %s: %v", u.Email, err)
		}

		_, err = pool.Exec(ctx, `
			INSERT INTO users (school_id, email, password_hash, role, is_active, is_first_login)
			VALUES ($1, $2, $3, $4::user_role, TRUE, TRUE)
			ON CONFLICT (email) DO NOTHING
		`, schoolID, u.Email, hash, u.Role)
		if err != nil {
			log.Fatalf("Failed to seed user %s: %v", u.Email, err)
		}
		fmt.Printf("✅ User: %s (%s)\n", u.Email, u.Role)
	}

	// 3. Seed academic year
	ayID := "b0000000-0000-0000-0000-000000000001"
	_, err = pool.Exec(ctx, `
		INSERT INTO academic_years (id, school_id, name, start_date, end_date, is_active)
		VALUES ($1, $2, '2025/2026', '2025-07-14', '2026-06-20', TRUE)
		ON CONFLICT (id) DO NOTHING
	`, ayID, schoolID)
	if err != nil {
		log.Fatalf("Failed to seed academic year: %v", err)
	}
	fmt.Println("✅ Academic Year: 2025/2026")

	// 4. Seed subjects
	subjects := []struct {
		ID   string
		Name string
		Code string
	}{
		{"c0000000-0000-0000-0000-000000000001", "Matematika", "MTK"},
		{"c0000000-0000-0000-0000-000000000002", "Bahasa Indonesia", "BIN"},
		{"c0000000-0000-0000-0000-000000000003", "IPA", "IPA"},
		{"c0000000-0000-0000-0000-000000000004", "Bahasa Inggris", "BIG"},
	}
	for _, s := range subjects {
		_, err = pool.Exec(ctx, `
			INSERT INTO subjects (id, school_id, name, code)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (id) DO NOTHING
		`, s.ID, schoolID, s.Name, s.Code)
		if err != nil {
			log.Fatalf("Failed to seed subject %s: %v", s.Name, err)
		}
	}
	fmt.Println("✅ Subjects: MTK, BIN, IPA, BIG")

	// 5. Seed classes
	classes := []struct {
		ID    string
		Name  string
		Grade int
	}{
		{"d0000000-0000-0000-0000-000000000001", "7A", 7},
		{"d0000000-0000-0000-0000-000000000002", "8A", 8},
		{"d0000000-0000-0000-0000-000000000003", "9A", 9},
	}
	for _, c := range classes {
		_, err = pool.Exec(ctx, `
			INSERT INTO classes (id, school_id, academic_year_id, name, grade_level)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (id) DO NOTHING
		`, c.ID, schoolID, ayID, c.Name, c.Grade)
		if err != nil {
			log.Fatalf("Failed to seed class %s: %v", c.Name, err)
		}
	}
	fmt.Println("✅ Classes: 7A, 8A, 9A")

	fmt.Println("\n🎉 Seed complete! Login credentials:")
	fmt.Println("   Admin:  admin@akubelajar.id / Admin@123!")
	fmt.Println("   Guru:   guru@akubelajar.id  / Guru@123!")
	fmt.Println("   Siswa:  siswa@akubelajar.id / Siswa@123!")
}
