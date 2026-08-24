package catalog

import (
	"fmt"
	"html/template"
	"io"
	"strings"

	"wedding-sign/internal/model"
)

type RenderModel struct {
	Title      string
	Date       string
	Venue      string
	Background string
	Fit        string
	Overlay    string
	Theme      ThemePalette
	Message    GuestMessage
}

func BuildRenderModel(record model.Record, profile model.Profile) RenderModel {
	background := "default-garden"
	if record.HasHeroImage() {
		background = record.HeroImage
	}
	overlay := "none"
	if profile.DimBackground {
		overlay = "dim"
	}
	fit := profile.FitMode
	if fit == "" {
		fit = "cover"
	}
	return RenderModel{Title: record.DisplayTitle(), Date: record.Date, Venue: record.Venue, Background: background, Fit: fit, Overlay: overlay, Theme: Theme(profile.Theme).WithAccent(profile.AccentColor), Message: BuildGuestMessage(record, profile)}
}

func (m RenderModel) DataAttributes() map[string]string {
	return map[string]string{"data-background": m.Background, "data-fit": m.Fit, "data-overlay": m.Overlay, "data-theme": m.Theme.Name}
}

func (m RenderModel) Text() string {
	return strings.Join([]string{m.Title, m.Date, m.Venue, m.Message.SearchText()}, " ")
}

func RenderStatic(writer io.Writer, record model.Record, profile model.Profile) error {
	model := BuildRenderModel(record, profile)
	tpl := template.Must(template.New("static").Parse(`<section data-background="{{.Background}}" data-fit="{{.Fit}}" data-overlay="{{.Overlay}}"><p>{{.Message.Greeting}}</p><h1>{{.Title}}</h1><p>{{.Date}}</p><p>{{.Venue}}</p><button>{{.Message.Action}}</button></section>`))
	if err := tpl.Execute(writer, model); err != nil {
		return fmt.Errorf("render static sign: %w", err)
	}
	return nil
}

func (m RenderModel) SearchTerms() []string { return strings.Fields(m.Text()) }
