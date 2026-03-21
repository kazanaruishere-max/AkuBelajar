# 🐳 Docker Setup — AkuBelajar

> File konfigurasi Docker lengkap untuk menjalankan seluruh stack di local machine.

---

## Quick Start

```bash
# 1. Clone dan setup environment
git clone https://github.com/your-org/akubelajar.git
cd akubelajar
cp .env.example .env

# 2. Jalankan semua service
make dev

# 3. Buka browser
# Frontend: http://localhost:3000
# API:      http://localhost:8080
# MinIO:    http://localhost:9001 (admin console)
# MailHog:  http://localhost:8025 (email catcher)
```

---

## 1. docker-compose.yml (Development)

```yaml
version: "3.9"

services:
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: ${DB_NAME:-akubelajar}
      POSTGRES_USER: ${DB_USER:-akubelajar}
      POSTGRES_PASSWORD: ${DB_PASSWORD:-akubelajar_dev}
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data
      - ./migrations/init.sql:/docker-entrypoint-initdb.d/init.sql
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER:-akubelajar}"]
      interval: 5s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    command: redis-server --appendonly yes
    ports:
      - "6379:6379"
    volumes:
      - redisdata:/data
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 3s
      retries: 5

  minio:
    image: minio/minio
    command: server /data --console-address ":9001"
    environment:
      MINIO_ROOT_USER: ${MINIO_ACCESS_KEY:-minioadmin}
      MINIO_ROOT_PASSWORD: ${MINIO_SECRET_KEY:-minioadmin}
    ports:
      - "9000:9000"
      - "9001:9001"
    volumes:
      - miniodata:/data

  minio-init:
    image: minio/mc
    depends_on:
      - minio
    entrypoint: >
      /bin/sh -c "
      mc alias set local http://minio:9000 minioadmin minioadmin;
      mc mb --ignore-existing local/akubelajar-files;
      mc mb --ignore-existing local/akubelajar-avatars;
      exit 0;
      "

  mailhog:
    image: mailhog/mailhog
    ports:
      - "1025:1025"   # SMTP
      - "8025:8025"   # Web UI

  api:
    build:
      context: ./backend
      dockerfile: Dockerfile
      target: development
    environment:
      - DB_URL=postgres://${DB_USER:-akubelajar}:${DB_PASSWORD:-akubelajar_dev}@postgres:5432/${DB_NAME:-akubelajar}?sslmode=disable
      - REDIS_URL=redis://redis:6379
      - MINIO_ENDPOINT=minio:9000
      - MINIO_ACCESS_KEY=${MINIO_ACCESS_KEY:-minioadmin}
      - MINIO_SECRET_KEY=${MINIO_SECRET_KEY:-minioadmin}
      - SMTP_HOST=mailhog
      - SMTP_PORT=1025
      - GEMINI_API_KEY=${GEMINI_API_KEY}
      - PASETO_KEY=${PASETO_KEY:-dev-key-32-bytes-long-placeholder!}
    ports:
      - "8080:8080"
    volumes:
      - ./backend:/app
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
      target: development
    environment:
      - NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
      - NEXT_PUBLIC_WS_URL=ws://localhost:8080/ws
    ports:
      - "3000:3000"
    volumes:
      - ./frontend:/app
      - /app/node_modules
    depends_on:
      - api

volumes:
  pgdata:
  redisdata:
  miniodata:
```

---

## 2. docker-compose.prod.yml

```yaml
version: "3.9"

services:
  api:
    build:
      context: ./backend
      dockerfile: Dockerfile
      target: production
    restart: always
    deploy:
      resources:
        limits:
          cpus: "2.0"
          memory: 512M

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
      target: production
    restart: always
    deploy:
      resources:
        limits:
          cpus: "1.0"
          memory: 256M

  postgres:
    restart: always
    deploy:
      resources:
        limits:
          cpus: "2.0"
          memory: 1G

  redis:
    restart: always
    deploy:
      resources:
        limits:
          cpus: "0.5"
          memory: 128M
```

---

## 3. Dockerfile — Go API

```dockerfile
# === Development ===
FROM golang:1.23-alpine AS development
WORKDIR /app
RUN go install github.com/air-verse/air@latest
COPY go.mod go.sum ./
RUN go mod download
COPY . .
CMD ["air", "-c", ".air.toml"]

# === Builder ===
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server ./cmd/api

# === Production ===
FROM alpine:3.19 AS production
RUN apk --no-cache add ca-certificates tzdata
RUN adduser -D -u 1000 appuser
WORKDIR /app
COPY --from=builder /server .
COPY --from=builder /app/migrations ./migrations
USER appuser
EXPOSE 8080
HEALTHCHECK CMD wget -qO- http://localhost:8080/health || exit 1
CMD ["./server"]
```

---

## 4. Dockerfile — Next.js

```dockerfile
# === Development ===
FROM node:22-alpine AS development
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
CMD ["npm", "run", "dev"]

# === Dependencies ===
FROM node:22-alpine AS deps
WORKDIR /app
COPY package*.json ./
RUN npm ci --production

# === Builder ===
FROM node:22-alpine AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .
RUN npm run build

# === Production ===
FROM node:22-alpine AS production
RUN adduser -D -u 1000 appuser
WORKDIR /app
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static
COPY --from=builder /app/public ./public
USER appuser
EXPOSE 3000
CMD ["node", "server.js"]
```

---

## 5. .env.example

```bash
# === Database ===
DB_NAME=akubelajar
DB_USER=akubelajar
DB_PASSWORD=CHANGE_ME_IN_PRODUCTION
DB_URL=postgres://akubelajar:CHANGE_ME@postgres:5432/akubelajar?sslmode=disable

# === Redis ===
REDIS_URL=redis://redis:6379

# === MinIO (Object Storage) ===
MINIO_ENDPOINT=minio:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=CHANGE_ME_IN_PRODUCTION

# === Auth ===
PASETO_KEY=CHANGE_ME_32_BYTES_LONG_SECRET_KEY

# === AI ===
GEMINI_API_KEY=YOUR_GEMINI_API_KEY_HERE

# === Notifications ===
FONNTE_API_KEY=YOUR_FONNTE_API_KEY_HERE
SMTP_HOST=mailhog
SMTP_PORT=1025
SMTP_USER=
SMTP_PASSWORD=

# === Frontend ===
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
NEXT_PUBLIC_WS_URL=ws://localhost:8080/ws
NEXT_PUBLIC_APP_URL=http://localhost:3000
```

---

## 6. Makefile

```makefile
.PHONY: dev build test lint migrate-up migrate-down seed logs

dev:
	docker compose up --build

dev-bg:
	docker compose up --build -d

build:
	docker compose -f docker-compose.yml -f docker-compose.prod.yml build

stop:
	docker compose down

reset:
	docker compose down -v && docker compose up --build

migrate-up:
	docker compose exec api go run cmd/migrate/main.go up

migrate-down:
	docker compose exec api go run cmd/migrate/main.go down

seed:
	docker compose exec api go run cmd/seed/main.go

test:
	docker compose exec api go test ./... -v -cover
	docker compose exec frontend npm test

lint:
	docker compose exec api golangci-lint run
	docker compose exec frontend npm run lint

logs:
	docker compose logs -f $(service)
```

---

*Terakhir diperbarui: 21 Maret 2026*
