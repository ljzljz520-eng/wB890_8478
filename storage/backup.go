package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"go.etcd.io/bbolt"
	"memorialstation/domain"
)

type Snapshot struct {
	Records     []*domain.Record
	AuditEvents []*domain.AuditEvent
	Workflows   []*domain.Workflow
	Attachments []*domain.Attachment
	Batches     []*domain.Batch
	Reviews     []*domain.Review
	Archives    []*domain.ArchiveEntry
}

func (s *Store) ExportSnapshot() (Snapshot, error) {
	records, err := s.ListRecords()
	if err != nil {
		return Snapshot{}, err
	}
	audits, err := s.ListAudit("")
	if err != nil {
		return Snapshot{}, err
	}
	workflows, err := s.ListWorkflows("")
	if err != nil {
		return Snapshot{}, err
	}
	attachments, err := s.ListAttachments("")
	if err != nil {
		return Snapshot{}, err
	}
	reviews, err := s.ListReviews("")
	if err != nil {
		return Snapshot{}, err
	}
	archives, err := s.ListArchives("")
	if err != nil {
		return Snapshot{}, err
	}
	batches := []*domain.Batch{}
	return Snapshot{Records: records, AuditEvents: audits, Workflows: workflows, Attachments: attachments, Batches: batches, Reviews: reviews, Archives: archives}, nil
}
func (s *Store) SnapshotJSON() ([]byte, error) {
	snapshot, err := s.ExportSnapshot()
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(snapshot, "", "  ")
}
func (s *Store) ImportSnapshot(data []byte) error {
	var snapshot Snapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return err
	}
	for _, record := range snapshot.Records {
		if err := s.SaveRecord(record); err != nil {
			return err
		}
	}
	for _, event := range snapshot.AuditEvents {
		if err := s.SaveAudit(event); err != nil {
			return err
		}
	}
	for _, workflow := range snapshot.Workflows {
		if err := s.SaveWorkflow(workflow); err != nil {
			return err
		}
	}
	for _, attachment := range snapshot.Attachments {
		if err := s.SaveAttachment(attachment); err != nil {
			return err
		}
	}
	for _, review := range snapshot.Reviews {
		if err := s.SaveReview(review); err != nil {
			return err
		}
	}
	for _, archive := range snapshot.Archives {
		if err := s.SaveArchive(archive); err != nil {
			return err
		}
	}
	return nil
}
func (s *Store) Compact() error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("storage is closed")
	}
	return s.db.View(func(tx *bbolt.Tx) error {
		for _, name := range bucketNames {
			if tx.Bucket(name) == nil {
				return fmt.Errorf("missing bucket %s", name)
			}
		}
		return nil
	})
}
