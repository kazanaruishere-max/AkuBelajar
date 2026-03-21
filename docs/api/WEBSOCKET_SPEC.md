# 🔌 WebSocket Specification — AkuBelajar

> Kontrak lengkap WebSocket untuk real-time communication: CBT heartbeat, notifikasi, dan live updates.

---

## 1. Endpoint & Autentikasi

| Item | Detail |
|:---|:---|
| URL | `wss://api.akubelajar.id/ws` |
| Auth | Query param `?token=<paseto_token>` |
| Max koneksi per user | 3 simultan |
| Idle timeout | 5 menit tanpa activity → disconnect |
| Heartbeat | Server kirim `ping` setiap 30 detik, client harus `pong` |

> **Mengapa query param?** WebSocket API tidak mendukung custom headers di browser. Token dikirim via query param, server memvalidasi dan menghapus dari log.

---

## 2. Message Envelope (Standar Semua Pesan)

```json
{
  "type": "event_name",
  "payload": {},
  "timestamp": "2026-03-21T10:00:00Z",
  "id": "019516a2-uuid-v7"
}
```

| Field | Type | Deskripsi |
|:---|:---|:---|
| `type` | string | Nama event (format: `domain:action`) |
| `payload` | object | Data event |
| `timestamp` | ISO 8601 | Waktu event dari server |
| `id` | UUID v7 | ID unik pesan (untuk dedup) |

---

## 3. Client → Server Events

### `quiz:heartbeat`
Dikirim setiap **5 detik** selama ujian CBT aktif.

```json
{ "type": "quiz:heartbeat", "payload": { "session_id": "uuid" } }
```
> Server tidak merespon. Jika tidak terima heartbeat > 30 detik → anggap disconnect.

### `quiz:autosave`
Dikirim setiap **30 detik** atau saat siswa pindah soal.

```json
{
  "type": "quiz:autosave",
  "payload": {
    "session_id": "uuid",
    "answers": [
      { "question_id": "uuid", "selected_key": "B" },
      { "question_id": "uuid", "selected_key": "A" }
    ]
  }
}
```
**Server response:** `quiz:autosave_ack`
```json
{ "type": "quiz:autosave_ack", "payload": { "saved_count": 15, "server_time": "..." } }
```

### `notification:mark_read`
```json
{ "type": "notification:mark_read", "payload": { "notification_id": "uuid" } }
```

---

## 4. Server → Client Events

### `quiz:time_update`
Dikirim setiap **60 detik** selama ujian.
```json
{ "type": "quiz:time_update", "payload": { "session_id": "uuid", "remaining_seconds": 1740 } }
```

### `quiz:force_submit`
Waktu habis — server auto-submit.
```json
{ "type": "quiz:force_submit", "payload": { "session_id": "uuid", "score": 85, "reason": "time_expired" } }
```

### `quiz:session_locked`
Anti-cheat triggered.
```json
{
  "type": "quiz:session_locked",
  "payload": {
    "session_id": "uuid",
    "reason": "tab_switch_exceeded",
    "cheat_count": 4,
    "message": "Sesi ujian dikunci karena terdeteksi berpindah tab terlalu sering."
  }
}
```

### `notification:new`
Notifikasi baru masuk.
```json
{
  "type": "notification:new",
  "payload": {
    "id": "uuid",
    "type": "assignment_new",
    "title": "Tugas Baru: Laporan Biologi",
    "body": "Deadline: 25 Maret 2026",
    "related_entity_type": "assignment",
    "related_entity_id": "uuid"
  }
}
```

### `attendance:updated`
Absensi berubah real-time (untuk dashboard guru).
```json
{
  "type": "attendance:updated",
  "payload": {
    "class_id": "uuid",
    "date": "2026-03-21",
    "summary": { "present": 28, "permission": 1, "sick": 1, "absent": 2 }
  }
}
```

### `grade:published`
Nilai baru keluar.
```json
{
  "type": "grade:published",
  "payload": {
    "subject_name": "Biologi",
    "score": 85,
    "grade_letter": "B"
  }
}
```

### `session:revoked`
Sesi di-revoke (ganti password, admin suspend, dll).
```json
{ "type": "session:revoked", "payload": { "reason": "password_changed", "message": "Silakan login ulang." } }
```

---

## 5. Reconnect Strategy (Client — TypeScript)

```typescript
class WebSocketManager {
  private ws: WebSocket | null = null;
  private retryCount = 0;
  private maxRetries = 10;
  private missedMessages: string[] = [];

  connect(token: string) {
    this.ws = new WebSocket(`wss://api.akubelajar.id/ws?token=${token}`);

    this.ws.onopen = () => {
      this.retryCount = 0; // Reset on successful connect
    };

    this.ws.onclose = (event) => {
      if (event.code === 4001) return; // Auth failed, don't retry
      this.scheduleReconnect(token);
    };

    this.ws.onmessage = (event) => {
      const msg = JSON.parse(event.data);
      this.handleMessage(msg);
    };
  }

  private scheduleReconnect(token: string) {
    if (this.retryCount >= this.maxRetries) {
      // Show error UI: "Koneksi terputus. Silakan refresh halaman."
      return;
    }
    const delay = Math.min(1000 * Math.pow(2, this.retryCount), 30000);
    // Backoff: 1s, 2s, 4s, 8s, 16s, 30s, 30s, 30s...
    setTimeout(() => {
      this.retryCount++;
      this.connect(token);
    }, delay);
  }
}
```

---

## 6. Error Handling

### Close Codes

| Code | Arti | Client Action |
|:---|:---|:---|
| `1000` | Normal close | — |
| `1001` | Server going away (restart) | Auto-reconnect |
| `1006` | Abnormal close (network) | Auto-reconnect |
| `4001` | Authentication failed | Redirect to login |
| `4002` | Token expired | Refresh token, reconnect |
| `4003` | Max connections exceeded | Show error |
| `4004` | Session revoked | Redirect to login |

### Token Expired Mid-Connection

```
1. Server detects token expired
2. Server sends: { "type": "auth:token_expired" }
3. Client calls POST /auth/refresh
4. Client reconnects with new token
```

---

*Terakhir diperbarui: 21 Maret 2026*
