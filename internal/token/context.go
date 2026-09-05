package token

import (
	"iter"
	"strings"
)

// Context mendefinisikan read-only query API ke Single Source of Truth (SSOT) token tema.
// Mengisolasi implementasi graph dan parser dari rules penilai (Layer 4).
// Menjamin immutability penuh (zero mutation leak) dan zero-alloc pada jalur evaluasi statis (hot-path).
type Context interface {
	// Path mengembalikan path berkas CSS sumber (jika dibaca dari filesystem).
	Path() string

	// TokenCount mengembalikan jumlah token terdaftar dalam graph (O(1), 0 alloc).
	TokenCount() int

	// HasToken memeriksa keberadaan token berdasarkan nama (O(1), 0 alloc).
	HasToken(name string) bool

	// LookupToken mengambil token pertama yang cocok dengan nama tertentu tanpa alokasi heap (O(1), 0 alloc).
	LookupToken(name string) (Token, bool)

	// TokenByID mengambil token spesifik berdasarkan identifier numeriknya (O(1), 0 alloc).
	TokenByID(id TokenID) (Token, bool)

	// ScopeCount mengembalikan jumlah scope unik dalam berkas tema (O(1), 0 alloc).
	ScopeCount() int

	// ScopeByID mengambil scope spesifik berdasarkan indeks (O(1), 0 alloc).
	ScopeByID(idx int) (Scope, bool)

	// HasScopeProperty memeriksa apakah ada scope yang mendeklarasikan properti standar tertentu (misal: "color-scheme").
	HasScopeProperty(property, value string) bool

	// FindCycles memeriksa siklus relasional var() dan mengembalikan daftar TokenID yang terlibat.
	FindCycles() [][]TokenID

	// AllTokens mengembalikan iterator Go 1.26 range-over-func untuk seluruh token (zero heap alloc).
	AllTokens() iter.Seq[Token]

	// AllScopes mengembalikan iterator Go 1.26 range-over-func untuk seluruh scope (zero heap alloc).
	AllScopes() iter.Seq[Scope]

	// TokensByName mengembalikan iterator Go 1.26 range-over-func untuk token dengan nama tertentu (zero heap alloc).
	TokensByName(name string) iter.Seq[Token]

	// TokensByPrefix mengembalikan iterator Go 1.26 range-over-func untuk token dengan prefix tertentu (zero heap alloc).
	TokensByPrefix(prefix string) iter.Seq[Token]

	// Tokens mengembalikan salinan defensif (defensive snapshot) seluruh token.
	// Aman dari mutasi eksternal, tetapi mengalokasikan slice baru di heap.
	Tokens() []Token

	// Scopes mengembalikan salinan defensif (defensive snapshot) seluruh scope.
	// Aman dari mutasi eksternal, tetapi mengalokasikan slice baru di heap.
	Scopes() []Scope

	// ByName mengembalikan salinan defensif seluruh deklarasi token dengan nama identik.
	ByName(name string) []Token

	// ByPrefix mengembalikan salinan defensif seluruh deklarasi token dengan prefix tertentu.
	ByPrefix(prefix string) []Token

	// Resolve meresolusi nilai akhir suatu token dengan mengevaluasi rantai dependensi var(--...).
	Resolve(id TokenID, opts ResolveOptions) (string, bool, error)
}

type themeContext struct {
	path            string
	graph           *TokenGraph
	scopes          []Scope
	scopeProperties map[string][]string // property -> list of values
}

// NewEmptyContext mengembalikan Context kosong yang aman saat tidak ada berkas tema yang ditemukan.
func NewEmptyContext() Context {
	return &themeContext{
		graph:           NewTokenGraph(),
		scopes:          make([]Scope, 0),
		scopeProperties: make(map[string][]string),
	}
}

