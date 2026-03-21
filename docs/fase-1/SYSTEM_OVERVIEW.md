# 🏗️ System Overview — AkuBelajar

> Blueprint arsitektur keseluruhan dalam satu halaman.

---

## Daftar Isi

1. [Filosofi Arsitektur](#filosofi-arsitektur)
2. [Diagram Arsitektur](#diagram-arsitektur)
3. [Layer Breakdown](#layer-breakdown)
4. [Alur Data (Data Flow)](#alur-data-data-flow)
5. [Komunikasi Antar Service](#komunikasi-antar-service)
6. [Keputusan Arsitektur Kunci](#keputusan-arsitektur-kunci)

---

## Filosofi Arsitektur

AkuBelajar dibangun dengan tiga prinsip utama:

| Prinsip | Implementasi |
|:---|:---|
| **Separation of Concerns** | Setiap layer (client, gateway, service, data) memiliki tanggung jawab yang jelas |
| **Defense in Depth** | Keamanan berlapis di setiap level — WAF, RBAC, RLS, encryption |
| **Horizontal Scalability** | Setiap service dapat di-scale secara independen via Kubernetes |

---

## Diagram Arsitektur

```
┌──────────────────────────────────────────────────────────────────────┐
│                         CLIENT LAYER                                  │
│                                                                       │
│   ┌─────────────────┐   ┌─────────────────┐   ┌──────────────────┐  │
│   │  Next.js 15     │   │  PWA (Mobile)   │   │  React Native    │  │
│   │  Web App        │   │  Service Worker │   │  (Expo) [Q3'26]  │  │
│   └────────┬────────┘   └────────┬────────┘   └────────┬─────────┘  │
└────────────┼────────────────────┼────────────────────┼───────────────┘
             │                    │                    │
             └────────────────────┼────────────────────┘
                                  │ HTTPS / WSS
┌─────────────────────────────────▼────────────────────────────────────┐
│                       GATEWAY LAYER                                   │
│                                                                       │
│   ┌─────────────┐  ┌──────────────┐  ┌───────────────┐              │
│   │  Cloudflare  │  │  Nginx       │  │  Rate Limiter │              │
│   │  WAF + CDN   │→ │  Reverse     │→ │  (Redis)      │              │
│   │              │  │  Proxy + TLS │  │               │              │
│   └─────────────┘  └──────────────┘  └───────┬───────┘              │
└──────────────────────────────────────────────┼───────────────────────┘
                                               │
                    ┌──────────────────────────┼──────────────┐
                    │                          │              │
┌───────────────────▼───┐  ┌───────────────────▼───┐  ┌──────▼──────────┐
│    AUTH SERVICE        │  │    CORE API SERVICE    │  │   AI SERVICE     │
│                        │  │                        │  │                  │
│  • JWT/Paseto          │  │  • Academic Module     │  │  • Gemini 2.0    │
│  • Login/Register      │  │  • Quiz/CBT Engine     │  │  • Quiz Gen      │
│  • OAuth 2.0           │  │  • Attendance          │  │  • Essay Grading  │
│  • Session Mgmt        │  │  • Grading             │  │  • Early Warning  │
│                        │  │  • Notification        │  │                  │
│  [Go + Gin]            │  │  [Go + Gin]            │  │  [Go + Gemini]   │
└───────────┬────────────┘  └───────────┬────────────┘  └────────┬────────┘
            │                           │                        │
            └───────────────────────────┼────────────────────────┘
                                        │
              ┌─────────────────────────┼─────────────────────┐
              │                         │                     │
┌─────────────▼──────┐  ┌──────────────▼──────┐  ┌──────────▼─────────┐
│   PostgreSQL 16+   │  │    Redis 7+         │  │   MinIO (S3)       │
│                    │  │                     │  │                    │
│  • Primary + Read  │  │  • Session Cache    │  │  • File Storage    │
│    Replica         │  │  • Rate Limit       │  │  • PDF Rapor       │
│  • RLS Policies    │  │  • Job Queue        │  │  • Upload Materi   │
│  • UUID v7 PKs     │  │  • Real-time Pub/   │  │                    │
│  • Immutable Audit │  │    Sub              │  │                    │
└────────────────────┘  └─────────────────────┘  └────────────────────┘
```

---

## Layer Breakdown

### 1. Client Layer

| Komponen | Teknologi | Peran |
|:---|:---|:---|
| Web App | Next.js 15 + TypeScript | SSR/ISR/CSR hybrid rendering |
| PWA | Service Worker + Workbox | Offline support, push notification |
| Mobile App | React Native (Expo) | Native access (kamera, biometrik) — Roadmap Q3 2026 |

### 2. Gateway Layer

| Komponen | Teknologi | Peran |
|:---|:---|:---|
| CDN + WAF | Cloudflare | DDoS protection, edge caching, bot mitigation |
| Reverse Proxy | Nginx | TLS termination, load balancing, request routing |
| Rate Limiter | Redis Sliding Window | Per-endpoint, per-user rate limiting |

### 3. Service Layer (Go Microservices)

| Service | Tanggung Jawab | Port Default |
|:---|:---|:---|
| Auth Service | Autentikasi, otorisasi, session | `:8081` |
| Core API | Semua endpoint bisnis (CRUD) | `:8080` |
| AI Service | Integrasi Gemini, quiz gen, essay grading | `:8082` |
| Worker | Background jobs (email, PDF, import) | `:8083` |
| WebSocket | Real-time notifikasi, CBT heartbeat | `:8084` |

### 4. Data Layer

| Komponen | Teknologi | Peran |
|:---|:---|:---|
| RDBMS | PostgreSQL 16+ | Source of truth, ACID, RLS |
| Cache | Redis 7+ | Session, rate limit, pub/sub, queue |
| Object Storage | MinIO (S3-compatible) | File, PDF, avatar, materi |
| Secret Store | HashiCorp Vault | API keys, DB credentials, certificates |

---

## Alur Data (Data Flow)

### Request Lifecycle

```
Client Request
    ↓
[1] Cloudflare WAF        → Filter malicious traffic
    ↓
[2] Nginx Reverse Proxy   → TLS termination, route to service
    ↓
[3] Rate Limiter (Redis)  → Check per-user/per-IP limits
    ↓
[4] Auth Middleware        → Validate JWT/Paseto token
    ↓
[5] RBAC Middleware        → Check role permissions
    ↓
[6] Input Validation       → Sanitize & validate with struct tags
    ↓
[7] Business Logic         → Service layer processing
    ↓
[8] Database Query         → RLS filter applied automatically
    ↓
[9] Audit Log              → Immutable event recorded
    ↓
[10] Response              → JSON response with proper status code
```

---

## Komunikasi Antar Service

| Pola | Implementasi | Use Case |
|:---|:---|:---|
| **Synchronous** | HTTP REST (JSON) | CRUD operations, auth validation |
| **Asynchronous** | Redis Queue (RPUSH/BLPOP) | Email sending, PDF generation, bulk import |
| **Real-time** | WebSocket (gorilla/websocket) | CBT heartbeat, live notifications |
| **Event-driven** | Redis Pub/Sub | Cache invalidation, broadcast notifications |

---

## Keputusan Arsitektur Kunci

Untuk detail lengkap, lihat folder **[ADR/](ADR/)**.

| Keputusan | Pilihan | Alasan Singkat |
|:---|:---|:---|
| Backend Language | Go (bukan Node.js) | Performa, type safety, binary deployment |
| Token Format | Paseto + JWT | Paseto lebih aman, JWT untuk kompatibilitas |
| Primary Key | UUID v7 (bukan auto-increment) | Anti-IDOR, time-sortable |
| ORM | Tanpa ORM (pgx raw) | Full SQL control, performa maksimal |
| Frontend Framework | Next.js 15 | SSR + ISR + CSR hybrid, React ecosystem |
| State Management | Zustand | Minimal boilerplate, type-safe |

---

## Referensi Terkait

- [ADR: Paseto vs JWT](ADR/001-paseto-vs-jwt.md)
- [Database Schema](DATABASE_SCHEMA.md)
- [Scalability Plan](SCALABILITY_PLAN.md)
- [Performance Budget](PERFORMANCE_BUDGET.md)

---

*Terakhir diperbarui: 21 Maret 2026*
