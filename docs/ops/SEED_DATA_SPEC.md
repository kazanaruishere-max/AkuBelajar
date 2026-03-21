# 🌱 Seed Data Specification — AkuBelajar

> Data yang harus ada saat sistem pertama kali dijalankan — production dan development.

---

## 1. Production Seed (Wajib)

Data minimal yang HARUS ada di production:

### SuperAdmin Default

```sql
INSERT INTO schools (id, name, code, is_active)
VALUES (gen_random_uuid(), '[NAMA SEKOLAH]', '[CODE]', true)
ON CONFLICT DO NOTHING;

INSERT INTO users (id, school_id, email, password, role, is_first_login)
VALUES (
    gen_random_uuid(),
    (SELECT id FROM schools LIMIT 1),
    'admin@[domain]',
    '[ARGON2ID_HASH_OF_TEMP_PASSWORD]',
    'super_admin',
    true    -- Force password change on first login
)
ON CONFLICT (email) DO NOTHING;
```

- Credentials dikirim ke owner sekolah via channel aman
- `is_first_login: true` → WAJIB ganti password
- **Tidak ada data dummy** di production seed

---

## 2. Development Seed (Realistis)

### Sekolah

```sql
INSERT INTO schools (id, name, code, config, is_active) VALUES
    ('019516a2-0001-7000-8000-000000000001', 'SMP AkuBelajar', 'AB0001',
     '{"grading_weights":{"assignment":60,"quiz":40},"kkm_default":70}', true)
ON CONFLICT DO NOTHING;
```

### Tahun Ajaran

```sql
INSERT INTO academic_years (id, school_id, name, start_date, end_date, is_active) VALUES
    ('019516a2-0002-7000-8000-000000000001',
     '019516a2-0001-7000-8000-000000000001',
     '2025/2026', '2025-07-14', '2026-06-30', true)
ON CONFLICT DO NOTHING;
```

### Users (password default: `Test@12345`)

```sql
-- SuperAdmin
INSERT INTO users (id, school_id, email, password, role, is_first_login) VALUES
    ('[uuid]', '[school_id]', 'admin@akubelajar.test', '[hash]', 'super_admin', false);

-- Guru (3)
INSERT INTO users VALUES
    ('[uuid]', '[school_id]', 'guru1@akubelajar.test', '[hash]', 'teacher', false),
    ('[uuid]', '[school_id]', 'guru2@akubelajar.test', '[hash]', 'teacher', false),
    ('[uuid]', '[school_id]', 'guru3@akubelajar.test', '[hash]', 'teacher', false);

-- Ketua Kelas (1)
INSERT INTO users VALUES
    ('[uuid]', '[school_id]', 'ketua@akubelajar.test', '[hash]', 'class_leader', false);

-- Siswa (10)
INSERT INTO users VALUES
    ('[uuid]', '[school_id]', 'siswa01@akubelajar.test', '[hash]', 'student', false),
    ('[uuid]', '[school_id]', 'siswa02@akubelajar.test', '[hash]', 'student', false),
    -- ... siswa03 sampai siswa10
    ('[uuid]', '[school_id]', 'siswa10@akubelajar.test', '[hash]', 'student', false);
```

### Kelas

```sql
INSERT INTO classes VALUES
    ('[uuid]', '[school_id]', '[ay_id]', '8A', 8, '[guru1_id]'),  -- 10 siswa + 1 ketua
    ('[uuid]', '[school_id]', '[ay_id]', '8B', 8, '[guru2_id]');  -- kosong
```

### Mata Pelajaran (5)

```sql
INSERT INTO subjects VALUES
    ('[uuid]', '[school_id]', 'Matematika', 'MTK', 'Matematika kelas 8'),
    ('[uuid]', '[school_id]', 'Bahasa Indonesia', 'BIN', 'Bahasa Indonesia'),
    ('[uuid]', '[school_id]', 'IPA', 'IPA', 'Ilmu Pengetahuan Alam'),
    ('[uuid]', '[school_id]', 'IPS', 'IPS', 'Ilmu Pengetahuan Sosial'),
    ('[uuid]', '[school_id]', 'Bahasa Inggris', 'BIG', 'Bahasa Inggris');
```

### Sample Content

**3 Tugas:**
1. Tugas aktif: "Laporan Praktikum IPA" — deadline 3 hari dari sekarang
2. Tugas deadline lewat: "Essay Bahasa Indonesia" — deadline kemarin, 5 submission
3. Tugas sudah dinilai: "Latihan Matematika" — semua siswa dinilai

**2 Kuis:**
1. Published: "Kuis Bab 1: Bilangan Bulat" — 10 soal, 30 menit
2. Draft: "Kuis Bab 2: Aljabar" — 15 soal

**Absensi:** Data 2 minggu terakhir untuk kelas 8A (10 record per hari × 10 hari)

---

## 3. Idempotent Seed

```go
// cmd/seed/main.go
func main() {
    reset := os.Getenv("RESET") == "true"
    
    if reset {
        log.Println("⚠️  Resetting all data...")
        db.Exec("TRUNCATE schools, users, ... CASCADE")
    }
    
    log.Println("🌱 Seeding...")
    seedSchools(db)
    seedUsers(db)
    seedAcademicData(db)
    seedSampleContent(db)
    log.Println("✅ Seed complete")
}
```

- Aman dijalankan berkali-kali: `INSERT ... ON CONFLICT DO NOTHING`
- Run: `make seed` atau `go run cmd/seed/main.go`
- Reset: `make seed RESET=true`

---

## 4. Factory Functions (untuk Test)

```go
// internal/testutil/factory.go

func CreateTestUser(t *testing.T, db *pgxpool.Pool, role string, overrides ...map[string]any) User {
    user := User{
        ID:       uuid.New(),
        SchoolID: TestSchoolID,
        Email:    fmt.Sprintf("test-%s@akubelajar.test", uuid.New().String()[:8]),
        Password: HashPassword("Test@12345"),
        Role:     role,
    }
    // Apply overrides
    for _, o := range overrides {
        if v, ok := o["email"]; ok { user.Email = v.(string) }
        if v, ok := o["school_id"]; ok { user.SchoolID = v.(uuid.UUID) }
    }
    db.QueryRow(ctx, "INSERT INTO users ... RETURNING *", ...).Scan(&user)
    t.Cleanup(func() { db.Exec(ctx, "DELETE FROM users WHERE id=$1", user.ID) })
    return user
}

func CreateTestAssignment(t *testing.T, db *pgxpool.Pool, classID, teacherID uuid.UUID, overrides ...map[string]any) Assignment {
    // Similar pattern...
}

func CreateTestQuiz(t *testing.T, db *pgxpool.Pool, classID, teacherID uuid.UUID, overrides ...map[string]any) Quiz {
    // Similar pattern...
}
```

### Contoh Penggunaan di Test

```go
func TestGradeAssignment(t *testing.T) {
    teacher := testutil.CreateTestUser(t, db, "teacher")
    student := testutil.CreateTestUser(t, db, "student")
    assignment := testutil.CreateTestAssignment(t, db, classID, teacher.ID)
    
    // Act: student submit + teacher grade
    sub := submitAssignment(t, student, assignment, []byte("file content"))
    gradeSubmission(t, teacher, sub.ID, 85, "Bagus!")
    
    // Assert
    assert.Equal(t, 85, sub.Grade)
    assert.Equal(t, "graded", sub.Status)
}
```

---

*Terakhir diperbarui: 21 Maret 2026*