// NewContext menginisialisasi implementasi Context read-only query facade.
func NewContext(path string, graph *TokenGraph, scopes []Scope, scopeProperties map[string][]string) Context {
	if graph == nil {
		graph = NewTokenGraph()
	}
	if scopeProperties == nil {
		scopeProperties = make(map[string][]string)
	}
	return &themeContext{
		path:            path,
		graph:           graph,
		scopes:          scopes,
		scopeProperties: scopeProperties,
	}
}

func (c *themeContext) Path() string {
	return c.path
}

func (c *themeContext) TokenCount() int {
	return len(c.graph.Nodes)
}

func (c *themeContext) HasToken(name string) bool {
	ids, ok := c.graph.ByName[name]
	return ok && len(ids) > 0
}

func (c *themeContext) LookupToken(name string) (Token, bool) {
	ids, ok := c.graph.ByName[name]
	if !ok || len(ids) == 0 {
		return Token{}, false
	}
	return c.graph.Nodes[ids[0]], true
}

func (c *themeContext) TokenByID(id TokenID) (Token, bool) {
	if int(id) < len(c.graph.Nodes) {
		return c.graph.Nodes[id], true
	}
	return Token{}, false
}

func (c *themeContext) ScopeCount() int {
	return len(c.scopes)
}

func (c *themeContext) ScopeByID(idx int) (Scope, bool) {
	if idx >= 0 && idx < len(c.scopes) {
		return c.scopes[idx], true
	}
	return Scope{}, false
}

func (c *themeContext) FindCycles() [][]TokenID {
	return c.graph.FindCycles()
}

func (c *themeContext) AllTokens() iter.Seq[Token] {
	return func(yield func(Token) bool) {
		for _, node := range c.graph.Nodes {
			if !yield(node) {
				return
			}
		}
	}
}

func (c *themeContext) AllScopes() iter.Seq[Scope] {
	return func(yield func(Scope) bool) {
		for _, s := range c.scopes {
			if !yield(s) {
				return
			}
		}
	}
}

func (c *themeContext) TokensByName(name string) iter.Seq[Token] {
	return func(yield func(Token) bool) {
		ids, ok := c.graph.ByName[name]
		if !ok {
			return
		}
		for _, id := range ids {
			if !yield(c.graph.Nodes[id]) {
				return
			}
		}
	}
}

func (c *themeContext) TokensByPrefix(prefix string) iter.Seq[Token] {
	return func(yield func(Token) bool) {
		for _, tok := range c.graph.Nodes {
			if strings.HasPrefix(tok.Name, prefix) {
				if !yield(tok) {
					return
				}
			}
		}
	}
}

func (c *themeContext) Tokens() []Token {
	if len(c.graph.Nodes) == 0 {
		return nil
	}
	res := make([]Token, len(c.graph.Nodes))
	copy(res, c.graph.Nodes)
	return res
}

func (c *themeContext) Scopes() []Scope {
	if len(c.scopes) == 0 {
		return nil
	}
	res := make([]Scope, len(c.scopes))
	copy(res, c.scopes)
	return res
}

func (c *themeContext) ByName(name string) []Token {
	ids, ok := c.graph.ByName[name]
	if !ok || len(ids) == 0 {
		return nil
	}
	res := make([]Token, len(ids))
	for i, id := range ids {
		res[i] = c.graph.Nodes[id]
	}
	return res
}

func (c *themeContext) ByPrefix(prefix string) []Token {
	var res []Token
	for _, tok := range c.graph.Nodes {
		if strings.HasPrefix(tok.Name, prefix) {
			res = append(res, tok)
		}
	}
	return res
}

func (c *themeContext) Resolve(id TokenID, opts ResolveOptions) (string, bool, error) {
	return c.graph.ResolveValue(id, opts)
}

func (c *themeContext) HasScopeProperty(property, value string) bool {
	vals, ok := c.scopeProperties[strings.ToLower(property)]
	if !ok {
		return false
	}
	if value == "" {
		return len(vals) > 0
	}
	for _, v := range vals {
		if strings.EqualFold(v, value) {
			return true
		}
	}
	return false
}
