package catalog

import (
	"sort"
	"sync"
	"time"
)

type ViewMetric struct {
	Screen    string    `json:"screen"`
	Opens     int       `json:"opens"`
	Starts    int       `json:"starts"`
	LastStart time.Time `json:"last_start"`
}

type Metrics struct {
	mu    sync.Mutex
	items map[string]ViewMetric
}

func NewMetrics() *Metrics { return &Metrics{items: make(map[string]ViewMetric)} }

func (m *Metrics) Open(screen string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	item := m.items[screen]
	item.Screen = screen
	item.Opens++
	m.items[screen] = item
}

func (m *Metrics) Start(screen string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	item := m.items[screen]
	item.Screen = screen
	item.Starts++
	item.LastStart = time.Now().UTC()
	m.items[screen] = item
}

func (m *Metrics) Snapshot() []ViewMetric {
	m.mu.Lock()
	defer m.mu.Unlock()
	result := make([]ViewMetric, 0, len(m.items))
	for _, item := range m.items {
		result = append(result, item)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Screen < result[j].Screen })
	return result
}

func (m *Metrics) TotalStarts() int {
	total := 0
	for _, item := range m.Snapshot() {
		total += item.Starts
	}
	return total
}
