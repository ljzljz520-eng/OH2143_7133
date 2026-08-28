package model

import (
	"strings"
	"time"
)

type Session struct {
	ID         string    `json:"id"`
	RecordID   string    `json:"record_id"`
	ProfileID  string    `json:"profile_id"`
	State      string    `json:"state"`
	StartedAt  time.Time `json:"started_at"`
	ClosedAt   time.Time `json:"closed_at"`
	LastScreen string    `json:"last_screen"`
}

func NewSession(id, recordID, profileID string) Session {
	return Session{ID: strings.TrimSpace(id), RecordID: strings.TrimSpace(recordID), ProfileID: strings.TrimSpace(profileID), State: "open", StartedAt: time.Now().UTC()}
}

func (s *Session) Close() {
	if s.State == "open" {
		s.State = "closed"
		s.ClosedAt = time.Now().UTC()
	}
}

func (s *Session) SetScreen(screen string) bool {
	screen = strings.TrimSpace(screen)
	if screen == "" || s.State != "open" {
		return false
	}
	s.LastScreen = screen
	return true
}

func (s Session) IsOpen() bool { return s.State == "open" }
