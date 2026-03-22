'use client';

import { useEffect, useState } from 'react';
import { useRouter, usePathname } from 'next/navigation';
import { useAuthStore } from '@/lib/store/authStore';
import { useUIStore } from '@/lib/store/uiStore';

// Role-based navigation items
const navItems = {
  super_admin: [
    { label: 'Dashboard', href: '/dashboard', icon: 'grid' },
    { label: 'Tahun Ajaran', href: '/dashboard/admin/academic-years', icon: 'calendar' },
    { label: 'Kelas', href: '/dashboard/admin/classes', icon: 'users' },
    { label: 'Mata Pelajaran', href: '/dashboard/admin/subjects', icon: 'book' },
    { label: 'Pengguna', href: '/dashboard/admin/users', icon: 'user-plus' },
    { label: 'Audit Log', href: '/dashboard/admin/audit', icon: 'shield' },
  ],
  teacher: [
    { label: 'Dashboard', href: '/dashboard', icon: 'grid' },
    { label: 'Kelas Saya', href: '/dashboard/teacher/classes', icon: 'users' },
    { label: 'Tugas', href: '/dashboard/teacher/assignments', icon: 'file-text' },
    { label: 'Kuis', href: '/dashboard/teacher/quizzes', icon: 'help-circle' },
    { label: 'Absensi', href: '/dashboard/teacher/attendance', icon: 'check-square' },
    { label: 'Nilai', href: '/dashboard/teacher/grades', icon: 'bar-chart-2' },
  ],
  class_leader: [
    { label: 'Dashboard', href: '/dashboard', icon: 'grid' },
    { label: 'Absensi', href: '/dashboard/student/attendance', icon: 'check-square' },
    { label: 'Tugas', href: '/dashboard/student/assignments', icon: 'file-text' },
    { label: 'Kuis', href: '/dashboard/student/quizzes', icon: 'help-circle' },
    { label: 'Nilai', href: '/dashboard/student/grades', icon: 'bar-chart-2' },
  ],
  student: [
    { label: 'Dashboard', href: '/dashboard', icon: 'grid' },
    { label: 'Tugas', href: '/dashboard/student/assignments', icon: 'file-text' },
    { label: 'Kuis', href: '/dashboard/student/quizzes', icon: 'help-circle' },
    { label: 'Nilai', href: '/dashboard/student/grades', icon: 'bar-chart-2' },
  ],
};

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const router = useRouter();
  const pathname = usePathname();
  const { user, accessToken, logout } = useAuthStore();
  const { sidebarOpen, toggleSidebar } = useUIStore();
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
    if (!accessToken) {
      router.push('/login');
    }
  }, [accessToken, router]);

  if (!mounted || !accessToken || !user) {
    return (
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', height: '100vh', background: 'var(--color-bg)' }}>
        <div className="loading-pulse" />
      </div>
    );
  }

  const role = user.role as keyof typeof navItems;
  const items = navItems[role] || navItems.student;

  const handleLogout = () => {
    logout();
    router.push('/login');
  };

  const roleLabel: Record<string, string> = {
    super_admin: 'Admin',
    teacher: 'Guru',
    class_leader: 'Ketua Kelas',
    student: 'Siswa',
  };

  return (
    <div className="dashboard-layout">
      {/* Mobile overlay */}
      {sidebarOpen && (
        <div className="sidebar-overlay" onClick={toggleSidebar} />
      )}

      {/* Sidebar */}
      <aside className={`sidebar ${sidebarOpen ? 'open' : ''}`}>
        <div className="sidebar-header">
          <div className="sidebar-logo">
            <svg width="32" height="32" viewBox="0 0 48 48" fill="none">
              <rect width="48" height="48" rx="12" fill="var(--color-primary)"/>
              <path d="M24 12L14 18V30L24 36L34 30V18L24 12Z" stroke="white" strokeWidth="2" fill="none"/>
            </svg>
            <span>AkuBelajar</span>
          </div>
          <button className="sidebar-close" onClick={toggleSidebar}>
            ✕
          </button>
        </div>

        <nav className="sidebar-nav">
          {items.map((item) => (
            <a
              key={item.href}
              href={item.href}
              className={`nav-item ${pathname === item.href ? 'active' : ''}`}
              onClick={(e) => {
                e.preventDefault();
                router.push(item.href);
                if (window.innerWidth < 768) toggleSidebar();
              }}
            >
              <span className="nav-icon">{getIcon(item.icon)}</span>
              <span>{item.label}</span>
            </a>
          ))}
        </nav>

        <div className="sidebar-footer">
          <div className="user-info">
            <div className="user-avatar">
              {user.email?.charAt(0).toUpperCase()}
            </div>
            <div className="user-details">
              <span className="user-name">{user.email?.split('@')[0]}</span>
              <span className="user-role">{roleLabel[user.role] || user.role}</span>
            </div>
          </div>
          <button className="btn-logout" onClick={handleLogout}>Keluar</button>
        </div>
      </aside>

      {/* Main content */}
      <main className="main-content">
        {/* Topbar */}
        <header className="topbar">
          <button className="topbar-menu" onClick={toggleSidebar}>
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <line x1="3" y1="6" x2="21" y2="6"/>
              <line x1="3" y1="12" x2="21" y2="12"/>
              <line x1="3" y1="18" x2="21" y2="18"/>
            </svg>
          </button>

          <div className="topbar-right">
            <button className="topbar-notif">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M18 8A6 6 0 006 8c0 7-3 9-3 9h18s-3-2-3-9"/>
                <path d="M13.73 21a2 2 0 01-3.46 0"/>
              </svg>
            </button>
          </div>
        </header>

        <div className="content-area">
          {children}
        </div>
      </main>

      <style jsx>{`
        .dashboard-layout {
          display: flex;
          min-height: 100vh;
          background: var(--color-bg);
        }

        /* Sidebar */
        .sidebar {
          width: 260px;
          background: var(--color-surface);
          border-right: 1px solid var(--color-border);
          display: flex;
          flex-direction: column;
          position: fixed;
          top: 0;
          left: 0;
          bottom: 0;
          z-index: 50;
          transform: translateX(-100%);
          transition: transform 0.3s ease;
        }

        .sidebar.open {
          transform: translateX(0);
        }

        @media (min-width: 768px) {
          .sidebar {
            position: sticky;
            top: 0;
            height: 100vh;
            transform: translateX(0);
          }
          .sidebar-close, .sidebar-overlay { display: none !important; }
          .topbar-menu { display: none !important; }
        }

        .sidebar-overlay {
          position: fixed;
          inset: 0;
          background: rgba(0,0,0,0.5);
          z-index: 40;
        }

        .sidebar-header {
          display: flex;
          align-items: center;
          justify-content: space-between;
          padding: 1.25rem 1rem;
          border-bottom: 1px solid var(--color-border);
        }

        .sidebar-logo {
          display: flex;
          align-items: center;
          gap: 0.75rem;
          font-weight: 700;
          font-size: 1.1rem;
          color: var(--color-text);
        }

        .sidebar-close {
          background: none;
          border: none;
          color: var(--color-text-secondary);
          font-size: 1.25rem;
          cursor: pointer;
        }

        .sidebar-nav {
          flex: 1;
          padding: 0.75rem 0.5rem;
          overflow-y: auto;
          display: flex;
          flex-direction: column;
          gap: 2px;
        }

        .nav-item {
          display: flex;
          align-items: center;
          gap: 0.75rem;
          padding: 0.625rem 0.75rem;
          border-radius: 8px;
          color: var(--color-text-secondary);
          text-decoration: none;
          font-size: 0.9rem;
          transition: all 0.15s;
        }

        .nav-item:hover {
          background: var(--color-border);
          color: var(--color-text);
        }

        .nav-item.active {
          background: rgba(59, 130, 246, 0.1);
          color: var(--color-primary);
          font-weight: 500;
        }

        .nav-icon {
          display: flex;
          align-items: center;
          width: 20px;
        }

        .sidebar-footer {
          padding: 1rem;
          border-top: 1px solid var(--color-border);
        }

        .user-info {
          display: flex;
          align-items: center;
          gap: 0.75rem;
          margin-bottom: 0.75rem;
        }

        .user-avatar {
          width: 36px;
          height: 36px;
          border-radius: 50%;
          background: var(--color-primary);
          color: white;
          display: flex;
          align-items: center;
          justify-content: center;
          font-weight: 600;
          font-size: 0.9rem;
        }

        .user-details {
          display: flex;
          flex-direction: column;
          min-width: 0;
        }

        .user-name {
          font-size: 0.85rem;
          font-weight: 500;
          color: var(--color-text);
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }

        .user-role {
          font-size: 0.75rem;
          color: var(--color-text-tertiary);
        }

        .btn-logout {
          width: 100%;
          padding: 0.5rem;
          border: 1px solid var(--color-border);
          border-radius: 6px;
          background: transparent;
          color: var(--color-text-secondary);
          font-size: 0.85rem;
          cursor: pointer;
          transition: all 0.15s;
        }

        .btn-logout:hover {
          border-color: hsl(0, 60%, 40%);
          color: hsl(0, 70%, 60%);
          background: hsl(0, 60%, 10%);
        }

        /* Main content */
        .main-content {
          flex: 1;
          display: flex;
          flex-direction: column;
          min-width: 0;
        }

        .topbar {
          display: flex;
          align-items: center;
          justify-content: space-between;
          padding: 0.75rem 1.5rem;
          border-bottom: 1px solid var(--color-border);
          background: var(--color-surface);
          position: sticky;
          top: 0;
          z-index: 30;
        }

        .topbar-menu {
          background: none;
          border: none;
          color: var(--color-text-secondary);
          cursor: pointer;
          padding: 4px;
        }

        .topbar-right {
          display: flex;
          align-items: center;
          gap: 0.5rem;
        }

        .topbar-notif {
          background: none;
          border: none;
          color: var(--color-text-secondary);
          cursor: pointer;
          padding: 6px;
          border-radius: 6px;
          transition: background 0.15s;
        }

        .topbar-notif:hover {
          background: var(--color-border);
        }

        .content-area {
          flex: 1;
          padding: 1.5rem;
          overflow-y: auto;
        }

        @media (max-width: 767px) {
          .content-area { padding: 1rem; }
        }

        .loading-pulse {
          width: 40px;
          height: 40px;
          border-radius: 50%;
          background: var(--color-primary);
          animation: pulse 1s ease-in-out infinite;
        }

        @keyframes pulse {
          0%, 100% { transform: scale(0.8); opacity: 0.5; }
          50% { transform: scale(1); opacity: 1; }
        }
      `}</style>
    </div>
  );
}

