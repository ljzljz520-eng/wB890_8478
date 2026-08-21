package review

import (
	"fmt"
	"memorialstation/domain"
	"memorialstation/storage"
)

type Service struct{ store *storage.Store }

func New(store *storage.Store) *Service { return &Service{store: store} }

func (s *Service) Inspect(recordID string) (*domain.Record, []*domain.Attachment, []*domain.Review, error) {
	record, err := s.store.GetRecord(recordID)
	if err != nil {
		return nil, nil, nil, err
	}
	attachments, err := s.store.ListAttachments(recordID)
	if err != nil {
		return nil, nil, nil, err
	}
	reviews, err := s.store.ListReviews(recordID)
	if err != nil {
		return nil, nil, nil, err
	}
	return record, attachments, reviews, nil
}

func (s *Service) Approve(recordID, reviewer, notes, now string) (*domain.Review, error) {
	record, err := s.store.GetRecord(recordID)
	if err != nil {
		return nil, err
	}
	if record.Status != domain.StatusSubmitted {
		return nil, fmt.Errorf("%w: review requires submitted", domain.ErrInvalidTransition)
	}
	reviews, err := s.store.ListReviews(recordID)
	if err != nil {
		return nil, err
	}
	review := &domain.Review{ID: fmt.Sprintf("review-%s-%d", recordID, len(reviews)+1), RecordID: recordID, Reviewer: reviewer, Decision: "approve", Notes: notes, Score: score(record), CreatedAt: now}
	if err := s.store.SaveReview(review); err != nil {
		return nil, err
	}
	return review, nil
}

func score(record *domain.Record) int {
	score := len(record.Message) * 3
	if score > 100 {
		score = 100
	}
	if record.Visibility == domain.VisibilityPublic {
		score += 5
	}
	if score > 100 {
		score = 100
	}
	return score
}

func (s *Service) Pending(batchID string) ([]*domain.Record, error) {
	records, err := s.store.ListRecords()
	if err != nil {
		return nil, err
	}
	pending := make([]*domain.Record, 0)
	for _, record := range records {
		if batchID != "" && record.BatchID != batchID {
			continue
		}
		if record.Status == domain.StatusSubmitted {
			pending = append(pending, record)
		}
	}
	return pending, nil
}

func (s *Service) Decision(recordID string) (string, error) {
	reviews, err := s.store.ListReviews(recordID)
	if err != nil {
		return "", err
	}
	if len(reviews) == 0 {
		return "pending", nil
	}
	latest := reviews[len(reviews)-1]
	if latest.Decision == "approve" {
		return "approved", nil
	}
	return "rejected", nil
}

func (s *Service) ValidateAttachments(recordID string) error {
	attachments, err := s.store.ListAttachments(recordID)
	if err != nil {
		return err
	}
	if len(attachments) == 0 {
		return nil
	}
	for _, attachment := range attachments {
		if attachment.Checksum == "" || len(attachment.Content) == 0 {
			return domain.ErrInvalidInput
		}
		if attachment.Approved {
			continue
		}
		attachment.Approved = true
		if err := s.store.SaveAttachment(attachment); err != nil {
			return err
		}
	}
	return nil
}
