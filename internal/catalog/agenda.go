package catalog

import (
	"errors"
	"sort"
	"strings"
	"time"
)

type AgendaItem struct {
	ID          string    `json:"id"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
	Label       string    `json:"label"`
	Description string    `json:"description"`
	Required    bool      `json:"required"`
}

type Agenda struct {
	Date  string       `json:"date"`
	Items []AgendaItem `json:"items"`
}

func NewAgenda(date string) Agenda {
	return Agenda{Date: strings.TrimSpace(date), Items: make([]AgendaItem, 0)}
}

func (a *Agenda) Add(item AgendaItem) error {
	item.ID = strings.TrimSpace(item.ID)
	item.Label = strings.TrimSpace(item.Label)
	if item.ID == "" || item.Label == "" {
		return errors.New("agenda item id and label are required")
	}
	if !item.End.After(item.Start) {
		return errors.New("agenda item must end after it starts")
	}
	for _, existing := range a.Items {
		if existing.ID == item.ID {
			return errors.New("agenda item already exists")
		}
		if item.Start.Before(existing.End) && existing.Start.Before(item.End) {
			return errors.New("agenda item overlaps another item")
		}
	}
	a.Items = append(a.Items, item)
	sort.Slice(a.Items, func(i, j int) bool { return a.Items[i].Start.Before(a.Items[j].Start) })
	return nil
}

func (a Agenda) Validate() error {
	if a.Date == "" {
		return errors.New("agenda date is required")
	}
	if _, err := time.Parse("2006-01-02", a.Date); err != nil {
		return errors.New("agenda date is invalid")
	}
	if len(a.Items) == 0 {
		return errors.New("agenda needs at least one item")
	}
	for index, item := range a.Items {
		if item.ID == "" || item.Label == "" || !item.End.After(item.Start) {
			return errors.New("agenda item at position " + string(rune(index+'0')) + " is invalid")
		}
	}
	return nil
}

func (a Agenda) Current(now time.Time) (AgendaItem, bool) {
	for _, item := range a.Items {
		if !now.Before(item.Start) && now.Before(item.End) {
			return item, true
		}
	}
	return AgendaItem{}, false
}

func (a Agenda) Next(now time.Time) (AgendaItem, bool) {
	for _, item := range a.Items {
		if item.Start.After(now) {
			return item, true
		}
	}
	return AgendaItem{}, false
}

func (a Agenda) RequiredCount() int {
	count := 0
	for _, item := range a.Items {
		if item.Required {
			count++
		}
	}
	return count
}

func (a Agenda) Labels() []string {
	labels := make([]string, 0, len(a.Items))
	for _, item := range a.Items {
		labels = append(labels, item.Label)
	}
	return labels
}
