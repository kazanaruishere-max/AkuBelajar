'use client';

import React, { useState, useEffect, useCallback } from 'react';
import { api } from '@/lib/api/client';

interface AdminStats { total_users: number; total_teachers: number; total_students: number; total_classes: number; total_subjects: number; }
interface TeacherStats { total_assignments: number; total_quizzes: number; pending_submissions: number; }
interface StudentStats { total_assignments: number; total_quizzes: number; average_score: number; unread_notifications: number; }

export default function DashboardPage() {
  const [stats, setStats] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [role, setRole] = useState('');

  const fetchStats = useCallback(async () => {
    try {
      setLoading(true);
      const res = await api.get('/api/v1/dashboard/stats');
      const data = (res.data as any).data;
      setStats(data);
      setRole(data?.role || '');
    } catch { /* ignore */ } finally { setLoading(false); }
  }, []);

  useEffect(() => { fetchStats(); }, [fetchStats]);

  if (loading) return <div style={{ padding: 32, color: '#94A3B8' }}>Memuat dashboard...</div>;

  return (
    <div style={{ padding: 32 }}>
      <h1 style={{ fontSize: 28, fontWeight: 800, marginBottom: 8, background: 'linear-gradient(135deg,#E2E8F0,#94A3B8)', WebkitBackgroundClip: 'text', WebkitTextFillColor: 'transparent' }}>
        Dashboard
      </h1>
      <p style={{ fontSize: 14, color: '#64748B', marginBottom: 32 }}>
        Selamat datang kembali! 👋
      </p>

      {role === 'super_admin' && <AdminDashboard stats={stats as AdminStats} />}
      {role === 'teacher' && <TeacherDashboard stats={stats as TeacherStats} />}
      {(role === 'student' || role === 'class_leader') && <StudentDashboard stats={stats as StudentStats} />}
    </div>
  );
}

function AdminDashboard({ stats }: { stats: AdminStats }) {
  const cards = [
    { label: 'Total Pengguna', value: stats.total_users, icon: '👥', color: '#3B82F6' },
    { label: 'Guru', value: stats.total_teachers, icon: '👨‍🏫', color: '#22C55E' },
    { label: 'Siswa', value: stats.total_students, icon: '👩‍🎓', color: '#A78BFA' },
    { label: 'Kelas', value: stats.total_classes, icon: '🏫', color: '#F59E0B' },
    { label: 'Mata Pelajaran', value: stats.total_subjects, icon: '📚', color: '#EC4899' },
  ];
  return <StatsGrid cards={cards} />;
}

function TeacherDashboard({ stats }: { stats: TeacherStats }) {
  const cards = [
    { label: 'Tugas Dibuat', value: stats.total_assignments, icon: '📝', color: '#3B82F6' },
    { label: 'Kuis Dibuat', value: stats.total_quizzes, icon: '❓', color: '#A78BFA' },
    { label: 'Submission Menunggu', value: stats.pending_submissions, icon: '⏳', color: '#F59E0B' },
  ];
  return <StatsGrid cards={cards} />;
}

function StudentDashboard({ stats }: { stats: StudentStats }) {
  const scoreColor = stats.average_score >= 80 ? '#22C55E' : stats.average_score >= 60 ? '#F59E0B' : '#EF4444';
  const cards = [
    { label: 'Tugas Dikerjakan', value: stats.total_assignments, icon: '📝', color: '#3B82F6' },
    { label: 'Kuis Dikerjakan', value: stats.total_quizzes, icon: '❓', color: '#A78BFA' },
    { label: 'Rata-rata Nilai', value: stats.average_score, icon: '📊', color: scoreColor },
    { label: 'Notifikasi Baru', value: stats.unread_notifications, icon: '🔔', color: '#EF4444' },
  ];
  return <StatsGrid cards={cards} />;
}

function StatsGrid({ cards }: { cards: { label: string; value: number; icon: string; color: string }[] }) {
  return (
    <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(220px, 1fr))', gap: 20 }}>
      {cards.map((c) => (
        <div key={c.label} style={{
          background: 'rgba(15,23,42,0.6)', borderRadius: 20, padding: 28,
          border: '1px solid rgba(148,163,184,0.08)', position: 'relative', overflow: 'hidden',
        }}>
          <div style={{ position: 'absolute', top: -20, right: -20, width: 80, height: 80, borderRadius: '50%', background: `${c.color}10` }} />
          <div style={{ fontSize: 28, marginBottom: 12 }}>{c.icon}</div>
          <p style={{ fontSize: 42, fontWeight: 800, color: c.color, lineHeight: 1 }}>{c.value}</p>
          <p style={{ fontSize: 13, color: '#94A3B8', marginTop: 8, fontWeight: 500 }}>{c.label}</p>
        </div>
      ))}
    </div>
  );
}
