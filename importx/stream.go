package importx

import (
	"bufio"
	"io"
	"memorialstation/domain"
	"strings"
)

func ScanLines(reader io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(reader)
	lines := []string{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines, scanner.Err()
}
func ParseTags(value string) []string {
	parts := strings.FieldsFunc(value, func(r rune) bool { return r == ',' || r == '|' || r == ';' })
	tags := []string{}
	seen := map[string]bool{}
	for _, part := range parts {
		tag := strings.ToLower(strings.TrimSpace(part))
		if tag != "" && !seen[tag] {
			tags = append(tags, tag)
			seen[tag] = true
		}
	}
	return tags
}
func CheckRecord(record *domain.Record) error {
	if record == nil {
		return domain.ErrInvalidInput
	}
	return domain.NormalizeRecord(record)
}
func BuildRecord(id, batch, name, message string, year int, tags []string) (*domain.Record, error) {
	record := &domain.Record{ID: id, BatchID: batch, StudentName: name, Message: message, GraduationYear: year, Tags: tags, Status: domain.StatusDraft, Visibility: domain.VisibilityPrivate}
	if err := CheckRecord(record); err != nil {
		return nil, err
	}
	return record, nil
}
