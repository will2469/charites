package token

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/will2469/charites/internal/token/theme"
)

// Standard CSS search paths relative to project root for automatic discovery.
var standardCSSPaths = []string{
	"src/style/global.css",
	"src/styles/global.css",
	"styles/global.css",
	"src/global.css",
	"src/index.css",
	"global.css",
	"index.css",
	"tests/fixtures/global.css",
}

var varRefRegex = regexp.MustCompile(`var\(\s*(--[a-zA-Z0-9_-]+)`)

// DiscoverAndLoad mencari berkas CSS tema berdasarkan discovery fallback engine Charites:
//  1. Jika customPath tidak kosong, periksa customPath.
//  2. Jika customPath kosong, cari di daftar standardCSSPaths dengan melakukan upward walk
//     mulai dari projectRoot (atau direktori kerja saat ini) hingga akar filesystem.
//  3. Jika berkas ditemukan, parse dan kembalikan Context terstruktur.
//  4. Jika tidak ditemukan berkas CSS manapun, kembalikan Context kosong tanpa error
//     (mematuhi invarian Zero-Config Default: YES).
func DiscoverAndLoad(projectRoot, customPath string) (Context, error) {
	var targetPath string

	if customPath != "" {
		resolved := customPath
		if !filepath.IsAbs(resolved) && projectRoot != "" {
			resolved = filepath.Join(projectRoot, resolved)
		}
		if _, err := os.Stat(resolved); err == nil {
			targetPath = resolved
		} else {
			return nil, fmt.Errorf("specified theme css file not found: %s: %w", customPath, err)
		}
	} else {
		startDir := projectRoot
		if startDir == "" {
			if cwd, err := os.Getwd(); err == nil {
				startDir = cwd
			}
		}
		if abs, err := filepath.Abs(startDir); err == nil {
			startDir = abs
		}

		curr := startDir
		for {
			for _, rel := range standardCSSPaths {
				cand := filepath.Join(curr, rel)
				if _, err := os.Stat(cand); err == nil {
					targetPath = cand
					break
				}
			}
			if targetPath != "" {
				break
			}
			parent := filepath.Dir(curr)
			if parent == curr {
				break
			}
			curr = parent
		}
	}

	if targetPath == "" {
		return NewEmptyContext(), nil
	}

	data, err := os.ReadFile(filepath.Clean(targetPath))
	if err != nil {
		return nil, fmt.Errorf("failed to read theme css file %s: %w", targetPath, err)
	}

	ctx, err := ParseCSS(data)
	if err != nil {
		return nil, err
	}
	if tc, ok := ctx.(*themeContext); ok {
		tc.path = targetPath
	}
	return ctx, nil
}

// ParseCSS mem-parse buffer CSS mentah menjadi Context terstruktur yang design-agnostic.
func ParseCSS(src []byte) (Context, error) {
	sheet, err := theme.Parse(src)
	if err != nil {
		return nil, fmt.Errorf("failed to parse css: %w", err)
	}

	graph := NewTokenGraph()
	var scopes []Scope
	scopeProps := make(map[string][]string)
	sourceOrder := 0

	var walkRules func(rules []theme.Rule, currentScope Scope)
	walkRules = func(rules []theme.Rule, currentScope Scope) {
		for _, r := range rules {
			switch node := r.(type) {
			case *theme.AtRule:
				at := AtRule{
					Name:       node.Name,
					Prelude:    node.Prelude,
					Conditions: parseConditions(node.Name, node.Prelude),
				}
				childScope := currentScope
				childScope.AtRules = append(childScope.AtRules, at)
				if strings.EqualFold(node.Name, "@layer") && node.Prelude != "" {
					childScope.Layers = append(childScope.Layers, node.Prelude)
				}

				// Deklarasi langsung di dalam AtRule (misal: @theme { --font-sans: ...; })
				for _, decl := range node.Declarations {
					processDeclaration(decl, childScope, graph, scopeProps, &sourceOrder)
				}

				if len(node.Rules) > 0 {
					walkRules(node.Rules, childScope)
				}

			case *theme.StyleRule:
				resolvedSelector := resolveNestedSelector(currentScope.Selector, node.Selector)
				childScope := currentScope
				childScope.Selector = resolvedSelector
				childScope.Specificity = ComputeSpecificity(resolvedSelector)
				sourceOrder++
				childScope.SourceOrder = sourceOrder
				scopes = append(scopes, childScope)

				for _, decl := range node.Declarations {
					processDeclaration(decl, childScope, graph, scopeProps, &sourceOrder)
				}

				if len(node.Rules) > 0 {
					walkRules(node.Rules, childScope)
				}
			}
		}
	}

	// Scope global awal
	initialScope := Scope{
		Selector:    "",
		AtRules:     make([]AtRule, 0),
		Layers:      make([]string, 0),
		SourceOrder: 0,
		Specificity: Specificity{},
	}

	// Deklarasi root level
	for _, decl := range sheet.Declarations {
		processDeclaration(decl, initialScope, graph, scopeProps, &sourceOrder)
	}

	// Aturan tingkat atas
	walkRules(sheet.Rules, initialScope)

	// Hubungkan dependensi relasional token
	graph.BuildDependencies()

	return NewContext("", graph, scopes, scopeProps), nil
}

func processDeclaration(
	decl theme.Declaration,
	scope Scope,
	graph *TokenGraph,
	scopeProps map[string][]string,
	sourceOrder *int,
) {
	prop := strings.TrimSpace(decl.Property)
	val := strings.TrimSpace(decl.Value)

	if strings.HasPrefix(prop, "--") {
		*sourceOrder++
		tokenScope := scope
		tokenScope.SourceOrder = *sourceOrder
		refs := extractVarReferences(val)
		graph.AddToken(prop, val, tokenScope, decl.Span, refs)
		return
	}

	// Properti non-custom (misal: color-scheme: light dark)
	scopeProps[strings.ToLower(prop)] = append(scopeProps[strings.ToLower(prop)], val)
}

func extractVarReferences(val string) []string {
	matches := varRefRegex.FindAllStringSubmatch(val, -1)
	if len(matches) == 0 {
		return nil
	}

	var refs []string
	seen := make(map[string]bool)
	for _, m := range matches {
		if len(m) >= 2 && !seen[m[1]] {
			seen[m[1]] = true
			refs = append(refs, m[1])
		}
	}
	return refs
}

func resolveNestedSelector(parent, child string) string {
	parent = strings.TrimSpace(parent)
	child = strings.TrimSpace(child)

	if parent == "" {
		return child
	}
	if strings.Contains(child, "&") {
		return strings.ReplaceAll(child, "&", parent)
	}
	return parent + " " + child
}

func parseConditions(atName, prelude string) []Condition {
	name := strings.ToLower(atName)
	var condType string
	switch name {
	case "@media":
		condType = "media"
	case "@supports":
		condType = "supports"
	case "@container":
		condType = "container"
	default:
		condType = strings.TrimPrefix(name, "@")
	}

	if prelude == "" {
		return nil
	}

	return []Condition{
		{Type: condType, Query: prelude},
	}
}
