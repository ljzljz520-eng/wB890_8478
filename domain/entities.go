package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const (
	StatusDraft       = "draft"
	StatusSubmitted   = "submitted"
	StatusApproved    = "approved"
	StatusRejected    = "rejected"
	StatusArchived    = "archived"
	VisibilityPrivate = "private"
	VisibilityClass   = "class"
	VisibilityPublic  = "public"
)

var (
	ErrMissingRecord     = errors.New("record not found")
	ErrInvalidTransition = errors.New("invalid status transition")
	ErrAlreadyExists     = errors.New("entity already exists")
	ErrInvalidInput      = errors.New("invalid input")
)

type Record struct {
	ID             string   `json:"id"`
	BatchID        string   `json:"batch_id"`
	StudentName    string   `json:"student_name"`
	GraduationYear int      `json:"graduation_year"`
	Message        string   `json:"message"`
	Status         string   `json:"status"`
	Visibility     string   `json:"visibility"`
	Editor         string   `json:"editor"`
	Version        int      `json:"version"`
	CreatedAt      string   `json:"created_at"`
	UpdatedAt      string   `json:"updated_at"`
	Tags           []string `json:"tags"`
}

type AuditEvent struct {
	ID         string `json:"id"`
	RecordID   string `json:"record_id"`
	BatchID    string `json:"batch_id"`
	Action     string `json:"action"`
	Actor      string `json:"actor"`
	FromStatus string `json:"from_status"`
	ToStatus   string `json:"to_status"`
	Detail     string `json:"detail"`
	Sequence   int    `json:"sequence"`
	CreatedAt  string `json:"created_at"`
}

type Workflow struct {
	ID          string   `json:"id"`
	BatchID     string   `json:"batch_id"`
	Name        string   `json:"name"`
	Owner       string   `json:"owner"`
	State       string   `json:"state"`
	CurrentStep int      `json:"current_step"`
	Steps       []string `json:"steps"`
	StartedAt   string   `json:"started_at"`
	CompletedAt string   `json:"completed_at"`
}

type Attachment struct {
	ID         string `json:"id"`
	RecordID   string `json:"record_id"`
	BatchID    string `json:"batch_id"`
	Name       string `json:"name"`
	MediaType  string `json:"media_type"`
	Content    []byte `json:"content"`
	Checksum   string `json:"checksum"`
	Approved   bool   `json:"approved"`
	UploadedBy string `json:"uploaded_by"`
}

type Batch struct {
	ID             string `json:"id"`
	Label          string `json:"label"`
	GraduationYear int    `json:"graduation_year"`
	Coordinator    string `json:"coordinator"`
	Status         string `json:"status"`
	RecordCount    int    `json:"record_count"`
	ApprovedCount  int    `json:"approved_count"`
	ArchivedAt     string `json:"archived_at"`
	CreatedAt      string `json:"created_at"`
}

type Review struct {
	ID        string `json:"id"`
	RecordID  string `json:"record_id"`
	Reviewer  string `json:"reviewer"`
	Decision  string `json:"decision"`
	Notes     string `json:"notes"`
	Score     int    `json:"score"`
	CreatedAt string `json:"created_at"`
}

type ArchiveEntry struct {
	ID         string `json:"id"`
	BatchID    string `json:"batch_id"`
	RecordID   string `json:"record_id"`
	Snapshot   string `json:"snapshot"`
	Confirmed  bool   `json:"confirmed"`
	Viewer     string `json:"viewer"`
	ArchivedAt string `json:"archived_at"`
}

type SearchQuery struct {
	BatchID     string
	Status      string
	Visibility  string
	StudentName string
	Tag         string
}

type SearchResult struct {
	Records []*Record
	Total   int
	Applied SearchQuery
}

