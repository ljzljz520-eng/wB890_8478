package workflow

import (
	"fmt"
	"memorialstation/domain"
	"memorialstation/storage"
)

type Service struct{ store *storage.Store }

func New(store *storage.Store) *Service { return &Service{store: store} }

func (s *Service) StartBatch(batch *domain.Batch, now string) error {
	if batch == nil {
		return domain.ErrInvalidInput
	}
	if batch.Status == "" {
		batch.Status = domain.StatusDraft
	}
	batch.CreatedAt = now
	return s.store.SaveBatch(batch)
}

func (s *Service) CreateRecord(record *domain.Record, now string) error {
	if record == nil {
		return domain.ErrInvalidInput
	}
	if record.Status == "" {
		record.Status = domain.StatusDraft
	}
	if record.CreatedAt == "" {
		record.CreatedAt = now
	}
	record.UpdatedAt = now
	return s.store.SaveRecord(record)
}

func (s *Service) Submit(recordID, actor, now string) error {
	return s.change(recordID, domain.StatusSubmitted, actor, "submit", now)
}

func (s *Service) change(recordID, target, actor, action, now string) error {
	record, err := s.store.GetRecord(recordID)
	if err != nil {
		return err
	}
	from := record.Status
	if err := record.Transition(target); err != nil {
		return err
	}
	record.Editor = actor
	record.UpdatedAt = now
	if err := s.store.SaveRecord(record); err != nil {
		return err
	}
	sequence, err := s.store.NextAuditSequence(record.BatchID)
	if err != nil {
		return err
	}
	event := &domain.AuditEvent{ID: fmt.Sprintf("%s-%03d", record.BatchID, sequence), RecordID: record.ID, BatchID: record.BatchID, Action: action, Actor: actor, FromStatus: from, ToStatus: target, Sequence: sequence, CreatedAt: now}
	return s.store.SaveAudit(event)
}

func (s *Service) Confirm(recordID, actor, now string) error {
	return s.change(recordID, domain.StatusApproved, actor, "confirm", now)
}
func (s *Service) Reject(recordID, actor, notes, now string) error {
	if err := s.change(recordID, domain.StatusRejected, actor, "reject", now); err != nil {
		return err
	}
	reviews, err := s.store.ListReviews(recordID)
	if err != nil {
		return err
	}
	review := &domain.Review{ID: fmt.Sprintf("review-%s-%d", recordID, len(reviews)+1), RecordID: recordID, Reviewer: actor, Decision: "reject", Notes: notes, Score: 0, CreatedAt: now}
	return s.store.SaveReview(review)
}

func (s *Service) Attach(attachment *domain.Attachment) error {
	return s.store.SaveAttachment(attachment)
}

func (s *Service) StartReviewWorkflow(batchID, owner, now string) (*domain.Workflow, error) {
	workflow := &domain.Workflow{ID: "review-" + batchID, BatchID: batchID, Name: "登记-审核-确认-归档", Owner: owner, State: "pending", Steps: []string{"登记", "审核", "确认", "归档"}, StartedAt: now}
	if err := workflow.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.SaveWorkflow(workflow); err != nil {
		return nil, err
	}
	return workflow, nil
}

func (s *Service) AdvanceWorkflow(id string, now string) (*domain.Workflow, error) {
	item, err := s.store.GetWorkflow(id)
	if err != nil {
		return nil, err
	}
	if item.State == "cancelled" || item.State == "completed" {
		return nil, domain.ErrInvalidTransition
	}
	if item.CurrentStep >= len(item.Steps)-1 {
		item.CurrentStep = len(item.Steps)
		item.State = "completed"
		item.CompletedAt = now
	} else {
		item.CurrentStep++
		item.State = "active"
	}
	if item.CurrentStep == 2 && item.State == "active" {
		item.State = "awaiting-confirmation"
	}
	if err := s.store.SaveWorkflow(item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) CancelWorkflow(id string) error {
	item, err := s.store.GetWorkflow(id)
	if err != nil {
		return err
	}
	if item.State == "completed" {
		return domain.ErrInvalidTransition
	}
	item.State = "cancelled"
	return s.store.SaveWorkflow(item)
}

func (s *Service) Publish(recordID, now string) error {
	record, err := s.store.GetRecord(recordID)
	if err != nil {
		return err
	}
	if record.Status != domain.StatusArchived {
		return fmt.Errorf("%w: publish requires archive", domain.ErrInvalidTransition)
	}
	if record.Visibility == domain.VisibilityPrivate {
		record.Visibility = domain.VisibilityPublic
	}
	record.UpdatedAt = now
	return s.store.SaveRecord(record)
}
