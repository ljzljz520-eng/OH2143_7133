package display

import (
	"strings"
	"testing"

	"wedding-sign/internal/model"
)

func TestDefaultAndCustomBackground(t *testing.T) {
	profile := model.NewProfile("p", "Ceremony")
	fallback := ResolveBackground(model.NewRecord("r", "A & B", "2026-08-24", "Hall", ""), profile)
	if !fallback.Fallback || fallback.Fit != "cover" || !strings.Contains(fallback.Overlay, "rgba") {
		t.Fatalf("unexpected fallback: %+v", fallback)
	}
	custom := ResolveBackground(model.NewRecord("r2", "A & B", "2026-08-24", "Hall", "photos/hero.jpg"), profile)
	if custom.Fallback || !custom.IsAccessible() {
		t.Fatalf("unexpected custom: %+v", custom)
	}
}

func TestRendererContainsWelcomeFields(t *testing.T) {
	renderer := NewRenderer()
	record := model.NewRecord("r", "Lin & Kai", "2026-08-24", "Jasmine Hall", "")
	profile := model.NewProfile("p", "Ceremony")
	raw, err := renderer.RenderJSON(record, profile)
	if err != nil || !strings.Contains(string(raw), "Jasmine Hall") {
		t.Fatalf("render failed: %s %v", raw, err)
	}
}
