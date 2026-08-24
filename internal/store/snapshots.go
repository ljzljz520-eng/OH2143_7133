package store

import (
	"encoding/json"
	"fmt"
	"time"

	"wedding-sign/internal/model"
)

type Snapshot struct {
	CapturedAt time.Time       `json:"captured_at"`
	Records    []model.Record  `json:"records"`
	Batches    []model.Batch   `json:"batches"`
	Profiles   []model.Profile `json:"profiles"`
}

func (s *Store) Snapshot() (Snapshot, error) {
	records, err := s.ListRecords()
	if err != nil {
		return Snapshot{}, err
	}
	batches, err := s.ListBatches()
	if err != nil {
		return Snapshot{}, err
	}
	profiles, err := s.ListProfiles()
	if err != nil {
		return Snapshot{}, err
	}
	return Snapshot{CapturedAt: time.Now().UTC(), Records: records, Batches: batches, Profiles: profiles}, nil
}

func (s *Store) ExportJSON() ([]byte, error) {
	snapshot, err := s.Snapshot()
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(snapshot, "", "  ")
}

func (s *Store) ImportJSON(raw []byte) error {
	var snapshot Snapshot
	if err := json.Unmarshal(raw, &snapshot); err != nil {
		return fmt.Errorf("decode snapshot: %w", err)
	}
	for _, record := range snapshot.Records {
		if err := s.PutRecord(record); err != nil {
			return err
		}
	}
	for _, batch := range snapshot.Batches {
		if err := s.PutBatch(batch); err != nil {
			return err
		}
	}
	for _, profile := range snapshot.Profiles {
		if err := s.PutProfile(profile); err != nil {
			return err
		}
	}
	return nil
}
