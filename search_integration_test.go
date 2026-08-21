package memorialstation

import (
	"memorialstation/domain"
	"testing"
)

func TestWorkflowSearchUpdatePublish(t *testing.T) {
	services := openTestServices(t)
	record := seedRecord(t, services, "record-2")
	if err := services.flow.Submit(record.ID, "student", "2024-05-01"); err != nil {
		t.Fatal(err)
	}
	if _, err := services.review.Approve(record.ID, "teacher", "ok", "2024-05-02"); err != nil {
		t.Fatal(err)
	}
	if err := services.flow.Confirm(record.ID, "teacher", "2024-05-03"); err != nil {
		t.Fatal(err)
	}
	if _, err := services.archive.Archive(record.ID, "teacher", "2024-05-04"); err != nil {
		t.Fatal(err)
	}
	item, err := services.store.GetRecord(record.ID)
	if err != nil {
		t.Fatal(err)
	}
	item.Visibility = domain.VisibilityPrivate
	if err := services.store.SaveRecord(item); err != nil {
		t.Fatal(err)
	}
	if err := services.flow.Publish(record.ID, "2024-05-05"); err != nil {
		t.Fatal(err)
	}
	result, err := services.search.Query(domain.SearchQuery{StudentName: "林"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Total != 1 || result.Records[0].Visibility != domain.VisibilityPublic {
		t.Fatalf("unexpected search result: %#v", result)
	}
}
