package memorialstation

import (
	"memorialstation/importx"
	"strings"
	"testing"
)

func TestWorkflowImportReport(t *testing.T) {
	services := openTestServices(t)
	batch := seedRecord(t, services, "seed").BatchID
	importer := importx.New(services.store)
	report, err := importer.ImportCSV(strings.NewReader("id,student,year,message,tags\nrecord-3,周宁,2024,毕业快乐大家,club|memory\n"), batch, "importer", "2024-05-01")
	if err != nil {
		t.Fatal(err)
	}
	if report.Imported != 1 || report.Rejected != 0 {
		t.Fatalf("unexpected report: %#v", report)
	}
	issues, err := importer.ValidateBatch(batch)
	if err != nil {
		t.Fatal(err)
	}
	if len(issues) != 0 {
		t.Fatalf("unexpected issues: %v", issues)
	}
}
