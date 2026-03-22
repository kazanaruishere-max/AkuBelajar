'use client';

import { useState, FormEvent } from 'react';
import { useRouter } from 'next/navigation';
import { useAuthStore } from '@/lib/store/authStore';
import { api } from '@/lib/api/client';

interface LoginResponse {
  access_token: string;
  refresh_token: string;
  user: {
    id: string;
    email: string;
    role: string;
    school_id: string;
    is_first_login: boolean;
  };
}

export default function LoginPage() {
  const router = useRouter();
  const { login: storeLogin } = useAuthStore();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [rememberMe, setRememberMe] = useState(false);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [showPassword, setShowPassword] = useState(false);

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      const res = await api.post<LoginResponse>('/auth/login', {
        email,
        password,
        remember_me: rememberMe,
      });

      const { user, access_token, refresh_token } = res.data;
      storeLogin(
        { ...user, role: user.role as 'super_admin' | 'teacher' | 'class_leader' | 'student' },
        access_token,
        refresh_token,
      );
      router.push('/dashboard');
    } catch (err: unknown) {
      const apiErr = err as { message?: string };
      setError(apiErr?.message || 'Terjadi kesalahan. Coba lagi.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="login-container">
      <div className="login-card">
        <div className="login-header">
          <div className="login-logo">
            <svg width="48" height="48" viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
              <rect width="48" height="48" rx="12" fill="var(--color-primary)"/>
              <path d="M24 12L14 18V30L24 36L34 30V18L24 12Z" stroke="white" strokeWidth="2" fill="none"/>
              <path d="M24 12V36M14 18L34 30M34 18L14 30" stroke="white" strokeWidth="1.5" opacity="0.5"/>
            </svg>
          </div>
          <h1>AkuBelajar</h1>
          <p>Masuk ke akun Anda</p>
        </div>

        {error && (
          <div className="login-error">
            <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
              <path d="M8 1a7 7 0 100 14A7 7 0 008 1zm0 10.5a.75.75 0 110-1.5.75.75 0 010 1.5zM8.75 8a.75.75 0 01-1.5 0V4.5a.75.75 0 011.5 0V8z"/>
            </svg>
            <span>{error}</span>
          </div>
        )}

        <form onSubmit={handleSubmit} className="login-form">
          <div className="form-group">
            <label htmlFor="email">Email</label>
            <div className="input-wrapper">
              <svg className="input-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"/>
                <polyline points="22,6 12,13 2,6"/>
              </svg>
              <input
                id="email"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="admin@sekolah.sch.id"
                required
                autoComplete="email"
                autoFocus
              />
            </div>
          </div>

          <div className="form-group">
            <label htmlFor="password">Password</label>
            <div className="input-wrapper">
              <svg className="input-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
                <path d="M7 11V7a5 5 0 0110 0v4"/>
              </svg>
              <input
                id="password"
                type={showPassword ? 'text' : 'password'}
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="Masukkan password"
                required
                autoComplete="current-password"
                minLength={8}
              />
              <button
                type="button"
                className="toggle-password"
                onClick={() => setShowPassword(!showPassword)}
                tabIndex={-1}
              >
                {showPassword ? (
                  <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                    <path d="M17.94 17.94A10.07 10.07 0 0112 20c-7 0-11-8-11-8a18.45 18.45 0 015.06-5.94M9.9 4.24A9.12 9.12 0 0112 4c7 0 11 8 11 8a18.5 18.5 0 01-2.16 3.19m-6.72-1.07a3 3 0 11-4.24-4.24"/>
                    <line x1="1" y1="1" x2="23" y2="23"/>
                  </svg>
                ) : (
                  <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                    <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                    <circle cx="12" cy="12" r="3"/>
                  </svg>
                )}
              </button>
            </div>
          </div>

          <div className="form-options">
            <label className="checkbox-label">
              <input
                type="checkbox"
                checked={rememberMe}
                onChange={(e) => setRememberMe(e.target.checked)}
              />
              <span>Ingat saya</span>
            </label>
            <a href="/forgot-password" className="forgot-link">Lupa password?</a>
          </div>

          <button type="submit" className="btn-login" disabled={loading}>
            {loading ? <span className="spinner" /> : 'Masuk'}
          </button>
        </form>

        <div className="login-footer">
          <p>Belum punya akun? Minta kode undangan dari admin sekolah.</p>
        </div>
      </div>

      <style jsx>{`
        .login-container {
          min-height: 100vh;
          display: flex;
          align-items: center;
          justify-content: center;
          padding: 1rem;
          background: linear-gradient(135deg, var(--color-bg) 0%, hsl(220, 30%, 12%) 100%);
        }
        .login-card {
          width: 100%;
          max-width: 420px;
          background: var(--color-surface);
          border: 1px solid var(--color-border);
          border-radius: 16px;
          padding: 2.5rem 2rem;
          box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
        }
        .login-header { text-align: center; margin-bottom: 2rem; }
        .login-logo { display: flex; justify-content: center; margin-bottom: 1rem; }
        .login-header h1 { font-size: 1.75rem; font-weight: 700; color: var(--color-text); margin: 0; }
        .login-header p { color: var(--color-text-secondary); margin: 0.25rem 0 0; font-size: 0.9rem; }
        .login-error {
          display: flex; align-items: center; gap: 0.5rem; padding: 0.75rem 1rem;
          border-radius: 8px; background: hsl(0, 70%, 15%); border: 1px solid hsl(0, 60%, 30%);
          color: hsl(0, 80%, 70%); font-size: 0.85rem; margin-bottom: 1.5rem;
          animation: slideDown 0.2s ease-out;
        }
        @keyframes slideDown { from { opacity: 0; transform: translateY(-8px); } to { opacity: 1; transform: translateY(0); } }
        .login-form { display: flex; flex-direction: column; gap: 1.25rem; }
        .form-group { display: flex; flex-direction: column; gap: 0.375rem; }
        .form-group label { font-size: 0.85rem; font-weight: 500; color: var(--color-text-secondary); }
        .input-wrapper { position: relative; display: flex; align-items: center; }
        .input-icon { position: absolute; left: 12px; color: var(--color-text-tertiary); pointer-events: none; }
        .input-wrapper input {
          width: 100%; padding: 0.75rem 0.75rem 0.75rem 2.75rem;
          border: 1px solid var(--color-border); border-radius: 8px;
          background: var(--color-bg); color: var(--color-text);
          font-size: 0.95rem; transition: border-color 0.2s, box-shadow 0.2s; outline: none;
        }
        .input-wrapper input:focus { border-color: var(--color-primary); box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.15); }
        .input-wrapper input::placeholder { color: var(--color-text-tertiary); }
        .toggle-password {
          position: absolute; right: 8px; background: none; border: none;
          color: var(--color-text-tertiary); cursor: pointer; padding: 4px; border-radius: 4px;
          display: flex; align-items: center;
        }
        .toggle-password:hover { color: var(--color-text-secondary); }
        .form-options { display: flex; align-items: center; justify-content: space-between; }
        .checkbox-label {
          display: flex; align-items: center; gap: 0.5rem;
          font-size: 0.85rem; color: var(--color-text-secondary); cursor: pointer;
        }
        .checkbox-label input[type="checkbox"] { width: 16px; height: 16px; accent-color: var(--color-primary); cursor: pointer; }
        .forgot-link { font-size: 0.85rem; color: var(--color-primary); text-decoration: none; }
        .forgot-link:hover { text-decoration: underline; }
        .btn-login {
          width: 100%; padding: 0.75rem; border: none; border-radius: 8px;
          background: var(--color-primary); color: white; font-size: 1rem; font-weight: 600;
          cursor: pointer; transition: background 0.2s, transform 0.1s;
          display: flex; align-items: center; justify-content: center; min-height: 44px;
        }
        .btn-login:hover:not(:disabled) { background: var(--color-primary-hover); transform: translateY(-1px); }
        .btn-login:active:not(:disabled) { transform: translateY(0); }
        .btn-login:disabled { opacity: 0.7; cursor: not-allowed; }
        .spinner {
          width: 20px; height: 20px; border: 2px solid rgba(255, 255, 255, 0.3);
          border-top-color: white; border-radius: 50%; animation: spin 0.6s linear infinite;
        }
        @keyframes spin { to { transform: rotate(360deg); } }
        .login-footer { text-align: center; margin-top: 1.5rem; padding-top: 1.5rem; border-top: 1px solid var(--color-border); }
        .login-footer p { font-size: 0.8rem; color: var(--color-text-tertiary); margin: 0; }
      `}</style>
    </div>
  );
}
