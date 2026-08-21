package storage

import (
	"errors"
	"go.etcd.io/bbolt"
	"memorialstation/domain"
)

func (s *Store) SaveReview(review *domain.Review) error {
	if err := review.Validate(); err != nil {
		return err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("storage is closed")
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return putJSON(tx, "reviews", review.ID, review) })
}
func (s *Store) ListReviews(recordID string) ([]*domain.Review, error) {
	items, err := listJSON(s, "reviews", func(data []byte) (*domain.Review, error) {
		var item domain.Review
		if err := decode(data, &item); err != nil {
			return nil, err
		}
		return &item, nil
	})
	if err != nil {
		return nil, err
	}
	filtered := make([]*domain.Review, 0, len(items))
	for _, item := range items {
		if recordID == "" || item.RecordID == recordID {
			filtered = append(filtered, item)
		}
	}
	return filtered, nil
}
