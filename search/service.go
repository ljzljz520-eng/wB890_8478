package search

import (
	"memorialstation/domain"
	"memorialstation/storage"
	"strings"
)

type Service struct{ store *storage.Store }

func New(store *storage.Store) *Service { return &Service{store: store} }

func (s *Service) Query(query domain.SearchQuery) (*domain.SearchResult, error) {
	records, err := s.store.ListRecords()
	if err != nil {
		return nil, err
	}
	filtered := make([]*domain.Record, 0, len(records))
	for _, record := range records {
		if query.BatchID != "" && record.BatchID != query.BatchID {
			continue
		}
		if query.Status != "" && record.Status != query.Status {
			continue
		}
		if query.Visibility != "" && record.Visibility != query.Visibility {
			continue
		}
		if query.StudentName != "" && !strings.Contains(strings.ToLower(record.StudentName), strings.ToLower(query.StudentName)) {
			continue
		}
		if query.Tag != "" && !contains(record.Tags, query.Tag) {
			continue
		}
		filtered = append(filtered, record)
	}
	return &domain.SearchResult{Records: filtered, Total: len(filtered), Applied: query}, nil
}

func contains(values []string, wanted string) bool {
	for _, value := range values {
		if strings.EqualFold(value, wanted) {
			return true
		}
	}
	return false
}
func (s *Service) Public(batchID string) ([]*domain.Record, error) {
	result, err := s.Query(domain.SearchQuery{BatchID: batchID, Visibility: domain.VisibilityPublic, Status: domain.StatusArchived})
	if err != nil {
		return nil, err
	}
	return result.Records, nil
}
func (s *Service) Facets(batchID string) (map[string]int, error) {
	records, err := s.store.ListRecords()
	if err != nil {
		return nil, err
	}
	facets := map[string]int{}
	for _, record := range records {
		if record.BatchID == batchID {
			facets[record.Visibility]++
			for _, tag := range record.Tags {
				facets["tag:"+tag]++
			}
		}
	}
	return facets, nil
}
