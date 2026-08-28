package catalog

import (
	"errors"
	"sort"
	"strings"
	"time"
)

type Schedule struct {
	Day      string   `json:"day"`
	Windows  []Window `json:"windows"`
	Timezone string   `json:"timezone"`
}

type Window struct {
	Name        string    `json:"name"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
	Capacity    int       `json:"capacity"`
	Description string    `json:"description"`
}

func NewSchedule(day, timezone string) Schedule {
	return Schedule{Day: strings.TrimSpace(day), Timezone: strings.TrimSpace(timezone), Windows: make([]Window, 0)}
}

func (s *Schedule) AddWindow(window Window) error {
	window.Name = strings.TrimSpace(window.Name)
	if window.Name == "" || !window.End.After(window.Start) {
		return errors.New("schedule window is invalid")
	}
	if window.Capacity < 1 {
		return errors.New("schedule capacity must be positive")
	}
	for _, existing := range s.Windows {
		if existing.Name == window.Name {
			return errors.New("schedule window already exists")
		}
		if window.Start.Before(existing.End) && existing.Start.Before(window.End) {
			return errors.New("schedule windows overlap")
		}
	}
	s.Windows = append(s.Windows, window)
	sort.Slice(s.Windows, func(i, j int) bool { return s.Windows[i].Start.Before(s.Windows[j].Start) })
	return nil
}

func (s Schedule) Validate() error {
	if s.Day == "" || s.Timezone == "" {
		return errors.New("schedule day and timezone are required")
	}
	if _, err := time.Parse("2006-01-02", s.Day); err != nil {
		return errors.New("schedule day is invalid")
	}
	if len(s.Windows) == 0 {
		return errors.New("schedule needs at least one window")
	}
	return nil
}

func (s Schedule) AvailableAt(now time.Time) []Window {
	available := make([]Window, 0)
	for _, window := range s.Windows {
		if !now.Before(window.Start) && now.Before(window.End) {
			available = append(available, window)
		}
	}
	return available
}

func (s Schedule) TotalCapacity() int {
	total := 0
	for _, window := range s.Windows {
		total += window.Capacity
	}
	return total
}

func (s Schedule) Names() []string {
	names := make([]string, 0, len(s.Windows))
	for _, window := range s.Windows {
		names = append(names, window.Name)
	}
	return names
}
