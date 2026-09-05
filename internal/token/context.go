package token

import (
	"strings"
)

// Context mendefinisikan read-only query API ke Single Source of Truth (SSOT) token tema.
// Mengisolasi implementasi graph dan parser dari rules penilai (Layer 4).
// Bersifat thread-safe dan zero-alloc pada jalur evaluasi statis (hot-path).
type Context interface {
	// Path mengembalikan path berkas CSS sumber (jika dibaca dari filesystem).
	Path() string

	// Tokens mengembalikan seluruh node token yang diekstrak sebagai fakta CSS.
	Tokens() []Token

	// TokenByID mengambil token spesifik berdasarkan identifier numeriknya.
	TokenByID(id TokenID) (Token, bool)

	// ByName mengembalikan seluruh deklarasi token yang memiliki nama identik (bisa multi-scope).
	ByName(name string) []Token

	// ByPrefix mengembalikan seluruh deklarasi token yang memiliki prefix nama tertentu (misal: "--color-").
	ByPrefix(prefix string) []Token

	// Scopes mengembalikan daftar seluruh scope yang teridentifikasi dalam berkas CSS.
	Scopes() []Scope

	// Graph mengembalikan instance TokenGraph yang mendasari untuk analisis relasional tingkat lanjut.
	Graph() *TokenGraph

	// Resolve meresolusi nilai akhir suatu token dengan mengevaluasi rantai dependensi var(--...).
	Resolve(id TokenID, opts ResolveOptions) (string, bool, error)

	// HasScopeProperty memeriksa apakah ada scope yang mendeklarasikan properti standar tertentu (misal: "color-scheme").
	HasScopeProperty(property, value string) bool
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

func (c *themeContext) Tokens() []Token {
	return c.graph.Nodes
}

func (c *themeContext) TokenByID(id TokenID) (Token, bool) {
	if int(id) < len(c.graph.Nodes) {
		return c.graph.Nodes[id], true
	}
	return Token{}, false
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

func (c *themeContext) Scopes() []Scope {
	return c.scopes
}

func (c *themeContext) Graph() *TokenGraph {
	return c.graph
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
