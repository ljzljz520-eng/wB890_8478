package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"go.etcd.io/bbolt"
	"memorialstation/domain"
)

var bucketNames = [][]byte{[]byte("records"), []byte("audit_events"), []byte("workflows"), []byte("attachments"), []byte("batches"), []byte("reviews"), []byte("archives")}

type Store struct {
	db   *bbolt.DB
	mu   sync.RWMutex
	path string
}

func Open(path string) (*Store, error) {
	if strings.TrimSpace(path) == "" {
		return nil, errors.New("storage path is required")
	}
	db, err := bbolt.Open(filepath.Clean(path), 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("open bolt: %w", err)
	}
	store := &Store{db: db, path: filepath.Clean(path)}
	if err := store.initialize(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return store, nil
}

func (s *Store) initialize() error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		for _, name := range bucketNames {
			if _, err := tx.CreateBucketIfNotExists(name); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db == nil {
		return nil
	}
	err := s.db.Close()
	s.db = nil
	return err
}
func (s *Store) Path() string          { s.mu.RLock(); defer s.mu.RUnlock(); return s.path }
func encode(value any) ([]byte, error) { return json.Marshal(value) }
func decode(data []byte, target any) error {
	if len(data) == 0 {
		return domain.ErrMissingRecord
	}
	return json.Unmarshal(data, target)
}
func putJSON(tx *bbolt.Tx, bucket, key string, value any) error {
	data, err := encode(value)
	if err != nil {
		return err
	}
	return tx.Bucket([]byte(bucket)).Put([]byte(key), data)
}
func getJSON(tx *bbolt.Tx, bucket, key string, target any) error {
	return decode(tx.Bucket([]byte(bucket)).Get([]byte(key)), target)
}
func deleteKey(tx *bbolt.Tx, bucket, key string) error {
	return tx.Bucket([]byte(bucket)).Delete([]byte(key))
}

func listJSON[T any](s *Store, bucket string, decodeValue func([]byte) (*T, error)) ([]*T, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return nil, errors.New("storage is closed")
	}
	values := make([]*T, 0)
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte(bucket)).ForEach(func(_, value []byte) error {
			item, err := decodeValue(value)
			if err != nil {
				return err
			}
			values = append(values, item)
			return nil
		})
	})
	return values, err
}

func sortRecords(records []*domain.Record) {
	sort.Slice(records, func(i, j int) bool {
		if records[i].BatchID == records[j].BatchID {
			return records[i].ID < records[j].ID
		}
		return records[i].BatchID < records[j].BatchID
	})
}
func sortEvents(events []*domain.AuditEvent) {
	sort.Slice(events, func(i, j int) bool {
		if events[i].Sequence == events[j].Sequence {
			return events[i].ID < events[j].ID
		}
		return events[i].Sequence < events[j].Sequence
	})
}
