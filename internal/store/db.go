package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	bolt "go.etcd.io/bbolt"
	"wedding-sign/internal/model"
)

var bucketNames = map[string][]byte{
	"records":  []byte("records"),
	"batches":  []byte("batches"),
	"audits":   []byte("audits"),
	"profiles": []byte("profiles"),
	"sessions": []byte("sessions"),
}

type Store struct {
	path string
	db   *bolt.DB
	mu   sync.RWMutex
}

func Open(path string) (*Store, error) {
	if path == "" {
		return nil, errors.New("store path is required")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	db, err := bolt.Open(path, 0o600, &bolt.Options{Timeout: time.Second})
	if err != nil {
		return nil, fmt.Errorf("open store: %w", err)
	}
	s := &Store{path: path, db: db}
	if err := s.initBuckets(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) initBuckets() error {
	return s.db.Update(func(tx *bolt.Tx) error {
		for _, name := range bucketNames {
			if _, err := tx.CreateBucketIfNotExists(name); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db == nil {
		return nil
	}
	err := s.db.Close()
	s.db = nil
	return err
}

func (s *Store) Path() string { return s.path }

func (s *Store) ensureOpen() error {
	if s == nil || s.db == nil {
		return errors.New("store is closed")
	}
	return nil
}

func encode(value any) ([]byte, error) { return json.Marshal(value) }

func decode(raw []byte, value any) error {
	if len(raw) == 0 {
		return errors.New("stored value is empty")
	}
	return json.Unmarshal(raw, value)
}

func (s *Store) put(bucket, key string, value any) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if err := s.ensureOpen(); err != nil {
		return err
	}
	raw, err := encode(value)
	if err != nil {
		return err
	}
	name, ok := bucketNames[bucket]
	if !ok {
		return fmt.Errorf("unknown bucket %s", bucket)
	}
	return s.db.Update(func(tx *bolt.Tx) error { return tx.Bucket(name).Put([]byte(key), raw) })
}

func (s *Store) get(bucket, key string, value any) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if err := s.ensureOpen(); err != nil {
		return err
	}
	name, ok := bucketNames[bucket]
	if !ok {
		return fmt.Errorf("unknown bucket %s", bucket)
	}
	return s.db.View(func(tx *bolt.Tx) error {
		raw := tx.Bucket(name).Get([]byte(key))
		if raw == nil {
			return fmt.Errorf("%s %s not found", bucket, key)
		}
		copyRaw := append([]byte(nil), raw...)
		return decode(copyRaw, value)
	})
}

func (s *Store) list(bucket string, decodeValue func([]byte) error) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if err := s.ensureOpen(); err != nil {
		return err
	}
	name, ok := bucketNames[bucket]
	if !ok {
		return fmt.Errorf("unknown bucket %s", bucket)
	}
	return s.db.View(func(tx *bolt.Tx) error {
		return tx.Bucket(name).ForEach(func(_, raw []byte) error {
			if raw == nil {
				return nil
			}
			return decodeValue(append([]byte(nil), raw...))
		})
	})
}

func (s *Store) remove(bucket, key string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if err := s.ensureOpen(); err != nil {
		return err
	}
	name, ok := bucketNames[bucket]
	if !ok {
		return fmt.Errorf("unknown bucket %s", bucket)
	}
	return s.db.Update(func(tx *bolt.Tx) error { return tx.Bucket(name).Delete([]byte(key)) })
}

func (s *Store) PutRecord(record model.Record) error { return s.put("records", record.ID, record) }

func (s *Store) GetRecord(id string) (model.Record, error) {
	var record model.Record
	err := s.get("records", id, &record)
	return record, err
}

func (s *Store) PutBatch(batch model.Batch) error { return s.put("batches", batch.ID, batch) }

func (s *Store) GetBatch(id string) (model.Batch, error) {
	var batch model.Batch
	err := s.get("batches", id, &batch)
	return batch, err
}

func (s *Store) PutAudit(audit model.Audit) error { return s.put("audits", audit.ID, audit) }

func (s *Store) GetAudit(id string) (model.Audit, error) {
	var audit model.Audit
	err := s.get("audits", id, &audit)
	return audit, err
}

func (s *Store) PutProfile(profile model.Profile) error {
	return s.put("profiles", profile.ID, profile)
}

func (s *Store) GetProfile(id string) (model.Profile, error) {
	var profile model.Profile
	err := s.get("profiles", id, &profile)
	return profile, err
}

func (s *Store) PutSession(session model.Session) error {
	return s.put("sessions", session.ID, session)
}

func (s *Store) GetSession(id string) (model.Session, error) {
	var session model.Session
	err := s.get("sessions", id, &session)
	return session, err
}
