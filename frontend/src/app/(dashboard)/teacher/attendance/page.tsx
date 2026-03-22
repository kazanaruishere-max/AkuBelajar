'use client';

import React, { useState, useEffect, useCallback } from 'react';
import { api } from '@/lib/api/client';

interface Student { id: string; email: string; name: string; }
interface AttendanceItem { id: string; student_id: string; student_email: string; status: string; note: string; }

export default function TeacherAttendancePage() {
  const [students, setStudents] = useState<Student[]>([]);
  const [records, setRecords] = useState<Record<string, { status: string; note: string }>>({});
  const [classId, setClassId] = useState('');
  const [subjectId, setSubjectId] = useState('');
  const [ayId, setAyId] = useState('');
  const [date, setDate] = useState(new Date().toISOString().split('T')[0]);
  const [saving, setSaving] = useState(false);
  const [loaded, setLoaded] = useState(false);

  const loadAttendance = async () => {
    if (!classId || !subjectId || !date) return;
    try {
      // Load students in class
      const sRes = await api.get(`/api/v1/academic/classes/${classId}/students`);
      const sList: Student[] = (sRes.data as any).data || [];
      setStudents(sList);

      // Load existing attendance
      const aRes = await api.get(`/api/v1/attendance/teacher?class_id=${classId}&subject_id=${subjectId}&date=${date}`);
      const aList: AttendanceItem[] = (aRes.data as any).data || [];

      const map: Record<string, { status: string; note: string }> = {};
      sList.forEach((s) => { map[s.id] = { status: 'present', note: '' }; });
      aList.forEach((a) => { map[a.student_id] = { status: a.status, note: a.note }; });
      setRecords(map);
      setLoaded(true);
    } catch { /* ignore */ }
  };

  const handleSave = async () => {
    setSaving(true);
    try {
      const recs = Object.entries(records).map(([studentId, rec]) => ({
        student_id: studentId, status: rec.status, note: rec.note,
      }));
      await api.post('/api/v1/attendance/teacher', {
        class_id: classId, subject_id: subjectId, academic_year_id: ayId, date, records: recs,
      });
      alert('Presensi berhasil disimpan!');
    } catch { /* ignore */ } finally { setSaving(false); }
  };

  const setStatus = (studentId: string, status: string) => {
    setRecords((prev) => ({ ...prev, [studentId]: { ...prev[studentId], status } }));
  };

  const statusColors: Record<string, { bg: string; fg: string; label: string }> = {
    present: { bg: 'rgba(34,197,94,0.15)', fg: '#22C55E', label: 'Hadir' },
    absent: { bg: 'rgba(239,68,68,0.15)', fg: '#EF4444', label: 'Absen' },
    late: { bg: 'rgba(251,191,36,0.15)', fg: '#FBBF24', label: 'Telat' },
    excused: { bg: 'rgba(59,130,246,0.15)', fg: '#60A5FA', label: 'Izin' },
  };

  return (
    <div style={{ padding: 32 }}>
      <h1 style={{ fontSize: 24, fontWeight: 700, color: '#E2E8F0', marginBottom: 24 }}>Presensi</h1>

      {/* Filters */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(200px, 1fr))', gap: 12, marginBottom: 24 }}>
        <label style={labelStyle}>Class ID<input style={inputStyle} value={classId} onChange={(e) => setClassId(e.target.value)} /></label>
        <label style={labelStyle}>Subject ID<input style={inputStyle} value={subjectId} onChange={(e) => setSubjectId(e.target.value)} /></label>
        <label style={labelStyle}>Academic Year ID<input style={inputStyle} value={ayId} onChange={(e) => setAyId(e.target.value)} /></label>
        <label style={labelStyle}>Tanggal<input type="date" style={inputStyle} value={date} onChange={(e) => setDate(e.target.value)} /></label>
      </div>
      <button onClick={loadAttendance} style={{ ...btnPrimary, marginBottom: 24 }}>Muat Data</button>

      {loaded && students.length > 0 && (
        <>
          <div style={{ overflowX: 'auto' }}>
            <table style={tableStyle}>
              <thead>
                <tr>
                  <th style={th}>No</th>
                  <th style={th}>Email</th>
                  <th style={th}>Status</th>
                  <th style={th}>Keterangan</th>
                </tr>
              </thead>
              <tbody>
                {students.map((s, i) => {
                  const rec = records[s.id] || { status: 'present', note: '' };
                  return (
                    <tr key={s.id} style={{ borderBottom: '1px solid rgba(148,163,184,0.1)' }}>
                      <td style={td}>{i + 1}</td>
                      <td style={{ ...td, fontWeight: 500 }}>{s.email}</td>
                      <td style={td}>
                        <div style={{ display: 'flex', gap: 6 }}>
                          {(['present', 'absent', 'late', 'excused'] as const).map((st) => {
                            const c = statusColors[st];
                            const isActive = rec.status === st;
                            return (
                              <button key={st} onClick={() => setStatus(s.id, st)}
                                style={{ padding: '4px 10px', borderRadius: 8, border: isActive ? `2px solid ${c.fg}` : '1px solid transparent',
                                  background: isActive ? c.bg : 'rgba(148,163,184,0.05)', color: isActive ? c.fg : '#64748B',
                                  cursor: 'pointer', fontSize: 12, fontWeight: 600, transition: 'all 0.15s' }}>
                                {c.label}
                              </button>
                            );
                          })}
                        </div>
                      </td>
                      <td style={td}>
                        <input style={{ ...inputStyle, padding: '6px 10px', fontSize: 13 }} value={rec.note}
                          onChange={(e) => setRecords((prev) => ({ ...prev, [s.id]: { ...prev[s.id], note: e.target.value } }))}
                          placeholder="Keterangan..." />
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
          <div style={{ display: 'flex', justifyContent: 'flex-end', marginTop: 20 }}>
            <button onClick={handleSave} disabled={saving} style={btnPrimary}>
              {saving ? 'Menyimpan...' : '💾 Simpan Presensi'}
            </button>
          </div>
        </>
      )}
    </div>
  );
}

const btnPrimary: React.CSSProperties = { padding: '10px 20px', borderRadius: 10, border: 'none', cursor: 'pointer', background: 'linear-gradient(135deg,#3B82F6,#2563EB)', color: '#fff', fontWeight: 600, fontSize: 14 };
const tableStyle: React.CSSProperties = { width: '100%', borderCollapse: 'collapse', background: 'rgba(15,23,42,0.6)', borderRadius: 12 };
const th: React.CSSProperties = { padding: '14px 16px', textAlign: 'left', fontSize: 12, fontWeight: 600, color: '#64748B', textTransform: 'uppercase', letterSpacing: '0.05em', borderBottom: '1px solid rgba(148,163,184,0.1)' };
const td: React.CSSProperties = { padding: '14px 16px', fontSize: 14, color: '#CBD5E1' };
const labelStyle: React.CSSProperties = { display: 'flex', flexDirection: 'column', gap: 6, color: '#CBD5E1', fontSize: 14, fontWeight: 500 };
const inputStyle: React.CSSProperties = { padding: '10px 14px', borderRadius: 10, border: '1px solid rgba(148,163,184,0.2)', background: 'rgba(15,23,42,0.8)', color: '#E2E8F0', fontSize: 14, outline: 'none' };
