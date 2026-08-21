package memorialstation

import "testing"

func TestWorkflowCreateReviewArchive(t *testing.T) {
	services := openTestServices(t)
	record := seedRecord(t, services, "record-1")
	if err := services.flow.Submit(record.ID, "student", "2024-05-01"); err != nil {
		t.Fatal(err)
	}
	if _, err := services.review.Approve(record.ID, "teacher", "looks good", "2024-05-02"); err != nil {
		t.Fatal(err)
	}
	if err := services.flow.Confirm(record.ID, "teacher", "2024-05-03"); err != nil {
		t.Fatal(err)
	}
	if _, err := services.archive.Archive(record.ID, "teacher", "2024-05-04"); err != nil {
		t.Fatal(err)
	}
	read, err := services.archive.ConfirmReadBack(record.BatchID, "teacher")
	if err != nil {
		t.Fatal(err)
	}
	if len(read) != 1 || read[0].Status != "archived" {
		t.Fatalf("unexpected readback: %#v", read)
	}
}
