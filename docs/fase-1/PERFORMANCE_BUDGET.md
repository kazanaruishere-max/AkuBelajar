# ⚡ Performance Budget — AkuBelajar

> Batasan performa yang **tidak boleh dilanggar**. Setiap PR yang menyebabkan regresi di bawah threshold ini wajib diperbaiki sebelum merge.

---

## Core Web Vitals (Frontend)

| Metrik | Target | Threshold Kritis | Tool Pengukuran |
|:---|:---|:---|:---|
| **LCP** (Largest Contentful Paint) | < 2.5s | > 4.0s = blocker | Lighthouse CI |
| **FID** (First Input Delay) | < 100ms | > 300ms = blocker | Web Vitals SDK |
| **CLS** (Cumulative Layout Shift) | < 0.1 | > 0.25 = blocker | Lighthouse CI |
| **TTFB** (Time to First Byte) | < 200ms | > 600ms = blocker | Server monitoring |
| **FCP** (First Contentful Paint) | < 1.8s | > 3.0s = blocker | Lighthouse CI |

---

## API Response Time (Backend)

| Endpoint Category | P50 Target | P99 Target | Threshold Kritis |
|:---|:---|:---|:---|
| Auth (login/register) | < 100ms | < 300ms | > 500ms |
| CRUD sederhana | < 50ms | < 200ms | > 400ms |
| List/Search (paginated) | < 100ms | < 500ms | > 1s |
| AI Generation (quiz) | < 5s | < 15s | > 30s |
| PDF Generation | < 2s | < 10s | > 20s |
| File Upload (10MB) | < 3s | < 10s | > 15s |

---

## Bundle Size (Frontend)

| Metric | Budget | Tool |
|:---|:---|:---|
| Initial JS bundle | < 150 KB (gzipped) | `@next/bundle-analyzer` |
| Per-route JS | < 50 KB (gzipped) | Next.js build output |
| Total CSS | < 30 KB (gzipped) | Build output |
| Largest image | < 200 KB (WebP) | `next/image` auto-optimization |

---

## Database Query Performance

| Query Type | Target | Alert |
|:---|:---|:---|
| Simple SELECT (by PK) | < 5ms | > 20ms |
| JOIN (2-3 tables) | < 20ms | > 100ms |
| Aggregation (COUNT, AVG) | < 50ms | > 200ms |
| Full-text search | < 100ms | > 500ms |
| Bulk insert (500 rows) | < 500ms | > 2s |

---

## Resource Limits

| Resource | Per Service Pod | Alert Threshold |
|:---|:---|:---|
| Memory (Go service) | 128 MB | > 256 MB |
| CPU | 0.25 vCPU | > 0.5 vCPU sustained |
| Goroutines | < 1000 | > 5000 |
| DB Connections | max 50 | > 40 (80%) |
| Redis Connections | max 20 | > 16 (80%) |

---

## Enforcement

### CI Pipeline Checks

```yaml
# .github/workflows/perf-budget.yml
- name: Lighthouse CI
  run: lhci autorun --assert.preset=lighthouse:recommended

- name: Bundle Size Check
  run: npx bundlesize

- name: API Benchmark
  run: make bench-api
```

### Monitoring (Production)

- **Grafana Dashboard** — Real-time P50/P99 latency per endpoint
- **Prometheus Alerts** — Auto-notify jika threshold terlampaui
- **Weekly Performance Report** — Tren mingguan untuk deteksi regresi gradual

---

*Terakhir diperbarui: 21 Maret 2026*
