package storage

import (
	"errors"
	"go.etcd.io/bbolt"
	"memorialstation/domain"
)

func (s *Store) SaveWorkflow(workflow *domain.Workflow) error {
	if err := workflow.Validate(); err != nil {
		return err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("storage is closed")
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return putJSON(tx, "workflows", workflow.ID, workflow) })
}

func (s *Store) GetWorkflow(id string) (*domain.Workflow, error) {
	if id == "" {
		return nil, domain.ErrInvalidInput
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return nil, errors.New("storage is closed")
	}
	var item domain.Workflow
	err := s.db.View(func(tx *bbolt.Tx) error { return getJSON(tx, "workflows", id, &item) })
	if err != nil {
		return nil, err
	}
	return &item, nil
}
func (s *Store) ListWorkflows(batchID string) ([]*domain.Workflow, error) {
	items, err := listJSON(s, "workflows", func(data []byte) (*domain.Workflow, error) {
		var item domain.Workflow
		if err := decode(data, &item); err != nil {
			return nil, err
		}
		return &item, nil
	})
	if err != nil {
		return nil, err
	}
	filtered := make([]*domain.Workflow, 0, len(items))
	for _, item := range items {
		if batchID == "" || item.BatchID == batchID {
			filtered = append(filtered, item)
		}
	}
	return filtered, nil
}
