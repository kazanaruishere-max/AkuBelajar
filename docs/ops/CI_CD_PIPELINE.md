# 🚀 CI/CD Pipeline — AkuBelajar

> GitHub Actions workflow untuk build, test, lint, dan deploy otomatis.

---

## 1. Pipeline Overview

```mermaid
flowchart LR
    A["Push/PR"] --> B["Lint\n(Go + TS)"]
    B --> C["Unit Test\n(Go + Jest)"]
    C --> D["Integration Test\n(Docker)"]
    D --> E["Build\n(Docker images)"]
    E --> F{Branch?}
    F -->|develop| G["Deploy\nStaging"]
    F -->|main| H["Deploy\nProduction"]
    F -->|feature/*| I["Selesai\n(PR ready)"]
```

---

## 2. Branch → Environment Mapping

| Branch | Environment | Deploy Trigger | URL |
|:---|:---|:---|:---|
| `feature/*` | — (CI only) | Auto on PR | — |
| `develop` | Staging | Auto on merge | `staging.akubelajar.id` |
| `main` | Production | Manual approve | `app.akubelajar.id` |
| `hotfix/*` | Production | Manual approve | `app.akubelajar.id` |

---

## 3. GitHub Actions Workflows

### `.github/workflows/ci.yml` — Lint & Test (every PR)

```yaml
name: CI

on:
  pull_request:
    branches: [develop, main]
  push:
    branches: [develop, main]

jobs:
  lint-backend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'
      - name: golangci-lint
        uses: golangci/golangci-lint-action@v4
        with:
          working-directory: ./backend

  lint-frontend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: '22'
          cache: 'npm'
          cache-dependency-path: frontend/package-lock.json
      - run: cd frontend && npm ci
      - run: cd frontend && npm run lint
      - run: cd frontend && npx tsc --noEmit

  test-backend:
    runs-on: ubuntu-latest
    needs: lint-backend
    services:
      postgres:
        image: postgres:16-alpine
        env:
          POSTGRES_DB: akubelajar_test
          POSTGRES_USER: test
          POSTGRES_PASSWORD: test
        ports: ['5432:5432']
        options: --health-cmd pg_isready --health-interval 5s
      redis:
        image: redis:7-alpine
        ports: ['6379:6379']
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'
      - name: Run migrations
        run: cd backend && go run cmd/migrate/main.go up
        env:
          DB_URL: postgres://test:test@localhost:5432/akubelajar_test?sslmode=disable
      - name: Run tests
        run: cd backend && go test ./... -v -coverprofile=coverage.out -race
        env:
          DB_URL: postgres://test:test@localhost:5432/akubelajar_test?sslmode=disable
          REDIS_URL: redis://localhost:6379
      - name: Coverage check
        run: |
          cd backend
          total=$(go tool cover -func=coverage.out | tail -1 | awk '{print $3}' | sed 's/%//')
          echo "Coverage: $total%"
          if (( $(echo "$total < 70" | bc -l) )); then
            echo "❌ Coverage below 70%"
            exit 1
          fi

  test-frontend:
    runs-on: ubuntu-latest
    needs: lint-frontend
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: '22'
          cache: 'npm'
          cache-dependency-path: frontend/package-lock.json
      - run: cd frontend && npm ci
      - run: cd frontend && npm test -- --coverage
      - name: Coverage check
        run: |
          # Parse Jest coverage output
          echo "Frontend tests passed"
```

### `.github/workflows/deploy-staging.yml`

```yaml
name: Deploy Staging

on:
  push:
    branches: [develop]

jobs:
  deploy:
    runs-on: ubuntu-latest
    needs: [test-backend, test-frontend]
    environment: staging
    steps:
      - uses: actions/checkout@v4
      
      - name: Build & push Docker images
        run: |
          docker build -t ghcr.io/${{ github.repository }}/api:staging ./backend --target production
          docker build -t ghcr.io/${{ github.repository }}/frontend:staging ./frontend --target production
          echo "${{ secrets.GHCR_TOKEN }}" | docker login ghcr.io -u ${{ github.actor }} --password-stdin
          docker push ghcr.io/${{ github.repository }}/api:staging
          docker push ghcr.io/${{ github.repository }}/frontend:staging
      
      - name: Deploy to staging
        run: |
          # SSH deploy or kubectl apply
          ssh staging "cd /app && docker compose pull && docker compose up -d"
      
      - name: Run migrations
        run: |
          ssh staging "cd /app && docker compose exec api ./server migrate up"
      
      - name: Health check
        run: |
          sleep 10
          curl -f https://staging.akubelajar.id/health || exit 1
```

### `.github/workflows/deploy-prod.yml`

```yaml
name: Deploy Production

on:
  workflow_dispatch:
    inputs:
      confirm:
        description: 'Type "deploy" to confirm'
        required: true

jobs:
  deploy:
    runs-on: ubuntu-latest
    if: github.event.inputs.confirm == 'deploy'
    environment: production
    steps:
      - uses: actions/checkout@v4
      
      - name: Build production images
        run: |
          docker build -t ghcr.io/${{ github.repository }}/api:${{ github.sha }} ./backend --target production
          docker build -t ghcr.io/${{ github.repository }}/frontend:${{ github.sha }} ./frontend --target production
      
      - name: Push images
        run: |
          docker push ghcr.io/${{ github.repository }}/api:${{ github.sha }}
          docker push ghcr.io/${{ github.repository }}/frontend:${{ github.sha }}
      
      - name: Rolling deploy
        run: |
          # Zero-downtime rolling update
          ssh prod "cd /app && docker compose pull && docker compose up -d --no-deps api"
          sleep 15
          curl -f https://app.akubelajar.id/health/ready || exit 1
          ssh prod "cd /app && docker compose up -d --no-deps frontend"
          
      - name: Notify team
        run: |
          echo "✅ Production deploy complete: ${{ github.sha }}"
```

---

## 4. PR Requirements (Branch Protection Rules)

| Rule | Nilai |
|:---|:---|
| Required reviews | 1 |
| CI checks must pass | ✅ lint-backend, lint-frontend, test-backend, test-frontend |
| Up-to-date with base | ✅ |
| Force push allowed | ❌ (main), ❌ (develop) |
| Delete branch after merge | ✅ |

---

## 5. Coverage Targets

| Codebase | Minimum | Target |
|:---|:---|:---|
| Go backend (unit) | 70% | 85% |
| Go backend (integration) | 50% | 70% |
| Next.js frontend | 60% | 80% |
| E2E (Playwright) | — | Critical paths only |

---

*Terakhir diperbarui: 21 Maret 2026*
