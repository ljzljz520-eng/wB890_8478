package memorialstation

import "testing"

func TestBusiness24Regression(t *testing.T) {
	services := openTestServices(t)
	record := seedRecord(t, services, "890-24")
	if err := services.flow.Submit(record.ID, "student", "1"); err != nil {
		t.Fatal(err)
	}
	if _, err := services.review.Approve(record.ID, "teacher", "confirmed", "2"); err != nil {
		t.Fatal(err)
	}
	if err := services.flow.Confirm(record.ID, "teacher", "3"); err != nil {
		t.Fatal(err)
	}
	if _, err := services.archive.Archive(record.ID, "teacher", "4"); err != nil {
		t.Fatal(err)
	}
	if err := services.store.DeleteRecord(record.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := services.archive.ArchiveBatch(record.BatchID, "teacher", "5"); err != nil {
		t.Fatal(err)
	}
	records, err := services.archive.ReadBack(record.BatchID, "teacher")
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 1 || records[0].Status != "archived" {
		t.Fatalf("batch ZX89024 reread returned incorrect state: %#v", records)
	}
}
