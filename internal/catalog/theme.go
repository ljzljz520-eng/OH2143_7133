package catalog

import (
	"fmt"
	"math"
	"strings"
)

type ThemePalette struct {
	Name       string            `json:"name"`
	Background string            `json:"background"`
	Surface    string            `json:"surface"`
	Text       string            `json:"text"`
	Muted      string            `json:"muted"`
	Accent     string            `json:"accent"`
	Variables  map[string]string `json:"variables"`
}

func Theme(name string) ThemePalette {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "night":
		return ThemePalette{Name: "night", Background: "#151a20", Surface: "#25313a", Text: "#fffaf2", Muted: "#c4cbd0", Accent: "#e0a978"}
	case "minimal":
		return ThemePalette{Name: "minimal", Background: "#f6f2ec", Surface: "#ffffff", Text: "#282521", Muted: "#726a61", Accent: "#9f694f"}
	default:
		return ThemePalette{Name: "garden", Background: "#273a35", Surface: "#3c554c", Text: "#fffdf8", Muted: "#d8e1da", Accent: "#d7a078"}
	}
}

func (p ThemePalette) WithAccent(accent string) ThemePalette {
	if strings.TrimSpace(accent) != "" {
		p.Accent = strings.TrimSpace(accent)
	}
	return p
}

func (p ThemePalette) CSSVariables() map[string]string {
	variables := map[string]string{"--background": p.Background, "--surface": p.Surface, "--text": p.Text, "--muted": p.Muted, "--accent": p.Accent}
	p.Variables = variables
	return variables
}

func (p ThemePalette) ContrastScore() float64 {
	return luminance(p.Text) / math.Max(luminance(p.Background), 0.001)
}

func luminance(hex string) float64 {
	value := strings.TrimPrefix(strings.TrimSpace(hex), "#")
	if len(value) != 6 {
		return 1
	}
	total := 0.0
	for index := 0; index < 6; index += 2 {
		var component int
		if _, err := fmt.Sscanf(value[index:index+2], "%02x", &component); err != nil {
			return 1
		}
		total += float64(component) / 255
	}
	return total / 3
}

func (p ThemePalette) Readable() bool { return p.ContrastScore() > 2.7 }
