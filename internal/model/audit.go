package model

import (
	"strings"
	"time"
)

type Audit struct {
	ID        string    `json:"id"`
	RecordID  string    `json:"record_id"`
	BatchID   string    `json:"batch_id"`
	Action    string    `json:"action"`
	Actor     string    `json:"actor"`
	Before    string    `json:"before"`
	After     string    `json:"after"`
	RequestID string    `json:"request_id"`
	At        time.Time `json:"at"`
}

func NewAudit(id, recordID, batchID, action, actor, before, after, requestID string) Audit {
	return Audit{ID: strings.TrimSpace(id), RecordID: strings.TrimSpace(recordID), BatchID: strings.TrimSpace(batchID), Action: strings.TrimSpace(action), Actor: strings.TrimSpace(actor), Before: before, After: after, RequestID: strings.TrimSpace(requestID), At: time.Now().UTC()}
}

func (a Audit) IsIdempotentWith(other Audit) bool {
	return a.RecordID == other.RecordID && a.Action == other.Action && a.RequestID != "" && a.RequestID == other.RequestID
}

func (a Audit) Summary() string {
	parts := []string{a.Action, a.RecordID}
	if a.BatchID != "" {
		parts = append(parts, a.BatchID)
	}
	if a.Actor != "" {
		parts = append(parts, "by "+a.Actor)
	}
	return strings.Join(parts, " ")
}
