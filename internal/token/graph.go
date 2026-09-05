package token

import (
	"errors"
	"strings"

	"github.com/will2469/charites/internal/parser/css"
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
func (g *Graph) AddToken(name, rawValue string, scope Scope, span css.SourceSpan, refs []string) ID {
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

func (g *Graph) computeCascadeRank(cID TokenID) CascadeRank {
	tok := g.Nodes[cID]
	candScope := tok.Scope

	// LayerRank (Unlayered styles win over layered styles in CSS Cascade 5)
	layerRank := 1000000
	if len(candScope.Layers) > 0 {
		topLayer := candScope.Layers[len(candScope.Layers)-1]
		if order, ok := g.LayerOrder[topLayer]; ok {
			layerRank = order
		} else {
			layerRank = 1
		}
	}

	return CascadeRank{
		LayerRank:   layerRank,
		Specificity: candScope.Specificity,
		SourceOrder: candScope.SourceOrder,
	}
}

// matchBestScope memilih TokenID pemenang dari daftar kandidat sesuai semantik W3C CSS Cascade 5:
// 1. Applicability Filtering: Deklarasi yang kondisinya tidak kompatibel dengan konteks difilter keluar.
// 2. Element Scoping (Inheritance): Deklarasi langsung pada elemen mengalahkan nilai warisan dari :root / ancestor.
// 3. Cascade Sorting (§6): Diurutkan murni berdasarkan LayerRank -> Specificity -> SourceOrder.
func (g *Graph) matchBestScope(candidates []TokenID, currentScope Scope) TokenID {
	if len(candidates) == 1 {
		return candidates[0]
	}

	applicable := g.filterApplicableScopes(candidates, currentScope)
	if currentScope.Selector != "" && !currentScope.IsRoot() {
		applicable = g.filterDirectElementScopes(applicable, currentScope.Selector)
	}

	bestID := applicable[0]
	bestRank := g.computeCascadeRank(bestID)

	for _, id := range applicable[1:] {
		rank := g.computeCascadeRank(id)
		if rank.GreaterThan(bestRank) {
			bestID = id
			bestRank = rank
		}
	}

	return bestID
}

func (g *Graph) filterApplicableScopes(candidates []TokenID, currentScope Scope) []TokenID {
	applicable := make([]TokenID, 0, len(candidates))
	for _, id := range candidates {
		if g.Nodes[id].Scope.MatchesConditions(currentScope) {
			applicable = append(applicable, id)
		}
	}
	if len(applicable) > 0 {
		return applicable
	}

	// Fallback alami ke deklarasi tanpa kondisi
	for _, id := range candidates {
		if len(g.Nodes[id].Scope.AtRules) == 0 {
			applicable = append(applicable, id)
		}
	}
	if len(applicable) == 0 {
		return candidates
	}
	return applicable
}

func (g *Graph) filterDirectElementScopes(candidates []TokenID, selector string) []TokenID {
	direct := make([]TokenID, 0, len(candidates))
	for _, id := range candidates {
		if g.Nodes[id].Scope.Selector == selector {
			direct = append(direct, id)
		}
	}
	if len(direct) > 0 {
		return direct
	}
	return candidates
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

	budgetCounter := 0
	activePath := make(map[TokenID]bool)
	return g.resolveRecursive(tokenID, activePath, &budgetCounter, maxBudget, opts)
}

func (g *Graph) resolveRecursive(
	tokenID TokenID,
	activePath map[TokenID]bool,
	budgetCounter *int,
	maxBudget int,
	opts ResolveOptions,
) (string, bool, error) {
	*budgetCounter++
	if *budgetCounter > maxBudget {
		return "", false, ErrEvaluationBudgetExceeded
	}

	if activePath[tokenID] {
		return "", false, ErrCycleDetected
	}

	activePath[tokenID] = true
	defer func() {
		delete(activePath, tokenID)
	}()

	node := g.Nodes[tokenID]
	return g.resolveString(node.RawValue, tokenID, activePath, budgetCounter, maxBudget, opts)
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

	calls := css.ExtractTopLevelVarCalls(input)
	if len(calls) == 0 {
		return input, true, nil
	}

	result := input
	// Substitusi dari kanan ke kiri (offset terbesar ke terkecil) agar slice offsets stabil
	for i := len(calls) - 1; i >= 0; i-- {
		call := calls[i]
		replacement, ok, err := g.resolveVarReplacement(call, originID, activePath, budgetCounter, maxBudget, opts)
		if err != nil {
			return "", false, err
		}
		if !ok {
			continue
		}
		result = result[:call.StartOffset] + replacement + result[call.EndOffset:]
	}

	return result, true, nil
}

func (g *Graph) resolveVarReplacement(
	call css.VarCall,
	originID TokenID,
	activePath map[TokenID]bool,
	budgetCounter *int,
	maxBudget int,
	opts ResolveOptions,
) (string, bool, error) {
	candidates, exists := g.ByName[call.Name]
	if !exists || len(candidates) == 0 {
		if !call.HasFallback {
			return "", false, nil
		}
		val, _, err := g.resolveString(call.Fallback, originID, activePath, budgetCounter, maxBudget, opts)
		if err != nil {
			return "", false, err
		}
		return val, true, nil
	}

	targetID := g.selectTargetID(candidates, originID, opts)
	val, ok, err := g.resolveRecursive(targetID, activePath, budgetCounter, maxBudget, opts)
	if err != nil {
		return "", false, err
	}
	if !ok {
		return "", false, nil
	}
	return val, true, nil
}

func (g *Graph) selectTargetID(candidates []TokenID, originID TokenID, opts ResolveOptions) TokenID {
	if opts.ScopeMatcher != nil {
		for _, cID := range candidates {
			if opts.ScopeMatcher(g.Nodes[cID].Scope) {
				return cID
			}
		}
	}
	return g.matchBestScope(candidates, g.Nodes[originID].Scope)
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
