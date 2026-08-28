package service

import (
	"fmt"
	"strings"

	"wedding-sign/internal/model"
)

func (s *Service) ConfirmRecord(id, actor string) (model.Record, error) {
	record, err := s.store.GetRecord(strings.TrimSpace(id))
	if err != nil {
		return model.Record{}, err
	}
	if record.Status != model.StatusValidated && record.Status != model.StatusDraft {
		return model.Record{}, fmt.Errorf("record %s cannot be confirmed from %s", id, record.Status)
	}
	before := record.Status
	record.Status = model.StatusConfirmed
	record.Normalize()
	if err := s.store.PutRecord(record); err != nil {
		return model.Record{}, err
	}
	_, err = s.logger.Record(record.ID, "", "confirm", actor, before, record.Status, s.nextID("request"))
	return record, err
}

func (s *Service) QueueRecord(id, batchID string) (model.Record, error) {
	record, err := s.store.GetRecord(strings.TrimSpace(id))
	if err != nil {
		return model.Record{}, err
	}
	if !record.CanQueue() {
		return model.Record{}, fmt.Errorf("record %s is not confirmed", id)
	}
	before := record.Status
	record.Status = model.StatusQueued
	record.Normalize()
	if err := s.store.PutRecord(record); err != nil {
		return model.Record{}, err
	}
	_, err = s.logger.Record(record.ID, batchID, "queue", "dispatcher", before, record.Status, s.nextID("request"))
	return record, err
}

func (s *Service) ArchiveRecord(id string) (model.Record, error) {
	record, err := s.store.GetRecord(strings.TrimSpace(id))
	if err != nil {
		return model.Record{}, err
	}
	if record.Status != model.StatusPublished && record.Status != model.StatusDispatched {
		return model.Record{}, fmt.Errorf("record %s is not ready to archive", id)
	}
	record.Status = model.StatusArchived
	record.Normalize()
	if err := s.store.PutRecord(record); err != nil {
		return model.Record{}, err
	}
	_, err = s.logger.Record(record.ID, "", "archive", "operator", model.StatusPublished, record.Status, s.nextID("request"))
	return record, err
}

func (s *Service) CreateBatch(id, eventDate string, recordIDs []string) (model.Batch, error) {
	batch := model.NewBatch(id, eventDate)
	for _, recordID := range recordIDs {
		if !batch.AddRecord(recordID) {
			return model.Batch{}, fmt.Errorf("invalid or duplicate record %s", recordID)
		}
	}
	if err := batch.Validate(); err != nil {
		return model.Batch{}, err
	}
	if err := s.store.PutBatch(batch); err != nil {
		return model.Batch{}, err
	}
	return batch, nil
}

func (s *Service) GetBatch(id string) (model.Batch, error) {
	return s.store.GetBatch(strings.TrimSpace(id))
}
