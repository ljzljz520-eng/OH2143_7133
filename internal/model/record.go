package model

import (
	"errors"
	"strings"
	"time"
)

const (
	StatusDraft      = "draft"
	StatusValidated  = "validated"
	StatusConfirmed  = "confirmed"
	StatusQueued     = "queued"
	StatusPublished  = "published"
	StatusDispatched = "dispatched"
	StatusArchived   = "archived"
)

type Record struct {
	ID           string    `json:"id"`
	CoupleNames  string    `json:"couple_names"`
	Date         string    `json:"date"`
	Venue        string    `json:"venue"`
	HeroImage    string    `json:"hero_image"`
	WelcomeText  string    `json:"welcome_text"`
	Status       string    `json:"status"`
	DisplayOrder int       `json:"display_order"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func NewRecord(id, coupleNames, date, venue, heroImage string) Record {
	now := time.Now().UTC()
	return Record{ID: strings.TrimSpace(id), CoupleNames: strings.TrimSpace(coupleNames), Date: strings.TrimSpace(date), Venue: strings.TrimSpace(venue), HeroImage: strings.TrimSpace(heroImage), WelcomeText: "Welcome to our celebration", Status: StatusDraft, CreatedAt: now, UpdatedAt: now}
}

func (r Record) Validate() error {
	if strings.TrimSpace(r.ID) == "" {
		return errors.New("record id is required")
	}
	if strings.TrimSpace(r.CoupleNames) == "" {
		return errors.New("couple names are required")
	}
	if _, err := time.Parse("2006-01-02", r.Date); err != nil {
		return errors.New("date must use YYYY-MM-DD")
	}
	if strings.TrimSpace(r.Venue) == "" {
		return errors.New("venue is required")
	}
	if r.Status == "" {
		return errors.New("status is required")
	}
	return nil
}

func (r *Record) Normalize() {
	r.ID = strings.TrimSpace(r.ID)
	r.CoupleNames = strings.Join(strings.Fields(r.CoupleNames), " ")
	r.Date = strings.TrimSpace(r.Date)
	r.Venue = strings.Join(strings.Fields(r.Venue), " ")
	r.HeroImage = strings.TrimSpace(r.HeroImage)
	r.WelcomeText = strings.TrimSpace(r.WelcomeText)
	if r.WelcomeText == "" {
		r.WelcomeText = "Welcome to our celebration"
	}
	r.UpdatedAt = time.Now().UTC()
}

func (r Record) HasHeroImage() bool { return strings.TrimSpace(r.HeroImage) != "" }

func (r Record) IsMutable() bool {
	return r.Status == StatusDraft || r.Status == StatusValidated || r.Status == StatusConfirmed
}

func (r Record) CanQueue() bool { return r.Status == StatusConfirmed }

func (r Record) CanPublish() bool { return r.Status == StatusQueued || r.Status == StatusConfirmed }

func (r Record) DisplayTitle() string {
	if strings.TrimSpace(r.CoupleNames) == "" {
		return "Wedding Welcome"
	}
	return r.CoupleNames
}
