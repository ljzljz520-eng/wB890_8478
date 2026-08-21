package memorialstation

import "testing"

func TestArchiveManifest(t *testing.T) {
	services := openTestServices(t)
	record := seedRecord(t, services, "record-manifest")
	if err := services.flow.Submit(record.ID, "student", "1"); err != nil {
		t.Fatal(err)
	}
	if _, err := services.review.Approve(record.ID, "teacher", "ok", "2"); err != nil {
		t.Fatal(err)
	}
	if err := services.flow.Confirm(record.ID, "teacher", "3"); err != nil {
		t.Fatal(err)
	}
	if _, err := services.archive.Archive(record.ID, "teacher", "4"); err != nil {
		t.Fatal(err)
	}
	manifest, err := services.archive.BuildManifest(record.BatchID)
	if err != nil {
		t.Fatal(err)
	}
	valid, err := services.archive.VerifyManifest(manifest)
	if err != nil || !valid {
		t.Fatalf("manifest valid=%v err=%v", valid, err)
	}
}
