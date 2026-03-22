'use client';

import React, { useState, useEffect, useCallback } from 'react';
import { api } from '@/lib/api/client';
import { useAuthStore } from '@/lib/store/authStore';

interface AcademicYear {
  id: string;
  name: string;
  start_date: string;
  end_date: string;
  is_active: boolean;
  created_at: string;
}

export default function AcademicYearsPage() {
  const [years, setYears] = useState<AcademicYear[]>([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [editing, setEditing] = useState<AcademicYear | null>(null);
  const [form, setForm] = useState({ name: '', start_date: '', end_date: '', is_active: false });
  const [saving, setSaving] = useState(false);
  const token = useAuthStore((s) => s.accessToken);

  const fetchYears = useCallback(async () => {
    try {
      setLoading(true);
      const res = await api.get('/api/v1/academic/years');
      setYears((res.data as any).data || []);
    } catch { /* ignore */ } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => { fetchYears(); }, [fetchYears]);

  const openCreate = () => {
    setEditing(null);
    setForm({ name: '', start_date: '', end_date: '', is_active: false });
    setShowModal(true);
  };

  const openEdit = (y: AcademicYear) => {
    setEditing(y);
    setForm({ name: y.name, start_date: y.start_date, end_date: y.end_date, is_active: y.is_active });
    setShowModal(true);
  };

  const handleSave = async () => {
    setSaving(true);
    try {
      if (editing) {
        await api.put(`/api/v1/academic/years/${editing.id}`, form);
      } else {
        await api.post('/api/v1/academic/years', form);
      }
      setShowModal(false);
      fetchYears();
    } catch { /* ignore */ } finally {
      setSaving(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm('Hapus tahun ajaran ini?')) return;
    try {
      await api.delete(`/api/v1/academic/years/${id}`);
      fetchYears();
    } catch { /* ignore */ }
  };

  return (
    <div style={{ padding: '32px' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
        <h1 style={{ fontSize: 24, fontWeight: 700, color: '#E2E8F0' }}>Tahun Ajaran</h1>
        <button onClick={openCreate} style={btnPrimary}>+ Tambah Tahun Ajaran</button>
      </div>

      {loading ? (
        <p style={{ color: '#94A3B8' }}>Memuat...</p>
      ) : years.length === 0 ? (
        <div style={emptyState}>
          <p style={{ color: '#94A3B8', fontSize: 16 }}>Belum ada tahun ajaran</p>
          <button onClick={openCreate} style={btnPrimary}>Buat Tahun Ajaran Pertama</button>
        </div>
      ) : (
        <div style={{ overflowX: 'auto' }}>
          <table style={tableStyle}>
            <thead>
              <tr>
                <th style={th}>Nama</th>
                <th style={th}>Mulai</th>
                <th style={th}>Selesai</th>
                <th style={th}>Status</th>
                <th style={th}>Aksi</th>
              </tr>
            </thead>
            <tbody>
              {years.map((y) => (
                <tr key={y.id} style={{ borderBottom: '1px solid rgba(148, 163, 184, 0.1)' }}>
                  <td style={td}>{y.name}</td>
                  <td style={td}>{y.start_date}</td>
                  <td style={td}>{y.end_date}</td>
                  <td style={td}>
                    <span style={{
                      padding: '4px 12px', borderRadius: 20, fontSize: 12, fontWeight: 600,
                      backgroundColor: y.is_active ? 'rgba(34,197,94,0.15)' : 'rgba(148,163,184,0.1)',
                      color: y.is_active ? '#22C55E' : '#94A3B8',
                    }}>{y.is_active ? 'Aktif' : 'Nonaktif'}</span>
                  </td>
                  <td style={td}>
                    <button onClick={() => openEdit(y)} style={btnSmall}>Edit</button>
                    {!y.is_active && (
                      <button onClick={() => handleDelete(y.id)} style={{ ...btnSmall, color: '#EF4444', marginLeft: 8 }}>Hapus</button>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {/* Modal */}
      {showModal && (
        <div style={overlay} onClick={() => setShowModal(false)}>
          <div style={modal} onClick={(e) => e.stopPropagation()}>
            <h2 style={{ fontSize: 20, fontWeight: 700, color: '#E2E8F0', marginBottom: 20 }}>
              {editing ? 'Edit Tahun Ajaran' : 'Tambah Tahun Ajaran'}
            </h2>
            <div style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
              <label style={labelStyle}>
                Nama
                <input style={inputStyle} value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })}
                  placeholder="Contoh: 2025/2026" />
              </label>
              <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 12 }}>
                <label style={labelStyle}>
                  Tanggal Mulai
                  <input type="date" style={inputStyle} value={form.start_date}
                    onChange={(e) => setForm({ ...form, start_date: e.target.value })} />
                </label>
                <label style={labelStyle}>
                  Tanggal Selesai
                  <input type="date" style={inputStyle} value={form.end_date}
                    onChange={(e) => setForm({ ...form, end_date: e.target.value })} />
                </label>
              </div>
              <label style={{ display: 'flex', alignItems: 'center', gap: 8, color: '#CBD5E1', fontSize: 14 }}>
                <input type="checkbox" checked={form.is_active}
                  onChange={(e) => setForm({ ...form, is_active: e.target.checked })} />
                Aktifkan tahun ajaran ini
              </label>
            </div>
            <div style={{ display: 'flex', justifyContent: 'flex-end', gap: 12, marginTop: 24 }}>
              <button onClick={() => setShowModal(false)} style={btnGhost}>Batal</button>
              <button onClick={handleSave} disabled={saving} style={btnPrimary}>
                {saving ? 'Menyimpan...' : 'Simpan'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

// ── Styles ────────────────────────────────────────────────

const btnPrimary: React.CSSProperties = {
  padding: '10px 20px', borderRadius: 10, border: 'none', cursor: 'pointer',
  background: 'linear-gradient(135deg, #3B82F6, #2563EB)', color: '#fff',
  fontWeight: 600, fontSize: 14, transition: 'all 0.2s',
};
const btnGhost: React.CSSProperties = {
  padding: '10px 20px', borderRadius: 10, border: '1px solid rgba(148,163,184,0.2)',
  cursor: 'pointer', background: 'transparent', color: '#94A3B8', fontWeight: 600, fontSize: 14,
};
const btnSmall: React.CSSProperties = {
  padding: '6px 12px', borderRadius: 8, border: 'none', cursor: 'pointer',
  background: 'rgba(148,163,184,0.1)', color: '#3B82F6', fontWeight: 500, fontSize: 13,
};
const tableStyle: React.CSSProperties = {
  width: '100%', borderCollapse: 'collapse',
  background: 'rgba(15, 23, 42, 0.6)', borderRadius: 12,
};
const th: React.CSSProperties = {
  padding: '14px 16px', textAlign: 'left', fontSize: 12, fontWeight: 600,
  color: '#64748B', textTransform: 'uppercase', letterSpacing: '0.05em',
  borderBottom: '1px solid rgba(148, 163, 184, 0.1)',
};
const td: React.CSSProperties = {
  padding: '14px 16px', fontSize: 14, color: '#CBD5E1',
};
const overlay: React.CSSProperties = {
  position: 'fixed', inset: 0, background: 'rgba(0, 0, 0, 0.6)',
  display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 50,
};
const modal: React.CSSProperties = {
  background: '#1E293B', borderRadius: 16, padding: 32, width: '100%', maxWidth: 480,
  border: '1px solid rgba(148, 163, 184, 0.1)', boxShadow: '0 25px 50px rgba(0,0,0,0.5)',
};
const labelStyle: React.CSSProperties = {
  display: 'flex', flexDirection: 'column', gap: 6, color: '#CBD5E1', fontSize: 14, fontWeight: 500,
};
const inputStyle: React.CSSProperties = {
  padding: '10px 14px', borderRadius: 10, border: '1px solid rgba(148, 163, 184, 0.2)',
  background: 'rgba(15, 23, 42, 0.8)', color: '#E2E8F0', fontSize: 14, outline: 'none',
};
const emptyState: React.CSSProperties = {
  display: 'flex', flexDirection: 'column', alignItems: 'center',
  justifyContent: 'center', gap: 16, padding: 64,
  background: 'rgba(15, 23, 42, 0.4)', borderRadius: 16,
  border: '1px dashed rgba(148, 163, 184, 0.2)',
};
