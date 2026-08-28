package workflow

import (
	"fmt"
	"strings"

	"wedding-sign/internal/model"
	"wedding-sign/internal/service"
)

type PublishResult struct {
	Batch   model.Batch
	Records []model.Record
	Ready   bool
}

func PublishWelcome(svc *service.Service, batchID, date string, recordIDs []string) (PublishResult, error) {
	batch, err := svc.CreateBatch(batchID, date, recordIDs)
	if err != nil {
		return PublishResult{}, err
	}
	records := make([]model.Record, 0, len(recordIDs))
	for _, recordID := range recordIDs {
		record, err := svc.QueueRecord(recordID, batch.ID)
		if err != nil {
			return PublishResult{}, err
		}
		records = append(records, record)
	}
	batch, err = svc.DispatchBatch(batch.ID, "publishing-desk")
	if err != nil {
		return PublishResult{}, err
	}
	for index, record := range records {
		updated, getErr := svc.GetRecord(record.ID)
		if getErr != nil {
			return PublishResult{}, getErr
		}
		records[index] = updated
	}
	ready := batch.Status == model.StatusDispatched && len(records) > 0
	return PublishResult{Batch: batch, Records: records, Ready: ready}, nil
}

func ValidatePublishInput(batchID, date string, ids []string) error {
	if strings.TrimSpace(batchID) == "" || strings.TrimSpace(date) == "" {
		return fmt.Errorf("batch id and date are required")
	}
	if len(ids) == 0 {
		return fmt.Errorf("at least one record is required")
	}
	return nil
}
