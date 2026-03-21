# 📈 Scalability Plan — AkuBelajar

> Target skala, roadmap teknis horizontal scaling, dan strategi pertumbuhan infrastruktur.

---

## Target Skala per Fase

| Fase | Timeline | Users Concurrent | Sekolah | Infra |
|:---|:---|:---|:---|:---|
| **MVP** | Q1 2026 | 500 | 1-3 | Single Server + Docker Compose |
| **Growth** | Q3 2026 | 5.000 | 10-30 | Multi-server + Load Balancer |
| **Scale** | Q1 2027 | 50.000 | 100-300 | Kubernetes Cluster |
| **Enterprise** | Q3 2027 | 200.000+ | 1.000+ | Multi-region K8s + CDN |

---

## Strategi Scaling per Layer

### Backend (Go Services)

| Strategi | Detail |
|:---|:---|
| **Horizontal Pod Autoscaler** | Scale berdasarkan CPU (>70%) dan RPS (>1000/pod) |
| **Stateless Design** | Session di Redis — setiap pod interchangeable |
| **Connection Pooling** | pgx pool: min=5, max=50 per pod |
| **Circuit Breaker** | Prevent cascade failure ke Gemini API |

### Database (PostgreSQL)

| Strategi | Detail |
|:---|:---|
| **Read Replicas** | 1 primary + 2 read replicas (query berat ke replica) |
| **Connection Pooling** | PgBouncer di depan PostgreSQL |
| **Partitioning** | Tabel `audit_logs` di-partition per bulan |
| **Materialized Views** | Dashboard analytics di-refresh setiap 5 menit |
| **Archival** | Data > 2 tahun dipindah ke cold storage |

### Frontend (Next.js)

| Strategi | Detail |
|:---|:---|
| **CDN** | Static assets via Cloudflare edge cache |
| **ISR** | Incremental Static Regeneration untuk halaman semi-static |
| **Code Splitting** | Dynamic import untuk komponen berat |
| **Image CDN** | next/image dengan loader Cloudflare Images |

### Cache (Redis)

| Strategi | Detail |
|:---|:---|
| **Redis Cluster** | 3 nodes minimum untuk HA |
| **TTL Strategy** | Session: 24h, Rate limit: sliding window, Cache: 5-30 min |
| **Eviction** | `allkeys-lru` (Least Recently Used) |

---

## Bottleneck yang Diantisipasi

| Bottleneck | Trigge | Mitigasi |
|:---|:---|:---|
| DB Connection exhaustion | >100 pods concurrent | PgBouncer + connection limit per service |
| Gemini API rate limit | >60 RPM (free tier) | Queue + retry with exponential backoff |
| PDF generation memory | >1000 rapor sekaligus | Worker pool + memory limit per goroutine |
| WebSocket connections | >10K concurrent | Sticky sessions + Redis Pub/Sub fan-out |

---

## Referensi Terkait

- [Performance Budget](PERFORMANCE_BUDGET.md)
- [Zero Downtime Deploy](../fase-4/ZERO_DOWNTIME_DEPLOY.md)
- [Multitenancy Plan](../fase-5/MULTITENANCY_PLAN.md)

---

*Terakhir diperbarui: 21 Maret 2026*
