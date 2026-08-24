package display

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"strings"

	"wedding-sign/internal/model"
)

type View struct {
	Background Background `json:"background"`
	Layout     Layout     `json:"layout"`
	RecordID   string     `json:"record_id"`
	ProfileID  string     `json:"profile_id"`
}

type Renderer struct{ template *template.Template }

func NewRenderer() *Renderer {
	return &Renderer{template: template.Must(template.New("welcome").Parse(`<!doctype html><html lang="en"><head><meta charset="utf-8"><title>{{.Layout.Title}}</title></head><body data-fit="{{.Background.Fit}}" data-dim="{{.Background.Overlay}}"><main><p>{{.Layout.WelcomeText}}</p><h1>{{.Layout.Title}}</h1><p>{{.Layout.DateLine}}</p><p>{{.Layout.VenueLine}}</p><button>{{.Layout.StartLabel}}</button></main></body></html>`))}
}

func (r *Renderer) View(record model.Record, profile model.Profile) View {
	return View{Background: ResolveBackground(record, profile), Layout: BuildLayout(record, profile), RecordID: record.ID, ProfileID: profile.ID}
}

func (r *Renderer) RenderHTML(writer io.Writer, record model.Record, profile model.Profile) error {
	if r == nil || r.template == nil {
		return fmt.Errorf("renderer is not initialized")
	}
	return r.template.Execute(writer, r.View(record, profile))
}

func (r *Renderer) RenderJSON(record model.Record, profile model.Profile) ([]byte, error) {
	view := r.View(record, profile)
	return json.MarshalIndent(view, "", "  ")
}

func (v View) SearchText() string {
	return strings.Join([]string{v.Layout.Title, v.Layout.DateLine, v.Layout.VenueLine, v.Layout.WelcomeText}, " ")
}
