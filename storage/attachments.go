package storage

import (
	"errors"
	"go.etcd.io/bbolt"
	"memorialstation/domain"
)

func (s *Store) SaveAttachment(attachment *domain.Attachment) error {
	if err := attachment.Validate(); err != nil {
		return err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("storage is closed")
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return putJSON(tx, "attachments", attachment.ID, attachment) })
}
func (s *Store) GetAttachment(id string) (*domain.Attachment, error) {
	if id == "" {
		return nil, domain.ErrInvalidInput
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return nil, errors.New("storage is closed")
	}
	var item domain.Attachment
	err := s.db.View(func(tx *bbolt.Tx) error { return getJSON(tx, "attachments", id, &item) })
	if err != nil {
		return nil, err
	}
	return &item, nil
}
func (s *Store) ListAttachments(recordID string) ([]*domain.Attachment, error) {
	items, err := listJSON(s, "attachments", func(data []byte) (*domain.Attachment, error) {
		var item domain.Attachment
		if err := decode(data, &item); err != nil {
			return nil, err
		}
		return &item, nil
	})
	if err != nil {
		return nil, err
	}
	filtered := make([]*domain.Attachment, 0, len(items))
	for _, item := range items {
		if recordID == "" || item.RecordID == recordID {
			filtered = append(filtered, item)
		}
	}
	return filtered, nil
}
