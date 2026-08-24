package display

import (
	"strings"

	"wedding-sign/internal/catalog"
	"wedding-sign/internal/model"
)

type DisplayCheck struct {
	Issues     []catalog.AccessibilityIssue `json:"issues"`
	GuestCopy  catalog.GuestMessage         `json:"guest_copy"`
	Theme      catalog.ThemePalette         `json:"theme"`
	SearchText string                       `json:"search_text"`
	Ready      bool                         `json:"ready"`
}

func CheckDisplay(record model.Record, profile model.Profile) DisplayCheck {
	theme := catalog.Theme(profile.Theme).WithAccent(profile.AccentColor)
	issues := catalog.CheckRecord(record, profile)
	copy := catalog.BuildGuestMessage(record, profile)
	text := strings.Join([]string{record.DisplayTitle(), record.Date, record.Venue, copy.SearchText()}, " ")
	return DisplayCheck{Issues: issues, GuestCopy: copy, Theme: theme, SearchText: text, Ready: !catalog.HasErrors(issues) && copy.Complete() && theme.Readable()}
}

func (c DisplayCheck) WarningCount() int {
	count := 0
	for _, issue := range c.Issues {
		if issue.Severity == "warning" {
			count++
		}
	}
	return count
}

func (c DisplayCheck) IssueSummary() string { return catalog.Summarize(c.Issues) }
