package rules

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
)

var (
	// ErrNilRule dikembalikan saat mencoba mendaftarkan rule bernilai nil.
	ErrNilRule = errors.New("cannot register nil rule")
	// ErrEmptyRuleID dikembalikan saat ID rule kosong.
	ErrEmptyRuleID = errors.New("cannot register rule with empty ID")
)

// Registry merepresentasikan katalog in-memory thread-safe seluruh aturan static analysis Charites.
type Registry struct {
	mu         sync.RWMutex
	rules      map[string]Rule
	categories map[string][]Rule
}

// NewRegistry membuat instance Registry baru yang kosong.
func NewRegistry() *Registry {
	return &Registry{
		rules:      make(map[string]Rule),
		categories: make(map[string][]Rule),
	}
}

// Register mendaftarkan rule ke dalam registry.
// Menolak rule bernilai nil, ID kosong, atau ID yang sudah terdaftar.
func (r *Registry) Register(rule Rule) error {
	if rule == nil {
		return ErrNilRule
	}

	id := rule.ID()
	if id == "" {
		return ErrEmptyRuleID
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.rules[id]; exists {
		return fmt.Errorf("rule with ID %q already registered", id)
	}

	r.rules[id] = rule
	cat := rule.Category()
	r.categories[cat] = append(r.categories[cat], rule)
	return nil
}

// Get mencari rule berdasarkan Charites Rule ID kanonikal (<category>.<slug>).
func (r *Registry) Get(id string) (Rule, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rule, ok := r.rules[id]
	return rule, ok
}

// All mengembalikan seluruh rule terdaftar yang diurutkan secara leksikografis berdasarkan Rule.ID().
// Mengembalikan salinan defensif (defensive copy) yang aman dari mutasi eksternal.
func (r *Registry) All() []Rule {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]Rule, 0, len(r.rules))
	for _, rule := range r.rules {
		list = append(list, rule)
	}

	slices.SortFunc(list, func(a, b Rule) int {
		return strings.Compare(a.ID(), b.ID())
	})
	return list
}

// ByCategory mengembalikan seluruh rule pada kategori tertentu yang diurutkan secara leksikografis berdasarkan Rule.ID().
// Mengembalikan salinan defensif (defensive copy).
func (r *Registry) ByCategory(category string) []Rule {
	r.mu.RLock()
	defer r.mu.RUnlock()

	rules := r.categories[category]
	out := make([]Rule, len(rules))
	copy(out, rules)

	slices.SortFunc(out, func(a, b Rule) int {
		return strings.Compare(a.ID(), b.ID())
	})
	return out
}

// Count mengembalikan total rule yang saat ini terdaftar di dalam registry.
func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.rules)
}

// defaultRegistry adalah instance singleton global bawaan.
var defaultRegistry = NewRegistry()

// DefaultRegistry mengembalikan instance global singleton registry.
func DefaultRegistry() *Registry {
	return defaultRegistry
}

// Register mendaftarkan rule ke default registry global.
func Register(rule Rule) error {
	return defaultRegistry.Register(rule)
}

// Get mencari rule di default registry global.
func Get(id string) (Rule, bool) {
	return defaultRegistry.Get(id)
}

// All mengembalikan seluruh rule dari default registry global dalam urutan leksikografis.
func All() []Rule {
	return defaultRegistry.All()
}

// ByCategory mengembalikan seluruh rule kategori tertentu dari default registry global dalam urutan leksikografis.
func ByCategory(category string) []Rule {
	return defaultRegistry.ByCategory(category)
}

// Count mengembalikan total rule di default registry global.
func Count() int {
	return defaultRegistry.Count()
}
