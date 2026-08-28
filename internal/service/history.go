package service

import (
	"sort"
	"strings"
	"time"

	"wedding-sign/internal/audit"
	"wedding-sign/internal/catalog"
	"wedding-sign/internal/model"
)

type RecordHistory struct {
	Record model.Record
	Audits []model.Audit
	Report audit.Report
	Valid  catalog.ValidationSummary
}

func (s *Service) History(recordID, profileID string) (RecordHistory, error) {
	record, err := s.store.GetRecord(strings.TrimSpace(recordID))
	if err != nil {
		return RecordHistory{}, err
	}
	profile, err := s.store.GetProfile(strings.TrimSpace(profileID))
	if err != nil {
		return RecordHistory{}, err
	}
	audits, err := s.logger.History(record.ID)
	if err != nil {
		return RecordHistory{}, err
	}
	return RecordHistory{Record: record, Audits: audits, Report: audit.BuildReport(audits), Valid: catalog.ValidateRecordAndProfile(record, profile)}, nil
}

func (s *Service) ActionsSince(recordID string, since time.Time) ([]model.Audit, error) {
	audits, err := s.logger.History(recordID)
	if err != nil {
		return nil, err
	}
	filtered := make([]model.Audit, 0)
	for _, item := range audits {
		if item.At.After(since) || item.At.Equal(since) {
			filtered = append(filtered, item)
		}
	}
	return filtered, nil
}

func (s *Service) LatestAction(recordID string) (model.Audit, bool, error) {
	audits, err := s.logger.History(recordID)
	if err != nil {
		return model.Audit{}, false, err
	}
	if len(audits) == 0 {
		return model.Audit{}, false, nil
	}
	sort.SliceStable(audits, func(i, j int) bool { return audits[i].At.Before(audits[j].At) })
	return audits[len(audits)-1], true, nil
}

func (s *Service) HasConfirmedHistory(recordID string) (bool, error) {
	history, err := s.logger.History(recordID)
	if err != nil {
		return false, err
	}
	for _, item := range history {
		if item.Action == "confirm" && item.After == model.StatusConfirmed {
			return true, nil
		}
	}
	return false, nil
}
