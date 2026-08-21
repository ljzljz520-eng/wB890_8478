package search

import (
	"memorialstation/domain"
	"sort"
	"strings"
)

func Rank(records []*domain.Record, query string) []*domain.Record {
	result := append([]*domain.Record(nil), records...)
	lower := strings.ToLower(query)
	sort.SliceStable(result, func(i, j int) bool { return relevance(result[i], lower) > relevance(result[j], lower) })
	return result
}
func relevance(record *domain.Record, query string) int {
	score := 0
	if strings.Contains(strings.ToLower(record.StudentName), query) {
		score += 5
	}
	if strings.Contains(strings.ToLower(record.Message), query) {
		score += 3
	}
	for _, tag := range record.Tags {
		if strings.Contains(strings.ToLower(tag), query) {
			score++
		}
	}
	if record.Status == domain.StatusArchived {
		score++
	}
	return score
}
func GroupByVisibility(records []*domain.Record) map[string][]*domain.Record {
	grouped := map[string][]*domain.Record{}
	for _, record := range records {
		grouped[record.Visibility] = append(grouped[record.Visibility], record)
	}
	return grouped
}
func Names(records []*domain.Record) []string {
	names := make([]string, 0, len(records))
	seen := map[string]bool{}
	for _, record := range records {
		if !seen[record.StudentName] {
			names = append(names, record.StudentName)
			seen[record.StudentName] = true
		}
	}
	sort.Strings(names)
	return names
}
