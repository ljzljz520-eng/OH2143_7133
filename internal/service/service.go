package service

import (
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"wedding-sign/internal/audit"
	"wedding-sign/internal/config"
	"wedding-sign/internal/display"
	"wedding-sign/internal/model"
	"wedding-sign/internal/store"
)

type Service struct {
	store    *store.Store
	logger   *audit.Logger
	renderer *display.Renderer
	settings config.Settings
	sequence uint64
}

func New(s *store.Store, settings config.Settings) (*Service, error) {
	if s == nil {
		return nil, fmt.Errorf("store is required")
	}
	if err := settings.Validate(); err != nil {
		return nil, err
	}
	return &Service{store: s, logger: audit.NewLogger(s), renderer: display.NewRenderer(), settings: settings}, nil
}

func (s *Service) nextID(prefix string) string {
	n := atomic.AddUint64(&s.sequence, 1)
	return fmt.Sprintf("%s-%d-%d", prefix, time.Now().UnixNano(), n)
}

func (s *Service) Store() *store.Store { return s.store }

func (s *Service) Settings() config.Settings { return s.settings }

func (s *Service) CreateProfile(id, name string) (model.Profile, error) {
	profile := model.NewProfile(id, name)
	profile.Theme = s.settings.DefaultTheme
	profile.Locale = s.settings.DefaultLocale
	profile.FitMode = s.settings.DefaultFitMode
	profile.DimBackground = s.settings.DimBackground
	profile.Fullscreen = s.settings.Fullscreen
	if err := profile.Validate(); err != nil {
		return model.Profile{}, err
	}
	if err := s.store.PutProfile(profile); err != nil {
		return model.Profile{}, err
	}
	return profile, nil
}

func (s *Service) CreateRecord(coupleNames, date, venue, heroImage string) (model.Record, error) {
	record := model.NewRecord(s.nextID("record"), coupleNames, date, venue, heroImage)
	record.Normalize()
	if err := record.Validate(); err != nil {
		return model.Record{}, err
	}
	if err := s.store.PutRecord(record); err != nil {
		return model.Record{}, err
	}
	_, err := s.logger.Record(record.ID, "", "create", "operator", "", model.StatusDraft, record.ID)
	return record, err
}

func (s *Service) ValidateRecord(id string) (model.Record, error) {
	record, err := s.store.GetRecord(strings.TrimSpace(id))
	if err != nil {
		return model.Record{}, err
	}
	if err := record.Validate(); err != nil {
		return model.Record{}, err
	}
	if record.Status == model.StatusDraft {
		record.Status = model.StatusValidated
		record.Normalize()
		if err := s.store.PutRecord(record); err != nil {
			return model.Record{}, err
		}
		_, err = s.logger.Record(record.ID, "", "validate", "operator", model.StatusDraft, model.StatusValidated, s.nextID("request"))
	}
	return record, err
}

func (s *Service) GetRecord(id string) (model.Record, error) {
	return s.store.GetRecord(strings.TrimSpace(id))
}

func (s *Service) ListRecords(date string) ([]model.Record, error) {
	if strings.TrimSpace(date) == "" {
		return s.store.ListRecords()
	}
	return s.store.FindRecordsByDate(date)
}
