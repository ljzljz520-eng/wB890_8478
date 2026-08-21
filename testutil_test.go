package memorialstation

import (
	"memorialstation/archive"
	"memorialstation/domain"
	"memorialstation/review"
	"memorialstation/search"
	"memorialstation/storage"
	"memorialstation/workflow"
	"testing"
)

type testServices struct {
	store   *storage.Store
	flow    *workflow.Service
	review  *review.Service
	archive *archive.Service
	search  *search.Service
}

func openTestServices(t *testing.T) testServices {
	t.Helper()
	store, err := storage.Open(t.TempDir() + "/memorial.db")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { store.Close() })
	return testServices{store: store, flow: workflow.New(store), review: review.New(store), archive: archive.New(store), search: search.New(store)}
}
func seedRecord(t *testing.T, services testServices, id string) *domain.Record {
	t.Helper()
	batch := &domain.Batch{ID: "ZX89024", Label: "2024毕业班", GraduationYear: 2024, Coordinator: "teacher", CreatedAt: "2024-01-01"}
	if err := services.flow.StartBatch(batch, "2024-01-01"); err != nil {
		t.Fatal(err)
	}
	record := &domain.Record{ID: id, BatchID: batch.ID, StudentName: "林晓", GraduationYear: 2024, Message: "一起走过的青春会发光", Visibility: domain.VisibilityClass, Tags: []string{"memory"}, CreatedAt: "2024-01-01", UpdatedAt: "2024-01-01"}
	if err := services.flow.CreateRecord(record, "2024-01-01"); err != nil {
		t.Fatal(err)
	}
	return record
}
