# 📴 Offline Strategy — AkuBelajar

> Strategi Service Worker dan cache untuk memastikan platform tetap fungsional di jaringan lemah atau offline — kritis untuk daerah dengan konektivitas terbatas di Indonesia.

---

## Prinsip Offline

1. **Graceful Degradation** — Fitur tetap berfungsi sebaik mungkin, bukan "halaman error"
2. **Cache-First untuk Shell** — App shell selalu tersedia dari cache
3. **Sync saat Online** — Data yang dibuat offline di-sync otomatis saat koneksi kembali
4. **Transparansi** — User selalu tahu status koneksi mereka (online/offline indicator)

---

## Strategi Caching per Tipe Konten

| Tipe Konten | Strategi | Penjelasan |
|:---|:---|:---|
| App Shell (HTML/JS/CSS) | **Cache-First** | Selalu dari cache, update di background |
| Jadwal & Pengumuman | **Stale-While-Revalidate** | Tampilkan cache, update diam-diam |
| Data real-time (nilai, notif) | **Network-First** | Prioritas jaringan, fallback ke cache |
| File materi (PDF) | **Cache-on-Demand** | Hanya simpan yang pernah diakses |
| API Responses (GET) | **Network-First + Cache** | Cache untuk 5-30 menit |
| API Mutations (POST/PUT) | **Queue Offline** | Simpan di IndexedDB, sync saat online |

---

## Implementasi Teknis

### Service Worker Registration

```javascript
// public/sw-register.js
if ('serviceWorker' in navigator) {
  window.addEventListener('load', () => {
    navigator.serviceWorker.register('/sw.js')
      .then(reg => console.log('SW registered:', reg.scope))
      .catch(err => console.error('SW registration failed:', err));
  });
}
```

### Workbox Configuration

```javascript
// next.config.js (via next-pwa)
const withPWA = require('next-pwa')({
  dest: 'public',
  disable: process.env.NODE_ENV === 'development',
  runtimeCaching: [
    {
      urlPattern: /^https:\/\/api\.akubelajar\.id\/v1\/.*/,
      handler: 'NetworkFirst',
      options: {
        cacheName: 'api-cache',
        expiration: { maxEntries: 200, maxAgeSeconds: 300 },
        networkTimeoutSeconds: 5,
      },
    },
    {
      urlPattern: /\.(png|jpg|jpeg|webp|svg|ico)$/,
      handler: 'CacheFirst',
      options: {
        cacheName: 'image-cache',
        expiration: { maxEntries: 100, maxAgeSeconds: 86400 },
      },
    },
  ],
});
```

### Background Sync (Offline Mutations)

```javascript
// Saat user submit absensi dalam keadaan offline
async function submitAttendance(data) {
  try {
    await fetch('/api/v1/attendances', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  } catch (error) {
    // Simpan ke IndexedDB untuk sync nanti
    await saveToOfflineQueue('attendance', data);
    showNotification('Data disimpan offline. Akan di-sync otomatis.');
  }
}
```

---

## Fitur yang Tersedia Offline

| Fitur | Offline | Keterangan |
|:---|:---|:---|
| Lihat jadwal | ✅ | Dari cache terakhir |
| Baca pengumuman | ✅ | Dari cache terakhir |
| Lihat materi (PDF) | ✅ | Jika pernah dibuka sebelumnya |
| Input absensi | ✅ | Queued, sync saat online |
| Kerjakan kuis | ❌ | Butuh server-side timer & validation |
| Lihat nilai | ⚠️ | Dari cache terakhir (mungkin tidak terbaru) |
| Login | ❌ | Butuh server verification |

---

## Offline UX Indicators

```
┌─────────────────────────────────────────────┐
│  ⚠️ Anda sedang offline                     │
│  Data yang ditampilkan mungkin tidak terbaru │
│  Perubahan akan di-sync otomatis saat online │
└─────────────────────────────────────────────┘
```

---

## Referensi

- [PWA Documentation — web.dev](https://web.dev/progressive-web-apps/)
- [Workbox — Google](https://developer.chrome.com/docs/workbox)
- [Frontend Guide](../fase-2/FRONTEND_GUIDE.md)

---

*Terakhir diperbarui: 21 Maret 2026*
