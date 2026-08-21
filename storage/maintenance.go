package storage

import (
	"errors"
	"fmt"
	"go.etcd.io/bbolt"
	"memorialstation/domain"
)

type Health struct {
	Open        bool
	Buckets     int
	Records     int
	Audits      int
	Workflows   int
	Attachments int
	Archives    int
}

func (s *Store) Health() (Health, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return Health{}, errors.New("storage is closed")
	}
	health := Health{Open: true}
	err := s.db.View(func(tx *bbolt.Tx) error {
		for _, bucket := range bucketNames {
			if tx.Bucket(bucket) == nil {
				return fmt.Errorf("missing bucket %s", bucket)
			}
		}
		health.Buckets = len(bucketNames)
		health.Records = tx.Bucket([]byte("records")).Stats().KeyN
		health.Audits = tx.Bucket([]byte("audit_events")).Stats().KeyN
		health.Workflows = tx.Bucket([]byte("workflows")).Stats().KeyN
		health.Attachments = tx.Bucket([]byte("attachments")).Stats().KeyN
		health.Archives = tx.Bucket([]byte("archives")).Stats().KeyN
		return nil
	})
	return health, err
}
func (s *Store) VerifyRecord(recordID string) error {
	record, err := s.GetRecord(recordID)
	if err != nil {
		return err
	}
	if err := record.Validate(); err != nil {
		return err
	}
	if record.Version < 0 {
		return domain.ErrInvalidInput
	}
	return nil
}
func (s *Store) CountByStatus(batchID string) (map[string]int, error) {
	records, err := s.ListRecords()
	if err != nil {
		return nil, err
	}
	counts := map[string]int{}
	for _, record := range records {
		if record.BatchID != batchID {
			continue
		}
		counts[record.Status]++
	}
	return counts, nil
}
func (s *Store) CopyBatch(batchID, destination string) (int, error) {
	records, err := s.ListRecords()
	if err != nil {
		return 0, err
	}
	copied := 0
	for _, record := range records {
		if record.BatchID != batchID {
			continue
		}
		clone := *record
		clone.ID = destination + "-" + record.ID
		clone.BatchID = destination
		clone.Version = 0
		if err := s.SaveRecord(&clone); err != nil {
			return copied, err
		}
		copied++
	}
	return copied, nil
}
func (s *Store) RemoveBatch(batchID string) (int, error) {
	records, err := s.ListRecords()
	if err != nil {
		return 0, err
	}
	removed := 0
	for _, record := range records {
		if record.BatchID != batchID {
			continue
		}
		if err := s.DeleteRecord(record.ID); err != nil {
			return removed, err
		}
		removed++
	}
	return removed, nil
}
