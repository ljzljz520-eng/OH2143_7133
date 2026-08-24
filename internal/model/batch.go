package model

import (
	"errors"
	"strings"
	"time"
)

type Batch struct {
	ID          string    `json:"id"`
	EventDate   string    `json:"event_date"`
	RecordIDs   []string  `json:"record_ids"`
	Status      string    `json:"status"`
	Cursor      int       `json:"cursor"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	PublishedAt time.Time `json:"published_at"`
}

func NewBatch(id, eventDate string) Batch {
	now := time.Now().UTC()
	return Batch{ID: strings.TrimSpace(id), EventDate: strings.TrimSpace(eventDate), RecordIDs: []string{}, Status: StatusDraft, CreatedAt: now, UpdatedAt: now}
}

func (b Batch) Validate() error {
	if b.ID == "" {
		return errors.New("batch id is required")
	}
	if _, err := time.Parse("2006-01-02", b.EventDate); err != nil {
		return errors.New("batch event date is invalid")
	}
	if len(b.RecordIDs) == 0 {
		return errors.New("batch requires records")
	}
	if b.Cursor < 0 || b.Cursor > len(b.RecordIDs) {
		return errors.New("batch cursor out of range")
	}
	return nil
}

func (b *Batch) AddRecord(id string) bool {
	id = strings.TrimSpace(id)
	if id == "" {
		return false
	}
	for _, existing := range b.RecordIDs {
		if existing == id {
			return false
		}
	}
	b.RecordIDs = append(b.RecordIDs, id)
	b.UpdatedAt = time.Now().UTC()
	return true
}

func (b *Batch) Next() (string, bool) {
	if b.Cursor >= len(b.RecordIDs) {
		return "", false
	}
	id := b.RecordIDs[b.Cursor]
	b.Cursor++
	b.UpdatedAt = time.Now().UTC()
	return id, true
}

func (b Batch) IsComplete() bool { return len(b.RecordIDs) > 0 && b.Cursor >= len(b.RecordIDs) }

func (b Batch) Progress() float64 {
	if len(b.RecordIDs) == 0 {
		return 0
	}
	return float64(b.Cursor) / float64(len(b.RecordIDs))
}
