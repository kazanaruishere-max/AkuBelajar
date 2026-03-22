'use client';

import { useAuthStore } from '@/lib/store/authStore';

export default function DashboardPage() {
  const { user } = useAuthStore();

  const roleGreeting: Record<string, string> = {
    super_admin: 'Admin',
    teacher: 'Bapak/Ibu Guru',
    class_leader: 'Ketua Kelas',
    student: 'Siswa',
  };

  const greeting = roleGreeting[user?.role || 'student'] || 'User';

  return (
    <div>
      <div className="dashboard-welcome">
        <h1>Selamat Datang, {greeting}! 👋</h1>
        <p>Berikut ringkasan aktivitas Anda hari ini.</p>
      </div>

      <div className="stats-grid">
        <div className="stat-card">
          <div className="stat-icon blue">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z"/>
              <polyline points="14 2 14 8 20 8"/>
            </svg>
          </div>
          <div className="stat-info">
            <span className="stat-value">—</span>
            <span className="stat-label">Tugas Aktif</span>
          </div>
        </div>

        <div className="stat-card">
          <div className="stat-icon green">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <polyline points="9 11 12 14 22 4"/>
              <path d="M21 12v7a2 2 0 01-2 2H5a2 2 0 01-2-2V5a2 2 0 012-2h11"/>
            </svg>
          </div>
          <div className="stat-info">
            <span className="stat-value">—</span>
            <span className="stat-label">Absensi Hari Ini</span>
          </div>
        </div>

        <div className="stat-card">
          <div className="stat-icon purple">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <circle cx="12" cy="12" r="10"/>
              <path d="M9.09 9a3 3 0 015.83 1c0 2-3 3-3 3"/>
              <line x1="12" y1="17" x2="12.01" y2="17"/>
            </svg>
          </div>
          <div className="stat-info">
            <span className="stat-value">—</span>
            <span className="stat-label">Kuis Mendatang</span>
          </div>
        </div>

        <div className="stat-card">
          <div className="stat-icon amber">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <line x1="18" y1="20" x2="18" y2="10"/>
              <line x1="12" y1="20" x2="12" y2="4"/>
              <line x1="6" y1="20" x2="6" y2="14"/>
            </svg>
          </div>
          <div className="stat-info">
            <span className="stat-value">—</span>
            <span className="stat-label">Rata-rata Nilai</span>
          </div>
        </div>
      </div>

      <style jsx>{`
        .dashboard-welcome {
          margin-bottom: 2rem;
        }

        .dashboard-welcome h1 {
          font-size: 1.5rem;
          font-weight: 700;
          color: var(--color-text);
          margin: 0 0 0.25rem;
        }

        .dashboard-welcome p {
          color: var(--color-text-secondary);
          margin: 0;
        }

        .stats-grid {
          display: grid;
          grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
          gap: 1rem;
        }

        .stat-card {
          display: flex;
          align-items: center;
          gap: 1rem;
          padding: 1.25rem;
          background: var(--color-surface);
          border: 1px solid var(--color-border);
          border-radius: 12px;
          transition: border-color 0.15s;
        }

        .stat-card:hover {
          border-color: var(--color-text-tertiary);
        }

        .stat-icon {
          width: 48px;
          height: 48px;
          border-radius: 10px;
          display: flex;
          align-items: center;
          justify-content: center;
          flex-shrink: 0;
        }

        .stat-icon.blue { background: rgba(59, 130, 246, 0.15); color: hsl(217, 91%, 60%); }
        .stat-icon.green { background: rgba(34, 197, 94, 0.15); color: hsl(142, 71%, 45%); }
        .stat-icon.purple { background: rgba(168, 85, 247, 0.15); color: hsl(271, 91%, 65%); }
        .stat-icon.amber { background: rgba(245, 158, 11, 0.15); color: hsl(38, 92%, 50%); }

        .stat-info {
          display: flex;
          flex-direction: column;
        }

        .stat-value {
          font-size: 1.5rem;
          font-weight: 700;
          color: var(--color-text);
        }

        .stat-label {
          font-size: 0.8rem;
          color: var(--color-text-secondary);
        }
      `}</style>
    </div>
  );
}
