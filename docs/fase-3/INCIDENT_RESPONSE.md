# 🚨 Incident Response — AkuBelajar

> Playbook: apa yang dilakukan jika server diserang pukul 3 pagi — siapa yang dihubungi, langkah apa yang diambil.

---

## Severity Levels

| Level | Deskripsi | Contoh | Response Time |
|:---|:---|:---|:---|
| **P0 — Critical** | Data breach, sistem down total | Database compromised, RCE | ≤ 15 menit |
| **P1 — High** | Fitur kritis down, kerentanan aktif dieksploitasi | Auth bypass, DDoS ongoing | ≤ 1 jam |
| **P2 — Medium** | Fitur non-kritis down, kerentanan terdeteksi | Stored XSS, minor DoS | ≤ 4 jam |
| **P3 — Low** | Anomali terdeteksi, belum ada dampak | Unusual traffic spike, failed login spike | ≤ 24 jam |

---

## Escalation Chain

```
Deteksi Alert (Grafana/Prometheus)
        │
        ▼
[1] On-Call Engineer (Primary)
    📱 Kazanaru — +62-xxx-xxx-xxxx
        │
        ├─── P3/P2: Handle sendiri
        │
        ▼ (jika P1/P0)
[2] Tech Lead
    📱 [Nama] — +62-xxx-xxx-xxxx
        │
        ▼ (jika data breach)
[3] Legal & Compliance
    📧 legal@akubelajar.id
        │
        ▼ (jika perlu komunikasi publik)
[4] PR / Communications
    📧 pr@akubelajar.id
```

---

## Response Playbook

### Fase 1: Deteksi & Triase (0-15 menit)

```
□ Konfirmasi alert bukan false positive
□ Tentukan severity level (P0/P1/P2/P3)
□ Buka incident channel di Slack/Discord
□ Assign incident commander
□ Mulai timeline log (catat semua tindakan + timestamp)
```

### Fase 2: Containment (15-60 menit)

```
□ Isolasi komponen yang terdampak
  - Jika auth breach → rotate semua JWT secrets
  - Jika database breach → revoke semua active sessions
  - Jika DDoS → aktifkan Cloudflare "Under Attack" mode
  - Jika malware di upload → isolasi MinIO bucket
□ Block IP penyerang di WAF
□ Preserve evidence (jangan hapus log)
□ Notifikasi stakeholder internal
```

### Fase 3: Eradication (1-4 jam)

```
□ Identifikasi root cause
□ Patch vulnerability
□ Scan sistem untuk indicator of compromise (IoC) lainnya
□ Verifikasi patch di staging
□ Deploy fix ke production
```

### Fase 4: Recovery (4-24 jam)

```
□ Restore service ke operasional normal
□ Monitor closely selama 24 jam
□ Verifikasi integritas data
□ Re-enable fitur yang di-disable saat containment
□ Konfirmasi semua alert sudah clear
```

### Fase 5: Post-Mortem (24-72 jam)

```
□ Tulis post-mortem document
  - Timeline lengkap
  - Root cause analysis (5 Whys)
  - Apa yang berjalan baik
  - Apa yang perlu diperbaiki
  - Action items dengan deadline
□ Share post-mortem ke tim
□ Update threat model
□ Update monitoring/alerting rules
□ Jadwalkan follow-up meeting
```

---

## Communication Templates

### Internal (Slack/Discord)

```
🚨 INCIDENT DECLARED — [SEVERITY]
━━━━━━━━━━━━━━━━━━━━━━━━━━━━
What: [Deskripsi singkat]
Impact: [Siapa/apa yang terdampak]
Commander: [Nama]
Status: INVESTIGATING / MITIGATING / RESOLVED
Channel: #incident-YYYY-MM-DD
━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

### Eksternal (jika diperlukan — data breach)

```
Kepada pengguna AkuBelajar,

Kami mendeteksi [deskripsi insiden] pada [tanggal].
[Langkah yang telah kami ambil].
[Langkah yang kami sarankan untuk pengguna].

Tim Keamanan AkuBelajar
```

---

## Tools & Access

| Tool | Akses | Purpose |
|:---|:---|:---|
| Grafana Dashboard | `monitoring.akubelajar.id` | Real-time metrics |
| Cloudflare Dashboard | Via SSO | WAF rules, DDoS mitigation |
| PostgreSQL (read-only) | Via bastion host | Data investigation |
| Server SSH | Via VPN + bastion | Log analysis |
| HashiCorp Vault | On-call only | Secret rotation |

---

## Referensi

- [Threat Model](THREAT_MODEL.md)
- [Audit Log Spec](AUDIT_LOG_SPEC.md)
- [Security Policy](../fase-0/SECURITY.md)

---

*Terakhir diperbarui: 21 Maret 2026*
