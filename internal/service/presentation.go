package service

import (
	"bytes"
	"fmt"

	"wedding-sign/internal/display"
	"wedding-sign/internal/model"
)

func (s *Service) BuildView(recordID, profileID string) (display.View, error) {
	record, err := s.store.GetRecord(recordID)
	if err != nil {
		return display.View{}, err
	}
	profile, err := s.store.GetProfile(profileID)
	if err != nil {
		return display.View{}, err
	}
	return s.renderer.View(record, profile), nil
}

func (s *Service) RenderWelcome(recordID, profileID string) (string, error) {
	record, err := s.store.GetRecord(recordID)
	if err != nil {
		return "", err
	}
	profile, err := s.store.GetProfile(profileID)
	if err != nil {
		return "", err
	}
	var buffer bytes.Buffer
	if err := s.renderer.RenderHTML(&buffer, record, profile); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func (s *Service) UpdateRecordText(recordID, text string) (model.Record, error) {
	record, err := s.store.GetRecord(recordID)
	if err != nil {
		return model.Record{}, err
	}
	if !record.IsMutable() {
		return model.Record{}, fmt.Errorf("record %s is locked", recordID)
	}
	if text == "" {
		return model.Record{}, fmt.Errorf("welcome text is required")
	}
	before := record.WelcomeText
	record.WelcomeText = text
	record.Normalize()
	if err := s.store.PutRecord(record); err != nil {
		return model.Record{}, err
	}
	_, err = s.logger.Record(record.ID, "", "update-text", "operator", before, text, s.nextID("request"))
	return record, err
}
