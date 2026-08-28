package workflow

import (
	"sort"

	"wedding-sign/internal/model"
	"wedding-sign/internal/service"
)

type DaySummary struct {
	Date       string
	Total      int
	Confirmed  int
	Dispatched int
	Published  int
}

func BuildDaySummary(svc *service.Service, date string) (DaySummary, error) {
	records, err := svc.ListRecords(date)
	if err != nil {
		return DaySummary{}, err
	}
	summary := DaySummary{Date: date, Total: len(records)}
	for _, record := range records {
		switch record.Status {
		case model.StatusConfirmed:
			summary.Confirmed++
		case model.StatusDispatched:
			summary.Dispatched++
		case model.StatusPublished:
			summary.Published++
		}
	}
	return summary, nil
}

func SortByDisplayOrder(records []model.Record) []model.Record {
	result := append([]model.Record(nil), records...)
	sort.SliceStable(result, func(i, j int) bool { return result[i].DisplayOrder < result[j].DisplayOrder })
	return result
}
