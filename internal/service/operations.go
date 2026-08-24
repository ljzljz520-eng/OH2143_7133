package service

import (
	"fmt"
	"strings"

	"wedding-sign/internal/catalog"
	"wedding-sign/internal/model"
	"wedding-sign/internal/store"
)

func (s *Service) Search(filter store.RecordFilter) ([]model.Record, error) {
	return s.store.SearchRecords(filter)
}

func (s *Service) CheckAccessibility(recordID, profileID string) (catalog.AccessibilityIssue, []catalog.AccessibilityIssue, error) {
	record, err := s.store.GetRecord(recordID)
	if err != nil {
		return catalog.AccessibilityIssue{}, nil, err
	}
	profile, err := s.store.GetProfile(profileID)
	if err != nil {
		return catalog.AccessibilityIssue{}, nil, err
	}
	issues := catalog.CheckRecord(record, profile)
	if len(issues) == 0 {
		return catalog.AccessibilityIssue{}, issues, nil
	}
	return issues[0], issues, nil
}

func (s *Service) UpdateProfile(profileID, theme, accent, fit string, dim, fullscreen bool) (model.Profile, error) {
	profile, err := s.store.GetProfile(strings.TrimSpace(profileID))
	if err != nil {
		return model.Profile{}, err
	}
	if strings.TrimSpace(theme) != "" {
		profile.Theme = strings.TrimSpace(theme)
	}
	if strings.TrimSpace(accent) != "" {
		profile.AccentColor = strings.TrimSpace(accent)
	}
	if fit != "" && !profile.SetFitMode(fit) {
		return model.Profile{}, fmt.Errorf("fit mode %s is invalid", fit)
	}
	profile.DimBackground = dim
	profile.Fullscreen = fullscreen
	if err := profile.Validate(); err != nil {
		return model.Profile{}, err
	}
	if err := s.store.PutProfile(profile); err != nil {
		return model.Profile{}, err
	}
	_, err = s.logger.Record("", "", "update-profile", "operator", profileID, profile.Theme, s.nextID("request"))
	return profile, err
}

func (s *Service) OpenSession(sessionID, recordID, profileID string) (model.Session, error) {
	session := model.NewSession(sessionID, recordID, profileID)
	if session.ID == "" || session.RecordID == "" || session.ProfileID == "" {
		return model.Session{}, fmt.Errorf("session references are required")
	}
	if _, err := s.store.GetRecord(recordID); err != nil {
		return model.Session{}, err
	}
	if _, err := s.store.GetProfile(profileID); err != nil {
		return model.Session{}, err
	}
	if err := s.store.PutSession(session); err != nil {
		return model.Session{}, err
	}
	_, err := s.logger.Record(recordID, "", "open-session", "display", "closed", "open", sessionID)
	return session, err
}

func (s *Service) CloseSession(sessionID string) (model.Session, error) {
	session, err := s.store.GetSession(sessionID)
	if err != nil {
		return model.Session{}, err
	}
	if !session.IsOpen() {
		return session, nil
	}
	session.Close()
	if err := s.store.PutSession(session); err != nil {
		return model.Session{}, err
	}
	_, err = s.logger.Record(session.RecordID, "", "close-session", "display", "open", "closed", sessionID)
	return session, err
}

func (s *Service) SetSessionScreen(sessionID, screen string) (model.Session, error) {
	session, err := s.store.GetSession(sessionID)
	if err != nil {
		return model.Session{}, err
	}
	if !session.SetScreen(screen) {
		return model.Session{}, fmt.Errorf("session %s is not open", sessionID)
	}
	if err := s.store.PutSession(session); err != nil {
		return model.Session{}, err
	}
	return session, nil
}
