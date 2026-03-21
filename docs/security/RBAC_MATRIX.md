# 🔐 RBAC Permission Matrix — AkuBelajar

> Matriks lengkap: Role × Resource × Action. Sumber kebenaran untuk middleware dan RLS.

---

## 1. Role Hierarchy

```
super_admin (level 4) — akses penuh per sekolah
  ├── teacher (level 3) — akses kelas & mapel sendiri
  │     └── class_leader (level 2) — siswa + input absensi draft
  └── student (level 1) — akses data sendiri saja
```

---

## 2. Permission Matrix

✅ = diizinkan | ❌ = ditolak | 🔒 = hanya data sendiri/kelasnya

### User Management

| Action | SuperAdmin | Teacher | ClassLeader | Student |
|:---|:---|:---|:---|:---|
| List all users | ✅ | ❌ | ❌ | ❌ |
| View user detail | ✅ | 🔒 siswa di kelasnya | ❌ | ❌ |
| Create user (single) | ✅ | ❌ | ❌ | ❌ |
| Bulk import users | ✅ | ❌ | ❌ | ❌ |
| Update any user | ✅ | ❌ | ❌ | ❌ |
| Update own profile | ✅ | ✅ | ✅ | ✅ |
| Change user role | ✅ | 🔒 siswa→KK di kelasnya | ❌ | ❌ |
| Suspend user | ✅ | ❌ | ❌ | ❌ |
| Delete user | ✅ | ❌ | ❌ | ❌ |
| Reset password | ✅ | 🔒 siswa di kelasnya | ❌ | ❌ |

### Academic Data

| Action | SuperAdmin | Teacher | ClassLeader | Student |
|:---|:---|:---|:---|:---|
| Manage academic years | ✅ | ❌ | ❌ | ❌ |
| Manage classes | ✅ | ❌ | ❌ | ❌ |
| Manage subjects | ✅ | ❌ | ❌ | ❌ |
| Assign teacher to class | ✅ | ❌ | ❌ | ❌ |
| Assign student to class | ✅ | ❌ | ❌ | ❌ |
| View class list | ✅ | 🔒 kelasnya | 🔒 kelasnya | 🔒 kelasnya |
| View class_subjects | ✅ | 🔒 kelasnya | 🔒 kelasnya | 🔒 kelasnya |

### Attendance

| Action | SuperAdmin | Teacher | ClassLeader | Student |
|:---|:---|:---|:---|:---|
| Input attendance (FINAL) | ✅ | 🔒 mapelnya | ❌ | ❌ |
| Input attendance (DRAFT) | ❌ | ❌ | 🔒 kelasnya | ❌ |
| Approve draft attendance | ❌ | 🔒 kelasnya | ❌ | ❌ |
| Edit attendance < 24h | ✅ | 🔒 miliknya | ❌ | ❌ |
| Edit attendance > 24h | ✅ | 🔒 miliknya + reason | ❌ | ❌ |
| View attendance (all) | ✅ | 🔒 kelasnya | 🔒 kelasnya | ❌ |
| View own attendance | — | — | ✅ | ✅ |
| Submit leave request | ❌ | ❌ | ❌ | ✅ |

### Assignments

| Action | SuperAdmin | Teacher | ClassLeader | Student |
|:---|:---|:---|:---|:---|
| Create assignment | ❌ | 🔒 mapelnya | ❌ | ❌ |
| Edit assignment | ❌ | 🔒 miliknya | ❌ | ❌ |
| Delete assignment | ✅ | 🔒 miliknya | ❌ | ❌ |
| View assignment list | ✅ | 🔒 kelasnya | 🔒 kelasnya | 🔒 kelasnya |
| Submit assignment | ❌ | ❌ | ✅ | ✅ |
| View submissions | ✅ | 🔒 kelasnya | ❌ | 🔒 miliknya |
| Grade submission | ❌ | 🔒 miliknya | ❌ | ❌ |

### Quizzes / CBT

