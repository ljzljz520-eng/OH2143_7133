package display

import (
	"fmt"
	"strings"

	"wedding-sign/internal/model"
)

type Layout struct {
	Title       string `json:"title"`
	DateLine    string `json:"date_line"`
	VenueLine   string `json:"venue_line"`
	WelcomeText string `json:"welcome_text"`
	StartLabel  string `json:"start_label"`
	Fullscreen  bool   `json:"fullscreen"`
	Direction   string `json:"direction"`
}

func BuildLayout(record model.Record, profile model.Profile) Layout {
	direction := "ltr"
	if profile.Locale == "ar" || profile.Locale == "he" {
		direction = "rtl"
	}
	dateLine := record.Date
	if dateLine == "" {
		dateLine = "Date to be announced"
	}
	venueLine := record.Venue
	if venueLine == "" {
		venueLine = "Venue to be announced"
	}
	welcome := strings.TrimSpace(record.WelcomeText)
	if welcome == "" {
		welcome = "Welcome to our celebration"
	}
	label := "Begin celebration"
	if profile.Locale == "zh-CN" {
		label = "开始"
	}
	return Layout{Title: record.DisplayTitle(), DateLine: dateLine, VenueLine: venueLine, WelcomeText: welcome, StartLabel: label, Fullscreen: profile.Fullscreen, Direction: direction}
}

func (l Layout) Lines() []string {
	return []string{l.Title, l.DateLine, l.VenueLine, l.WelcomeText, l.StartLabel}
}

func (l Layout) AccessibleLabel() string {
	return fmt.Sprintf("%s, %s, %s", l.Title, l.DateLine, l.VenueLine)
}

func (l Layout) Valid() bool {
	return l.Title != "" && l.DateLine != "" && l.VenueLine != "" && l.StartLabel != ""
}
