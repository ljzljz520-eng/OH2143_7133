package model

import "testing"

func TestBatchProgression(t *testing.T) {
	batch := NewBatch("b1", "2026-08-24")
	if !batch.AddRecord("r1") || batch.AddRecord("r1") {
		t.Fatal("record should be added once")
	}
	if err := batch.Validate(); err != nil {
		t.Fatal(err)
	}
	id, ok := batch.Next()
	if !ok || id != "r1" || !batch.IsComplete() || batch.Progress() != 1 {
		t.Fatalf("unexpected progression: %+v", batch)
	}
}
