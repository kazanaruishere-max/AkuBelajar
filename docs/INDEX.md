# 📚 AKUBELAJAR — Indeks Dokumentasi Pengembangan

> **Platform Edukasi Digital Generasi Berikutnya**  
> AI-First • Cross-Platform • Enterprise-Grade Security

| Metadata | Detail |
|:---|:---|
| **Proyek** | AkuBelajar |
| **Maintainer** | Kazanaru |
| **Dibuat** | 21 Maret 2026 |
| **Versi Docs** | 3.0.0 |

---

## 🤖 File Root — AI Agent Context (Baca Pertama!)

| # | File | Lokasi | Deskripsi |
|:-:|:---|:---|:---|
| 1 | [AGENTS.md](../AGENTS.md) | ★ Root | Instruksi langsung ke AI agent: rules, scope, urutan kerja |
| 2 | [PROJECT_CONTEXT.md](../PROJECT_CONTEXT.md) | ★ Root | Ringkasan 1 halaman — baca setiap awal sesi |
| 3 | [DECISIONS_LOG.md](../DECISIONS_LOG.md) | ★ Root | 31 keputusan final — JANGAN ubah |
| 4 | [FORBIDDEN_PATTERNS.md](../FORBIDDEN_PATTERNS.md) | ★ Root | Anti-pattern yang DILARANG (kode contoh ❌/✅) |

---

## 📋 Produk & Prioritas

| # | Dokumen | Deskripsi |
|:-:|:---|:---|
| 1 | [FEATURE_PRIORITY_MATRIX.md](FEATURE_PRIORITY_MATRIX.md) | ★ Prioritas 40+ fitur: P0 (MVP) → P3 (Future) |
| 2 | [USER_PERSONAS.md](USER_PERSONAS.md) | ★ 4 persona: guru daerah, admin kota, siswa, ketua kelas |

---

## Navigasi Cepat per Fase

---

### 🏗️ Fase 0 — Pondasi & Perencanaan

| # | Dokumen | Deskripsi |
|:-:|:---|:---|
| 1 | [README.md](fase-0/README.md) | Overview, quick install, link ke semua docs |
| 2 | [CONTRIBUTING.md](fase-0/CONTRIBUTING.md) | Aturan kontribusi |
| 3 | [CHANGELOG.md](fase-0/CHANGELOG.md) | Catatan perubahan per versi |
| 4 | [SECURITY.md](fase-0/SECURITY.md) | Responsible disclosure |
| 5 | [CODE_OF_CONDUCT.md](fase-0/CODE_OF_CONDUCT.md) | Kode etik komunitas |
| 6 | [LICENSE.md](fase-0/LICENSE.md) | Lisensi legalitas |
| 7 | [AGENTS.md](fase-0/AGENTS.md) | 🤖 Instruksi untuk AI agent |

---

### 🧩 Fase 1 — Arsitektur & Desain Sistem

| # | Dokumen | Deskripsi |
|:-:|:---|:---|
| 1 | [SYSTEM_OVERVIEW.md](fase-1/SYSTEM_OVERVIEW.md) | Blueprint arsitektur |
| 2 | [ADR/](fase-1/ADR/) | Architecture Decision Records |
| 3 | [DATABASE_SCHEMA.md](fase-1/DATABASE_SCHEMA.md) | Skema PostgreSQL + RLS |
| 4 | [API_SPEC.md](fase-1/API_SPEC.md) | Kontrak API ringkas |
| 5 | [DATA_MODELS.md](fase-1/DATA_MODELS.md) | Mapping PostgreSQL ↔ Go ↔ TypeScript |
| 6 | [SCALABILITY_PLAN.md](fase-1/SCALABILITY_PLAN.md) | Target skala & roadmap |
| 7 | [PERFORMANCE_BUDGET.md](fase-1/PERFORMANCE_BUDGET.md) | Batasan performa |
| 8 | [OFFLINE_STRATEGY.md](fase-1/OFFLINE_STRATEGY.md) | Service Worker & cache |

---

### ⚙️ Fase 2 — Development & Standar Kode

