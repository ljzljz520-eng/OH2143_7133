package model

import "testing"

func TestRecordValidationAndNormalization(t *testing.T) {
	record := NewRecord(" r1 ", "  Lin   & Kai ", "2026-08-24", " Jasmine Hall ", "")
	record.Normalize()
	if err := record.Validate(); err != nil {
		t.Fatal(err)
	}
	if record.CoupleNames != "Lin & Kai" || record.HasHeroImage() {
		t.Fatalf("unexpected normalized record: %+v", record)
	}
}

func TestRecordMutability(t *testing.T) {
	record := NewRecord("r1", "A & B", "2026-08-24", "Hall", "hero.jpg")
	if !record.IsMutable() || record.CanQueue() {
		t.Fatal("draft should be mutable but not queueable")
	}
	record.Status = StatusConfirmed
	if !record.CanQueue() || !record.CanPublish() {
		t.Fatal("confirmed record should be queueable")
	}
}
