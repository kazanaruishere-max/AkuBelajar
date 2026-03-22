package attendance

import "time"

type MarkAttendanceRequest struct {
	ClassID        string `json:"class_id" validate:"required,uuid"`
	SubjectID      string `json:"subject_id" validate:"required,uuid"`
	AcademicYearID string `json:"academic_year_id" validate:"required,uuid"`
	Date           string `json:"date" validate:"required"`
	Records        []AttendanceRecord `json:"records" validate:"required,min=1,dive"`
}

type AttendanceRecord struct {
	StudentID string `json:"student_id" validate:"required,uuid"`
	Status    string `json:"status" validate:"required,oneof=present absent late excused"`
	Note      string `json:"note"`
}

type AttendanceResponse struct {
	ID        string    `json:"id"`
	StudentID string    `json:"student_id"`
	Email     string    `json:"student_email,omitempty"`
	ClassID   string    `json:"class_id"`
	SubjectID string    `json:"subject_id"`
	Date      string    `json:"date"`
	Status    string    `json:"status"`
	Note      string    `json:"note,omitempty"`
	MarkedBy  string    `json:"marked_by"`
	CreatedAt time.Time `json:"created_at"`
}

type AttendanceSummary struct {
	StudentID string `json:"student_id"`
	Email     string `json:"student_email"`
	Present   int    `json:"present"`
	Absent    int    `json:"absent"`
	Late      int    `json:"late"`
	Excused   int    `json:"excused"`
	Total     int    `json:"total"`
}
