# 📱 Mobile Guide — AkuBelajar

> Panduan khusus lapisan mobile: PWA configuration, React Native (Expo) roadmap, dan native features.

---

## Strategi Mobile: Two-Tier

| Tier | Platform | Status | Prioritas |
|:---|:---|:---|:---|
| **Tier 1** | PWA (Progressive Web App) | ✅ Active | Utama |
| **Tier 2** | React Native (Expo) | 📅 Roadmap Q3 2026 | Pelengkap |

---

## Tier 1: PWA

### Kapabilitas

| Fitur | Status | Keterangan |
|:---|:---|:---|
| Install ke Home Screen | ✅ | Via browser "Add to Home Screen" |
| Offline Support | ✅ | Service Worker + Workbox |
| Push Notification | ✅ | Web Push API |
| Fullscreen Mode | ✅ | `display: standalone` di manifest |
| Splash Screen | ✅ | Custom splash dengan logo |
| Background Sync | ✅ | Sync data saat kembali online |
| Kamera | ⚠️ | Via `getUserMedia` (terbatas) |
| File System | ⚠️ | Via File System Access API |
| Biometric Auth | ❌ | Tidak tersedia di web |

### Web App Manifest

```json
{
  "name": "AkuBelajar — Platform Edukasi Digital",
  "short_name": "AkuBelajar",
  "description": "Platform manajemen edukasi AI-first",
  "start_url": "/dashboard",
  "display": "standalone",
  "background_color": "#0A0A0A",
  "theme_color": "#3B82F6",
  "orientation": "portrait-primary",
  "icons": [
    { "src": "/icons/icon-192.png", "sizes": "192x192", "type": "image/png" },
    { "src": "/icons/icon-512.png", "sizes": "512x512", "type": "image/png" },
    { "src": "/icons/icon-maskable.png", "sizes": "512x512", "type": "image/png", "purpose": "maskable" }
  ]
}
```

### Responsive Breakpoints

| Breakpoint | Target | Action |
|:---|:---|:---|
| `< 640px` | Mobile phone | Single column, bottom nav |
| `640-1024px` | Tablet | Two column, sidebar collapsible |
| `> 1024px` | Desktop | Full layout, sidebar persistent |

---

## Tier 2: React Native (Expo) — Roadmap

### Timeline

| Milestone | Target | Deliverable |
|:---|:---|:---|
| Setup & Auth | Q3 2026 | Login, register, profile |
| Core Features | Q4 2026 | Dashboard, kuis, absensi |
| Native Features | Q1 2027 | QR scan, biometric, camera |
| Store Release | Q2 2027 | Play Store + App Store |

### Tech Stack

| Komponen | Teknologi |
|:---|:---|
| Framework | React Native + Expo SDK 52+ |
| Navigation | Expo Router (file-based) |
| State | Zustand (shared dengan web) |
| UI Kit | React Native Paper / Tamagui |
| Push Notification | Expo Notifications |
| Storage | AsyncStorage + SQLite (offline) |

### Shared Code dengan Web

```
packages/
├── shared/
│   ├── types/        # TypeScript interfaces (shared)
│   ├── validation/   # Zod schemas (shared)
│   ├── utils/        # Date formatting, etc. (shared)
│   └── api/          # API client (shared)
├── web/              # Next.js app
└── mobile/           # Expo app
```

---

## Testing di Mobile

| Tipe | Tool | Coverage Target |
|:---|:---|:---|
| PWA Audit | Lighthouse | Score > 90 |
| Responsiveness | Chrome DevTools | Semua breakpoints |
| Offline | Chrome DevTools (Network: Offline) | Core flows berfungsi |
| Performance | WebPageTest (mobile profile) | LCP < 3s on 3G |

---

## Referensi

- [Offline Strategy](../fase-1/OFFLINE_STRATEGY.md)
- [Frontend Guide](FRONTEND_GUIDE.md)
- [Performance Budget](../fase-1/PERFORMANCE_BUDGET.md)

---

*Terakhir diperbarui: 21 Maret 2026*
