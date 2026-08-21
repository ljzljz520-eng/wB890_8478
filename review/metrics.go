package review

import (
	"memorialstation/domain"
)

type Metrics struct {
	Total        int
	Pending      int
	Approved     int
	Rejected     int
	AverageScore int
}

func (s *Service) Metrics(batchID string) (Metrics, error) {
	records, err := s.store.ListRecords()
	if err != nil {
		return Metrics{}, err
	}
	metrics := Metrics{}
	totalScore := 0
	scoreCount := 0
	for _, record := range records {
		if record.BatchID != batchID {
			continue
		}
		metrics.Total++
		switch record.Status {
		case domain.StatusSubmitted:
			metrics.Pending++
		case domain.StatusApproved, domain.StatusArchived:
			metrics.Approved++
		case domain.StatusRejected:
			metrics.Rejected++
		}
		reviews, reviewErr := s.store.ListReviews(record.ID)
		if reviewErr != nil {
			return metrics, reviewErr
		}
		for _, review := range reviews {
			totalScore += review.Score
			scoreCount++
		}
	}
	if scoreCount > 0 {
		metrics.AverageScore = totalScore / scoreCount
	}
	return metrics, nil
}
func (s *Service) EligibleForArchive(recordID string) (bool, error) {
	record, err := s.store.GetRecord(recordID)
	if err != nil {
		return false, err
	}
	if record.Status != domain.StatusApproved {
		return false, nil
	}
	attachments, err := s.store.ListAttachments(recordID)
	if err != nil {
		return false, err
	}
	for _, attachment := range attachments {
		if !attachment.Approved {
			return false, nil
		}
	}
	return true, nil
}
func (s *Service) ReviewTrail(recordID string) ([]string, error) {
	reviews, err := s.store.ListReviews(recordID)
	if err != nil {
		return nil, err
	}
	trail := make([]string, 0, len(reviews))
	for _, review := range reviews {
		trail = append(trail, review.Reviewer+":"+review.Decision)
	}
	return trail, nil
}
