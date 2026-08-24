package service

import (
	"fmt"
	"strings"

	"wedding-sign/internal/model"
)

type QueueReport struct {
	BatchID   string
	Total     int
	Processed int
	Rejected  int
	Current   string
	Complete  bool
}

func (s *Service) ProcessBatch(batchID string) (QueueReport, error) {
	batch, err := s.store.GetBatch(strings.TrimSpace(batchID))
	if err != nil {
		return QueueReport{}, err
	}
	report := QueueReport{BatchID: batch.ID, Total: len(batch.RecordIDs), Current: ""}
	for _, recordID := range batch.RecordIDs {
		record, getErr := s.store.GetRecord(recordID)
		if getErr != nil {
			report.Rejected++
			continue
		}
		if record.Status == model.StatusQueued || record.Status == model.StatusConfirmed {
			record.Status = model.StatusPublished
			record.Normalize()
			if putErr := s.store.PutRecord(record); putErr != nil {
				return QueueReport{}, putErr
			}
			report.Processed++
			report.Current = record.ID
			if _, auditErr := s.logger.Record(record.ID, batch.ID, "publish", "queue", model.StatusQueued, model.StatusPublished, s.nextID("request")); auditErr != nil {
				return QueueReport{}, auditErr
			}
		} else {
			report.Rejected++
		}
	}
	batch.Status = model.StatusPublished
	if report.Processed == report.Total && report.Total > 0 {
		report.Complete = true
	}
	if err := s.store.PutBatch(batch); err != nil {
		return QueueReport{}, err
	}
	return report, nil
}

func (s *Service) ReviewBatch(batchID string) (QueueReport, error) {
	batch, err := s.store.GetBatch(batchID)
	if err != nil {
		return QueueReport{}, err
	}
	report := QueueReport{BatchID: batch.ID, Total: len(batch.RecordIDs)}
	for _, recordID := range batch.RecordIDs {
		record, getErr := s.store.GetRecord(recordID)
		if getErr != nil {
			report.Rejected++
			continue
		}
		switch record.Status {
		case model.StatusPublished, model.StatusDispatched:
			report.Processed++
		case model.StatusQueued, model.StatusConfirmed:
			report.Current = record.ID
		default:
			report.Rejected++
		}
	}
	report.Complete = report.Processed == report.Total && report.Total > 0
	return report, nil
}

func (s *Service) PublishRecord(id string) (model.Record, error) {
	record, err := s.store.GetRecord(id)
	if err != nil {
		return model.Record{}, err
	}
	if !record.CanPublish() {
		return model.Record{}, fmt.Errorf("record %s cannot publish from %s", id, record.Status)
	}
	before := record.Status
	record.Status = model.StatusPublished
	record.Normalize()
	if err := s.store.PutRecord(record); err != nil {
		return model.Record{}, err
	}
	_, err = s.logger.Record(record.ID, "", "publish", "operator", before, record.Status, s.nextID("request"))
	return record, err
}
