package workflow

import (
	"fmt"
	"strings"

	"wedding-sign/internal/catalog"
	"wedding-sign/internal/model"
	"wedding-sign/internal/service"
	"wedding-sign/internal/store"
)

type ChecklistItem struct {
	Name   string
	Passed bool
	Detail string
}

type Checklist struct{ Items []ChecklistItem }

func BuildChecklist(svc *service.Service, recordID, profileID string) (Checklist, error) {
	record, err := svc.GetRecord(recordID)
	if err != nil {
		return Checklist{}, err
	}
	profile, err := svc.Store().GetProfile(profileID)
	if err != nil {
		return Checklist{}, err
	}
	check := catalog.CheckRecord(record, profile)
	items := []ChecklistItem{{Name: "record-valid", Passed: record.Validate() == nil, Detail: record.Status}, {Name: "display-readable", Passed: !catalog.HasErrors(check), Detail: catalog.Summarize(check)}, {Name: "background-ready", Passed: record.HasHeroImage() || record.HeroImage == "", Detail: "fallback available"}}
	return Checklist{Items: items}, nil
}

func (c Checklist) Passed() bool {
	for _, item := range c.Items {
		if !item.Passed {
			return false
		}
	}
	return len(c.Items) > 0
}

func (c Checklist) Failures() []ChecklistItem {
	failures := make([]ChecklistItem, 0)
	for _, item := range c.Items {
		if !item.Passed {
			failures = append(failures, item)
		}
	}
	return failures
}

func AdvanceBatch(svc *service.Service, batchID string) (model.Batch, string, error) {
	batch, err := svc.GetBatch(batchID)
	if err != nil {
		return model.Batch{}, "", err
	}
	next, ok := batch.Next()
	if !ok {
		return batch, "complete", nil
	}
	if err := svc.Store().PutBatch(batch); err != nil {
		return model.Batch{}, "", err
	}
	return batch, next, nil
}

func DispatchDay(svc *service.Service, date string, actor string) (int, error) {
	records, err := svc.Search(store.RecordFilter{Date: strings.TrimSpace(date), Status: model.StatusConfirmed})
	if err != nil {
		return 0, err
	}
	if len(records) == 0 {
		return 0, fmt.Errorf("no confirmed records for %s", date)
	}
	ids := make([]string, 0, len(records))
	for _, record := range records {
		ids = append(ids, record.ID)
	}
	batch, err := svc.CreateBatch("day-"+strings.ReplaceAll(date, "-", ""), date, ids)
	if err != nil {
		return 0, err
	}
	if _, err := svc.DispatchBatch(batch.ID, actor); err != nil {
		return 0, err
	}
	return len(ids), nil
}