| Action | SuperAdmin | Teacher | ClassLeader | Student |
|:---|:---|:---|:---|:---|
| Create quiz | ❌ | ✅ | ❌ | ❌ |
| Generate AI questions | ❌ | ✅ | ❌ | ❌ |
| Publish quiz | ❌ | 🔒 miliknya | ❌ | ❌ |
| Take quiz | ❌ | ❌ | ✅ | ✅ |
| View quiz results | ✅ | 🔒 miliknya | ❌ | 🔒 miliknya |
| Reset locked session | ❌ | 🔒 miliknya | ❌ | ❌ |

### Grades & Report Cards

| Action | SuperAdmin | Teacher | ClassLeader | Student |
|:---|:---|:---|:---|:---|
| View all grades | ✅ | 🔒 mapelnya | ❌ | ❌ |
| View own grades | — | — | ✅ | ✅ |
| Lock/unlock grades | ✅ | ❌ | ❌ | ❌ |
| Generate report card PDF | ✅ | ❌ | ❌ | ❌ |
| View own report card | — | — | ✅ | ✅ |

### Notifications

| Action | SuperAdmin | Teacher | ClassLeader | Student |
|:---|:---|:---|:---|:---|
| View own notifications | ✅ | ✅ | ✅ | ✅ |
| Mark read | ✅ | ✅ | ✅ | ✅ |
| Update preferences | ✅ | ✅ | ✅ | ✅ |
| Send broadcast | ✅ | ❌ | ❌ | ❌ |

### System

| Action | SuperAdmin | Teacher | ClassLeader | Student |
|:---|:---|:---|:---|:---|
| View audit logs | ✅ | ❌ | ❌ | ❌ |
| Manage invite tokens | ✅ | ✅ | ❌ | ❌ |
| System settings | ✅ | ❌ | ❌ | ❌ |

---

## 3. RLS Policy per Tabel

```sql
-- Users: user hanya lihat sesama sekolah
CREATE POLICY users_school_isolation ON users
    USING (school_id = current_setting('app.school_id')::UUID);

-- Attendances: guru hanya lihat kelas yang dia ajar
CREATE POLICY att_teacher ON attendances FOR ALL TO teacher_role
    USING (class_id IN (
        SELECT class_id FROM class_subjects 
        WHERE teacher_id = current_setting('app.user_id')::UUID
    ));

-- Attendances: siswa hanya lihat milik sendiri
CREATE POLICY att_student ON attendances FOR SELECT TO student_role
    USING (student_id = current_setting('app.user_id')::UUID);

-- Grades: siswa hanya lihat milik sendiri
CREATE POLICY grades_student ON grades FOR SELECT TO student_role
    USING (student_id = current_setting('app.user_id')::UUID);

-- Notifications: user hanya lihat milik sendiri
CREATE POLICY notif_own ON notifications FOR ALL
    USING (user_id = current_setting('app.user_id')::UUID);
```

---

## 4. Go Middleware Implementation

```go
// Endpoint → minimum role mapping
var endpointRoles = map[string]map[string][]string{
    "GET": {
        "/api/v1/users":                {"super_admin"},
        "/api/v1/users/me":             {"super_admin", "teacher", "class_leader", "student"},
        "/api/v1/attendances":          {"super_admin", "teacher"},
        "/api/v1/attendances/my":       {"student", "class_leader"},
        "/api/v1/assignments":          {"super_admin", "teacher", "class_leader", "student"},
        "/api/v1/grades":               {"super_admin", "teacher"},
        "/api/v1/grades/my":            {"student", "class_leader"},
    },
    "POST": {
        "/api/v1/users":                {"super_admin"},
        "/api/v1/attendances":          {"teacher", "class_leader"},
        "/api/v1/assignments":          {"teacher"},
        "/api/v1/quizzes":              {"teacher"},
    },
}

func RBACMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        role := c.GetString("role")
        allowedRoles := endpointRoles[c.Request.Method][c.FullPath()]
        
        if !slices.Contains(allowedRoles, role) {
            c.JSON(403, AppError{Code: "AUTH_005", Message: "Tidak memiliki izin"})
            c.Abort()
            return
        }
        c.Next()
    }
}
```

---

*Terakhir diperbarui: 21 Maret 2026*
