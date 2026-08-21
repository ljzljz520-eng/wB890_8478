package memorialstation

import (
	"memorialstation/domain"
	"memorialstation/storage"
	"testing"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := t.TempDir() + "/persistent.db"
	store, err := storage.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.SaveRecord(&domain.Record{ID: "persist-record", BatchID: "persist-batch", StudentName: "赵一", GraduationYear: 2024, Message: "永远记得我们的教室", Status: domain.StatusDraft, Visibility: domain.VisibilityClass, Version: 1, CreatedAt: "2024", UpdatedAt: "2024"}); err != nil {
		t.Fatal(err)
	}
	if err := store.SaveAudit(&domain.AuditEvent{ID: "persist-audit", RecordID: "persist-record", BatchID: "persist-batch", Action: "create", Actor: "tester", FromStatus: "", ToStatus: "draft", Sequence: 1, CreatedAt: "2024"}); err != nil {
		t.Fatal(err)
	}
	if err := store.SaveWorkflow(&domain.Workflow{ID: "persist-workflow", BatchID: "persist-batch", Name: "闭环", Owner: "tester", State: "pending", Steps: []string{"登记", "审核", "确认", "归档"}}); err != nil {
		t.Fatal(err)
	}
	if err := store.SaveAttachment(&domain.Attachment{ID: "persist-attachment", RecordID: "persist-record", BatchID: "persist-batch", Name: "photo.txt", MediaType: "text/plain", Content: []byte("photo"), Checksum: "hash", UploadedBy: "tester"}); err != nil {
		t.Fatal(err)
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	reopened, err := storage.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	if _, err := reopened.GetRecord("persist-record"); err != nil {
		t.Fatal(err)
	}
	if _, err := reopened.GetWorkflow("persist-workflow"); err != nil {
		t.Fatal(err)
	}
	if _, err := reopened.GetAttachment("persist-attachment"); err != nil {
		t.Fatal(err)
	}
	events, err := reopened.ListAudit("persist-batch")
	if err != nil || len(events) != 1 {
		t.Fatalf("events: %v %#v", err, events)
	}
}
