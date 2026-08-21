package workflow

import (
	"fmt"
	"memorialstation/domain"
	"strings"
)

type Report struct {
	BatchID       string
	WorkflowCount int
	RecordCount   int
	Approved      int
	Archived      int
	Issues        []string
	Narrative     string
}

func (s *Service) BuildReport(batchID string) (Report, error) {
	report := Report{BatchID: batchID}
	records, err := s.store.ListRecords()
	if err != nil {
		return report, err
	}
	for _, record := range records {
		if record.BatchID != batchID {
			continue
		}
		report.RecordCount++
		if record.Status == domain.StatusApproved {
			report.Approved++
		}
		if record.Status == domain.StatusArchived {
			report.Archived++
		}
		if record.Message == "" {
			report.Issues = append(report.Issues, record.ID+":empty message")
		}
		if record.Version < 1 {
			report.Issues = append(report.Issues, record.ID+":never changed")
		}
	}
	workflows, err := s.store.ListWorkflows(batchID)
	if err != nil {
		return report, err
	}
	report.WorkflowCount = len(workflows)
	report.Narrative = fmt.Sprintf("batch %s contains %d records, %d approved and %d archived", batchID, report.RecordCount, report.Approved, report.Archived)
	return report, nil
}

func (s *Service) Checklist(batchID string) ([]string, error) {
	report, err := s.BuildReport(batchID)
	if err != nil {
		return nil, err
	}
	checks := []string{}
	if report.RecordCount == 0 {
		checks = append(checks, "no records")
	}
	if report.Approved < report.RecordCount {
		checks = append(checks, "not every record approved")
	}
	if report.Archived < report.Approved {
		checks = append(checks, "approved records remain unarchived")
	}
	checks = append(checks, report.Issues...)
	return checks, nil
}
func (s *Service) FormatReport(report Report) string {
	lines := []string{report.Narrative}
	if len(report.Issues) > 0 {
		lines = append(lines, "issues: "+strings.Join(report.Issues, ", "))
	}
	return strings.Join(lines, "\n")
}
