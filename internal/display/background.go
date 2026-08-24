package display

import (
	"path/filepath"
	"strings"

	"wedding-sign/internal/model"
)

type Background struct {
	Source      string `json:"source"`
	Fallback    bool   `json:"fallback"`
	Overlay     string `json:"overlay"`
	Fit         string `json:"fit"`
	Description string `json:"description"`
}

func ResolveBackground(record model.Record, profile model.Profile) Background {
	source := strings.TrimSpace(record.HeroImage)
	fallback := source == ""
	if fallback {
		source = "default-garden"
	} else {
		source = filepath.Clean(source)
	}
	overlay := "none"
	if profile.DimBackground {
		overlay = "rgba(18, 21, 28, 0.42)"
	}
	fit := profile.FitMode
	if fit != "contain" {
		fit = "cover"
	}
	description := "custom wedding image"
	if fallback {
		description = "elegant default garden background"
	}
	return Background{Source: source, Fallback: fallback, Overlay: overlay, Fit: fit, Description: description}
}

func (b Background) CSS() map[string]string {
	return map[string]string{"background-image": "url('" + b.Source + "')", "background-size": b.Fit, "background-position": "center", "background-overlay": b.Overlay}
}

func (b Background) IsAccessible() bool { return b.Fallback || b.Source != "" }
