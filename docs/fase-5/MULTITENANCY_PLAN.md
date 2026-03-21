# 🏢 Multitenancy Plan — AkuBelajar

> Roadmap teknis menuju SaaS multi-sekolah agar tidak perlu refactor total di kemudian hari.

---

## Status Saat Ini

| Aspek | Status |
|:---|:---|
| Tabel `schools` | ✅ Sudah ada |
| `school_id` di semua tabel utama | ✅ Sudah ada |
| RLS per sekolah | ✅ Sudah ada |
| Onboarding per sekolah | 🟡 Manual |
| Billing & subscription | ❌ Belum ada |
| Custom domain per sekolah | ❌ Belum ada |
| Config per sekolah | 🟡 Partial (JSONB `config`) |

---

## Strategi: Shared Database, Schema Isolation

```
┌───────────────────────────────────────────┐
│              Shared Infrastructure         │
│                                            │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐  │
│  │ Sekolah A│ │ Sekolah B│ │ Sekolah C│  │
│  │ (tenant) │ │ (tenant) │ │ (tenant) │  │
│  └────┬─────┘ └────┬─────┘ └────┬─────┘  │
│       │             │             │        │
│  ┌────▼─────────────▼─────────────▼────┐  │
│  │          PostgreSQL + RLS            │  │
│  │   Semua tenant dalam 1 database     │  │
│  │   Isolasi via school_id + RLS       │  │
│  └─────────────────────────────────────┘  │
└───────────────────────────────────────────┘
```

### Mengapa Shared Database?

| Pendekatan | Pro | Kontra | AkuBelajar |
|:---|:---|:---|:---|
| DB per tenant | Isolasi total | Rumit maintain 1000+ DB | ❌ |
| Schema per tenant | Isolasi baik | Migration rumit | ❌ |
| **Shared DB + RLS** | Simple, scalable | Perlu RLS ketat | ✅ Dipilih |

---

## Roadmap Implementasi

### Phase 1: Foundation (✅ Done — v2.0.0)

- [x] Tabel `schools` sebagai tenant anchor
- [x] `school_id` foreign key di semua tabel utama
- [x] RLS policies per `school_id`
- [x] Per-school config via JSONB column

### Phase 2: Self-Service Onboarding (Q3 2026)

- [ ] Wizard setup sekolah baru (via UI admin)
- [ ] Otomatis: buat admin sekolah, seed data master
- [ ] Custom branding per sekolah (logo, warna, nama)
- [ ] Subdomain atau path-based routing (`sekolaha.akubelajar.id`)

### Phase 3: Subscription & Billing (Q4 2026)

- [ ] Tier pricing (Free, Basic, Premium, Enterprise)
- [ ] Feature gating per tier
- [ ] Usage quota (jumlah siswa, storage, AI calls)
- [ ] Payment integration (Midtrans/Xendit)

### Phase 4: Advanced Isolation (Q1 2027)

- [ ] Data export per tenant (full JSON/SQL dump)
- [ ] Tenant suspension/archival
- [ ] Cross-tenant analytics (untuk platform admin)
- [ ] Dedicated resources untuk tenant Enterprise

---

## Tier Pricing (Proposal)

| Feature | 🆓 Free | 💼 Basic | ⭐ Premium | 🏢 Enterprise |
|:---|:---|:---|:---|:---|
| Siswa | ≤ 100 | ≤ 500 | ≤ 2.000 | Unlimited |
| Storage | 1 GB | 10 GB | 50 GB | Custom |
| AI Quiz/bulan | 50 | 500 | 5.000 | Unlimited |
| Custom branding | ❌ | ✅ | ✅ | ✅ |
| Custom domain | ❌ | ❌ | ✅ | ✅ |
| Priority support | ❌ | ❌ | ✅ | ✅ (dedicated) |
| SLA | — | 99.5% | 99.9% | 99.99% |

---

## Data Isolation Checklist

Setiap fitur baru harus melewati checklist ini:

- [ ] Query menggunakan `school_id` filter?
- [ ] RLS policy sudah di-apply?
- [ ] API response tidak leak data tenant lain?
- [ ] File upload di-isolasi per `school_id`?
- [ ] Cache key menyertakan `school_id`?

---

## Referensi

- [Database Schema](../fase-1/DATABASE_SCHEMA.md)
- [Scalability Plan](../fase-1/SCALABILITY_PLAN.md)
- [User Stories](USER_STORIES.md)

---

*Terakhir diperbarui: 21 Maret 2026*
