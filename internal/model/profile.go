package model

import (
	"errors"
	"strings"
	"time"
)

type Profile struct {
	ID            string    `json:"id"`
	DisplayName   string    `json:"display_name"`
	Theme         string    `json:"theme"`
	AccentColor   string    `json:"accent_color"`
	Locale        string    `json:"locale"`
	Fullscreen    bool      `json:"fullscreen"`
	DimBackground bool      `json:"dim_background"`
	FitMode       string    `json:"fit_mode"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func NewProfile(id, name string) Profile {
	return Profile{ID: strings.TrimSpace(id), DisplayName: strings.TrimSpace(name), Theme: "garden", AccentColor: "#c88b6b", Locale: "zh-CN", Fullscreen: true, DimBackground: true, FitMode: "cover", UpdatedAt: time.Now().UTC()}
}

func (p Profile) Validate() error {
	if p.ID == "" || p.DisplayName == "" {
		return errors.New("profile identity is required")
	}
	if p.FitMode != "cover" && p.FitMode != "contain" {
		return errors.New("fit mode must be cover or contain")
	}
	if p.Locale == "" {
		return errors.New("locale is required")
	}
	return nil
}

func (p *Profile) ToggleFullscreen() { p.Fullscreen = !p.Fullscreen; p.UpdatedAt = time.Now().UTC() }

func (p *Profile) SetFitMode(mode string) bool {
	if mode != "cover" && mode != "contain" {
		return false
	}
	p.FitMode = mode
	p.UpdatedAt = time.Now().UTC()
	return true
}

func (p Profile) UsesDarkOverlay() bool { return p.DimBackground }
