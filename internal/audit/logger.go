package audit

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"wedding-sign/internal/model"
	"wedding-sign/internal/store"
)

type Logger struct {
	store *store.Store
	mu    sync.Mutex
	clock func() time.Time
}

func NewLogger(s *store.Store) *Logger {
	return &Logger{store: s, clock: func() time.Time { return time.Now().UTC() }}
}

func (l *Logger) Record(recordID, batchID, action, actor, before, after, requestID string) (model.Audit, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.store == nil {
		return model.Audit{}, fmt.Errorf("audit store is missing")
	}
	audit := model.NewAudit(fmt.Sprintf("audit-%d", l.clock().UnixNano()), recordID, batchID, action, actor, before, after, requestID)
	audit.At = l.clock()
	if err := l.store.PutAudit(audit); err != nil {
		return model.Audit{}, err
	}
	return audit, nil
}

func (l *Logger) FindByRequest(recordID, requestID string) (model.Audit, bool, error) {
	audits, err := l.store.ListAudits(recordID)
	if err != nil {
		return model.Audit{}, false, err
	}
	for _, item := range audits {
		if item.RequestID == strings.TrimSpace(requestID) && requestID != "" {
			return item, true, nil
		}
	}
	return model.Audit{}, false, nil
}

func (l *Logger) History(recordID string) ([]model.Audit, error) {
	audits, err := l.store.ListAudits(recordID)
	if err != nil {
		return nil, err
	}
	sort.SliceStable(audits, func(i, j int) bool { return audits[i].At.Before(audits[j].At) })
	return audits, nil
}
