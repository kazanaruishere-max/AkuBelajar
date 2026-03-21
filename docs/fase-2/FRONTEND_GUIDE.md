# 🎨 Frontend Guide — AkuBelajar

> Panduan khusus lapisan frontend Next.js 15: rendering strategy, state management, dan component patterns.

---

## Rendering Strategy

| Halaman | Strategi | Alasan |
|:---|:---|:---|
| Landing page | **SSG** (Static) | SEO, performa maksimal |
| Dashboard | **SSR** (Server) | Data real-time, role-based |
| Jadwal (semi-static) | **ISR** (revalidate: 300s) | Jarang berubah, cache efisien |
| CBT Interface | **CSR** (Client) | Interaktivitas penuh, timer |
| Settings | **Server Component** | Fetch sekali, render di server |

---

## Component Patterns

### Struktur File per Fitur

```
components/features/quiz/
├── QuizCard.tsx          # Presentational component
├── QuizList.tsx          # Container component
├── QuizForm.tsx          # Form component
├── useQuiz.ts            # Custom hook (data fetching)
├── quiz.types.ts         # TypeScript interfaces
└── quiz.test.tsx         # Unit tests
```

### Server vs Client Component

```tsx
// ✅ Server Component (default di App Router)
// Tidak perlu "use client" — fetch data langsung di server
async function QuizList({ classId }: { classId: string }) {
  const quizzes = await fetchQuizzes(classId);
  return <div>{quizzes.map(q => <QuizCard key={q.id} quiz={q} />)}</div>;
}

// ✅ Client Component — hanya untuk interaktivitas
"use client";
function QuizTimer({ endTime }: { endTime: Date }) {
  const [remaining, setRemaining] = useState(calculateRemaining(endTime));
  useEffect(() => { /* timer logic */ }, []);
  return <span>{formatTime(remaining)}</span>;
}
```

---

## State Management

| Scope | Tool | Use Case |
|:---|:---|:---|
| Server state | **TanStack Query** | API data, caching, sync |
| Client state (global) | **Zustand** | Theme, sidebar, user prefs |
| Form state | **React Hook Form + Zod** | Input validation |
| URL state | **Next.js searchParams** | Filters, pagination |

```typescript
// Zustand store example
import { create } from 'zustand';

interface AuthStore {
  user: User | null;
  setUser: (user: User) => void;
  logout: () => void;
}

export const useAuthStore = create<AuthStore>((set) => ({
  user: null,
  setUser: (user) => set({ user }),
  logout: () => set({ user: null }),
}));
```

---

## API Client Pattern

```typescript
// lib/api/client.ts
const API_BASE = process.env.NEXT_PUBLIC_API_URL;

async function request<T>(endpoint: string, options?: RequestInit): Promise<T> {
  const response = await fetch(`${API_BASE}${endpoint}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${getToken()}`,
      ...options?.headers,
    },
  });

  if (!response.ok) {
    const error = await response.json();
    throw new APIError(response.status, error.message);
  }

  return response.json();
}

// Type-safe API functions
export const api = {
  quizzes: {
    list: (classId: string) => request<Quiz[]>(`/v1/quizzes?class_id=${classId}`),
    get: (id: string) => request<Quiz>(`/v1/quizzes/${id}`),
    submit: (id: string, answers: Answer[]) =>
      request<QuizResult>(`/v1/quizzes/${id}/submit`, {
        method: 'POST',
        body: JSON.stringify({ answers }),
      }),
  },
};
```

---

## Performance Checklist

- [ ] Dynamic import untuk komponen berat (chart, editor)
- [ ] `next/image` untuk semua gambar
- [ ] `next/font` untuk font loading
- [ ] Tidak ada `useEffect` fetch di client — gunakan Server Component
- [ ] Bundle size per route < 50 KB gzipped

---

## Referensi

- [Coding Standards](CODING_STANDARDS.md)
- [Offline Strategy](../fase-1/OFFLINE_STRATEGY.md)
- [Performance Budget](../fase-1/PERFORMANCE_BUDGET.md)

---

*Terakhir diperbarui: 21 Maret 2026*
