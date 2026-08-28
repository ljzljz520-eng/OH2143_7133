package catalog

import (
	"regexp"
	"strings"

	"wedding-sign/internal/model"
)

var hexColor = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)

type ValidationSummary struct {
	Valid    bool     `json:"valid"`
	Errors   []string `json:"errors"`
	Warnings []string `json:"warnings"`
}

func ValidateProfile(profile model.Profile) ValidationSummary {
	summary := ValidationSummary{Errors: make([]string, 0), Warnings: make([]string, 0)}
	if err := profile.Validate(); err != nil {
		summary.Errors = append(summary.Errors, err.Error())
	}
	if !hexColor.MatchString(profile.AccentColor) {
		summary.Errors = append(summary.Errors, "accent color must be a six-digit hex color")
	}
	if profile.Theme == "" {
		summary.Errors = append(summary.Errors, "theme is required")
	}
	if !profile.Fullscreen {
		summary.Warnings = append(summary.Warnings, "fullscreen is disabled")
	}
	summary.Valid = len(summary.Errors) == 0
	return summary
}

func ValidateRecordAndProfile(record model.Record, profile model.Profile) ValidationSummary {
	summary := ValidationSummary{Errors: make([]string, 0), Warnings: make([]string, 0)}
	if err := record.Validate(); err != nil {
		summary.Errors = append(summary.Errors, err.Error())
	}
	profileSummary := ValidateProfile(profile)
	summary.Errors = append(summary.Errors, profileSummary.Errors...)
	summary.Warnings = append(summary.Warnings, profileSummary.Warnings...)
	if !record.HasHeroImage() {
		summary.Warnings = append(summary.Warnings, "default background will be used")
	}
	summary.Valid = len(summary.Errors) == 0
	return summary
}

func (s ValidationSummary) Message() string {
	parts := make([]string, 0, len(s.Errors)+len(s.Warnings))
	for _, item := range s.Errors {
		parts = append(parts, "error: "+item)
	}
	for _, item := range s.Warnings {
		parts = append(parts, "warning: "+item)
	}
	if len(parts) == 0 {
		return "valid"
	}
	return strings.Join(parts, "; ")
}

func (s ValidationSummary) HasWarnings() bool { return len(s.Warnings) > 0 }
