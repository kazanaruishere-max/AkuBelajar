# 💾 Backup & Recovery — AkuBelajar

> Jadwal backup, enkripsi, dan prosedur restore drill.

---

## Strategi Backup

| Komponen | Metode | Frekuensi | Retensi | Enkripsi |
|:---|:---|:---|:---|:---|
| PostgreSQL (full) | `pg_dump` (logical) | Harian 02:00 WIB | 30 hari | AES-256 |
| PostgreSQL (WAL) | WAL archiving (continuous) | Real-time | 7 hari | AES-256 |
| Redis (snapshot) | RDB dump | Setiap 6 jam | 7 hari | AES-256 |
| MinIO (files) | rsync ke backup storage | Harian 03:00 WIB | 90 hari | AES-256 |
| Konfigurasi (Vault) | Vault snapshot | Harian | 30 hari | Built-in |

---

## Prosedur Backup

### PostgreSQL Full Backup

```bash
#!/bin/bash
# scripts/backup-db.sh

TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="akubelajar_${TIMESTAMP}.sql.gz.enc"
BACKUP_DIR="/backups/postgres"

# 1. Dump database
pg_dump -h $DB_HOST -U $DB_USER -d $DB_NAME \
  --format=custom \
  --compress=9 \
  | openssl enc -aes-256-cbc -salt \
    -pass file:/etc/backup/encryption.key \
    > "${BACKUP_DIR}/${BACKUP_FILE}"

# 2. Upload ke backup storage (S3/MinIO)
aws s3 cp "${BACKUP_DIR}/${BACKUP_FILE}" \
  s3://akubelajar-backups/postgres/

# 3. Cleanup lokal (keep 7 hari)
find ${BACKUP_DIR} -name "*.enc" -mtime +7 -delete

# 4. Verifikasi
echo "Backup completed: ${BACKUP_FILE} ($(du -h ${BACKUP_DIR}/${BACKUP_FILE} | cut -f1))"
```

### Cron Schedule

```crontab
# PostgreSQL full backup — setiap hari jam 2 pagi
0 2 * * * /opt/akubelajar/scripts/backup-db.sh >> /var/log/backup.log 2>&1

# Redis snapshot — setiap 6 jam
0 */6 * * * /opt/akubelajar/scripts/backup-redis.sh >> /var/log/backup.log 2>&1

# MinIO files — setiap hari jam 3 pagi
0 3 * * * /opt/akubelajar/scripts/backup-files.sh >> /var/log/backup.log 2>&1
```

---

## Prosedur Restore

### PostgreSQL Restore

```bash
#!/bin/bash
# scripts/restore-db.sh

BACKUP_FILE=$1  # Nama file backup

# 1. Decrypt
openssl enc -d -aes-256-cbc \
  -pass file:/etc/backup/encryption.key \
  -in "/backups/postgres/${BACKUP_FILE}" \
  | pg_restore -h $DB_HOST -U $DB_USER -d $DB_NAME \
    --clean --if-exists

# 2. Verifikasi
psql -h $DB_HOST -U $DB_USER -d $DB_NAME \
  -c "SELECT COUNT(*) FROM users;"

echo "Restore completed from: ${BACKUP_FILE}"
```

### Point-in-Time Recovery (PITR)

```bash
# Restore ke waktu spesifik menggunakan WAL
recovery_target_time = '2026-03-21 14:30:00+07'
```

---

## Restore Drill Schedule

| Drill | Frekuensi | Penanggung Jawab | Hasil Terakhir |
|:---|:---|:---|:---|
| Full DB restore | Bulanan | DevOps Lead | ✅ 15 Mar 2026 — 12 menit |
| PITR test | Kuartalan | DBA | ✅ 1 Mar 2026 — 8 menit |
| Full disaster recovery | Semesteran | Tim Infra | 📅 Dijadwalkan Jun 2026 |

### Restore Drill Checklist

```
□ Download backup terbaru dari storage
□ Decrypt backup
□ Restore ke database test (BUKAN production!)
□ Verifikasi row count pada tabel kritis
□ Test login dengan user dari backup
□ Catat waktu total (target: < 30 menit)
□ Dokumentasikan hasil
□ Laporkan ke tech lead
```

---

## RPO & RTO

| Metrik | Target | Keterangan |
|:---|:---|:---|
| **RPO** (Recovery Point Objective) | ≤ 1 jam | Data loss maksimal 1 jam (via WAL archiving) |
| **RTO** (Recovery Time Objective) | ≤ 30 menit | Waktu restore sampai service kembali online |

---

## Monitoring Backup

- **Alert: Backup gagal** → Notifikasi Slack + email ke DevOps
- **Alert: Backup size anomaly** → Size berubah > 50% dari rata-rata
- **Dashboard** → Grafana panel: backup success rate, size trend, duration

---

*Terakhir diperbarui: 21 Maret 2026*
