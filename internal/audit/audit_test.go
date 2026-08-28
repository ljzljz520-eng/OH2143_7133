package audit

import (
	"path/filepath"
	"testing"

	"wedding-sign/internal/model"
	"wedding-sign/internal/store"
)

func TestAuditReport(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "audit.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	logger := NewLogger(db)
	if _, err = logger.Record("r1", "", "create", "operator", "", model.StatusDraft, "req-1"); err != nil {
		t.Fatal(err)
	}
	if _, err = logger.Record("r1", "", "confirm", "operator", model.StatusValidated, model.StatusConfirmed, "req-2"); err != nil {
		t.Fatal(err)
	}
	history, err := logger.History("r1")
	if err != nil {
		t.Fatal(err)
	}
	report := BuildReport(history)
	if !report.HasAction("confirm") || report.RecordID != "r1" {
		t.Fatalf("unexpected report: %+v", report)
	}
}
