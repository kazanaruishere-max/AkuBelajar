'use client';

import React, { useState, useEffect, useCallback } from 'react';
import { api } from '@/lib/api/client';

interface Subject {
  id: string;
  name: string;
  code: string | null;
  description: string | null;
  created_at: string;
}

export default function SubjectsPage() {
  const [subjects, setSubjects] = useState<Subject[]>([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [editing, setEditing] = useState<Subject | null>(null);
  const [form, setForm] = useState({ name: '', code: '', description: '' });
  const [saving, setSaving] = useState(false);

  const fetchSubjects = useCallback(async () => {
    try {
      setLoading(true);
      const res = await api.get('/api/v1/academic/subjects');
      setSubjects((res.data as any).data || []);
    } catch { /* ignore */ } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => { fetchSubjects(); }, [fetchSubjects]);

  const openCreate = () => {
    setEditing(null);
    setForm({ name: '', code: '', description: '' });
    setShowModal(true);
  };

  const openEdit = (s: Subject) => {
    setEditing(s);
    setForm({ name: s.name, code: s.code || '', description: s.description || '' });
    setShowModal(true);
  };

  const handleSave = async () => {
    setSaving(true);
    try {
      const payload = {
        name: form.name,
        code: form.code || null,
        description: form.description || null,
      };
      if (editing) {
        await api.put(`/api/v1/academic/subjects/${editing.id}`, payload);
      } else {
        await api.post('/api/v1/academic/subjects', payload);
      }
      setShowModal(false);
      fetchSubjects();
    } catch { /* ignore */ } finally {
      setSaving(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm('Hapus mata pelajaran ini?')) return;
    try {
      await api.delete(`/api/v1/academic/subjects/${id}`);
      fetchSubjects();
    } catch { /* ignore */ }
  };

  return (
    <div style={{ padding: '32px' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
        <h1 style={{ fontSize: 24, fontWeight: 700, color: '#E2E8F0' }}>Mata Pelajaran</h1>
        <button onClick={openCreate} style={btnPrimary}>+ Tambah Mapel</button>
      </div>

      {loading ? (
        <p style={{ color: '#94A3B8' }}>Memuat...</p>
      ) : subjects.length === 0 ? (
        <div style={emptyState}>
          <p style={{ color: '#94A3B8', fontSize: 16 }}>Belum ada mata pelajaran</p>
          <button onClick={openCreate} style={btnPrimary}>Buat Mapel Pertama</button>
        </div>
      ) : (
        <div style={{ overflowX: 'auto' }}>
          <table style={tableStyle}>
            <thead>
              <tr>
                <th style={th}>Kode</th>
                <th style={th}>Nama</th>
                <th style={th}>Deskripsi</th>
                <th style={th}>Aksi</th>
              </tr>
            </thead>
            <tbody>
              {subjects.map((s) => (
                <tr key={s.id} style={{ borderBottom: '1px solid rgba(148,163,184,0.1)' }}>
                  <td style={td}>
                    <span style={{ padding: '4px 10px', borderRadius: 6, fontSize: 12, fontWeight: 600,
                      background: 'rgba(99,102,241,0.15)', color: '#818CF8' }}>
                      {s.code || '-'}
                    </span>
                  </td>
                  <td style={{ ...td, fontWeight: 500 }}>{s.name}</td>
                  <td style={{ ...td, color: '#94A3B8', maxWidth: 300, overflow: 'hidden', textOverflow: 'ellipsis' }}>
                    {s.description || '-'}
                  </td>
                  <td style={td}>
                    <button onClick={() => openEdit(s)} style={btnSmall}>Edit</button>
                    <button onClick={() => handleDelete(s.id)} style={{ ...btnSmall, color: '#EF4444', marginLeft: 8 }}>Hapus</button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {showModal && (
        <div style={overlay} onClick={() => setShowModal(false)}>
          <div style={modal} onClick={(e) => e.stopPropagation()}>
            <h2 style={{ fontSize: 20, fontWeight: 700, color: '#E2E8F0', marginBottom: 20 }}>
              {editing ? 'Edit Mata Pelajaran' : 'Tambah Mata Pelajaran'}
            </h2>
            <div style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
              <div style={{ display: 'grid', gridTemplateColumns: '100px 1fr', gap: 12 }}>
                <label style={labelStyle}>
                  Kode
                  <input style={inputStyle} value={form.code} onChange={(e) => setForm({ ...form, code: e.target.value })}
                    placeholder="MTK" maxLength={10} />
                </label>
                <label style={labelStyle}>
                  Nama
                  <input style={inputStyle} value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })}
                    placeholder="Matematika" />
                </label>
              </div>
              <label style={labelStyle}>
                Deskripsi (opsional)
                <textarea style={{ ...inputStyle, minHeight: 80, resize: 'vertical' }} value={form.description}
                  onChange={(e) => setForm({ ...form, description: e.target.value })}
                  placeholder="Deskripsi singkat..." />
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

const btnPrimary: React.CSSProperties = {
  padding: '10px 20px', borderRadius: 10, border: 'none', cursor: 'pointer',
  background: 'linear-gradient(135deg, #3B82F6, #2563EB)', color: '#fff', fontWeight: 600, fontSize: 14,
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
  width: '100%', borderCollapse: 'collapse', background: 'rgba(15, 23, 42, 0.6)', borderRadius: 12,
};
const th: React.CSSProperties = {
  padding: '14px 16px', textAlign: 'left', fontSize: 12, fontWeight: 600, color: '#64748B',
  textTransform: 'uppercase', letterSpacing: '0.05em', borderBottom: '1px solid rgba(148, 163, 184, 0.1)',
};
const td: React.CSSProperties = { padding: '14px 16px', fontSize: 14, color: '#CBD5E1' };
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
  display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center',
  gap: 16, padding: 64, background: 'rgba(15, 23, 42, 0.4)', borderRadius: 16,
  border: '1px dashed rgba(148, 163, 184, 0.2)',
};
