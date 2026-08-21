package storage

import (
	"errors"
	"go.etcd.io/bbolt"
	"memorialstation/domain"
)

func (s *Store) SaveArchive(entry *domain.ArchiveEntry) error {
	if err := entry.Validate(); err != nil {
		return err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("storage is closed")
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return putJSON(tx, "archives", entry.ID, entry) })
}
func (s *Store) GetArchive(id string) (*domain.ArchiveEntry, error) {
	if id == "" {
		return nil, domain.ErrInvalidInput
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return nil, errors.New("storage is closed")
	}
	var item domain.ArchiveEntry
	err := s.db.View(func(tx *bbolt.Tx) error { return getJSON(tx, "archives", id, &item) })
	if err != nil {
		return nil, err
	}
	return &item, nil
}
func (s *Store) ListArchives(batchID string) ([]*domain.ArchiveEntry, error) {
	items, err := listJSON(s, "archives", func(data []byte) (*domain.ArchiveEntry, error) {
		var item domain.ArchiveEntry
		if err := decode(data, &item); err != nil {
			return nil, err
		}
		return &item, nil
	})
	if err != nil {
		return nil, err
	}
	filtered := make([]*domain.ArchiveEntry, 0, len(items))
	for _, item := range items {
		if batchID == "" || item.BatchID == batchID {
			filtered = append(filtered, item)
		}
	}
	return filtered, nil
}
