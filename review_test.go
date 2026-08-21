package memorialstation

import (
	"memorialstation/review"
	"testing"
)

func TestReviewMetricsAndRules(t *testing.T) {
	services := openTestServices(t)
	record := seedRecord(t, services, "record-review")
	if err := services.flow.Submit(record.ID, "student", "1"); err != nil {
		t.Fatal(err)
	}
	reviewer := review.New(services.store)
	if _, err := reviewer.Approve(record.ID, "teacher", "good", "2"); err != nil {
		t.Fatal(err)
	}
	metrics, err := reviewer.Metrics(record.BatchID)
	if err != nil {
		t.Fatal(err)
	}
	if metrics.Pending != 1 {
		t.Fatalf("metrics=%#v", metrics)
	}
	item, err := services.store.GetRecord(record.ID)
	if err != nil {
		t.Fatal(err)
	}
	item.Tags = nil
	item.Version++
	if err := services.store.SaveRecord(item); err != nil {
		t.Fatal(err)
	}
	allowed, err := reviewer.CanApprove(record.ID, review.DefaultRules())
	if err != nil {
		t.Fatal(err)
	}
	if allowed {
		t.Fatal("record without required tag should be blocked by default rule")
	}
}