| # | Dokumen | Deskripsi |
|:-:|:---|:---|
| 1 | [SETUP_LOCAL.md](fase-2/SETUP_LOCAL.md) | Setup environment < 30 menit |
| 2 | [CODING_STANDARDS.md](fase-2/CODING_STANDARDS.md) | Style guide Go, TypeScript, SQL |
| 3 | [GIT_WORKFLOW.md](fase-2/GIT_WORKFLOW.md) | Branching & commit format |
| 4 | [FOLDER_STRUCTURE.md](fase-2/FOLDER_STRUCTURE.md) | Di mana file harus diletakkan |
| 5 | [BACKEND_GUIDE.md](fase-2/BACKEND_GUIDE.md) | Panduan backend (Go) |
| 6 | [FRONTEND_GUIDE.md](fase-2/FRONTEND_GUIDE.md) | Panduan frontend (Next.js) |
| 7 | [MOBILE_GUIDE.md](fase-2/MOBILE_GUIDE.md) | Panduan mobile (PWA + RN) |
| 8 | [BUSINESS_LOGIC.md](fase-2/BUSINESS_LOGIC.md) | Aturan bisnis akademik |
| 9 | [TESTING_STRATEGY.md](fase-2/TESTING_STRATEGY.md) | Coverage target & tools |
| 10 | [TEST_CASES.md](fase-2/TEST_CASES.md) | Test case per fitur |
| 11 | [ACCEPTANCE_CRITERIA.md](fase-2/ACCEPTANCE_CRITERIA.md) | Kapan fitur "selesai" |
| 12 | [FEATURE_FLAG.md](fase-2/FEATURE_FLAG.md) | Rollout fitur bertahap |

---

### 🔒 Fase 3 — Keamanan & Compliance

| # | Dokumen | Deskripsi |
|:-:|:---|:---|
| 1 | [THREAT_MODEL.md](fase-3/THREAT_MODEL.md) | Attack surface (STRIDE) |
| 2 | [OWASP_CHECKLIST.md](fase-3/OWASP_CHECKLIST.md) | Audit OWASP Top 10 |
| 3 | [SECURE_CODING.md](fase-3/SECURE_CODING.md) | Rules kode aman |
| 4 | [INCIDENT_RESPONSE.md](fase-3/INCIDENT_RESPONSE.md) | Playbook insiden |
| 5 | [AUDIT_LOG_SPEC.md](fase-3/AUDIT_LOG_SPEC.md) | Audit log immutable |
| 6 | [PRIVACY_POLICY.md](fase-3/PRIVACY_POLICY.md) | UU PDP Indonesia |

---

### 🚀 Fase 4 — Deployment, CI/CD & Ops

| # | Dokumen | Deskripsi |
|:-:|:---|:---|
| 1 | [DEPLOYMENT_GUIDE.md](fase-4/DEPLOYMENT_GUIDE.md) | Docker → K8s deploy |
| 2 | [ZERO_DOWNTIME_DEPLOY.md](fase-4/ZERO_DOWNTIME_DEPLOY.md) | Rolling update |
| 3 | [ENVIRONMENT_VARIABLES.md](fase-4/ENVIRONMENT_VARIABLES.md) | Env var + Vault |
| 4 | [BACKUP_AND_RECOVERY.md](fase-4/BACKUP_AND_RECOVERY.md) | Backup & restore |

---

### 🌱 Fase 5 — Pertumbuhan Jangka Panjang

| # | Dokumen | Deskripsi |
|:-:|:---|:---|
| 1 | [UPGRADE_GUIDE.md](fase-5/UPGRADE_GUIDE.md) | Panduan upgrade versi |
| 2 | [MULTITENANCY_PLAN.md](fase-5/MULTITENANCY_PLAN.md) | Roadmap SaaS multi-sekolah |
| 3 | [USER_STORIES.md](fase-5/USER_STORIES.md) | Kebutuhan nyata pengguna |

---

## 📂 Dokumentasi Tambahan (dari AKUBELAJAR_PROMPTS.md)

### User Flows (`docs/flows/`)

