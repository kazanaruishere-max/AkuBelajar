# 🤖 AI Quiz Generation Flow — AkuBelajar

> Flowchart lengkap: dari guru klik "Generate AI" sampai soal tersimpan di database. Termasuk error handling Gemini.

---

## 1. Flow Lengkap

```mermaid
flowchart TD
    A["Guru: klik 'Generate Soal AI'"] --> B["Frontend: tampilkan form"]
    B --> C["Guru isi:\n- Mapel\n- Topik/Bab\n- Jumlah soal\n- Difficulty"]
    C --> D["Frontend: POST /api/v1/quizzes/:id/generate-ai"]
    D --> E["Backend: validate input"]
    E --> F["Backend: check rate limit\n(10 req/guru/jam)"]
    F --> G{Rate limit OK?}
    G -->|Tidak| H["Return 429\nQUIZ_008: Rate limit exceeded"]
    G -->|Ya| I["Service: build prompt\n(system + user template)"]
    I --> J["Service: call Gemini API\n(gemini-2.0-flash)"]
    J --> K{Gemini response OK?}
    K -->|Error 429| L["Backoff 60s → retry"]
    K -->|Error 500| M["Backoff exponential\n(1s, 5s, 15s)"]
    K -->|Timeout| N["Retry (max 3x)"]
    K -->|Success| O["Parser: extract JSON\nfrom response"]
    L --> J
    M --> J
    N --> J
    O --> P{JSON valid?}
    P -->|Tidak| Q["Retry dengan prompt\n'HANYA JSON, tanpa penjelasan'"]
    P -->|Ya| R["Validator: check per soal\n- 4 options?\n- correct_key valid?\n- text length OK?"]
    Q --> J
    R --> S{Semua soal valid?}
    S -->|Ada yang invalid| T["Filter: skip soal invalid\nlog warning"]
    S -->|Ya| U["Content Safety Filter\n(blocked keywords check)"]
    T --> U
    U --> V["INSERT quiz_questions\n(batch, dalam transaction)"]
    V --> W["Return 200\n{generated: N, filtered: M}"]
    W --> X["Frontend: tampilkan soal\nGuru bisa edit/delete"]
```

---

## 2. Sequence Diagram

```mermaid
sequenceDiagram
    participant Guru
    participant FE as Frontend
    participant API as Go Backend
    participant AI as Gemini API
    participant DB as PostgreSQL

    Guru->>FE: Klik "Generate Soal AI"
    FE->>FE: Tampilkan form (mapel, topik, jumlah, difficulty)
    Guru->>FE: Submit form
    FE->>API: POST /quizzes/:id/generate-ai
    API->>API: Validate + rate limit check
    
    loop Max 3 retry
        API->>AI: Generate content (prompt)
        alt Success
            AI-->>API: JSON response
        else Rate limited
            AI-->>API: 429
            API->>API: Wait 60s
        else Server error
            AI-->>API: 500
            API->>API: Exponential backoff
        end
    end
    
    API->>API: Parse JSON + Validate soal
    API->>API: Content safety filter
    API->>DB: INSERT quiz_questions (batch)
    DB-->>API: OK
    API-->>FE: 200 {generated: 10, filtered: 1}
    FE->>FE: Tampilkan soal (editable)
    FE-->>Guru: "11 soal dihasilkan, 1 difilter"
```

---

## 3. Error Handling Matrix

| Error | HTTP | Error Code | User Message | Aksi Backend |
|:---|:---|:---|:---|:---|
| Rate limit guru (10/jam) | 429 | `QUIZ_008` | "Anda sudah generate 10x dalam 1 jam. Tunggu X menit." | — |
| Gemini quota habis | 503 | `QUIZ_007` | "Layanan AI sedang sibuk. Coba 5 menit lagi." | Retry 60s |
| Gemini return error | 502 | `QUIZ_007` | "Gagal generate soal. Coba lagi." | Retry 3x |
| Response bukan JSON | 502 | `QUIZ_007` | "Gagal generate soal. Coba lagi." | Retry + modified prompt |
| Jumlah soal tidak sesuai | 200 | — | "Diminta 10, berhasil generate 8." | Return partial |
| Konten tidak pantas | 200 | — | "11 soal dihasilkan, 1 difilter karena tidak sesuai." | Filter + log |
| API key invalid | 500 | `SYS_004` | "Konfigurasi sistem error. Hubungi admin." | Alert admin |
| Semua retry gagal | 503 | `QUIZ_007` | "Tidak bisa generate soal saat ini." | Log + alert |

---

## 4. Rate Limit Detail

```go
// Redis key: ai_quiz_rate:{teacher_id}
// Limit: 10 requests per jam
// Window: sliding window 1 hour

func CheckAIRateLimit(ctx context.Context, teacherID string) error {
    key := fmt.Sprintf("ai_quiz_rate:%s", teacherID)
    count, _ := redis.Incr(ctx, key)
    if count == 1 {
        redis.Expire(ctx, key, 1*time.Hour)
    }
    if count > 10 {
        ttl, _ := redis.TTL(ctx, key)
        return fmt.Errorf("rate limited, retry in %v", ttl)
    }
    return nil
}
```

---

## 5. Post-Generation: Guru Review

Setelah AI generate, guru WAJIB review sebelum publish:

```
Generate → Review → Edit (opsional) → Save Draft → Publish
```

- Guru bisa **edit** teks soal, opsi, dan jawaban benar
- Guru bisa **delete** soal yang tidak sesuai
- Guru bisa **tambah** soal manual ke kuis yang sama
- Kuis tetap `status: draft` sampai guru klik Publish

---

*Terakhir diperbarui: 21 Maret 2026*
