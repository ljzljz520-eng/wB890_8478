package review

import (
	"memorialstation/domain"
	"strings"
)

type RuleSet struct {
	BlockedWords []string
	MinimumScore int
	RequireTag   string
}

func DefaultRules() RuleSet {
	return RuleSet{BlockedWords: []string{"password", "secret"}, MinimumScore: 20, RequireTag: "memory"}
}
func (rules RuleSet) Check(record *domain.Record) []string {
	issues := []string{}
	if record == nil {
		return []string{"missing record"}
	}
	lower := strings.ToLower(record.Message)
	for _, word := range rules.BlockedWords {
		if strings.Contains(lower, word) {
			issues = append(issues, "blocked word:"+word)
		}
	}
	if len(record.Message) < rules.MinimumScore {
		issues = append(issues, "score below minimum")
	}
	if rules.RequireTag != "" && !hasTag(record.Tags, rules.RequireTag) {
		issues = append(issues, "required tag missing")
	}
	return issues
}
func hasTag(tags []string, wanted string) bool {
	for _, tag := range tags {
		if strings.EqualFold(tag, wanted) {
			return true
		}
	}
	return false
}
func (s *Service) RuleCheck(recordID string, rules RuleSet) ([]string, error) {
	record, err := s.store.GetRecord(recordID)
	if err != nil {
		return nil, err
	}
	return rules.Check(record), nil
}
func (s *Service) CanApprove(recordID string, rules RuleSet) (bool, error) {
	issues, err := s.RuleCheck(recordID, rules)
	if err != nil {
		return false, err
	}
	return len(issues) == 0, nil
}
