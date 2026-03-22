'use client';

import React, { useState, useEffect, useCallback } from 'react';
import { api } from '@/lib/api/client';

interface Assignment {
  id: string;
  class_name: string;
  subject_name: string;
  teacher_email: string;
  title: string;
  description: string | null;
  deadline_at: string;
  allow_late: boolean;
  max_late_days: number;
  late_penalty_pct: number;
  status: string;
  weight_pct: number;
}

interface Submission {
  id: string;
  submitted_at: string | null;
  is_late: boolean;
  late_days: number;
  status: string;
  grade: number | null;
  grade_after_penalty: number | null;
  feedback: string | null;
  files: { id: string; file_name: string }[];
}

export default function StudentAssignmentsPage() {
  const [assignments, setAssignments] = useState<Assignment[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedAssignment, setSelectedAssignment] = useState<Assignment | null>(null);
  const [submission, setSubmission] = useState<Submission | null>(null);
  const [submitting, setSubmitting] = useState(false);

  // This would come from student's class context — for now use query
  const classId = typeof window !== 'undefined'
    ? new URLSearchParams(window.location.search).get('class_id') || ''
    : '';

  const fetchAssignments = useCallback(async () => {
    if (!classId) { setLoading(false); return; }
    try {
      setLoading(true);
      const res = await api.get(`/api/v1/assignments/student?class_id=${classId}`);
      setAssignments((res.data as any).data || []);
    } catch { /* ignore */ } finally {
      setLoading(false);
    }
  }, [classId]);

  useEffect(() => { fetchAssignments(); }, [fetchAssignments]);

  const viewAssignment = async (a: Assignment) => {
    setSelectedAssignment(a);
    try {
      const res = await api.get(`/api/v1/assignments/student/${a.id}/my-submission`);
      setSubmission((res.data as any).data);
    } catch {
      setSubmission(null);
    }
  };

  const handleSubmit = async () => {
    if (!selectedAssignment) return;
    setSubmitting(true);
    try {
      const res = await api.post(`/api/v1/assignments/student/${selectedAssignment.id}/submit`);
      setSubmission((res.data as any).data);
    } catch { /* ignore */ } finally {
      setSubmitting(false);
    }
  };

  const formatDate = (d: string) => new Date(d).toLocaleDateString('id-ID', {
    day: 'numeric', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit',
  });

  const isOverdue = (d: string) => new Date(d) < new Date();

  return (
    <div style={{ padding: '32px' }}>
      <h1 style={{ fontSize: 24, fontWeight: 700, color: '#E2E8F0', marginBottom: 24 }}>Tugas</h1>

      {!classId ? (
        <div style={emptyState}>
          <p style={{ color: '#94A3B8' }}>Pilih kelas untuk melihat tugas (tambahkan ?class_id=... di URL)</p>
        </div>
      ) : loading ? (
        <p style={{ color: '#94A3B8' }}>Memuat...</p>
      ) : assignments.length === 0 ? (
        <div style={emptyState}>
          <p style={{ color: '#94A3B8', fontSize: 16 }}>🎉 Belum ada tugas!</p>
        </div>
      ) : (
        <div style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
          {assignments.map((a) => {
            const overdue = isOverdue(a.deadline_at);
            return (
              <div key={a.id} style={card} onClick={() => viewAssignment(a)}>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
                  <div>
                    <h3 style={{ fontSize: 16, fontWeight: 600, color: '#E2E8F0' }}>{a.title}</h3>
                    <p style={{ fontSize: 13, color: '#94A3B8', marginTop: 4 }}>
                      📚 {a.subject_name} · 🧑‍🏫 {a.teacher_email}
                    </p>
                  </div>
                  <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'flex-end', gap: 4 }}>
                    <span style={{
                      padding: '4px 10px', borderRadius: 20, fontSize: 11, fontWeight: 600,
                      background: overdue ? 'rgba(239,68,68,0.15)' : 'rgba(34,197,94,0.15)',
                      color: overdue ? '#EF4444' : '#22C55E',
                    }}>{overdue ? 'Lewat deadline' : 'Masih aktif'}</span>
                    <span style={{ fontSize: 12, color: '#64748B' }}>⏰ {formatDate(a.deadline_at)}</span>
                  </div>
                </div>
                {a.description && (
                  <p style={{ fontSize: 13, color: '#94A3B8', marginTop: 12, lineHeight: 1.5 }}>{a.description}</p>
                )}
              </div>
            );
          })}
        </div>
      )}

      {/* Detail + Submit Modal */}
      {selectedAssignment && (
        <div style={overlay} onClick={() => setSelectedAssignment(null)}>
          <div style={modal} onClick={(e) => e.stopPropagation()}>
            <h2 style={{ fontSize: 20, fontWeight: 700, color: '#E2E8F0' }}>{selectedAssignment.title}</h2>
            <p style={{ fontSize: 13, color: '#94A3B8', marginTop: 4 }}>
              {selectedAssignment.subject_name} · Deadline: {formatDate(selectedAssignment.deadline_at)}
            </p>
            {selectedAssignment.description && (
              <p style={{ fontSize: 14, color: '#CBD5E1', marginTop: 16, lineHeight: 1.6 }}>
                {selectedAssignment.description}
              </p>
            )}

            <div style={{ marginTop: 20, padding: 16, borderRadius: 12, background: 'rgba(15,23,42,0.5)',
              border: '1px solid rgba(148,163,184,0.1)' }}>
              {submission ? (
                <div>
                  <p style={{ fontSize: 14, color: '#E2E8F0', fontWeight: 600 }}>Status: {submission.status}</p>
                  {submission.grade !== null && (
                    <div style={{ marginTop: 12 }}>
                      <p style={{ color: '#CBD5E1' }}>Nilai: <strong style={{ color: '#22C55E', fontSize: 20 }}>{submission.grade_after_penalty ?? submission.grade}</strong>/100</p>
                      {submission.is_late && <p style={{ color: '#F59E0B', fontSize: 12, marginTop: 4 }}>Terlambat {submission.late_days} hari (penalti diterapkan)</p>}
                      {submission.feedback && <p style={{ color: '#94A3B8', marginTop: 8, fontSize: 13 }}>💬 {submission.feedback}</p>}
                    </div>
                  )}
                </div>
              ) : (
                <div style={{ textAlign: 'center' }}>
                  <p style={{ color: '#94A3B8', marginBottom: 16 }}>Belum dikumpulkan</p>
                  <button onClick={handleSubmit} disabled={submitting} style={btnPrimary}>
                    {submitting ? 'Mengirim...' : '📤 Kumpulkan Tugas'}
                  </button>
                </div>
              )}
            </div>

            <div style={{ display: 'flex', justifyContent: 'flex-end', marginTop: 20 }}>
              <button onClick={() => setSelectedAssignment(null)} style={btnGhost}>Tutup</button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

const card: React.CSSProperties = {
  background: 'rgba(15,23,42,0.6)', borderRadius: 16, padding: 20,
  border: '1px solid rgba(148,163,184,0.1)', cursor: 'pointer', transition: 'all 0.2s',
};
const btnPrimary: React.CSSProperties = {
  padding: '10px 20px', borderRadius: 10, border: 'none', cursor: 'pointer',
  background: 'linear-gradient(135deg, #3B82F6, #2563EB)', color: '#fff', fontWeight: 600, fontSize: 14,
};
const btnGhost: React.CSSProperties = {
  padding: '10px 20px', borderRadius: 10, border: '1px solid rgba(148,163,184,0.2)',
  cursor: 'pointer', background: 'transparent', color: '#94A3B8', fontWeight: 600, fontSize: 14,
};
const overlay: React.CSSProperties = {
  position: 'fixed', inset: 0, background: 'rgba(0,0,0,0.6)',
  display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 50,
};
const modal: React.CSSProperties = {
  background: '#1E293B', borderRadius: 16, padding: 32, width: '100%', maxWidth: 560,
  border: '1px solid rgba(148,163,184,0.1)', boxShadow: '0 25px 50px rgba(0,0,0,0.5)',
  maxHeight: '90vh', overflowY: 'auto',
};
const emptyState: React.CSSProperties = {
  display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center',
  gap: 16, padding: 64, background: 'rgba(15,23,42,0.4)', borderRadius: 16,
  border: '1px dashed rgba(148,163,184,0.2)',
};