func (r *Record) Validate() error {
	if r == nil {
		return ErrInvalidInput
	}
	if strings.TrimSpace(r.ID) == "" || strings.TrimSpace(r.BatchID) == "" {
		return fmt.Errorf("%w: id and batch are required", ErrInvalidInput)
	}
	if strings.TrimSpace(r.StudentName) == "" {
		return fmt.Errorf("%w: student name is required", ErrInvalidInput)
	}
	if r.GraduationYear < 1900 || r.GraduationYear > 2200 {
		return fmt.Errorf("%w: graduation year", ErrInvalidInput)
	}
	if len(strings.TrimSpace(r.Message)) < 3 {
		return fmt.Errorf("%w: message too short", ErrInvalidInput)
	}
	if r.Status == "" {
		r.Status = StatusDraft
	}
	if r.Visibility == "" {
		r.Visibility = VisibilityPrivate
	}
	if !validStatus(r.Status) || !validVisibility(r.Visibility) {
		return fmt.Errorf("%w: unsupported state", ErrInvalidInput)
	}
	return nil
}

func validStatus(value string) bool {
	switch value {
	case StatusDraft, StatusSubmitted, StatusApproved, StatusRejected, StatusArchived:
		return true
	default:
		return false
	}
}

func validVisibility(value string) bool {
	switch value {
	case VisibilityPrivate, VisibilityClass, VisibilityPublic:
		return true
	default:
		return false
	}
}

func CanTransition(from, to string) bool {
	if from == to {
		return true
	}
	switch from {
	case StatusDraft:
		return to == StatusSubmitted
	case StatusSubmitted:
		return to == StatusApproved || to == StatusRejected
	case StatusRejected:
		return to == StatusDraft || to == StatusSubmitted
	case StatusApproved:
		return to == StatusArchived
	default:
		return false
	}
}

func (r *Record) Transition(to string) error {
	if r == nil {
		return ErrMissingRecord
	}
	if !CanTransition(r.Status, to) {
		return fmt.Errorf("%w: %s -> %s", ErrInvalidTransition, r.Status, to)
	}
	r.Status = to
	r.Version++
	return nil
}

func (b *Batch) Validate() error {
	if b == nil || strings.TrimSpace(b.ID) == "" || strings.TrimSpace(b.Label) == "" {
		return ErrInvalidInput
	}
	if b.GraduationYear < 1900 || b.GraduationYear > 2200 {
		return ErrInvalidInput
	}
	if b.Status == "" {
		b.Status = StatusDraft
	}
	if b.Status != StatusDraft && b.Status != StatusArchived && b.Status != StatusApproved {
		return ErrInvalidInput
	}
	return nil
}

func (w *Workflow) Validate() error {
	if w == nil || strings.TrimSpace(w.ID) == "" || strings.TrimSpace(w.BatchID) == "" || len(w.Steps) < 4 {
		return ErrInvalidInput
	}
	if w.State == "" {
		w.State = "pending"
	}
	if w.CurrentStep < 0 || w.CurrentStep > len(w.Steps) {
		return ErrInvalidInput
	}
	return nil
}

func (a *Attachment) Validate() error {
	if a == nil || a.ID == "" || a.RecordID == "" || a.BatchID == "" || a.Name == "" {
		return ErrInvalidInput
	}
	if len(a.Content) == 0 {
		return ErrInvalidInput
	}
	return nil
}

func (r *Review) Validate() error {
	if r == nil || r.ID == "" || r.RecordID == "" || r.Reviewer == "" {
		return ErrInvalidInput
	}
	if r.Decision != "approve" && r.Decision != "reject" {
		return ErrInvalidInput
	}
	if r.Score < 0 || r.Score > 100 {
		return ErrInvalidInput
	}
	return nil
}

func (e *ArchiveEntry) Validate() error {
	if e == nil || e.ID == "" || e.BatchID == "" || e.RecordID == "" || e.Snapshot == "" {
		return ErrInvalidInput
	}
	return nil
}

// Restore reconstructs the business record that was frozen at archive time.
// When the live record has since been deleted, readback must surface the real
// archived result rather than fabricating a draft state, so repeated
// archive/readback cycles never drift away from the confirmed business
// record.
func (e *ArchiveEntry) Restore() (*Record, error) {
	if e == nil || e.Snapshot == "" {
		return nil, ErrInvalidInput
	}
	var record Record
	if err := json.Unmarshal([]byte(e.Snapshot), &record); err != nil {
		return nil, fmt.Errorf("%w: decode archive snapshot: %v", ErrInvalidInput, err)
	}
	return &record, nil
}
