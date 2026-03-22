'use client';

import React, { useState, useEffect, useCallback } from 'react';
import { api } from '@/lib/api/client';

interface Quiz {
  id: string;
  class_name: string;
  subject_name: string;
  title: string;
  time_limit: number;
  status: string;
  question_count: number;
  session_count: number;
  max_attempts: number;
  randomize_questions: boolean;
  created_at: string;
}

export default function TeacherQuizzesPage() {
  const [quizzes, setQuizzes] = useState<Quiz[]>([]);
  const [loading, setLoading] = useState(true);
  const [showCreate, setShowCreate] = useState(false);
  const [showQuestions, setShowQuestions] = useState<string | null>(null);
  const [form, setForm] = useState({
    class_id: '', subject_id: '', title: '', time_limit: 30,
    randomize_questions: true, randomize_options: true, max_attempts: 1, allow_review: true,
    start_at: '', end_at: '',
  });
  const [qForm, setQForm] = useState({ question_text: '', question_type: 'multiple_choice' as string, options: '', order_num: 1 });
  const [saving, setSaving] = useState(false);

  const fetchQuizzes = useCallback(async () => {
    try {
      setLoading(true);
      const res = await api.get('/api/v1/quizzes/teacher');
      setQuizzes((res.data as any).data || []);
    } catch { /* ignore */ } finally { setLoading(false); }
  }, []);

  useEffect(() => { fetchQuizzes(); }, [fetchQuizzes]);

  const handleCreate = async () => {
    setSaving(true);
    try {
      await api.post('/api/v1/quizzes/teacher', form);
      setShowCreate(false);
      fetchQuizzes();
    } catch { /* ignore */ } finally { setSaving(false); }
  };

  const handlePublish = async (id: string) => {
    await api.post(`/api/v1/quizzes/teacher/${id}/publish`);
    fetchQuizzes();
  };

  const handleDelete = async (id: string) => {
    if (!confirm('Hapus kuis ini?')) return;
    await api.delete(`/api/v1/quizzes/teacher/${id}`);
    fetchQuizzes();
  };

  const handleAddQuestion = async () => {
    if (!showQuestions) return;
    setSaving(true);
    try {
      let opts: any = null;
      if (qForm.question_type === 'multiple_choice' && qForm.options) {
        try { opts = JSON.parse(qForm.options); } catch { opts = null; }
      }
      await api.post(`/api/v1/quizzes/teacher/${showQuestions}/questions`, {
        ...qForm, options: opts ? JSON.stringify(opts) : null,
      });
      setQForm({ question_text: '', question_type: 'multiple_choice', options: '', order_num: qForm.order_num + 1 });
    } catch { /* ignore */ } finally { setSaving(false); }
  };

  const sc = (s: string) => {
    const map: Record<string, { bg: string; fg: string }> = {
      draft: { bg: 'rgba(251,191,36,0.15)', fg: '#FBBF24' },
      published: { bg: 'rgba(34,197,94,0.15)', fg: '#22C55E' },
      active: { bg: 'rgba(59,130,246,0.15)', fg: '#60A5FA' },
      ended: { bg: 'rgba(148,163,184,0.15)', fg: '#94A3B8' },
    };
    return map[s] || { bg: 'rgba(148,163,184,0.1)', fg: '#94A3B8' };
  };

  return (
    <div style={{ padding: 32 }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
        <h1 style={{ fontSize: 24, fontWeight: 700, color: '#E2E8F0' }}>Kuis Saya</h1>
        <button onClick={() => setShowCreate(true)} style={btnPrimary}>+ Buat Kuis</button>
      </div>

      {loading ? <p style={{ color: '#94A3B8' }}>Memuat...</p> : quizzes.length === 0 ? (
        <div style={emptyState}><p style={{ color: '#94A3B8' }}>Belum ada kuis</p></div>
      ) : (
        <div style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
          {quizzes.map((q) => {
            const color = sc(q.status);
            return (
              <div key={q.id} style={card}>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
                  <div>
                    <h3 style={{ fontSize: 16, fontWeight: 600, color: '#E2E8F0' }}>{q.title}</h3>
                    <p style={{ fontSize: 13, color: '#94A3B8', marginTop: 4 }}>
                      📚 {q.class_name} · {q.subject_name} · ⏱️ {q.time_limit} menit
                    </p>
                  </div>
                  <span style={{ padding: '4px 12px', borderRadius: 20, fontSize: 12, fontWeight: 600,
                    background: color.bg, color: color.fg, textTransform: 'capitalize' }}>{q.status}</span>
                </div>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginTop: 16 }}>
                  <div style={{ display: 'flex', gap: 16 }}>
                    <span style={{ fontSize: 13, color: '#64748B' }}>📝 {q.question_count} soal</span>
                    <span style={{ fontSize: 13, color: '#64748B' }}>👥 {q.session_count} peserta</span>
                  </div>
                  <div style={{ display: 'flex', gap: 8 }}>
                    <button onClick={() => setShowQuestions(q.id)} style={btnSmall}>+ Soal</button>
                    {q.status === 'draft' && <button onClick={() => handlePublish(q.id)} style={{ ...btnSmall, color: '#22C55E' }}>Publish</button>}
                    <button onClick={() => handleDelete(q.id)} style={{ ...btnSmall, color: '#EF4444' }}>Hapus</button>
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      )}

      {/* Create Quiz Modal */}
      {showCreate && (
        <div style={overlay} onClick={() => setShowCreate(false)}>
          <div style={modal} onClick={(e) => e.stopPropagation()}>
            <h2 style={{ fontSize: 20, fontWeight: 700, color: '#E2E8F0', marginBottom: 20 }}>Buat Kuis Baru</h2>
            <div style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
              <label style={labelStyle}>Judul<input style={inputStyle} value={form.title} onChange={(e) => setForm({ ...form, title: e.target.value })} /></label>
              <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 12 }}>
                <label style={labelStyle}>Class ID<input style={inputStyle} value={form.class_id} onChange={(e) => setForm({ ...form, class_id: e.target.value })} /></label>
                <label style={labelStyle}>Subject ID<input style={inputStyle} value={form.subject_id} onChange={(e) => setForm({ ...form, subject_id: e.target.value })} /></label>
              </div>
              <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr 1fr', gap: 12 }}>
                <label style={labelStyle}>Waktu (menit)<input type="number" style={inputStyle} value={form.time_limit} onChange={(e) => setForm({ ...form, time_limit: +e.target.value })} /></label>
                <label style={labelStyle}>Max Attempt<input type="number" style={inputStyle} value={form.max_attempts} onChange={(e) => setForm({ ...form, max_attempts: +e.target.value })} /></label>
                <label style={{ ...labelStyle, justifyContent: 'flex-end' }}>
                  <span style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
                    <input type="checkbox" checked={form.randomize_questions} onChange={(e) => setForm({ ...form, randomize_questions: e.target.checked })} /> Acak soal
                  </span>
                </label>
              </div>
            </div>
            <div style={{ display: 'flex', justifyContent: 'flex-end', gap: 12, marginTop: 24 }}>
              <button onClick={() => setShowCreate(false)} style={btnGhost}>Batal</button>
              <button onClick={handleCreate} disabled={saving} style={btnPrimary}>{saving ? 'Menyimpan...' : 'Simpan'}</button>
            </div>
          </div>
        </div>
      )}

      {/* Add Question Modal */}
      {showQuestions && (
        <div style={overlay} onClick={() => setShowQuestions(null)}>
          <div style={modal} onClick={(e) => e.stopPropagation()}>
            <h2 style={{ fontSize: 20, fontWeight: 700, color: '#E2E8F0', marginBottom: 20 }}>Tambah Soal</h2>
            <div style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
              <label style={labelStyle}>No. Soal<input type="number" style={inputStyle} value={qForm.order_num} onChange={(e) => setQForm({ ...qForm, order_num: +e.target.value })} /></label>
              <label style={labelStyle}>
                Tipe
                <select style={inputStyle} value={qForm.question_type} onChange={(e) => setQForm({ ...qForm, question_type: e.target.value })}>
                  <option value="multiple_choice">Pilihan Ganda</option>
                  <option value="essay">Essay</option>
                </select>
              </label>
              <label style={labelStyle}>Pertanyaan<textarea style={{ ...inputStyle, minHeight: 80 }} value={qForm.question_text} onChange={(e) => setQForm({ ...qForm, question_text: e.target.value })} /></label>
              {qForm.question_type === 'multiple_choice' && (
                <label style={labelStyle}>
                  Opsi (JSON array)
                  <textarea style={{ ...inputStyle, minHeight: 60, fontFamily: 'monospace', fontSize: 12 }}
                    value={qForm.options} onChange={(e) => setQForm({ ...qForm, options: e.target.value })}
                    placeholder='[{"key":"A","text":"...","is_correct":false}]' />
                </label>
              )}
            </div>
            <div style={{ display: 'flex', justifyContent: 'flex-end', gap: 12, marginTop: 24 }}>
              <button onClick={() => setShowQuestions(null)} style={btnGhost}>Tutup</button>
              <button onClick={handleAddQuestion} disabled={saving} style={btnPrimary}>{saving ? 'Menambahkan...' : 'Tambah Soal'}</button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

const btnPrimary: React.CSSProperties = { padding: '10px 20px', borderRadius: 10, border: 'none', cursor: 'pointer', background: 'linear-gradient(135deg,#3B82F6,#2563EB)', color: '#fff', fontWeight: 600, fontSize: 14 };
const btnGhost: React.CSSProperties = { padding: '10px 20px', borderRadius: 10, border: '1px solid rgba(148,163,184,0.2)', cursor: 'pointer', background: 'transparent', color: '#94A3B8', fontWeight: 600, fontSize: 14 };
const btnSmall: React.CSSProperties = { padding: '6px 12px', borderRadius: 8, border: 'none', cursor: 'pointer', background: 'rgba(148,163,184,0.1)', color: '#3B82F6', fontWeight: 500, fontSize: 13 };
const card: React.CSSProperties = { background: 'rgba(15,23,42,0.6)', borderRadius: 16, padding: 20, border: '1px solid rgba(148,163,184,0.1)' };
const overlay: React.CSSProperties = { position: 'fixed', inset: 0, background: 'rgba(0,0,0,0.6)', display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 50 };
const modal: React.CSSProperties = { background: '#1E293B', borderRadius: 16, padding: 32, width: '100%', maxWidth: 560, border: '1px solid rgba(148,163,184,0.1)', boxShadow: '0 25px 50px rgba(0,0,0,0.5)', maxHeight: '90vh', overflowY: 'auto' };
const labelStyle: React.CSSProperties = { display: 'flex', flexDirection: 'column', gap: 6, color: '#CBD5E1', fontSize: 14, fontWeight: 500 };
const inputStyle: React.CSSProperties = { padding: '10px 14px', borderRadius: 10, border: '1px solid rgba(148,163,184,0.2)', background: 'rgba(15,23,42,0.8)', color: '#E2E8F0', fontSize: 14, outline: 'none' };
const emptyState: React.CSSProperties = { display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', gap: 16, padding: 64, background: 'rgba(15,23,42,0.4)', borderRadius: 16, border: '1px dashed rgba(148,163,184,0.2)' };
