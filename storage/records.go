package storage

import (
	"errors"
	"go.etcd.io/bbolt"
	"memorialstation/domain"
)

func (s *Store) SaveRecord(record *domain.Record) error {
	if err := record.Validate(); err != nil {
		return err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("storage is closed")
	}
	return s.db.Update(func(tx *bbolt.Tx) error {
		existing := tx.Bucket([]byte("records")).Get([]byte(record.ID))
		if existing != nil && record.Version == 0 {
			return domain.ErrAlreadyExists
		}
		return putJSON(tx, "records", record.ID, record)
	})
}

func (s *Store) GetRecord(id string) (*domain.Record, error) {
	if id == "" {
		return nil, domain.ErrInvalidInput
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return nil, errors.New("storage is closed")
	}
	var result domain.Record
	err := s.db.View(func(tx *bbolt.Tx) error { return getJSON(tx, "records", id, &result) })
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *Store) DeleteRecord(id string) error {
	if id == "" {
		return domain.ErrInvalidInput
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("storage is closed")
	}
	return s.db.Update(func(tx *bbolt.Tx) error {
		if tx.Bucket([]byte("records")).Get([]byte(id)) == nil {
			return domain.ErrMissingRecord
		}
		return deleteKey(tx, "records", id)
	})
}

func (s *Store) ListRecords() ([]*domain.Record, error) {
	records, err := listJSON(s, "records", func(data []byte) (*domain.Record, error) {
		var item domain.Record
		if err := decode(data, &item); err != nil {
			return nil, err
		}
		return &item, nil
	})
	sortRecords(records)
	return records, err
}

func (s *Store) SaveBatch(batch *domain.Batch) error {
	if err := batch.Validate(); err != nil {
		return err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("storage is closed")
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return putJSON(tx, "batches", batch.ID, batch) })
}

func (s *Store) GetBatch(id string) (*domain.Batch, error) {
	if id == "" {
		return nil, domain.ErrInvalidInput
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return nil, errors.New("storage is closed")
	}
	var item domain.Batch
	err := s.db.View(func(tx *bbolt.Tx) error { return getJSON(tx, "batches", id, &item) })
	if err != nil {
		return nil, err
	}
	return &item, nil
}
