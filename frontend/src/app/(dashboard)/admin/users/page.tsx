'use client';

import React, { useState, useEffect, useCallback } from 'react';
import { api } from '@/lib/api/client';

interface User {
  id: string;
  email: string;
  full_name: string;
  role: string;
  is_active: boolean;
  is_first_login: boolean;
  created_at: string;
}

export default function AdminUsersPage() {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [showCreate, setShowCreate] = useState(false);
  const [editUser, setEditUser] = useState<User | null>(null);
  const [form, setForm] = useState({ email: '', full_name: '', role: 'student', password: '' });
  const [saving, setSaving] = useState(false);
  const [search, setSearch] = useState('');

  const fetchUsers = useCallback(async () => {
    try {
      setLoading(true);
      const res = await api.get('/api/v1/admin/users');
      setUsers((res.data as any).data || []);
    } catch { /* ignore */ } finally { setLoading(false); }
  }, []);

  useEffect(() => { fetchUsers(); }, [fetchUsers]);

  const handleCreate = async () => {
    setSaving(true);
    try {
      await api.post('/api/v1/admin/users', form);
      setShowCreate(false);
      setForm({ email: '', full_name: '', role: 'student', password: '' });
      fetchUsers();
    } catch { /* ignore */ } finally { setSaving(false); }
  };

  const handleUpdate = async () => {
    if (!editUser) return;
    setSaving(true);
    try {
      await api.put(`/api/v1/admin/users/${editUser.id}`, {
        full_name: editUser.full_name, role: editUser.role, is_active: editUser.is_active,
      });
      setEditUser(null);
      fetchUsers();
    } catch { /* ignore */ } finally { setSaving(false); }
  };

  const handleDelete = async (id: string) => {
    if (!confirm('Hapus pengguna ini?')) return;
    await api.delete(`/api/v1/admin/users/${id}`);
    fetchUsers();
  };

  const roleColors: Record<string, { bg: string; fg: string }> = {
    super_admin:  { bg: 'rgba(239,68,68,0.15)', fg: '#EF4444' },
    teacher:      { bg: 'rgba(59,130,246,0.15)', fg: '#60A5FA' },
    student:      { bg: 'rgba(34,197,94,0.15)', fg: '#22C55E' },
    class_leader: { bg: 'rgba(251,191,36,0.15)', fg: '#FBBF24' },
  };

  const filtered = users.filter((u) =>
    u.email.toLowerCase().includes(search.toLowerCase()) ||
    u.full_name.toLowerCase().includes(search.toLowerCase()) ||
    u.role.toLowerCase().includes(search.toLowerCase())
  );

  return (
    <div style={{ padding: 32 }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
        <h1 style={{ fontSize: 24, fontWeight: 700, color: '#E2E8F0' }}>Manajemen Pengguna</h1>
        <button onClick={() => setShowCreate(true)} style={btnPrimary}>+ Tambah Pengguna</button>
      </div>

      {/* Search */}
      <input style={{ ...inputStyle, width: '100%', maxWidth: 400, marginBottom: 20 }}
        placeholder="🔍 Cari email, nama, atau role..."
        value={search} onChange={(e) => setSearch(e.target.value)} />

      {/* Stats */}
      <div style={{ display: 'flex', gap: 12, marginBottom: 24, flexWrap: 'wrap' }}>
        {['super_admin', 'teacher', 'student', 'class_leader'].map((r) => {
          const c = roleColors[r];
          const count = users.filter((u) => u.role === r).length;
          return (
            <span key={r} style={{ padding: '6px 14px', borderRadius: 20, fontSize: 13, fontWeight: 600, background: c.bg, color: c.fg }}>
              {r.replace('_', ' ')} ({count})
            </span>
          );
        })}
      </div>

      {/* Table */}
      {loading ? <p style={{ color: '#94A3B8' }}>Memuat...</p> : (
        <div style={{ overflowX: 'auto' }}>
          <table style={tableStyle}>
            <thead>
              <tr>
                <th style={th}>Email</th>
                <th style={th}>Nama</th>
                <th style={th}>Role</th>
                <th style={th}>Status</th>
                <th style={th}>Aksi</th>
              </tr>
            </thead>
            <tbody>
              {filtered.map((u) => {
                const rc = roleColors[u.role] || roleColors.student;
                return (
                  <tr key={u.id} style={{ borderBottom: '1px solid rgba(148,163,184,0.1)' }}>
                    <td style={{ ...td, fontWeight: 500 }}>{u.email}</td>
                    <td style={td}>{u.full_name}</td>
                    <td style={td}>
                      <span style={{ padding: '3px 10px', borderRadius: 8, fontSize: 12, fontWeight: 600, background: rc.bg, color: rc.fg, textTransform: 'capitalize' }}>
                        {u.role.replace('_', ' ')}
                      </span>
                    </td>
                    <td style={td}>
                      <span style={{ width: 8, height: 8, borderRadius: '50%', display: 'inline-block', marginRight: 8,
                        background: u.is_active ? '#22C55E' : '#EF4444' }} />
                      {u.is_active ? 'Aktif' : 'Nonaktif'}
                    </td>
                    <td style={td}>
                      <div style={{ display: 'flex', gap: 8 }}>
                        <button onClick={() => setEditUser({ ...u })} style={{ ...btnSmall, color: '#60A5FA' }}>Edit</button>
                        <button onClick={() => handleDelete(u.id)} style={{ ...btnSmall, color: '#EF4444' }}>Hapus</button>
                      </div>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
          <p style={{ fontSize: 12, color: '#64748B', marginTop: 12 }}>Menampilkan {filtered.length} dari {users.length} pengguna</p>
        </div>
      )}

      {/* Create Modal */}
      {showCreate && (
        <div style={overlay} onClick={() => setShowCreate(false)}>
          <div style={modal} onClick={(e) => e.stopPropagation()}>
            <h2 style={{ fontSize: 20, fontWeight: 700, color: '#E2E8F0', marginBottom: 20 }}>Tambah Pengguna</h2>
            <div style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
              <label style={labelStyle}>Email<input style={inputStyle} value={form.email} onChange={(e) => setForm({ ...form, email: e.target.value })} /></label>
              <label style={labelStyle}>Nama Lengkap<input style={inputStyle} value={form.full_name} onChange={(e) => setForm({ ...form, full_name: e.target.value })} /></label>
              <label style={labelStyle}>
                Role
                <select style={inputStyle} value={form.role} onChange={(e) => setForm({ ...form, role: e.target.value })}>
                  <option value="student">Student</option>
                  <option value="class_leader">Class Leader</option>
                  <option value="teacher">Teacher</option>
                  <option value="super_admin">Super Admin</option>
                </select>
              </label>
              <label style={labelStyle}>Password<input type="password" style={inputStyle} value={form.password} onChange={(e) => setForm({ ...form, password: e.target.value })} /></label>
            </div>
            <div style={{ display: 'flex', justifyContent: 'flex-end', gap: 12, marginTop: 24 }}>
              <button onClick={() => setShowCreate(false)} style={btnGhost}>Batal</button>
              <button onClick={handleCreate} disabled={saving} style={btnPrimary}>{saving ? 'Menyimpan...' : 'Simpan'}</button>
            </div>
          </div>
        </div>
      )}

      {/* Edit Modal */}
      {editUser && (
        <div style={overlay} onClick={() => setEditUser(null)}>
          <div style={modal} onClick={(e) => e.stopPropagation()}>
            <h2 style={{ fontSize: 20, fontWeight: 700, color: '#E2E8F0', marginBottom: 20 }}>Edit Pengguna</h2>
            <div style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
              <label style={labelStyle}>Nama Lengkap<input style={inputStyle} value={editUser.full_name} onChange={(e) => setEditUser({ ...editUser, full_name: e.target.value })} /></label>
              <label style={labelStyle}>
                Role
                <select style={inputStyle} value={editUser.role} onChange={(e) => setEditUser({ ...editUser, role: e.target.value })}>
                  <option value="student">Student</option>
                  <option value="class_leader">Class Leader</option>
                  <option value="teacher">Teacher</option>
                  <option value="super_admin">Super Admin</option>
                </select>
              </label>
              <label style={{ display: 'flex', alignItems: 'center', gap: 8, color: '#CBD5E1', fontSize: 14 }}>
                <input type="checkbox" checked={editUser.is_active} onChange={(e) => setEditUser({ ...editUser, is_active: e.target.checked })} />
                Aktif
              </label>
            </div>
            <div style={{ display: 'flex', justifyContent: 'flex-end', gap: 12, marginTop: 24 }}>
              <button onClick={() => setEditUser(null)} style={btnGhost}>Batal</button>
              <button onClick={handleUpdate} disabled={saving} style={btnPrimary}>{saving ? 'Menyimpan...' : 'Simpan'}</button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

const btnPrimary: React.CSSProperties = { padding: '10px 20px', borderRadius: 10, border: 'none', cursor: 'pointer', background: 'linear-gradient(135deg,#3B82F6,#2563EB)', color: '#fff', fontWeight: 600, fontSize: 14 };
const btnGhost: React.CSSProperties = { padding: '10px 20px', borderRadius: 10, border: '1px solid rgba(148,163,184,0.2)', cursor: 'pointer', background: 'transparent', color: '#94A3B8', fontWeight: 600, fontSize: 14 };
const btnSmall: React.CSSProperties = { padding: '6px 12px', borderRadius: 8, border: 'none', cursor: 'pointer', background: 'rgba(148,163,184,0.1)', fontWeight: 500, fontSize: 13 };
const tableStyle: React.CSSProperties = { width: '100%', borderCollapse: 'collapse', background: 'rgba(15,23,42,0.6)', borderRadius: 12 };
const th: React.CSSProperties = { padding: '14px 16px', textAlign: 'left', fontSize: 12, fontWeight: 600, color: '#64748B', textTransform: 'uppercase', letterSpacing: '0.05em', borderBottom: '1px solid rgba(148,163,184,0.1)' };
const td: React.CSSProperties = { padding: '14px 16px', fontSize: 14, color: '#CBD5E1' };
const overlay: React.CSSProperties = { position: 'fixed', inset: 0, background: 'rgba(0,0,0,0.6)', display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 50 };
const modal: React.CSSProperties = { background: '#1E293B', borderRadius: 16, padding: 32, width: '100%', maxWidth: 480, border: '1px solid rgba(148,163,184,0.1)', boxShadow: '0 25px 50px rgba(0,0,0,0.5)' };
const labelStyle: React.CSSProperties = { display: 'flex', flexDirection: 'column', gap: 6, color: '#CBD5E1', fontSize: 14, fontWeight: 500 };
const inputStyle: React.CSSProperties = { padding: '10px 14px', borderRadius: 10, border: '1px solid rgba(148,163,184,0.2)', background: 'rgba(15,23,42,0.8)', color: '#E2E8F0', fontSize: 14, outline: 'none' };
