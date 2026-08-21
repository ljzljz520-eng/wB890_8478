package domain

import (
	"fmt"
	"sort"
	"strings"
)

func DescribeRecord(record *Record) string {
	if record == nil {
		return "<missing record>"
	}
	return fmt.Sprintf("%s (%s) %s: %s", record.ID, record.Status, record.StudentName, record.Message)
}
func CloneRecord(record *Record) *Record {
	if record == nil {
		return nil
	}
	clone := *record
	clone.Tags = append([]string(nil), record.Tags...)
	return &clone
}
func NormalizeRecord(record *Record) error {
	if record == nil {
		return ErrInvalidInput
	}
	record.StudentName = strings.TrimSpace(record.StudentName)
	record.Message = strings.TrimSpace(record.Message)
	record.BatchID = strings.TrimSpace(record.BatchID)
	record.Tags = sortedUnique(record.Tags)
	return record.Validate()
}
func sortedUnique(values []string) []string {
	seen := map[string]bool{}
	output := []string{}
	for _, value := range values {
		clean := strings.ToLower(strings.TrimSpace(value))
		if clean != "" && !seen[clean] {
			seen[clean] = true
			output = append(output, clean)
		}
	}
	sort.Strings(output)
	return output
}
func CompareRecords(left, right *Record) bool {
	if left == nil || right == nil {
		return left == right
	}
	return left.ID == right.ID && left.BatchID == right.BatchID && left.Status == right.Status && left.Version == right.Version && left.Message == right.Message
}
func IsVisible(record *Record, viewer string) bool {
	if record == nil {
		return false
	}
	switch record.Visibility {
	case VisibilityPublic:
		return true
	case VisibilityClass:
		return viewer != ""
	default:
		return viewer == record.Editor
	}
}
func StatusLabel(status string) string {
	labels := map[string]string{StatusDraft: "草稿", StatusSubmitted: "待审核", StatusApproved: "已确认", StatusRejected: "已退回", StatusArchived: "已归档"}
	if value, ok := labels[status]; ok {
		return value
	}
	return "未知状态"
}
func VisibilityLabel(value string) string {
	labels := map[string]string{VisibilityPrivate: "仅本人", VisibilityClass: "班级可见", VisibilityPublic: "公开展示"}
	if label, ok := labels[value]; ok {
		return label
	}
	return "未知范围"
}
