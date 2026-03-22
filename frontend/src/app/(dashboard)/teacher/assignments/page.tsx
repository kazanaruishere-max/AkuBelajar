'use client';

import React, { useState, useEffect, useCallback } from 'react';
import { api } from '@/lib/api/client';

interface Assignment {
  id: string;
  class_name: string;
  subject_name: string;
  title: string;
  description: string | null;
  deadline_at: string;
  status: string;
  submission_count: number;
  graded_count: number;
  weight_pct: number;
  created_at: string;
}

export default function TeacherAssignmentsPage() {
  const [assignments, setAssignments] = useState<Assignment[]>([]);
  const [loading, setLoading] = useState(true);
  const [showCreate, setShowCreate] = useState(false);
  const [form, setForm] = useState({
    class_id: '', subject_id: '', title: '', description: '',
    deadline_at: '', allow_late: true, late_penalty_pct: 10,
    max_late_days: 5, max_file_count: 5, max_file_size_mb: 20, weight_pct: 100,
  });
  const [saving, setSaving] = useState(false);

  const fetchAssignments = useCallback(async () => {
    try {
      setLoading(true);
      const res = await api.get('/api/v1/assignments/teacher');
      setAssignments((res.data as any).data || []);
    } catch { /* ignore */ } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => { fetchAssignments(); }, [fetchAssignments]);

  const handleCreate = async () => {
    setSaving(true);
    try {
      await api.post('/api/v1/assignments/teacher', form);
      setShowCreate(false);
      fetchAssignments();
    } catch { /* ignore */ } finally {
      setSaving(false);
    }
  };

  const handlePublish = async (id: string) => {
    try {
      await api.post(`/api/v1/assignments/teacher/${id}/publish`);
      fetchAssignments();
    } catch { /* ignore */ }
  };

  const handleClose = async (id: string) => {
    try {
      await api.post(`/api/v1/assignments/teacher/${id}/close`);
      fetchAssignments();
    } catch { /* ignore */ }
  };

  const handleDelete = async (id: string) => {
    if (!confirm('Hapus tugas ini?')) return;
    try {
      await api.delete(`/api/v1/assignments/teacher/${id}`);
      fetchAssignments();
    } catch { /* ignore */ }
  };

  const statusColor = (s: string) => {
    switch (s) {
      case 'draft': return { bg: 'rgba(251,191,36,0.15)', fg: '#FBBF24' };
      case 'published': return { bg: 'rgba(34,197,94,0.15)', fg: '#22C55E' };
      case 'closed': return { bg: 'rgba(148,163,184,0.15)', fg: '#94A3B8' };
      default: return { bg: 'rgba(148,163,184,0.1)', fg: '#94A3B8' };
    }
  };

  const formatDate = (d: string) => new Date(d).toLocaleDateString('id-ID', {
    day: 'numeric', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit',
  });

  return (
    <div style={{ padding: '32px' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
        <h1 style={{ fontSize: 24, fontWeight: 700, color: '#E2E8F0' }}>Tugas Saya</h1>
        <button onClick={() => setShowCreate(true)} style={btnPrimary}>+ Buat Tugas</button>
      </div>

      {loading ? (
        <p style={{ color: '#94A3B8' }}>Memuat...</p>
      ) : assignments.length === 0 ? (
        <div style={emptyState}>
          <p style={{ color: '#94A3B8', fontSize: 16 }}>Belum ada tugas dibuat</p>
          <button onClick={() => setShowCreate(true)} style={btnPrimary}>Buat Tugas Pertama</button>
        </div>
      ) : (
        <div style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
          {assignments.map((a) => {
            const sc = statusColor(a.status);
            return (
              <div key={a.id} style={card}>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
                  <div style={{ flex: 1 }}>
                    <h3 style={{ fontSize: 16, fontWeight: 600, color: '#E2E8F0' }}>{a.title}</h3>
                    <div style={{ display: 'flex', gap: 12, marginTop: 6, flexWrap: 'wrap' }}>
                      <span style={{ fontSize: 13, color: '#94A3B8' }}>📚 {a.class_name} · {a.subject_name}</span>
                      <span style={{ fontSize: 13, color: '#94A3B8' }}>⏰ {formatDate(a.deadline_at)}</span>
                    </div>
                  </div>
                  <span style={{ padding: '4px 12px', borderRadius: 20, fontSize: 12, fontWeight: 600,
                    background: sc.bg, color: sc.fg, textTransform: 'capitalize' }}>{a.status}</span>
                </div>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginTop: 16 }}>
                  <div style={{ display: 'flex', gap: 16 }}>
                    <span style={{ fontSize: 13, color: '#64748B' }}>📝 {a.submission_count} dikumpulkan</span>
                    <span style={{ fontSize: 13, color: '#64748B' }}>✅ {a.graded_count} dinilai</span>
                  </div>
                  <div style={{ display: 'flex', gap: 8 }}>
                    {a.status === 'draft' && (
                      <button onClick={() => handlePublish(a.id)} style={{ ...btnSmall, color: '#22C55E' }}>Publish</button>
                    )}
                    {a.status === 'published' && (
                      <button onClick={() => handleClose(a.id)} style={{ ...btnSmall, color: '#F59E0B' }}>Tutup</button>
                    )}
                    <button onClick={() => handleDelete(a.id)} style={{ ...btnSmall, color: '#EF4444' }}>Hapus</button>
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      )}

      {/* Create Modal */}
      {showCreate && (
        <div style={overlay} onClick={() => setShowCreate(false)}>
          <div style={modal} onClick={(e) => e.stopPropagation()}>
            <h2 style={{ fontSize: 20, fontWeight: 700, color: '#E2E8F0', marginBottom: 20 }}>Buat Tugas Baru</h2>
            <div style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
              <label style={labelStyle}>
                Judul Tugas
                <input style={inputStyle} value={form.title}
                  onChange={(e) => setForm({ ...form, title: e.target.value })} placeholder="Contoh: PR Matematika Bab 3" />
              </label>
              <label style={labelStyle}>
                Deskripsi
                <textarea style={{ ...inputStyle, minHeight: 80, resize: 'vertical' }} value={form.description}
                  onChange={(e) => setForm({ ...form, description: e.target.value })}
                  placeholder="Instruksi pengerjaan..." />
              </label>
              <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 12 }}>
                <label style={labelStyle}>
                  Class ID
                  <input style={inputStyle} value={form.class_id}
                    onChange={(e) => setForm({ ...form, class_id: e.target.value })} placeholder="UUID kelas" />
                </label>
                <label style={labelStyle}>
                  Subject ID
                  <input style={inputStyle} value={form.subject_id}
                    onChange={(e) => setForm({ ...form, subject_id: e.target.value })} placeholder="UUID mapel" />
                </label>
              </div>
              <label style={labelStyle}>
                Deadline
                <input type="datetime-local" style={inputStyle} value={form.deadline_at}
                  onChange={(e) => setForm({ ...form, deadline_at: e.target.value })} />
              </label>
              <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr 1fr', gap: 12 }}>
                <label style={labelStyle}>
                  Bobot (%)
                  <input type="number" style={inputStyle} value={form.weight_pct}
                    onChange={(e) => setForm({ ...form, weight_pct: parseInt(e.target.value) || 0 })} />
                </label>
                <label style={labelStyle}>
                  Penalty (%/hari)
                  <input type="number" style={inputStyle} value={form.late_penalty_pct}
                    onChange={(e) => setForm({ ...form, late_penalty_pct: parseInt(e.target.value) || 0 })} />
                </label>
                <label style={labelStyle}>
                  Max telat (hari)
                  <input type="number" style={inputStyle} value={form.max_late_days}
                    onChange={(e) => setForm({ ...form, max_late_days: parseInt(e.target.value) || 0 })} />
                </label>
              </div>
            </div>
            <div style={{ display: 'flex', justifyContent: 'flex-end', gap: 12, marginTop: 24 }}>
              <button onClick={() => setShowCreate(false)} style={btnGhost}>Batal</button>
              <button onClick={handleCreate} disabled={saving} style={btnPrimary}>
                {saving ? 'Menyimpan...' : 'Simpan Draft'}
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
  border: '1px solid rgba(148, 163, 184, 0.1)',
};
const overlay: React.CSSProperties = {
  position: 'fixed', inset: 0, background: 'rgba(0, 0, 0, 0.6)',
  display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 50,
};
const modal: React.CSSProperties = {
  background: '#1E293B', borderRadius: 16, padding: 32, width: '100%', maxWidth: 560,
  border: '1px solid rgba(148, 163, 184, 0.1)', boxShadow: '0 25px 50px rgba(0,0,0,0.5)',
  maxHeight: '90vh', overflowY: 'auto',
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
