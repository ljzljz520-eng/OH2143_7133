package catalog

import (
	"strings"

	"wedding-sign/internal/model"
)

type AccessibilityIssue struct {
	Field    string `json:"field"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

func CheckRecord(record model.Record, profile model.Profile) []AccessibilityIssue {
	issues := make([]AccessibilityIssue, 0)
	if strings.TrimSpace(record.CoupleNames) == "" {
		issues = append(issues, AccessibilityIssue{Field: "couple_names", Severity: "error", Message: "Names must be visible"})
	}
	if strings.TrimSpace(record.Venue) == "" {
		issues = append(issues, AccessibilityIssue{Field: "venue", Severity: "error", Message: "Venue must be visible"})
	}
	if !profile.DimBackground && record.HasHeroImage() {
		issues = append(issues, AccessibilityIssue{Field: "background", Severity: "warning", Message: "A dim overlay improves text contrast"})
	}
	if !profile.Fullscreen {
		issues = append(issues, AccessibilityIssue{Field: "fullscreen", Severity: "warning", Message: "Fullscreen helps distant guests read the sign"})
	}
	return issues
}

func HasErrors(issues []AccessibilityIssue) bool {
	for _, issue := range issues {
		if issue.Severity == "error" {
			return true
		}
	}
	return false
}

func Summarize(issues []AccessibilityIssue) string {
	if len(issues) == 0 {
		return "No accessibility issues"
	}
	fields := make([]string, 0, len(issues))
	for _, issue := range issues {
		fields = append(fields, issue.Field+": "+issue.Message)
	}
	return strings.Join(fields, "; ")
}
