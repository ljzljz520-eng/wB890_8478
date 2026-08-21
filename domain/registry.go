package domain

import (
	"sort"
	"strings"
)

type Registry struct {
	batches      map[string]*Batch
	coordinators map[string]string
	aliases      map[string]string
}

func NewRegistry() *Registry {
	return &Registry{batches: map[string]*Batch{}, coordinators: map[string]string{}, aliases: map[string]string{}}
}
func (r *Registry) AddBatch(batch *Batch) error {
	if r == nil || batch == nil {
		return ErrInvalidInput
	}
	if err := batch.Validate(); err != nil {
		return err
	}
	if _, exists := r.batches[batch.ID]; exists {
		return ErrAlreadyExists
	}
	r.batches[batch.ID] = batch
	r.coordinators[batch.ID] = batch.Coordinator
	return nil
}
func (r *Registry) RenameBatch(id, label string) error {
	if r == nil || strings.TrimSpace(label) == "" {
		return ErrInvalidInput
	}
	batch, ok := r.batches[id]
	if !ok {
		return ErrMissingRecord
	}
	batch.Label = strings.TrimSpace(label)
	return nil
}
func (r *Registry) SetAlias(alias, batchID string) error {
	if r == nil || strings.TrimSpace(alias) == "" || strings.TrimSpace(batchID) == "" {
		return ErrInvalidInput
	}
	if _, ok := r.batches[batchID]; !ok {
		return ErrMissingRecord
	}
	r.aliases[strings.ToLower(alias)] = batchID
	return nil
}
func (r *Registry) Resolve(value string) (*Batch, bool) {
	if r == nil {
		return nil, false
	}
	if batch, ok := r.batches[value]; ok {
		return batch, true
	}
	id, ok := r.aliases[strings.ToLower(strings.TrimSpace(value))]
	if !ok {
		return nil, false
	}
	batch, ok := r.batches[id]
	return batch, ok
}
func (r *Registry) RemoveBatch(id string) error {
	if r == nil {
		return ErrInvalidInput
	}
	if _, ok := r.batches[id]; !ok {
		return ErrMissingRecord
	}
	delete(r.batches, id)
	delete(r.coordinators, id)
	for alias, target := range r.aliases {
		if target == id {
			delete(r.aliases, alias)
		}
	}
	return nil
}
func (r *Registry) List() []*Batch {
	if r == nil {
		return nil
	}
	ids := make([]string, 0, len(r.batches))
	for id := range r.batches {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	output := make([]*Batch, 0, len(ids))
	for _, id := range ids {
		output = append(output, r.batches[id])
	}
	return output
}
func (r *Registry) Coordinator(id string) (string, bool) {
	if r == nil {
		return "", false
	}
	coordinator, ok := r.coordinators[id]
	return coordinator, ok
}
func (r *Registry) Count() int {
	if r == nil {
		return 0
	}
	return len(r.batches)
}
