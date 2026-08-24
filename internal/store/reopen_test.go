package store

import (
	"path/filepath"
	"testing"

	"wedding-sign/internal/model"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "persistent.db")
	first, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	record := model.NewRecord("persisted", "A & B", "2026-08-24", "Hall", "hero.jpg")
	profile := model.NewProfile("ceremony", "Ceremony")
	batch := model.NewBatch("batch", "2026-08-24")
	batch.AddRecord(record.ID)
	if err := first.PutRecord(record); err != nil {
		t.Fatal(err)
	}
	if err := first.PutProfile(profile); err != nil {
		t.Fatal(err)
	}
	if err := first.PutBatch(batch); err != nil {
		t.Fatal(err)
	}
	if err := first.Close(); err != nil {
		t.Fatal(err)
	}
	second, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer second.Close()
	got, err := second.GetRecord(record.ID)
	if err != nil || got.Venue != "Hall" {
		t.Fatalf("reopen record: %+v %v", got, err)
	}
	if _, err := second.GetProfile(profile.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := second.GetBatch(batch.ID); err != nil {
		t.Fatal(err)
	}
}
