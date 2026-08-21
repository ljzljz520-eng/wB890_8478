package archive

import (
	"encoding/json"
	"fmt"
	"memorialstation/domain"
	"memorialstation/storage"
)

type Service struct{ store *storage.Store }

func New(store *storage.Store) *Service { return &Service{store: store} }

func (s *Service) Archive(recordID, actor, now string) (*domain.ArchiveEntry, error) {
	record, err := s.store.GetRecord(recordID)
	if err != nil {
		return nil, err
	}
	if record.Status != domain.StatusApproved {
		return nil, fmt.Errorf("%w: only approved records can be archived", domain.ErrInvalidTransition)
	}
	if err := record.Transition(domain.StatusArchived); err != nil {
		return nil, err
	}
	record.Editor = actor
	record.UpdatedAt = now
	if err := s.store.SaveRecord(record); err != nil {
		return nil, err
	}
	snapshot, err := json.Marshal(record)
	if err != nil {
		return nil, err
	}
	entry := &domain.ArchiveEntry{ID: fmt.Sprintf("archive-%s-%d", record.BatchID, record.Version), BatchID: record.BatchID, RecordID: record.ID, Snapshot: string(snapshot), Confirmed: true, Viewer: actor, ArchivedAt: now}
	if err := s.store.SaveArchive(entry); err != nil {
		return nil, err
	}
	return entry, nil
}

func (s *Service) ReadBack(batchID, actor string) ([]*domain.Record, error) {
	entries, err := s.store.ListArchives(batchID)
	if err != nil {
		return nil, err
	}
	result := make([]*domain.Record, 0, len(entries))
	for _, entry := range entries {
		record, getErr := s.store.GetRecord(entry.RecordID)
		if getErr == domain.ErrMissingRecord {
			record = nil
		} else if getErr != nil {
			return nil, getErr
		}
		if record == nil {
			record = &domain.Record{ID: entry.RecordID, BatchID: entry.BatchID, Status: domain.StatusDraft, Editor: actor}
		}
		result = append(result, record)
	}
	return result, nil
}

func (s *Service) ConfirmReadBack(batchID, actor string) ([]*domain.Record, error) {
	records, err := s.ReadBack(batchID, actor)
	if err != nil {
		return nil, err
	}
	for _, record := range records {
		if record.Status != domain.StatusArchived {
			return records, fmt.Errorf("%w: %s is not archived", domain.ErrInvalidTransition, record.ID)
		}
	}
	return records, nil
}

func (s *Service) ArchiveBatch(batchID, actor, now string) (int, error) {
	records, err := s.store.ListRecords()
	if err != nil {
		return 0, err
	}
	count := 0
	for _, record := range records {
		if record.BatchID != batchID || record.Status != domain.StatusApproved {
			continue
		}
		if _, err := s.Archive(record.ID, actor, now); err != nil {
			return count, err
		}
		count++
	}
	batch, err := s.store.GetBatch(batchID)
	if err != nil {
		return count, err
	}
	batch.Status = domain.StatusArchived
	batch.ArchivedAt = now
	batch.ApprovedCount = count
	if err := s.store.SaveBatch(batch); err != nil {
		return count, err
	}
	return count, nil
}

func (s *Service) Summary(batchID string) (map[string]int, error) {
	records, err := s.store.ListRecords()
	if err != nil {
		return nil, err
	}
	summary := map[string]int{"draft": 0, "submitted": 0, "approved": 0, "rejected": 0, "archived": 0}
	for _, record := range records {
		if record.BatchID == batchID {
			summary[record.Status]++
		}
	}
	return summary, nil
}
