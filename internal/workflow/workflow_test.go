package workflow

import (
	"path/filepath"
	"testing"

	"wedding-sign/internal/config"
	"wedding-sign/internal/service"
	"wedding-sign/internal/store"
)

func workflowService(t *testing.T) *service.Service {
	t.Helper()
	settings := config.DefaultSettings(t.TempDir())
	db, err := store.Open(filepath.Join(t.TempDir(), "workflow.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	svc, err := service.New(db, settings)
	if err != nil {
		t.Fatal(err)
	}
	return svc
}

func TestWorkflowAccept(t *testing.T) {
	result, err := AcceptWelcome(workflowService(t), "ceremony", "Lin & Kai", "2026-08-24", "Jasmine Hall", "")
	if err != nil || !result.ViewReady || result.Record.Status != "confirmed" {
		t.Fatalf("accept failed: %+v %v", result, err)
	}
}

func TestWorkflowPublish(t *testing.T) {
	svc := workflowService(t)
	accepted, err := AcceptWelcome(svc, "ceremony", "Lin & Kai", "2026-08-24", "Jasmine Hall", "")
	if err != nil {
		t.Fatal(err)
	}
	result, err := PublishWelcome(svc, "publish-batch", "2026-08-24", []string{accepted.Record.ID})
	if err != nil || !result.Ready || result.Batch.Status != "dispatched" {
		t.Fatalf("publish failed: %+v %v", result, err)
	}
}

func TestWorkflowReopen(t *testing.T) {
	svc := workflowService(t)
	accepted, err := AcceptWelcome(svc, "ceremony", "Lin & Kai", "2026-08-24", "Jasmine Hall", "")
	if err != nil {
		t.Fatal(err)
	}
	session, err := OpenDisplaySession(svc, "session-1", accepted.Record.ID, accepted.Profile.ID)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = UpdateDisplaySession(svc, session.ID, "welcome"); err != nil {
		t.Fatal(err)
	}
	session, err = CloseDisplaySession(svc, session.ID)
	if err != nil || session.State != "closed" {
		t.Fatalf("close failed: %+v %v", session, err)
	}
}
