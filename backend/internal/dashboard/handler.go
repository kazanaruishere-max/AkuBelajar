package dashboard

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kazanaruishere-max/akubelajar/backend/pkg/response"
)

type Handler struct {
	db *pgxpool.Pool
}

func NewHandler(db *pgxpool.Pool) *Handler {
	return &Handler{db: db}
}

func uid(c *gin.Context) string { v, _ := c.Get("user_id"); return v.(string) }
func urole(c *gin.Context) string { v, _ := c.Get("user_role"); return v.(string) }

// Stats returns role-based dashboard statistics.
func (h *Handler) Stats(c *gin.Context) {
	ctx := c.Request.Context()
	role := urole(c)
	userID := uid(c)

	switch role {
	case "super_admin":
		h.adminStats(c, ctx)
	case "teacher":
		h.teacherStats(c, ctx, userID)
	case "student", "class_leader":
		h.studentStats(c, ctx, userID)
	default:
		response.OK(c, gin.H{"role": role})
	}
}

func (h *Handler) adminStats(c *gin.Context, ctx context.Context) {
	var totalUsers, totalTeachers, totalStudents, totalClasses, totalSubjects int
	h.db.QueryRow(ctx, `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL`).Scan(&totalUsers)
	h.db.QueryRow(ctx, `SELECT COUNT(*) FROM users WHERE role = 'teacher' AND deleted_at IS NULL`).Scan(&totalTeachers)
	h.db.QueryRow(ctx, `SELECT COUNT(*) FROM users WHERE role IN ('student','class_leader') AND deleted_at IS NULL`).Scan(&totalStudents)
	h.db.QueryRow(ctx, `SELECT COUNT(*) FROM classes WHERE deleted_at IS NULL`).Scan(&totalClasses)
	h.db.QueryRow(ctx, `SELECT COUNT(*) FROM subjects WHERE deleted_at IS NULL`).Scan(&totalSubjects)

	response.OK(c, gin.H{
		"role":           "super_admin",
		"total_users":    totalUsers,
		"total_teachers": totalTeachers,
		"total_students": totalStudents,
		"total_classes":  totalClasses,
		"total_subjects": totalSubjects,
	})
}

func (h *Handler) teacherStats(c *gin.Context, ctx context.Context, teacherID string) {
	var totalAssignments, totalQuizzes, pendingSubmissions int
	h.db.QueryRow(ctx, `SELECT COUNT(*) FROM assignments WHERE teacher_id = $1 AND deleted_at IS NULL`, teacherID).Scan(&totalAssignments)
	h.db.QueryRow(ctx, `SELECT COUNT(*) FROM quizzes WHERE teacher_id = $1 AND deleted_at IS NULL`, teacherID).Scan(&totalQuizzes)
	h.db.QueryRow(ctx, `SELECT COUNT(*) FROM assignment_submissions s JOIN assignments a ON a.id = s.assignment_id WHERE a.teacher_id = $1 AND s.status = 'submitted'`, teacherID).Scan(&pendingSubmissions)

	response.OK(c, gin.H{
		"role":                "teacher",
		"total_assignments":   totalAssignments,
		"total_quizzes":       totalQuizzes,
		"pending_submissions": pendingSubmissions,
	})
}

func (h *Handler) studentStats(c *gin.Context, ctx context.Context, studentID string) {
	var totalAssignments, totalQuizzes, avgScore int
	h.db.QueryRow(ctx, `SELECT COUNT(*) FROM assignment_submissions WHERE student_id = $1`, studentID).Scan(&totalAssignments)
	h.db.QueryRow(ctx, `SELECT COUNT(*) FROM quiz_sessions WHERE student_id = $1`, studentID).Scan(&totalQuizzes)
	h.db.QueryRow(ctx, `SELECT COALESCE(AVG(score),0)::INT FROM grades WHERE student_id = $1`, studentID).Scan(&avgScore)

	var unread int
	h.db.QueryRow(ctx, `SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND is_read = FALSE`, studentID).Scan(&unread)

	response.OK(c, gin.H{
		"role":               "student",
		"total_assignments":  totalAssignments,
		"total_quizzes":      totalQuizzes,
		"average_score":      avgScore,
		"unread_notifications": unread,
	})
}
