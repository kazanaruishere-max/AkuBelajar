'use client';

import React, { useState, useEffect, useCallback, useRef } from 'react';
import { api } from '@/lib/api/client';

interface Quiz { id: string; title: string; subject_name: string; time_limit: number; question_count: number; }
interface Question { id: string; question_text: string; question_type: string; options: any; order_num: number; }
interface Session { id: string; expires_at: string; status: string; score: number | null; }

export default function StudentQuizzesPage() {
  const [quizzes, setQuizzes] = useState<Quiz[]>([]);
  const [loading, setLoading] = useState(true);
  const [activeSession, setActiveSession] = useState<Session | null>(null);
  const [questions, setQuestions] = useState<Question[]>([]);
  const [currentQ, setCurrentQ] = useState(0);
  const [answers, setAnswers] = useState<Record<string, string>>({});
  const [timeLeft, setTimeLeft] = useState(0);
  const timerRef = useRef<NodeJS.Timeout | null>(null);

  const classId = typeof window !== 'undefined' ? new URLSearchParams(window.location.search).get('class_id') || '' : '';

  const fetchQuizzes = useCallback(async () => {
    if (!classId) { setLoading(false); return; }
    try {
      setLoading(true);
      const res = await api.get(`/api/v1/quizzes/student?class_id=${classId}`);
      setQuizzes((res.data as any).data || []);
    } catch { /* ignore */ } finally { setLoading(false); }
  }, [classId]);

  useEffect(() => { fetchQuizzes(); }, [fetchQuizzes]);

  const startQuiz = async (quizId: string) => {
    try {
      const res = await api.post(`/api/v1/quizzes/student/${quizId}/start`);
      const session = (res.data as any).data;
      setActiveSession(session);

      // Fetch questions
      // Note: In production, questions would come from the session endpoint
      // For CBT, we use the quiz questions endpoint
      const qRes = await api.get(`/api/v1/quizzes/teacher/${quizId}/questions`);
      setQuestions((qRes.data as any).data || []);
      setCurrentQ(0);
      setAnswers({});

      // Start timer
      const expires = new Date(session.expires_at).getTime();
      const updateTimer = () => {
        const remaining = Math.max(0, Math.floor((expires - Date.now()) / 1000));
        setTimeLeft(remaining);
        if (remaining <= 0) { handleSubmit(session.id); }
      };
      updateTimer();
      timerRef.current = setInterval(updateTimer, 1000);
    } catch { /* ignore */ }
  };

  const saveAnswer = async (questionId: string, key: string) => {
    if (!activeSession) return;
    setAnswers((prev) => ({ ...prev, [questionId]: key }));
    try {
      await api.post(`/api/v1/quizzes/student/sessions/${activeSession.id}/answer`, {
        question_id: questionId, selected_key: key,
      });
    } catch { /* ignore */ }
  };

  const handleSubmit = async (sessionId?: string) => {
    const sid = sessionId || activeSession?.id;
    if (!sid) return;
    if (timerRef.current) clearInterval(timerRef.current);
    try {
      const res = await api.post(`/api/v1/quizzes/student/sessions/${sid}/submit`);
      setActiveSession((res.data as any).data);
      setQuestions([]);
    } catch { /* ignore */ }
  };

  const formatTime = (s: number) => `${Math.floor(s / 60).toString().padStart(2, '0')}:${(s % 60).toString().padStart(2, '0')}`;

  // ── CBT Interface ───────────────────────────────────────
  if (activeSession && activeSession.status !== 'submitted' && questions.length > 0) {
    const q = questions[currentQ];
    let options: { key: string; text: string }[] = [];
    try { options = typeof q.options === 'string' ? JSON.parse(q.options) : q.options || []; } catch { options = []; }

    return (
      <div style={{ padding: 0, height: '100vh', display: 'flex', flexDirection: 'column', background: '#0F172A' }}>
        {/* Timer Bar */}
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', padding: '12px 24px',
          background: timeLeft < 60 ? 'rgba(239,68,68,0.2)' : 'rgba(15,23,42,0.8)', borderBottom: '1px solid rgba(148,163,184,0.1)' }}>
          <span style={{ color: '#E2E8F0', fontWeight: 600 }}>Soal {currentQ + 1}/{questions.length}</span>
          <span style={{ fontSize: 24, fontWeight: 700, fontFamily: 'monospace',
            color: timeLeft < 60 ? '#EF4444' : timeLeft < 300 ? '#F59E0B' : '#22C55E' }}>
            ⏱️ {formatTime(timeLeft)}
          </span>
          <button onClick={() => { if (confirm('Submit kuis?')) handleSubmit(); }}
            style={{ padding: '8px 20px', borderRadius: 8, border: 'none', cursor: 'pointer',
              background: '#EF4444', color: '#fff', fontWeight: 600, fontSize: 14 }}>Submit</button>
        </div>

        {/* Question */}
        <div style={{ flex: 1, display: 'flex', flexDirection: 'column', padding: 32, overflowY: 'auto' }}>
          <div style={{ background: 'rgba(30,41,59,0.8)', borderRadius: 16, padding: 24,
            border: '1px solid rgba(148,163,184,0.1)', marginBottom: 24 }}>
            <p style={{ fontSize: 18, color: '#E2E8F0', lineHeight: 1.6 }}>{q.question_text}</p>
          </div>

          <div style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
            {options.map((opt) => (
              <button key={opt.key} onClick={() => saveAnswer(q.id, opt.key)}
                style={{ display: 'flex', alignItems: 'center', gap: 16, padding: '16px 20px', borderRadius: 12,
                  border: answers[q.id] === opt.key ? '2px solid #3B82F6' : '1px solid rgba(148,163,184,0.15)',
                  background: answers[q.id] === opt.key ? 'rgba(59,130,246,0.1)' : 'rgba(15,23,42,0.4)',
                  cursor: 'pointer', textAlign: 'left', transition: 'all 0.2s' }}>
                <span style={{ width: 36, height: 36, borderRadius: '50%', display: 'flex', alignItems: 'center', justifyContent: 'center',
                  background: answers[q.id] === opt.key ? '#3B82F6' : 'rgba(148,163,184,0.15)',
                  color: answers[q.id] === opt.key ? '#fff' : '#94A3B8', fontWeight: 700, fontSize: 14 }}>{opt.key}</span>
                <span style={{ color: '#CBD5E1', fontSize: 15 }}>{opt.text}</span>
              </button>
            ))}
          </div>
        </div>

        {/* Navigation */}
        <div style={{ display: 'flex', justifyContent: 'space-between', padding: '16px 24px',
          borderTop: '1px solid rgba(148,163,184,0.1)', background: 'rgba(15,23,42,0.8)' }}>
          <button onClick={() => setCurrentQ(Math.max(0, currentQ - 1))} disabled={currentQ === 0}
            style={{ ...navBtn, opacity: currentQ === 0 ? 0.3 : 1 }}>← Sebelumnya</button>
          <div style={{ display: 'flex', gap: 6, flexWrap: 'wrap', justifyContent: 'center' }}>
            {questions.map((_, i) => (
              <button key={i} onClick={() => setCurrentQ(i)}
                style={{ width: 32, height: 32, borderRadius: 8, border: 'none', cursor: 'pointer', fontSize: 12, fontWeight: 600,
                  background: i === currentQ ? '#3B82F6' : answers[questions[i].id] ? 'rgba(34,197,94,0.3)' : 'rgba(148,163,184,0.1)',
                  color: i === currentQ ? '#fff' : answers[questions[i].id] ? '#22C55E' : '#64748B' }}>{i + 1}</button>
            ))}
          </div>
          <button onClick={() => setCurrentQ(Math.min(questions.length - 1, currentQ + 1))} disabled={currentQ === questions.length - 1}
            style={{ ...navBtn, opacity: currentQ === questions.length - 1 ? 0.3 : 1 }}>Selanjutnya →</button>
        </div>
      </div>
    );
  }

  // ── Result View ─────────────────────────────────────────
  if (activeSession && activeSession.status === 'submitted') {
    return (
      <div style={{ padding: 32, display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', minHeight: '60vh' }}>
        <div style={{ background: 'rgba(15,23,42,0.6)', borderRadius: 24, padding: 48, textAlign: 'center',
          border: '1px solid rgba(148,163,184,0.1)', maxWidth: 400 }}>
          <div style={{ fontSize: 64, marginBottom: 16 }}>🎉</div>
          <h2 style={{ fontSize: 24, fontWeight: 700, color: '#E2E8F0' }}>Kuis Selesai!</h2>
          <p style={{ fontSize: 48, fontWeight: 800, marginTop: 16,
            background: 'linear-gradient(135deg,#3B82F6,#22C55E)', WebkitBackgroundClip: 'text', WebkitTextFillColor: 'transparent' }}>
            {activeSession.score ?? 0}/100
          </p>
          <button onClick={() => { setActiveSession(null); fetchQuizzes(); }}
            style={{ ...btnPrimary, marginTop: 24 }}>Kembali ke Daftar</button>
        </div>
      </div>
    );
  }

  // ── Quiz List ───────────────────────────────────────────
  return (
    <div style={{ padding: 32 }}>
      <h1 style={{ fontSize: 24, fontWeight: 700, color: '#E2E8F0', marginBottom: 24 }}>Kuis</h1>
      {!classId ? (
        <div style={emptyState}><p style={{ color: '#94A3B8' }}>Pilih kelas (?class_id=...)</p></div>
      ) : loading ? <p style={{ color: '#94A3B8' }}>Memuat...</p> : quizzes.length === 0 ? (
        <div style={emptyState}><p style={{ color: '#94A3B8' }}>🎉 Tidak ada kuis saat ini</p></div>
      ) : (
        <div style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
          {quizzes.map((q) => (
            <div key={q.id} style={card}>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <div>
                  <h3 style={{ fontSize: 16, fontWeight: 600, color: '#E2E8F0' }}>{q.title}</h3>
                  <p style={{ fontSize: 13, color: '#94A3B8', marginTop: 4 }}>
                    📚 {q.subject_name} · 📝 {q.question_count} soal · ⏱️ {q.time_limit} menit
                  </p>
                </div>
                <button onClick={() => startQuiz(q.id)} style={btnPrimary}>Mulai Kuis →</button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

const btnPrimary: React.CSSProperties = { padding: '10px 20px', borderRadius: 10, border: 'none', cursor: 'pointer', background: 'linear-gradient(135deg,#3B82F6,#2563EB)', color: '#fff', fontWeight: 600, fontSize: 14 };
const navBtn: React.CSSProperties = { padding: '10px 20px', borderRadius: 10, border: '1px solid rgba(148,163,184,0.2)', cursor: 'pointer', background: 'transparent', color: '#CBD5E1', fontWeight: 500, fontSize: 14 };
const card: React.CSSProperties = { background: 'rgba(15,23,42,0.6)', borderRadius: 16, padding: 20, border: '1px solid rgba(148,163,184,0.1)' };
const emptyState: React.CSSProperties = { display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', gap: 16, padding: 64, background: 'rgba(15,23,42,0.4)', borderRadius: 16, border: '1px dashed rgba(148,163,184,0.2)' };
