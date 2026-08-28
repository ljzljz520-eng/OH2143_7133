package service

import "testing"

func TestWorkflow24(t *testing.T) {
	svc := testService(t)
	_, err := svc.CreateProfile("p24", "Ceremony")
	if err != nil {
		t.Fatal(err)
	}
	record, err := svc.CreateRecord("Lin & Kai", "2026-08-24", "Jasmine Hall", "hero.jpg")
	if err != nil {
		t.Fatal(err)
	}
	record, err = svc.ValidateRecord(record.ID)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = svc.ConfirmRecord(record.ID, "desk"); err != nil {
		t.Fatal(err)
	}
	batch, err := svc.CreateBatch("batch-24", record.Date, []string{record.ID})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = svc.DispatchRecord(record.ID, batch.ID, "dispatch-24"); err != nil {
		t.Fatal(err)
	}
	if _, err = svc.DispatchRecord(record.ID, batch.ID, "dispatch-24"); err != nil {
		t.Fatal(err)
	}
	count, err := svc.DispatchAuditCount(record.ID)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("dispatch should be idempotent, got %d audit entries", count)
	}
}
