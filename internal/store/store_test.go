package store

import (
	"path/filepath"
	"testing"

	"wedding-sign/internal/model"
)

func TestStoreRoundTrip(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "sign.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	record := model.NewRecord("r1", "A & B", "2026-08-24", "Hall", "")
	if err := db.PutRecord(record); err != nil {
		t.Fatal(err)
	}
	got, err := db.GetRecord("r1")
	if err != nil {
		t.Fatal(err)
	}
	if got.CoupleNames != record.CoupleNames {
		t.Fatalf("got %+v", got)
	}
}
