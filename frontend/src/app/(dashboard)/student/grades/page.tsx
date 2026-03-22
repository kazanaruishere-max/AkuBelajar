'use client';

import React, { useState, useEffect, useCallback } from 'react';
import { api } from '@/lib/api/client';

interface GradeItem {
  id: string;
  subject_name: string;
  category: string;
  label: string;
  score: number;
  weight_pct: number;
  weighted_score: number;
  created_at: string;
}

export default function StudentGradesPage() {
  const [grades, setGrades] = useState<GradeItem[]>([]);
  const [loading, setLoading] = useState(true);
  const classId = typeof window !== 'undefined' ? new URLSearchParams(window.location.search).get('class_id') || '' : '';
  const subjectId = typeof window !== 'undefined' ? new URLSearchParams(window.location.search).get('subject_id') || '' : '';

  const fetchGrades = useCallback(async () => {
    if (!classId || !subjectId) { setLoading(false); return; }
    try {
      setLoading(true);
      const res = await api.get(`/api/v1/grades/student?class_id=${classId}&subject_id=${subjectId}`);
      setGrades((res.data as any).data || []);
    } catch { /* ignore */ } finally { setLoading(false); }
  }, [classId, subjectId]);

  useEffect(() => { fetchGrades(); }, [fetchGrades]);

  const catIcon: Record<string, string> = { assignment: '📝', quiz: '❓', exam: '📄', midterm: '📋', final: '🏆' };
  const avg = grades.length > 0 ? Math.round(grades.reduce((s, g) => s + g.score, 0) / grades.length) : 0;
  const letterGrade = avg >= 90 ? 'A' : avg >= 80 ? 'B' : avg >= 70 ? 'C' : avg >= 60 ? 'D' : 'E';
  const gradeColor = avg >= 80 ? '#22C55E' : avg >= 60 ? '#F59E0B' : '#EF4444';

  return (
    <div style={{ padding: 32 }}>
      <h1 style={{ fontSize: 24, fontWeight: 700, color: '#E2E8F0', marginBottom: 24 }}>Nilai Saya</h1>

      {(!classId || !subjectId) ? (
        <div style={emptyState}><p style={{ color: '#94A3B8' }}>Tambah ?class_id=...&subject_id=... di URL</p></div>
      ) : loading ? <p style={{ color: '#94A3B8' }}>Memuat...</p> : (
        <>
          {/* Summary Card */}
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: 16, marginBottom: 32 }}>
            <div style={summaryCard}>
              <p style={{ fontSize: 13, color: '#64748B', fontWeight: 500 }}>Rata-rata</p>
              <p style={{ fontSize: 36, fontWeight: 800, color: gradeColor, marginTop: 4 }}>{avg}</p>
            </div>
            <div style={summaryCard}>
              <p style={{ fontSize: 13, color: '#64748B', fontWeight: 500 }}>Grade</p>
              <p style={{ fontSize: 36, fontWeight: 800, color: gradeColor, marginTop: 4 }}>{letterGrade}</p>
            </div>
            <div style={summaryCard}>
              <p style={{ fontSize: 13, color: '#64748B', fontWeight: 500 }}>Total Penilaian</p>
              <p style={{ fontSize: 36, fontWeight: 800, color: '#60A5FA', marginTop: 4 }}>{grades.length}</p>
            </div>
          </div>

          {/* Grades Table */}
          {grades.length === 0 ? (
            <div style={emptyState}><p style={{ color: '#94A3B8' }}>Belum ada nilai</p></div>
          ) : (
            <div style={{ overflowX: 'auto' }}>
              <table style={tableStyle}>
                <thead>
                  <tr>
                    <th style={th}>Kategori</th>
                    <th style={th}>Keterangan</th>
                    <th style={th}>Nilai</th>
                    <th style={th}>Bobot</th>
                    <th style={th}>Nilai Terbobot</th>
                  </tr>
                </thead>
                <tbody>
                  {grades.map((g) => (
                    <tr key={g.id} style={{ borderBottom: '1px solid rgba(148,163,184,0.1)' }}>
                      <td style={td}>
                        <span style={{ padding: '4px 10px', borderRadius: 8, fontSize: 12, fontWeight: 600,
                          background: 'rgba(99,102,241,0.15)', color: '#818CF8' }}>
                          {catIcon[g.category] || '📌'} {g.category}
                        </span>
                      </td>
                      <td style={{ ...td, fontWeight: 500 }}>{g.label}</td>
                      <td style={td}>
                        <span style={{ fontWeight: 700, fontSize: 16,
                          color: g.score >= 80 ? '#22C55E' : g.score >= 60 ? '#F59E0B' : '#EF4444' }}>{g.score}</span>
                      </td>
                      <td style={td}>{g.weight_pct}%</td>
                      <td style={td}>{g.weighted_score.toFixed(1)}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </>
      )}
    </div>
  );
}

const summaryCard: React.CSSProperties = { background: 'rgba(15,23,42,0.6)', borderRadius: 16, padding: 24, textAlign: 'center', border: '1px solid rgba(148,163,184,0.1)' };
const tableStyle: React.CSSProperties = { width: '100%', borderCollapse: 'collapse', background: 'rgba(15,23,42,0.6)', borderRadius: 12 };
const th: React.CSSProperties = { padding: '14px 16px', textAlign: 'left', fontSize: 12, fontWeight: 600, color: '#64748B', textTransform: 'uppercase', letterSpacing: '0.05em', borderBottom: '1px solid rgba(148,163,184,0.1)' };
const td: React.CSSProperties = { padding: '14px 16px', fontSize: 14, color: '#CBD5E1' };
const emptyState: React.CSSProperties = { display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', gap: 16, padding: 64, background: 'rgba(15,23,42,0.4)', borderRadius: 16, border: '1px dashed rgba(148,163,184,0.2)' };
