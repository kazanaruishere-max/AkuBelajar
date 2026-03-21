# 🚀 Deployment Guide — AkuBelajar

> Panduan langkah demi langkah dari Docker build hingga Kubernetes deploy.

---

## Overview Pipeline

```
Code Push → GitHub Actions → Build & Test → Docker Image → Registry → K8s Deploy
```

---

## 1. Docker Build

### Backend (Go)

```dockerfile
# Dockerfile.backend
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /api ./cmd/api

FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /api /api
EXPOSE 8080
ENTRYPOINT ["/api"]
```

### Frontend (Next.js)

```dockerfile
# Dockerfile.frontend
FROM node:22-alpine AS builder
WORKDIR /app
COPY package.json pnpm-lock.yaml ./
RUN corepack enable && pnpm install --frozen-lockfile
COPY . .
RUN pnpm build

FROM node:22-alpine AS runner
WORKDIR /app
ENV NODE_ENV=production
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static
COPY --from=builder /app/public ./public
EXPOSE 3000
CMD ["node", "server.js"]
```

### Build Commands

```bash
# Build semua images
docker-compose build

# Atau individual
docker build -f Dockerfile.backend -t akubelajar-api:latest .
docker build -f Dockerfile.frontend -t akubelajar-web:latest .
```

---

## 2. Docker Compose (Staging/Development)

```yaml
# docker-compose.yml
services:
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: akubelajar
      POSTGRES_USER: akubelajar
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - pgdata:/var/lib/postgresql/data
    ports:
      - "5432:5432"

  redis:
    image: redis:7-alpine
    command: redis-server --requirepass ${REDIS_PASSWORD}
    ports:
      - "6379:6379"

  api:
    build:
      context: .
      dockerfile: Dockerfile.backend
    env_file: .env
    depends_on:
      - postgres
      - redis
    ports:
      - "8080:8080"

  web:
    build:
      context: ./frontend
      dockerfile: Dockerfile.frontend
    environment:
      - NEXT_PUBLIC_API_URL=http://api:8080
    depends_on:
      - api
    ports:
      - "3000:3000"

volumes:
  pgdata:
```

---

## 3. Kubernetes Deploy (Production)

### Namespace

```bash
kubectl create namespace akubelajar-prod
```

### Deployment (Backend)

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: akubelajar-api
  namespace: akubelajar-prod
spec:
  replicas: 3
  selector:
    matchLabels:
      app: akubelajar-api
  template:
    spec:
      containers:
        - name: api
          image: registry.akubelajar.id/api:v2.0.0
          ports:
            - containerPort: 8080
          envFrom:
            - secretRef:
                name: akubelajar-secrets
          resources:
            requests:
              cpu: 250m
              memory: 128Mi
            limits:
              cpu: 500m
              memory: 256Mi
          livenessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 10
          readinessProbe:
            httpGet:
              path: /ready
              port: 8080
            initialDelaySeconds: 5
```

---

## 4. Health Check Endpoints

| Endpoint | Purpose | Expected Response |
|:---|:---|:---|
| `GET /health` | Liveness — apakah service hidup | `200: {"status":"ok"}` |
| `GET /ready` | Readiness — apakah service siap terima traffic | `200: {"status":"ready","db":"ok","redis":"ok"}` |

---

## 5. Post-Deploy Checklist

- [ ] Health check semua endpoints return 200
- [ ] Database migration sudah dijalankan
- [ ] SSL certificate valid
- [ ] Monitoring dashboard menunjukkan metrics normal
- [ ] Smoke test: login → dashboard → create quiz → submit
- [ ] Rollback plan sudah disiapkan

---

## Referensi

- [Zero Downtime Deploy](ZERO_DOWNTIME_DEPLOY.md)
- [Environment Variables](ENVIRONMENT_VARIABLES.md)
- [Backup & Recovery](BACKUP_AND_RECOVERY.md)

---

*Terakhir diperbarui: 21 Maret 2026*
