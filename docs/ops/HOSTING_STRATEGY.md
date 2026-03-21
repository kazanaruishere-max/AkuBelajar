# ☁️ Hosting Strategy — AkuBelajar (Demo/Portfolio)

> Stack hosting 100% gratis untuk project demo. Total biaya: Rp 0/bulan.

---

## Arsitektur

```mermaid
flowchart LR
    U["User Browser"] --> V["Vercel\nNext.js Frontend\nakubelajar.vercel.app"]
    V --> R["Render\nGo API + WebSocket\nakubelajar-api.onrender.com"]
    R --> S["Supabase\nPostgreSQL 500MB"]
    R --> UP["Upstash\nRedis Serverless"]
    R --> SS["Supabase Storage\nFile Upload 1GB"]
    R --> G["Gemini AI\nFree Tier"]
    R --> RE["Resend\nEmail 100/day"]
```

---

## Stack Detail

| Layer | Platform | URL | Free Tier |
|:---|:---|:---|:---|
| **Frontend** | Vercel | `akubelajar.vercel.app` | 100GB bandwidth, unlimited deploys |
| **Backend API** | Render | `akubelajar-api.onrender.com` | 750 jam/bulan, 512MB RAM |
| **Database** | Supabase | Direct connection string | PostgreSQL 500MB, 50K rows |
| **Cache** | Upstash | REST API Redis | 10K commands/day, 256MB |
| **File Storage** | Supabase Storage | Via API | 1GB, 50MB max file |
| **AI** | Google Gemini | `gemini-2.0-flash` | 15 req/min, 1500 req/day |
| **Email** | Resend | API | 100 emails/day |
| **WA Notifikasi** | ❌ Tidak dipakai | — | Diganti in-app + email |
| **DNS** | Vercel built-in | — | Auto SSL |

---

## Environment Variables

### Vercel (Frontend)

```bash
NEXT_PUBLIC_API_URL=https://akubelajar-api.onrender.com/api/v1
NEXT_PUBLIC_WS_URL=wss://akubelajar-api.onrender.com/ws
NEXT_PUBLIC_APP_URL=https://akubelajar.vercel.app
```

### Render (Backend)

```bash
# Database — dari Supabase dashboard
DATABASE_URL=postgresql://postgres.[ref]:[password]@aws-0-ap-southeast-1.pooler.supabase.com:6543/postgres

# Redis — dari Upstash dashboard
REDIS_URL=rediss://default:[password]@[region].upstash.io:6379

# Supabase Storage
SUPABASE_URL=https://[ref].supabase.co
SUPABASE_ANON_KEY=[anon_key]

# Auth
PASETO_KEY=[generate: openssl rand -hex 16]

# AI
GEMINI_API_KEY=[dari Google AI Studio — gratis]

# Email
RESEND_API_KEY=[dari resend.com — gratis]

# App
PORT=8080
ENV=production
CORS_ORIGIN=https://akubelajar.vercel.app
```

---

## Deploy Flow

### Frontend (Vercel) — Auto Deploy

```
1. Push ke GitHub branch `main`
2. Vercel otomatis build & deploy (< 1 menit)
3. Preview deploy otomatis untuk setiap PR
```

Vercel settings:
- Framework: Next.js
- Build command: `npm run build`
- Output directory: `.next`
- Root directory: `frontend/`

### Backend (Render) — Auto Deploy

```
1. Push ke GitHub branch `main`
2. Render otomatis build & deploy (3-5 menit)
3. Health check: GET /health
```

Render settings:
- Environment: Docker
- Dockerfile path: `backend/Dockerfile`
- Health check path: `/health`
- Auto-deploy: Yes (from `main` branch)

---

## Keterbatasan Free Tier & Solusi

| Keterbatasan | Impact | Solusi |
|:---|:---|:---|
| Render cold start ~30 detik | API lambat setelah idle 15 menit | Buka URL API dulu sebelum demo |
| Supabase pause setelah 7 hari idle | DB tidak bisa diakses | Login ke Supabase dashboard 1×/minggu |
| Upstash 10K commands/day | Rate limit & session overflow | Cukup untuk < 20 concurrent users |
| Resend 100 email/day | Limit notifikasi email | Fokus ke in-app notification |
| Gemini 15 req/min | AI quiz generation terbatas | Cukup untuk demo (1 guru generate) |
| Supabase Storage 1GB | File upload limit | Cukup untuk ~500 file submission |

---

## Upgrade Path (Jika Mau Production)

| Komponen | Free → Paid | Biaya/Bulan |
|:---|:---|:---|
| Render | Free → Starter | $7 (~Rp 110.000) |
| Supabase | Free → Pro | $25 (~Rp 400.000) |
| Upstash | Free → Pay-as-you-go | ~$5 (~Rp 80.000) |
| Fonnte (WA) | — → Active | Rp 25.000 |
| Custom domain | — → `.id` | Rp 200.000/tahun |
| **Total** | | **~Rp 630.000/bulan** |

---

*Terakhir diperbarui: 21 Maret 2026*
