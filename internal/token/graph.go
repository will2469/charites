package token

import (
	"errors"
	"fmt"
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
	LayerOrder map[string]int
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
		LayerOrder: make(map[string]int),
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

func (g *Graph) computeCascadeRank(cID TokenID, currentScope Scope) CascadeRank {
	tok := g.Nodes[cID]
	candScope := tok.Scope

	// 1. ConditionScore
	condScore := 1 // default unconditional
	hasConditions := false
	for _, at := range candScope.AtRules {
		if len(at.Conditions) > 0 {
			hasConditions = true
			break
		}
	}
	if hasConditions {
		if candScope.MatchesConditions(currentScope) {
			condScore = 2
		} else {
			condScore = 0
		}
	}

	// 2. LayerRank (Unlayered styles win over layered styles in CSS)
	layerRank := 1000000
	if len(candScope.Layers) > 0 {
		topLayer := candScope.Layers[len(candScope.Layers)-1]
		if order, ok := g.LayerOrder[topLayer]; ok {
			layerRank = order
		} else {
			layerRank = 1
		}
	}

	// 3. SelectorAffinity
	selAffinity := 0
	if candScope.Selector == currentScope.Selector && currentScope.Selector != "" {
		selAffinity = 2
	} else if candScope.IsRoot() {
		selAffinity = 1
	}

	return CascadeRank{
		ConditionScore:   condScore,
		LayerRank:        layerRank,
		Specificity:      candScope.Specificity,
		SelectorAffinity: selAffinity,
		SourceOrder:      candScope.SourceOrder,
	}
}

// matchBestScope memilih TokenID pemenang dari daftar kandidat sesuai semantik CSS Cascade penuh:
// 1. Conditional At-Rule applicability (@media, @supports, @container)
// 2. CSS Cascade Layer precedence (unlayered > later layers > earlier layers)
// 3. Selector Specificity (A, B, C)
// 4. Selector Affinity (exact selector match > root fallback)
// 5. Source Order (order of appearance in file: later declaration wins)
func (g *Graph) matchBestScope(candidates []TokenID, currentScope Scope) TokenID {
	if len(candidates) == 1 {
		return candidates[0]
	}

	bestID := candidates[0]
	bestRank := g.computeCascadeRank(bestID, currentScope)

	for _, id := range candidates[1:] {
		rank := g.computeCascadeRank(id, currentScope)
		if rank.GreaterThan(bestRank) {
			bestID = id
			bestRank = rank
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
	return g.resolveString(raw, currID, activePath, budgetCounter, maxBudget, opts)
}

func (g *Graph) resolveString(
	input string,
	originID TokenID,
	activePath map[TokenID]bool,
	budgetCounter *int,
	maxBudget int,
	opts ResolveOptions,
) (string, bool, error) {
	if !strings.Contains(input, "var(") && !strings.Contains(input, "VAR(") {
		return input, true, nil
	}

	calls := theme.ExtractTopLevelVarCalls(input)
	if len(calls) == 0 {
		return input, true, nil
	}

	result := input
	// Substitusi dari kanan ke kiri (offset terbesar ke terkecil) agar slice offsets stabil
	for i := len(calls) - 1; i >= 0; i-- {
		call := calls[i]
		refName := call.Name
		candidates, exists := g.ByName[refName]

		var replacement string
		if !exists || len(candidates) == 0 {
			if call.HasFallback {
				// Evaluasi jika fallback mengandung var() bersarang
				fallbackResolved, _, err := g.resolveString(call.Fallback, originID, activePath, budgetCounter, maxBudget, opts)
				if err != nil {
					return "", false, err
				}
				replacement = fallbackResolved
			} else {
				// Tanpa fallback dan token tidak ditemukan, pertahankan panggilan var asli
				continue
			}
		} else {
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
				targetID = g.matchBestScope(candidates, g.Nodes[originID].Scope)
			}

			val, ok, err := g.resolveRecursive(targetID, activePath, budgetCounter, maxBudget, opts)
			if err != nil {
				return "", false, err
			}
			if !ok {
				continue
			}
			replacement = val
		}

		result = result[:call.StartOffset] + replacement + result[call.EndOffset:]
	}

	return result, true, nil
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
