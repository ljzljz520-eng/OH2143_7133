package store

import (
	"sort"
	"strings"

	"wedding-sign/internal/model"
)

func (s *Store) ListRecords() ([]model.Record, error) {
	records := make([]model.Record, 0)
	err := s.list("records", func(raw []byte) error {
		var record model.Record
		if err := decode(raw, &record); err != nil {
			return err
		}
		records = append(records, record)
		return nil
	})
	sort.Slice(records, func(i, j int) bool {
		if records[i].Date == records[j].Date {
			return records[i].DisplayOrder < records[j].DisplayOrder
		}
		return records[i].Date < records[j].Date
	})
	return records, err
}

func (s *Store) FindRecordsByDate(date string) ([]model.Record, error) {
	all, err := s.ListRecords()
	if err != nil {
		return nil, err
	}
	found := make([]model.Record, 0)
	for _, record := range all {
		if record.Date == strings.TrimSpace(date) {
			found = append(found, record)
		}
	}
	return found, nil
}

func (s *Store) DeleteRecord(id string) error { return s.remove("records", id) }

func (s *Store) ListBatches() ([]model.Batch, error) {
	batches := make([]model.Batch, 0)
	err := s.list("batches", func(raw []byte) error {
		var batch model.Batch
		if err := decode(raw, &batch); err != nil {
			return err
		}
		batches = append(batches, batch)
		return nil
	})
	sort.Slice(batches, func(i, j int) bool { return batches[i].EventDate < batches[j].EventDate })
	return batches, err
}

func (s *Store) ListAudits(recordID string) ([]model.Audit, error) {
	audits := make([]model.Audit, 0)
	err := s.list("audits", func(raw []byte) error {
		var audit model.Audit
		if err := decode(raw, &audit); err != nil {
			return err
		}
		if strings.TrimSpace(recordID) == "" || audit.RecordID == recordID {
			audits = append(audits, audit)
		}
		return nil
	})
	sort.Slice(audits, func(i, j int) bool { return audits[i].At.Before(audits[j].At) })
	return audits, err
}
