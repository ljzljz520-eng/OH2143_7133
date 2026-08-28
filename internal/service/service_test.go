package service

import (
	"path/filepath"
	"testing"

	"wedding-sign/internal/config"
	"wedding-sign/internal/model"
	"wedding-sign/internal/store"
)

func testService(t *testing.T) *Service {
	t.Helper()
	settings := config.DefaultSettings(t.TempDir())
	db, err := store.Open(filepath.Join(t.TempDir(), "service.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	svc, err := New(db, settings)
	if err != nil {
		t.Fatal(err)
	}
	return svc
}

func TestServiceLifecycle(t *testing.T) {
	svc := testService(t)
	profile, err := svc.CreateProfile("p1", "Ceremony")
	if err != nil {
		t.Fatal(err)
	}
	record, err := svc.CreateRecord("Lin & Kai", "2026-08-24", "Jasmine Hall", "")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = svc.ValidateRecord(record.ID); err != nil {
		t.Fatal(err)
	}
	record, err = svc.ConfirmRecord(record.ID, "tester")
	if err != nil || record.Status != model.StatusConfirmed {
		t.Fatalf("confirmation failed: %+v %v", record, err)
	}
	view, err := svc.BuildView(record.ID, profile.ID)
	if err != nil || !view.Layout.Valid() || !view.Background.Fallback {
		t.Fatalf("view failed: %+v %v", view, err)
	}
}