// Simple icon component using inline SVGs
function getIcon(name: string): React.ReactNode {
  const icons: Record<string, React.ReactNode> = {
    'grid': <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/><rect x="14" y="14" width="7" height="7"/><rect x="3" y="14" width="7" height="7"/></svg>,
    'calendar': <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><rect x="3" y="4" width="18" height="18" rx="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/></svg>,
    'users': <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 00-3-3.87"/><path d="M16 3.13a4 4 0 010 7.75"/></svg>,
    'book': <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M4 19.5A2.5 2.5 0 016.5 17H20"/><path d="M6.5 2H20v20H6.5A2.5 2.5 0 014 19.5v-15A2.5 2.5 0 016.5 2z"/></svg>,
    'user-plus': <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M16 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2"/><circle cx="8.5" cy="7" r="4"/><line x1="20" y1="8" x2="20" y2="14"/><line x1="23" y1="11" x2="17" y2="11"/></svg>,
    'shield': <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg>,
    'file-text': <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/></svg>,
    'help-circle': <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><circle cx="12" cy="12" r="10"/><path d="M9.09 9a3 3 0 015.83 1c0 2-3 3-3 3"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>,
    'check-square': <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><polyline points="9 11 12 14 22 4"/><path d="M21 12v7a2 2 0 01-2 2H5a2 2 0 01-2-2V5a2 2 0 012-2h11"/></svg>,
    'bar-chart-2': <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><line x1="18" y1="20" x2="18" y2="10"/><line x1="12" y1="20" x2="12" y2="4"/><line x1="6" y1="20" x2="6" y2="14"/></svg>,
  };
  return icons[name] || <span>•</span>;
}
