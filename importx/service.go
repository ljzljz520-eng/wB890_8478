package importx

import (
	"encoding/csv"
	"fmt"
	"io"
	"memorialstation/domain"
	"memorialstation/storage"
	"strconv"
	"strings"
)

type Service struct{ store *storage.Store }

func New(store *storage.Store) *Service { return &Service{store: store} }

type Report struct {
	Imported int
	Rejected int
	Errors   []string
}

func (s *Service) ImportCSV(reader io.Reader, batchID, actor, now string) (*Report, error) {
	if strings.TrimSpace(batchID) == "" {
		return nil, domain.ErrInvalidInput
	}
	csvReader := csv.NewReader(reader)
	report := &Report{}
	row := 0
	for {
		fields, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		row++
		if err != nil {
			report.Rejected++
			report.Errors = append(report.Errors, fmt.Sprintf("row %d: %v", row, err))
			continue
		}
		if row == 1 && strings.EqualFold(fields[0], "id") {
			continue
		}
		record, err := parse(fields, batchID, actor, now)
		if err != nil {
			report.Rejected++
			report.Errors = append(report.Errors, fmt.Sprintf("row %d: %v", row, err))
			continue
		}
		if err := s.store.SaveRecord(record); err != nil {
			report.Rejected++
			report.Errors = append(report.Errors, fmt.Sprintf("row %d: %v", row, err))
			continue
		}
		report.Imported++
	}
	return report, nil
}

func parse(fields []string, batchID, actor, now string) (*domain.Record, error) {
	if len(fields) < 5 {
		return nil, domain.ErrInvalidInput
	}
	year, err := strconv.Atoi(fields[2])
	if err != nil {
		return nil, err
	}
	tags := strings.Split(fields[4], "|")
	return &domain.Record{ID: fields[0], BatchID: batchID, StudentName: fields[1], GraduationYear: year, Message: fields[3], Tags: tags, Editor: actor, Status: domain.StatusDraft, Visibility: domain.VisibilityPrivate, CreatedAt: now, UpdatedAt: now}, nil
}

func (s *Service) ExportCSV(writer io.Writer, batchID string) error {
	records, err := s.store.ListRecords()
	if err != nil {
		return err
	}
	output := csv.NewWriter(writer)
	if err := output.Write([]string{"id", "student", "year", "message", "tags", "status"}); err != nil {
		return err
	}
	for _, record := range records {
		if record.BatchID != batchID {
			continue
		}
		if err := output.Write([]string{record.ID, record.StudentName, strconv.Itoa(record.GraduationYear), record.Message, strings.Join(record.Tags, "|"), record.Status}); err != nil {
			return err
		}
	}
	output.Flush()
	return output.Error()
}

func (s *Service) ValidateBatch(batchID string) ([]string, error) {
	records, err := s.store.ListRecords()
	if err != nil {
		return nil, err
	}
	issues := make([]string, 0)
	seen := map[string]bool{}
	for _, record := range records {
		if record.BatchID != batchID {
			continue
		}
		if seen[record.ID] {
			issues = append(issues, "duplicate:"+record.ID)
		}
		seen[record.ID] = true
		if err := record.Validate(); err != nil {
			issues = append(issues, record.ID+":invalid")
		}
		if len(record.Tags) == 0 {
			issues = append(issues, record.ID+":untagged")
		}
	}
	return issues, nil
}
