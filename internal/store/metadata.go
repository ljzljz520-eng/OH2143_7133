package store

import (
	"sort"
	"strings"

	"wedding-sign/internal/model"
)

func (s *Store) ListProfiles() ([]model.Profile, error) {
	profiles := make([]model.Profile, 0)
	err := s.list("profiles", func(raw []byte) error {
		var profile model.Profile
		if err := decode(raw, &profile); err != nil {
			return err
		}
		profiles = append(profiles, profile)
		return nil
	})
	sort.Slice(profiles, func(i, j int) bool { return profiles[i].DisplayName < profiles[j].DisplayName })
	return profiles, err
}

func (s *Store) FindProfileByName(name string) (model.Profile, error) {
	profiles, err := s.ListProfiles()
	if err != nil {
		return model.Profile{}, err
	}
	for _, profile := range profiles {
		if strings.EqualFold(profile.DisplayName, strings.TrimSpace(name)) {
			return profile, nil
		}
	}
	return model.Profile{}, fmtNotFound("profile", name)
}

func (s *Store) ListSessions() ([]model.Session, error) {
	sessions := make([]model.Session, 0)
	err := s.list("sessions", func(raw []byte) error {
		var session model.Session
		if err := decode(raw, &session); err != nil {
			return err
		}
		sessions = append(sessions, session)
		return nil
	})
	sort.Slice(sessions, func(i, j int) bool { return sessions[i].StartedAt.Before(sessions[j].StartedAt) })
	return sessions, err
}

func (s *Store) DeleteSession(id string) error { return s.remove("sessions", id) }

type notFoundError string

func (e notFoundError) Error() string { return string(e) + " not found" }

func fmtNotFound(kind, id string) error { return notFoundError(kind + " " + strings.TrimSpace(id)) }
