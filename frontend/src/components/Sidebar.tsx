'use client';

import React from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';

interface SidebarProps {
  role: string;
}

const menuItems: Record<string, { label: string; href: string; icon: string }[]> = {
  super_admin: [
    { label: 'Dashboard', href: '/dashboard', icon: '📊' },
    { label: 'Pengguna', href: '/admin/users', icon: '👥' },
    { label: 'Tahun Ajaran', href: '/admin/academic-years', icon: '📅' },
    { label: 'Kelas', href: '/admin/classes', icon: '🏫' },
    { label: 'Mata Pelajaran', href: '/admin/subjects', icon: '📚' },
    { label: 'Notifikasi', href: '/notifications', icon: '🔔' },
  ],
  teacher: [
    { label: 'Dashboard', href: '/dashboard', icon: '📊' },
    { label: 'Tugas', href: '/teacher/assignments', icon: '📝' },
    { label: 'Kuis', href: '/teacher/quizzes', icon: '❓' },
    { label: 'Presensi', href: '/teacher/attendance', icon: '📋' },
    { label: 'Notifikasi', href: '/notifications', icon: '🔔' },
  ],
  student: [
    { label: 'Dashboard', href: '/dashboard', icon: '📊' },
    { label: 'Tugas', href: '/student/assignments', icon: '📝' },
    { label: 'Kuis', href: '/student/quizzes', icon: '❓' },
    { label: 'Nilai', href: '/student/grades', icon: '📈' },
    { label: 'Notifikasi', href: '/notifications', icon: '🔔' },
  ],
};

export default function Sidebar({ role }: SidebarProps) {
  const pathname = usePathname();
  const items = menuItems[role] || menuItems.student;

  return (
    <aside style={{
      width: 260, minHeight: '100vh', background: 'linear-gradient(180deg, #0F172A 0%, #1E293B 100%)',
      borderRight: '1px solid rgba(148,163,184,0.08)', padding: '24px 16px', display: 'flex', flexDirection: 'column',
    }}>
      {/* Logo */}
      <div style={{ padding: '8px 12px', marginBottom: 32 }}>
        <h1 style={{ fontSize: 22, fontWeight: 800, background: 'linear-gradient(135deg,#3B82F6,#8B5CF6)', WebkitBackgroundClip: 'text', WebkitTextFillColor: 'transparent' }}>
          AkuBelajar
        </h1>
        <p style={{ fontSize: 11, color: '#64748B', marginTop: 4, textTransform: 'uppercase', letterSpacing: '0.1em' }}>
          {role === 'super_admin' ? 'Admin Panel' : role === 'teacher' ? 'Teacher Panel' : 'Student Portal'}
        </p>
      </div>

      {/* Menu */}
      <nav style={{ display: 'flex', flexDirection: 'column', gap: 4, flex: 1 }}>
        {items.map((item) => {
          const isActive = pathname === item.href;
          return (
            <Link key={item.href} href={item.href} style={{
              display: 'flex', alignItems: 'center', gap: 12, padding: '12px 16px', borderRadius: 12,
              textDecoration: 'none', fontSize: 14, fontWeight: isActive ? 600 : 400, transition: 'all 0.2s',
              background: isActive ? 'rgba(59,130,246,0.1)' : 'transparent',
              color: isActive ? '#60A5FA' : '#94A3B8',
              borderLeft: isActive ? '3px solid #3B82F6' : '3px solid transparent',
            }}>
              <span style={{ fontSize: 18 }}>{item.icon}</span>
              <span>{item.label}</span>
            </Link>
          );
        })}
      </nav>

      {/* Footer */}
      <div style={{ padding: '16px 12px', borderTop: '1px solid rgba(148,163,184,0.08)', marginTop: 'auto' }}>
        <p style={{ fontSize: 11, color: '#475569' }}>© 2026 AkuBelajar</p>
      </div>
    </aside>
  );
}
