'use client';

import React, { useState, useEffect, useCallback } from 'react';
import { api } from '@/lib/api/client';

interface Notification {
  id: string;
  title: string;
  message: string;
  type: string;
  link: string;
  is_read: boolean;
  created_at: string;
}

export default function NotificationsPage() {
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [loading, setLoading] = useState(true);
  const [unreadCount, setUnreadCount] = useState(0);

  const fetchNotifications = useCallback(async () => {
    try {
      setLoading(true);
      const [nRes, cRes] = await Promise.all([
        api.get('/api/v1/notifications?limit=50'),
        api.get('/api/v1/notifications/unread-count'),
      ]);
      setNotifications((nRes.data as any).data || []);
      setUnreadCount((cRes.data as any).data?.count || 0);
    } catch { /* ignore */ } finally { setLoading(false); }
  }, []);

  useEffect(() => { fetchNotifications(); }, [fetchNotifications]);

  const markRead = async (id: string) => {
    await api.post(`/api/v1/notifications/${id}/read`);
    fetchNotifications();
  };

  const markAllRead = async () => {
    await api.post('/api/v1/notifications/read-all');
    fetchNotifications();
  };

  const typeConfig: Record<string, { icon: string; color: string }> = {
    info:       { icon: 'ℹ️', color: '#60A5FA' },
    assignment: { icon: '📝', color: '#22C55E' },
    quiz:       { icon: '❓', color: '#A78BFA' },
    grade:      { icon: '📊', color: '#FBBF24' },
    attendance: { icon: '📋', color: '#F97316' },
    system:     { icon: '⚙️', color: '#94A3B8' },
  };

  const timeAgo = (d: string) => {
    const s = Math.floor((Date.now() - new Date(d).getTime()) / 1000);
    if (s < 60) return 'baru saja';
    if (s < 3600) return `${Math.floor(s / 60)} menit lalu`;
    if (s < 86400) return `${Math.floor(s / 3600)} jam lalu`;
    return `${Math.floor(s / 86400)} hari lalu`;
  };

  return (
    <div style={{ padding: 32 }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
          <h1 style={{ fontSize: 24, fontWeight: 700, color: '#E2E8F0' }}>Notifikasi</h1>
          {unreadCount > 0 && (
            <span style={{ padding: '4px 12px', borderRadius: 20, fontSize: 13, fontWeight: 700,
              background: 'rgba(239,68,68,0.2)', color: '#EF4444' }}>{unreadCount} belum dibaca</span>
          )}
        </div>
        {unreadCount > 0 && (
          <button onClick={markAllRead} style={btnGhost}>✓ Tandai semua dibaca</button>
        )}
      </div>

      {loading ? <p style={{ color: '#94A3B8' }}>Memuat...</p> : notifications.length === 0 ? (
        <div style={emptyState}>
          <div style={{ fontSize: 48 }}>🔔</div>
          <p style={{ color: '#94A3B8' }}>Belum ada notifikasi</p>
        </div>
      ) : (
        <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
          {notifications.map((n) => {
            const tc = typeConfig[n.type] || typeConfig.info;
            return (
              <div key={n.id} onClick={() => !n.is_read && markRead(n.id)}
                style={{
                  padding: '16px 20px', borderRadius: 12, display: 'flex', alignItems: 'flex-start', gap: 14, cursor: 'pointer',
                  background: n.is_read ? 'rgba(15,23,42,0.3)' : 'rgba(59,130,246,0.06)',
                  border: n.is_read ? '1px solid rgba(148,163,184,0.08)' : '1px solid rgba(59,130,246,0.15)',
                  transition: 'all 0.2s',
                }}>
                <div style={{ width: 40, height: 40, borderRadius: 10, display: 'flex', alignItems: 'center', justifyContent: 'center',
                  background: `${tc.color}15`, fontSize: 18, flexShrink: 0 }}>{tc.icon}</div>
                <div style={{ flex: 1, minWidth: 0 }}>
                  <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                    <h3 style={{ fontSize: 14, fontWeight: n.is_read ? 500 : 700, color: n.is_read ? '#94A3B8' : '#E2E8F0' }}>{n.title}</h3>
                    <span style={{ fontSize: 12, color: '#64748B', flexShrink: 0 }}>{timeAgo(n.created_at)}</span>
                  </div>
                  <p style={{ fontSize: 13, color: '#94A3B8', marginTop: 4, lineHeight: 1.5 }}>{n.message}</p>
                  {n.link && <a href={n.link} style={{ fontSize: 13, color: '#3B82F6', marginTop: 6, display: 'inline-block' }}>Lihat detail →</a>}
                </div>
                {!n.is_read && <div style={{ width: 8, height: 8, borderRadius: '50%', background: '#3B82F6', marginTop: 6, flexShrink: 0 }} />}
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}

const btnGhost: React.CSSProperties = { padding: '8px 16px', borderRadius: 10, border: '1px solid rgba(148,163,184,0.2)', cursor: 'pointer', background: 'transparent', color: '#94A3B8', fontWeight: 500, fontSize: 13 };
const emptyState: React.CSSProperties = { display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', gap: 16, padding: 64, background: 'rgba(15,23,42,0.4)', borderRadius: 16, border: '1px dashed rgba(148,163,184,0.2)' };