| # | Dokumen | Deskripsi | Prioritas |
|:-:|:---|:---|:---|
| 1 | [ACCOUNT_CREATION_FLOW.md](flows/ACCOUNT_CREATION_FLOW.md) | Alur pembuatan akun per role | 🔴 KRITIS |
| 2 | [FIRST_LOGIN_FLOW.md](flows/FIRST_LOGIN_FLOW.md) | Onboarding: password change + profil wizard | 🔴 KRITIS |
| 3 | [PASSWORD_RESET_FLOW.md](flows/PASSWORD_RESET_FLOW.md) | Self-service OTP + admin reset | 🔴 KRITIS |
| 4 | [ATTENDANCE_FLOW.md](flows/ATTENDANCE_FLOW.md) | Absensi guru, ketua kelas, izin/sakit | 🔴 KRITIS |
| 5 | [ASSIGNMENT_FLOW.md](flows/ASSIGNMENT_FLOW.md) | Lifecycle tugas: create → submit → grade | 🔴 KRITIS |
| 6 | [ROLE_MANAGEMENT_FLOW.md](flows/ROLE_MANAGEMENT_FLOW.md) | Upgrade/downgrade/suspend/delete | 🔴 KRITIS |
| 7 | [CBT_DETAIL_FLOW.md](flows/CBT_DETAIL_FLOW.md) | Anti-cheat, reconnect, pengacakan soal | 🟡 PENTING |
| 8 | [NOTIFICATION_FLOW.md](flows/NOTIFICATION_FLOW.md) | Queue, template, failure handling | 🟡 PENTING |
| 9 | [ACADEMIC_YEAR_TRANSITION.md](flows/ACADEMIC_YEAR_TRANSITION.md) | ★ Transisi tahun ajaran, naik kelas, kelulusan | 🔴 KRITIS |
| 10 | [BULK_OPERATIONS_SPEC.md](flows/BULK_OPERATIONS_SPEC.md) | ★ Import Excel, batch processing, progress | 🟡 PENTING |
| 11 | [REPORT_CARD_GENERATION.md](flows/REPORT_CARD_GENERATION.md) | ★ Generate PDF rapor, QR verification | 🟡 PENTING |
| 12 | [OFFLINE_SYNC_FLOW.md](flows/OFFLINE_SYNC_FLOW.md) | ★ PWA offline sync, IndexedDB, conflict resolution | ⚪ P3 |
| 13 | [AI_QUIZ_GENERATION_FLOW.md](flows/AI_QUIZ_GENERATION_FLOW.md) | ★ Flowchart lengkap AI quiz: Gemini → DB | 🔴 KRITIS |
| 14 | [GRADE_CALCULATION_FLOW.md](flows/GRADE_CALCULATION_FLOW.md) | ★ Formula nilai, remedial, grade letter | 🔴 KRITIS |

### API Specifications (`docs/api/`)

| # | Dokumen | Deskripsi | Prioritas |
|:-:|:---|:---|:---|
| 1 | [API_SPEC_FULL.md](api/API_SPEC_FULL.md) | Kontrak API semua modul lengkap | 🔴 KRITIS |
| 2 | [WEBSOCKET_SPEC.md](api/WEBSOCKET_SPEC.md) | WebSocket events + reconnect | 🔴 KRITIS |
| 3 | [ERROR_CODE_CATALOGUE.md](api/ERROR_CODE_CATALOGUE.md) | Semua error code + pesan BI | 🔴 KRITIS |
| 4 | [AI_INTEGRATION_SPEC.md](api/AI_INTEGRATION_SPEC.md) | ★ Gemini AI: prompt, parsing, cost, safety | 🔴 KRITIS |
| 5 | [API_VERSIONING_STRATEGY.md](api/API_VERSIONING_STRATEGY.md) | ★ v1 stability, breaking changes, deprecation | 🔴 KRITIS |

### Database (`docs/database/`)

| # | Dokumen | Deskripsi | Prioritas |
|:-:|:---|:---|:---|
| 1 | [DATABASE_SCHEMA_FULL.md](database/DATABASE_SCHEMA_FULL.md) | DDL SQL 26 tabel lengkap | 🔴 KRITIS |
| 2 | [ERD.md](database/ERD.md) | Entity Relationship Diagram (Mermaid) | 🔴 KRITIS |
| 3 | [BUSINESS_LOGIC_FULL.md](database/BUSINESS_LOGIC_FULL.md) | Aturan bisnis: nilai, kehadiran, deadline | 🔴 KRITIS |
| 4 | [MIGRATION_STRATEGY.md](database/MIGRATION_STRATEGY.md) | ★ golang-migrate, expand-contract, rollback | 🔴 KRITIS |

### Security (`docs/security/`)

