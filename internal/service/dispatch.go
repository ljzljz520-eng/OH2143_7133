package service

import (
	"fmt"
	"strings"

	"wedding-sign/internal/model"
)

func (s *Service) DispatchRecord(recordID, batchID, requestID string) (model.Record, error) {
	record, err := s.store.GetRecord(strings.TrimSpace(recordID))
	if err != nil {
		return model.Record{}, err
	}
	batch, err := s.store.GetBatch(strings.TrimSpace(batchID))
	if err != nil {
		return model.Record{}, err
	}
	if !contains(batch.RecordIDs, record.ID) {
		return model.Record{}, fmt.Errorf("record %s is not in batch %s", record.ID, batch.ID)
	}
	if record.Status != model.StatusConfirmed && record.Status != model.StatusQueued && record.Status != model.StatusDispatched {
		return model.Record{}, fmt.Errorf("record %s cannot be dispatched from %s", record.ID, record.Status)
	}
	before := record.Status
	record.Status = model.StatusDispatched
	record.Normalize()
	if err := s.store.PutRecord(record); err != nil {
		return model.Record{}, err
	}
	_, err = s.logger.Record(record.ID, batch.ID, "dispatch", "dispatcher", before, record.Status, requestID)
	return record, err
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func (s *Service) DispatchBatch(batchID, actor string) (model.Batch, error) {
	batch, err := s.store.GetBatch(strings.TrimSpace(batchID))
	if err != nil {
		return model.Batch{}, err
	}
	if len(batch.RecordIDs) == 0 {
		return model.Batch{}, fmt.Errorf("batch %s is empty", batch.ID)
	}
	batch.Status = model.StatusQueued
	for _, recordID := range batch.RecordIDs {
		if _, err := s.DispatchRecord(recordID, batch.ID, s.nextID("dispatch")); err != nil {
			return model.Batch{}, err
		}
	}
	batch.Status = model.StatusDispatched
	if err := s.store.PutBatch(batch); err != nil {
		return model.Batch{}, err
	}
	_, err = s.logger.Record("", batch.ID, "dispatch-batch", actor, model.StatusQueued, batch.Status, s.nextID("request"))
	return batch, err
}

func (s *Service) DispatchAuditCount(recordID string) (int, error) {
	audits, err := s.logger.History(recordID)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, item := range audits {
		if item.Action == "dispatch" {
			count++
		}
	}
	return count, nil
}
