package importx

import (
	"memorialstation/domain"
	"strings"
)

func ValidateHeaders(headers []string) error {
	required := []string{"id", "student", "year", "message", "tags"}
	if len(headers) < len(required) {
		return domain.ErrInvalidInput
	}
	for index, wanted := range required {
		if !strings.EqualFold(strings.TrimSpace(headers[index]), wanted) {
			return domain.ErrInvalidInput
		}
	}
	return nil
}
func NormalizeFields(fields []string) []string {
	normalized := make([]string, len(fields))
	for index, field := range fields {
		normalized[index] = strings.TrimSpace(strings.ReplaceAll(field, "\r", ""))
	}
	return normalized
}
func MergeReports(first, second Report) Report {
	merged := Report{Imported: first.Imported + second.Imported, Rejected: first.Rejected + second.Rejected, Errors: append(append([]string{}, first.Errors...), second.Errors...)}
	return merged
}