| # | Dokumen | Deskripsi | Prioritas |
|:-:|:---|:---|:---|
| 1 | [INPUT_VALIDATION_RULES.md](security/INPUT_VALIDATION_RULES.md) | Validasi per field (Zod + Go) | 🔴 KRITIS |
| 2 | [SESSION_AND_TOKEN_SPEC.md](security/SESSION_AND_TOKEN_SPEC.md) | Paseto v4 lifecycle + Redis blocklist | 🔴 KRITIS |
| 3 | [FILE_UPLOAD_SECURITY.md](security/FILE_UPLOAD_SECURITY.md) | Magic bytes, signed URL, malware scan | 🔴 KRITIS |
| 4 | [RATE_LIMIT_POLICY.md](security/RATE_LIMIT_POLICY.md) | Limit per endpoint, 3-layer defense | 🟡 PENTING |
| 5 | [PASSWORD_POLICY.md](security/PASSWORD_POLICY.md) | Argon2id, password history, temp password | 🟡 PENTING |
| 6 | [RBAC_MATRIX.md](security/RBAC_MATRIX.md) | ★ Permission matrix Role×Resource×Action + RLS | 🔴 KRITIS |

### Frontend (`docs/frontend/`)

| # | Dokumen | Deskripsi | Prioritas |
|:-:|:---|:---|:---|
| 1 | [SITEMAP.md](frontend/SITEMAP.md) | Semua route + role access + rendering | 🔴 KRITIS |
| 2 | [STATE_MANAGEMENT.md](frontend/STATE_MANAGEMENT.md) | Zustand + TanStack Query patterns | 🔴 KRITIS |
| 3 | [UI_COMPONENT_SPEC.md](frontend/UI_COMPONENT_SPEC.md) | ★ Design tokens, component taxonomy, a11y, dark mode | 🟡 PENTING |
| 4 | [WIREFRAMES.md](frontend/WIREFRAMES.md) | ★ ASCII wireframes per halaman per role | 🔴 KRITIS |

### Testing (`docs/testing/`)

| # | Dokumen | Deskripsi | Prioritas |
|:-:|:---|:---|:---|
| 1 | [ACCEPTANCE_CRITERIA_FULL.md](testing/ACCEPTANCE_CRITERIA_FULL.md) | Given/When/Then per fitur | 🔴 KRITIS |
| 2 | [EDGE_CASES.md](testing/EDGE_CASES.md) | 29 skenario ekstrem | 🔴 KRITIS |

### Operations (`docs/ops/`)

| # | Dokumen | Deskripsi | Prioritas |
|:-:|:---|:---|:---|
| 1 | [DOCKER_SETUP.md](ops/DOCKER_SETUP.md) | Docker Compose + Dockerfile + Makefile | 🔴 KRITIS |
| 2 | [SEED_DATA_SPEC.md](ops/SEED_DATA_SPEC.md) | Seed production + development + factory | 🟡 PENTING |
| 3 | [MONITORING_OBSERVABILITY.md](ops/MONITORING_OBSERVABILITY.md) | ★ Prometheus, Grafana, alerting rules | 🔴 KRITIS |
| 4 | [CI_CD_PIPELINE.md](ops/CI_CD_PIPELINE.md) | ★ GitHub Actions, branch strategy, coverage | 🔴 KRITIS |
| 5 | [LOGGING_STANDARD.md](ops/LOGGING_STANDARD.md) | ★ Structured JSON, masking, correlation ID | 🟡 PENTING |
| 6 | [HOSTING_STRATEGY.md](ops/HOSTING_STRATEGY.md) | ★ Vercel + Render + Supabase (Rp 0 demo) | 🔴 KRITIS |

---

## Statistik Dokumentasi

| Kategori | Jumlah File |
|:---|:---|
| Root AI Agent Context | 4 file |
| Produk & Prioritas | 2 file |
| Fase 0–5 (foundation) | 41 file |
| User Flows | 14 file |
| API Specs | 5 file |
| Database | 4 file |
| Security | 6 file |
| Frontend | 4 file |
| Testing | 2 file |
| Operations | 6 file |
| **Total** | **88 file** |

> ★ = file baru (22 file ditambahkan dari gap analysis + audit)

---

*Terakhir diperbarui: 21 Maret 2026*
