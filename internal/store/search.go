package store

import (
	"strings"

	"wedding-sign/internal/model"
)

type RecordFilter struct {
	Date         string
	Status       string
	Venue        string
	Query        string
	RequireImage bool
}

func (s *Store) SearchRecords(filter RecordFilter) ([]model.Record, error) {
	records, err := s.ListRecords()
	if err != nil {
		return nil, err
	}
	result := make([]model.Record, 0, len(records))
	for _, record := range records {
		if filter.Date != "" && record.Date != filter.Date {
			continue
		}
		if filter.Status != "" && record.Status != filter.Status {
			continue
		}
		if filter.Venue != "" && !strings.Contains(strings.ToLower(record.Venue), strings.ToLower(filter.Venue)) {
			continue
		}
		if filter.RequireImage && !record.HasHeroImage() {
			continue
		}
		if filter.Query != "" && !containsRecordText(record, filter.Query) {
			continue
		}
		result = append(result, record)
	}
	return result, nil
}

func containsRecordText(record model.Record, query string) bool {
	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" {
		return true
	}
	fields := []string{record.ID, record.CoupleNames, record.Date, record.Venue, record.WelcomeText}
	for _, field := range fields {
		if strings.Contains(strings.ToLower(field), query) {
			return true
		}
	}
	return false
}

func GroupRecordsByStatus(records []model.Record) map[string][]model.Record {
	groups := make(map[string][]model.Record)
	for _, record := range records {
		groups[record.Status] = append(groups[record.Status], record)
	}
	return groups
}

func CountStatus(records []model.Record, status string) int {
	count := 0
	for _, record := range records {
		if record.Status == status {
			count++
		}
	}
	return count
}
