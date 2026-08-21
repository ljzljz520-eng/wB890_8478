package memorialstation

import (
	"memorialstation/domain"
	"testing"
)

func TestDomainTransitions(t *testing.T) {
	record := &domain.Record{ID: "r", BatchID: "b", StudentName: "人", GraduationYear: 2024, Message: "一段足够长的留言", Status: domain.StatusDraft, Visibility: domain.VisibilityPrivate}
	for _, status := range []string{domain.StatusSubmitted, domain.StatusApproved, domain.StatusArchived} {
		if err := record.Transition(status); err != nil {
			t.Fatal(err)
		}
	}
	record.Visibility = domain.VisibilityPublic
	if domain.CanTransition(domain.StatusArchived, domain.StatusDraft) {
		t.Fatal("archived record can not return to draft")
	}
	if !domain.IsVisible(record, "viewer") {
		t.Fatal("public visibility expected after change")
	}
}
