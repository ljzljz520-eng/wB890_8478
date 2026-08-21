package storage

import (
	"errors"
	"fmt"
	"go.etcd.io/bbolt"
	"memorialstation/domain"
)

func (s *Store) SaveAudit(event *domain.AuditEvent) error {
	if event == nil || event.ID == "" || event.RecordID == "" || event.BatchID == "" || event.Action == "" {
		return domain.ErrInvalidInput
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("storage is closed")
	}
	return s.db.Update(func(tx *bbolt.Tx) error {
		if event.Sequence < 1 {
			return fmt.Errorf("%w: sequence", domain.ErrInvalidInput)
		}
		return putJSON(tx, "audit_events", event.ID, event)
	})
}

func (s *Store) ListAudit(batchID string) ([]*domain.AuditEvent, error) {
	events, err := listJSON(s, "audit_events", func(data []byte) (*domain.AuditEvent, error) {
		var item domain.AuditEvent
		if err := decode(data, &item); err != nil {
			return nil, err
		}
		return &item, nil
	})
	if err != nil {
		return nil, err
	}
	filtered := make([]*domain.AuditEvent, 0, len(events))
	for _, event := range events {
		if batchID == "" || event.BatchID == batchID {
			filtered = append(filtered, event)
		}
	}
	sortEvents(filtered)
	return filtered, nil
}

func (s *Store) NextAuditSequence(batchID string) (int, error) {
	events, err := s.ListAudit(batchID)
	if err != nil {
		return 0, err
	}
	next := 1
	for _, event := range events {
		if event.Sequence >= next {
			next = event.Sequence + 1
		}
	}
	return next, nil
}
