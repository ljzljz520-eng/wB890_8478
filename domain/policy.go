package domain

import "strings"

type Policy struct {
	Name                 string
	AllowedTags          []string
	RequiredVisibility   string
	MinimumMessageLength int
	RequireAttachment    bool
}
type PolicyResult struct {
	Accepted       bool
	Reasons        []string
	NormalizedTags []string
}

func DefaultPolicy() Policy {
	return Policy{Name: "graduation-memory", AllowedTags: []string{"class", "club", "memory", "photo", "quote"}, RequiredVisibility: VisibilityClass, MinimumMessageLength: 8}
}
func (p Policy) Evaluate(record *Record, attachments []*Attachment) PolicyResult {
	result := PolicyResult{Accepted: true, NormalizedTags: normalizeTags(record.Tags)}
	if record == nil {
		return PolicyResult{Accepted: false, Reasons: []string{"missing record"}}
	}
	if len([]rune(strings.TrimSpace(record.Message))) < p.MinimumMessageLength {
		result.Accepted = false
		result.Reasons = append(result.Reasons, "message is too short")
	}
	if p.RequiredVisibility != "" && record.Visibility == VisibilityPrivate {
		result.Accepted = false
		result.Reasons = append(result.Reasons, "private records need class visibility")
	}
	if len(result.NormalizedTags) == 0 {
		result.Accepted = false
		result.Reasons = append(result.Reasons, "at least one tag is required")
	}
	for _, tag := range result.NormalizedTags {
		if !containsString(p.AllowedTags, tag) {
			result.Accepted = false
			result.Reasons = append(result.Reasons, "unsupported tag:"+tag)
		}
	}
	if p.RequireAttachment && len(attachments) == 0 {
		result.Accepted = false
		result.Reasons = append(result.Reasons, "attachment required")
	}
	return result
}

func normalizeTags(tags []string) []string {
	output := make([]string, 0, len(tags))
	seen := map[string]bool{}
	for _, tag := range tags {
		clean := strings.ToLower(strings.TrimSpace(tag))
		if clean == "" || seen[clean] {
			continue
		}
		seen[clean] = true
		output = append(output, clean)
	}
	return output
}
func containsString(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}
func (p Policy) Normalize(record *Record) {
	if record == nil {
		return
	}
	record.Tags = normalizeTags(record.Tags)
	if record.Visibility == "" {
		record.Visibility = VisibilityPrivate
	}
	if record.Status == "" {
		record.Status = StatusDraft
	}
}
func (p Policy) Explain(result PolicyResult) string {
	if result.Accepted {
		return "accepted"
	}
	return strings.Join(result.Reasons, "; ")
}
func (p Policy) Merge(other Policy) Policy {
	merged := p
	if other.Name != "" {
		merged.Name = other.Name
	}
	if other.RequiredVisibility != "" {
		merged.RequiredVisibility = other.RequiredVisibility
	}
	if other.MinimumMessageLength > 0 {
		merged.MinimumMessageLength = other.MinimumMessageLength
	}
	if other.AllowedTags != nil {
		merged.AllowedTags = append([]string(nil), other.AllowedTags...)
	}
	if other.RequireAttachment {
		merged.RequireAttachment = true
	}
	return merged
}
