package token

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/will2469/charites/internal/token/theme"
)

var (
	// ErrCycleDetected dikembalikan saat terdapat siklus dependensi sirkular (misal: --a -> --b -> --a).
	ErrCycleDetected = errors.New("cycle detected in token dependency graph")

	// ErrEvaluationBudgetExceeded dikembalikan saat evaluasi token melebihi batas MaxNodes (DoS protection).
	ErrEvaluationBudgetExceeded = errors.New("evaluation budget exceeded during token resolution")

	// ErrTokenNotFound dikembalikan saat token yang dicari tidak ada dalam graph.
	ErrTokenNotFound = errors.New("token declaration not found")
)

var varRegex = regexp.MustCompile(`var\(\s*(--[a-zA-Z0-9_-]+)(?:\s*,\s*([^)]+))?\s*\)`)

// ResolveOptions mendefinisikan opsi evaluasi saat me-resolve token dependency graph.
type ResolveOptions struct {
	// MaxNodes membatasi total kunjungan node untuk mencegah infinite loops atau serangan DoS.
	// Jika 0, default bernilai 1000.
	MaxNodes int

	// ScopeMatcher menyaring deklarasi token target jika terdapat banyak deklarasi bernama sama.
	ScopeMatcher func(Scope) bool
}

// Graph merepresentasikan Directed Token Dependency Graph.
// Menampung seluruh node deklarasi token yang teridentifikasi oleh ID,
// mengindeks berdasarkan nama token, dan melacak relasi dependensi dua arah (DependsOn & Dependents).
type Graph struct {
	Nodes      []Token
	ByName     map[string][]ID
	DependsOn  map[ID][]ID
	Dependents map[ID][]ID
}

// TokenGraph adalah alias untuk Graph demi kenyamanan penamaan.
//
//nolint:revive // alias retained for backward compatibility
type TokenGraph = Graph

// NewTokenGraph menginisialisasi Graph baru.
func NewTokenGraph() *Graph {
	return &Graph{
		Nodes:      make([]Token, 0),
		ByName:     make(map[string][]ID),
		DependsOn:  make(map[ID][]ID),
		Dependents: make(map[ID][]ID),
	}
}

// AddToken mendaftarkan node token baru ke dalam graph dan mengembalikan ID uniknya.
func (g *Graph) AddToken(name, rawValue string, scope Scope, span theme.SourceSpan, refs []string) ID {
	id := ID(len(g.Nodes)) // #nosec G115 -- length of slice in single CSS file is bounded
	tok := Token{
		ID:         id,
		Name:       name,
		RawValue:   rawValue,
		Scope:      scope,
		Span:       span,
		References: refs,
	}
	g.Nodes = append(g.Nodes, tok)
	g.ByName[name] = append(g.ByName[name], id)
	return id
}

// BuildDependencies menghubungkan relasi DependsOn dan Dependents antar token
// berdasarkan nama token yang direferensikan dalam var(--...).
func (g *Graph) BuildDependencies() {
	for _, node := range g.Nodes {
		for _, refName := range node.References {
			targetIDs, exists := g.ByName[refName]
			if !exists || len(targetIDs) == 0 {
				continue
			}
			bestTarget := g.matchBestScope(targetIDs, node.Scope)
			g.DependsOn[node.ID] = append(g.DependsOn[node.ID], bestTarget)
			g.Dependents[bestTarget] = append(g.Dependents[bestTarget], node.ID)
		}
	}
}

// matchBestScope memilih TokenID terbaik dari daftar kandidat yang cocok dengan scope referer.
func (g *Graph) matchBestScope(candidates []TokenID, currentScope Scope) TokenID {
	if len(candidates) == 1 {
		return candidates[0]
	}

	// 1. Prioritaskan selector yang sama persis
	for _, id := range candidates {
		if g.Nodes[id].Scope.Selector == currentScope.Selector {
			return id
		}
	}

	// 2. Prioritaskan kandidat dengan spesifisitas tertinggi
	bestID := candidates[0]
	for _, id := range candidates {
		if g.Nodes[id].Scope.Specificity.GreaterThan(g.Nodes[bestID].Scope.Specificity) {
			bestID = id
		}
	}

	return bestID
}

// ResolveValue mengevaluasi dan meresolusi nilai akhir dari sebuah token (mengganti var(--...)).
// Menggunakan visited set untuk cycle detection deterministik dan evaluation budget untuk perlindungan DoS.
func (g *Graph) ResolveValue(tokenID TokenID, opts ResolveOptions) (string, bool, error) {
	if int(tokenID) >= len(g.Nodes) {
		return "", false, ErrTokenNotFound
	}

	maxBudget := opts.MaxNodes
	if maxBudget <= 0 {
		maxBudget = 1000
	}

	visited := make(map[TokenID]bool)
	budgetCounter := 0

	return g.resolveRecursive(tokenID, visited, &budgetCounter, maxBudget, opts)
}

func (g *Graph) resolveRecursive(
	currID TokenID,
	activePath map[TokenID]bool,
	budgetCounter *int,
	maxBudget int,
	opts ResolveOptions,
) (string, bool, error) {
	*budgetCounter++
	if *budgetCounter > maxBudget {
		return "", false, fmt.Errorf("%w: limit %d nodes", ErrEvaluationBudgetExceeded, maxBudget)
	}

	if activePath[currID] {
		return "", false, fmt.Errorf("%w: token %s (ID %d)", ErrCycleDetected, g.Nodes[currID].Name, currID)
	}

	activePath[currID] = true
	defer func() {
		activePath[currID] = false
	}()

	raw := g.Nodes[currID].RawValue
	if !strings.Contains(raw, "var(") {
		return raw, true, nil
	}

	var resolveErr error
	resolved := varRegex.ReplaceAllStringFunc(raw, func(match string) string {
		if resolveErr != nil {
			return match
		}

		submatches := varRegex.FindStringSubmatch(match)
		if len(submatches) < 2 {
			return match
		}

		refName := submatches[1]
		candidates, exists := g.ByName[refName]
		if !exists || len(candidates) == 0 {
			// Jika ada fallback value: var(--foo, fallback)
			if len(submatches) >= 3 && submatches[2] != "" {
				return strings.TrimSpace(submatches[2])
			}
			return match
		}

		// Pilih kandidat berdasarkan ScopeMatcher atau fallback
		var targetID TokenID
		found := false
		if opts.ScopeMatcher != nil {
			for _, cID := range candidates {
				if opts.ScopeMatcher(g.Nodes[cID].Scope) {
					targetID = cID
					found = true
					break
				}
			}
		}
		if !found {
			targetID = g.matchBestScope(candidates, g.Nodes[currID].Scope)
		}

		val, ok, err := g.resolveRecursive(targetID, activePath, budgetCounter, maxBudget, opts)
		if err != nil {
			resolveErr = err
			return match
		}
		if !ok {
			return match
		}
		return val
	})

	if resolveErr != nil {
		return "", false, resolveErr
	}

	return resolved, true, nil
}

// FindCycles memeriksa seluruh token dalam graph dan mengembalikan daftar TokenID yang terlibat dalam siklus.
func (g *Graph) FindCycles() [][]TokenID {
	var cycles [][]TokenID
	opts := ResolveOptions{MaxNodes: len(g.Nodes) * 2}

	for _, node := range g.Nodes {
		_, _, err := g.ResolveValue(node.ID, opts)
		if errors.Is(err, ErrCycleDetected) {
			cycles = append(cycles, []TokenID{node.ID})
		}
	}

	return cycles
}
