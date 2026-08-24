package audit

import (
	"sort"
	"strings"

	"wedding-sign/internal/model"
)

type Report struct {
	RecordID string
	Actions  []string
	Actors   []string
	Last     model.Audit
}

func BuildReport(audits []model.Audit) Report {
	report := Report{Actions: make([]string, 0), Actors: make([]string, 0)}
	if len(audits) == 0 {
		return report
	}
	sort.SliceStable(audits, func(i, j int) bool { return audits[i].At.Before(audits[j].At) })
	report.RecordID = audits[0].RecordID
	report.Last = audits[len(audits)-1]
	seenActions := map[string]bool{}
	seenActors := map[string]bool{}
	for _, audit := range audits {
		if !seenActions[audit.Action] {
			report.Actions = append(report.Actions, audit.Action)
			seenActions[audit.Action] = true
		}
		actor := strings.TrimSpace(audit.Actor)
		if actor != "" && !seenActors[actor] {
			report.Actors = append(report.Actors, actor)
			seenActors[actor] = true
		}
	}
	return report
}

func (r Report) HasAction(action string) bool {
	for _, item := range r.Actions {
		if item == action {
			return true
		}
	}
	return false
}

func (r Report) Complete() bool {
	return r.HasAction("create") && r.HasAction("confirm") && r.HasAction("publish")
}
