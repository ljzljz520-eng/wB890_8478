package workflow

import (
	"fmt"
	"memorialstation/domain"
)

type Command struct {
	Name     string
	RecordID string
	Actor    string
	Notes    string
	At       string
}

func (s *Service) Execute(command Command) error {
	if command.RecordID == "" || command.Actor == "" {
		return domain.ErrInvalidInput
	}
	switch command.Name {
	case "submit":
		return s.Submit(command.RecordID, command.Actor, command.At)
	case "confirm":
		return s.Confirm(command.RecordID, command.Actor, command.At)
	case "reject":
		return s.Reject(command.RecordID, command.Actor, command.Notes, command.At)
	default:
		return fmt.Errorf("%w: command %s", domain.ErrInvalidInput, command.Name)
	}
}
func (s *Service) AllowedCommands(recordID string) ([]string, error) {
	record, err := s.store.GetRecord(recordID)
	if err != nil {
		return nil, err
	}
	commands := []string{}
	switch record.Status {
	case domain.StatusDraft:
		commands = append(commands, "submit")
	case domain.StatusSubmitted:
		commands = append(commands, "confirm", "reject")
	case domain.StatusRejected:
		commands = append(commands, "submit")
	case domain.StatusApproved:
		commands = append(commands, "archive")
	}
	return commands, nil
}
func (s *Service) Progress(id string) (float64, error) {
	item, err := s.store.GetWorkflow(id)
	if err != nil {
		return 0, err
	}
	if len(item.Steps) == 0 {
		return 0, nil
	}
	return float64(item.CurrentStep) / float64(len(item.Steps)), nil
}
func (s *Service) ResetWorkflow(id string, now string) (*domain.Workflow, error) {
	item, err := s.store.GetWorkflow(id)
	if err != nil {
		return nil, err
	}
	if item.State == "completed" {
		return nil, domain.ErrInvalidTransition
	}
	item.CurrentStep = 0
	item.State = "pending"
	item.StartedAt = now
	item.CompletedAt = ""
	if err := s.store.SaveWorkflow(item); err != nil {
		return nil, err
	}
	return item, nil
}
