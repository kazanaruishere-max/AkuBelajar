'use client';

import React, { useState, useEffect, useCallback } from 'react';
import { api } from '@/lib/api/client';

interface ClassItem {
  id: string;
  name: string;
  grade_level: number;
  academic_year_id: string;
  academic_year_name?: string;
  homeroom_teacher_id?: string;
  homeroom_teacher?: string;
  student_count: number;
  created_at: string;
}

interface AcademicYear {
  id: string;
  name: string;
  is_active: boolean;
}

export default function ClassesPage() {
  const [classes, setClasses] = useState<ClassItem[]>([]);
  const [years, setYears] = useState<AcademicYear[]>([]);
  const [selectedYear, setSelectedYear] = useState('');
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [editing, setEditing] = useState<ClassItem | null>(null);
  const [form, setForm] = useState({ name: '', grade_level: 7, academic_year_id: '' });
  const [saving, setSaving] = useState(false);

  // Fetch years on mount
  useEffect(() => {
    (async () => {
      try {
        const res = await api.get('/api/v1/academic/years');
        const yrs: AcademicYear[] = (res.data as any).data || [];
        setYears(yrs);
        const active = yrs.find((y) => y.is_active);
        if (active) setSelectedYear(active.id);
        else if (yrs.length > 0) setSelectedYear(yrs[0].id);
      } catch { /* ignore */ }
    })();
  }, []);

  const fetchClasses = useCallback(async () => {
    if (!selectedYear) return;
    try {
      setLoading(true);
      const res = await api.get(`/api/v1/academic/classes?academic_year_id=${selectedYear}`);
      setClasses((res.data as any).data || []);
    } catch { /* ignore */ } finally {
      setLoading(false);
    }
  }, [selectedYear]);

  useEffect(() => { fetchClasses(); }, [fetchClasses]);

  const openCreate = () => {
    setEditing(null);
    setForm({ name: '', grade_level: 7, academic_year_id: selectedYear });
    setShowModal(true);
  };

  const openEdit = (c: ClassItem) => {
    setEditing(c);
    setForm({ name: c.name, grade_level: c.grade_level, academic_year_id: c.academic_year_id });
    setShowModal(true);
  };

  const handleSave = async () => {
    setSaving(true);
    try {
      if (editing) {
        await api.put(`/api/v1/academic/classes/${editing.id}`, form);
      } else {
        await api.post('/api/v1/academic/classes', form);
      }
      setShowModal(false);
      fetchClasses();
    } catch { /* ignore */ } finally {
      setSaving(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm('Hapus kelas ini?')) return;
    try {
      await api.delete(`/api/v1/academic/classes/${id}`);
      fetchClasses();
    } catch { /* ignore */ }
  };

  return (
    <div style={{ padding: '32px' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
        <div>
          <h1 style={{ fontSize: 24, fontWeight: 700, color: '#E2E8F0', marginBottom: 8 }}>Kelas</h1>
          <select value={selectedYear} onChange={(e) => setSelectedYear(e.target.value)} style={selectStyle}>
            {years.map((y) => (
              <option key={y.id} value={y.id}>{y.name}{y.is_active ? ' (Aktif)' : ''}</option>
            ))}
          </select>
        </div>
        <button onClick={openCreate} style={btnPrimary}>+ Tambah Kelas</button>
      </div>

      {loading ? (
        <p style={{ color: '#94A3B8' }}>Memuat...</p>
      ) : classes.length === 0 ? (
        <div style={emptyState}>
          <p style={{ color: '#94A3B8', fontSize: 16 }}>Belum ada kelas di tahun ajaran ini</p>
          <button onClick={openCreate} style={btnPrimary}>Buat Kelas Pertama</button>
        </div>
      ) : (
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(280px, 1fr))', gap: 16 }}>
          {classes.map((c) => (
            <div key={c.id} style={card}>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
                <div>
                  <h3 style={{ fontSize: 20, fontWeight: 700, color: '#E2E8F0' }}>{c.name}</h3>
                  <p style={{ fontSize: 13, color: '#64748B', marginTop: 4 }}>Kelas {c.grade_level}</p>
                </div>
                <span style={badge}>{c.student_count} siswa</span>
              </div>
              {c.homeroom_teacher && (
                <p style={{ fontSize: 13, color: '#94A3B8', marginTop: 12 }}>
                  🧑‍🏫 {c.homeroom_teacher}
                </p>
              )}
              <div style={{ display: 'flex', gap: 8, marginTop: 16 }}>
                <button onClick={() => openEdit(c)} style={btnSmall}>Edit</button>
                <button onClick={() => handleDelete(c.id)} style={{ ...btnSmall, color: '#EF4444' }}>Hapus</button>
              </div>
            </div>
          ))}
        </div>
      )}

      {showModal && (
        <div style={overlay} onClick={() => setShowModal(false)}>
          <div style={modal} onClick={(e) => e.stopPropagation()}>
            <h2 style={{ fontSize: 20, fontWeight: 700, color: '#E2E8F0', marginBottom: 20 }}>
              {editing ? 'Edit Kelas' : 'Tambah Kelas'}
            </h2>
            <div style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
              <label style={labelStyle}>
                Nama Kelas
                <input style={inputStyle} value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })}
                  placeholder="Contoh: 7A" />
              </label>
              <label style={labelStyle}>
                Tingkat
                <select style={selectStyle} value={form.grade_level}
                  onChange={(e) => setForm({ ...form, grade_level: parseInt(e.target.value) })}>
                  {[7, 8, 9, 10, 11, 12].map((g) => (
                    <option key={g} value={g}>Kelas {g}</option>
                  ))}
                </select>
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
const card: React.CSSProperties = {
  background: 'rgba(15, 23, 42, 0.6)', borderRadius: 16, padding: 20,
  border: '1px solid rgba(148, 163, 184, 0.1)', transition: 'all 0.2s',
};
const badge: React.CSSProperties = {
  padding: '4px 12px', borderRadius: 20, fontSize: 12, fontWeight: 600,
  background: 'rgba(59, 130, 246, 0.15)', color: '#60A5FA',
};
const selectStyle: React.CSSProperties = {
  padding: '10px 14px', borderRadius: 10, border: '1px solid rgba(148, 163, 184, 0.2)',
  background: 'rgba(15, 23, 42, 0.8)', color: '#E2E8F0', fontSize: 14, outline: 'none',
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
  display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center',
  gap: 16, padding: 64, background: 'rgba(15, 23, 42, 0.4)', borderRadius: 16,
  border: '1px dashed rgba(148, 163, 184, 0.2)',
};
