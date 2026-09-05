package token

import (
	"fmt"
	"os"
	"path/filepath"
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

func findSSOTInDir(dir string) string {
	for _, rel := range standardCSSPaths {
		cand := filepath.Join(dir, rel)
		if _, err := os.Stat(cand); err == nil {
			return cand
		}
	}
	return ""
}

func resolveCSSPath(projectRoot, customPath string) (string, error) {
	if customPath != "" {
		resolved := customPath
		if !filepath.IsAbs(resolved) && projectRoot != "" {
			resolved = filepath.Join(projectRoot, resolved)
		}
		if _, err := os.Stat(resolved); err != nil {
			return "", fmt.Errorf("specified theme css file not found: %s: %w", customPath, err)
		}
		return resolved, nil
	}

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
		if found := findSSOTInDir(curr); found != "" {
			return found, nil
		}
		parent := filepath.Dir(curr)
		if parent == curr {
			break
		}
		curr = parent
	}
	return "", nil
}

// DiscoverAndLoad mencari berkas CSS tema berdasarkan discovery fallback engine Charites:
//  1. Jika customPath tidak kosong, periksa customPath.
//  2. Jika customPath kosong, cari di daftar standardCSSPaths dengan melakukan upward walk
//     mulai dari projectRoot (atau direktori kerja saat ini) hingga akar filesystem.
//  3. Jika berkas ditemukan, parse dan kembalikan Context terstruktur.
//  4. Jika tidak ditemukan berkas CSS manapun, kembalikan Context kosong tanpa error
//     (mematuhi invarian Zero-Config Default: YES).
func DiscoverAndLoad(projectRoot, customPath string) (Context, error) {
	targetPath, err := resolveCSSPath(projectRoot, customPath)
	if err != nil {
		return nil, err
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

func registerLayerOrder(prelude string, graph *TokenGraph) {
	for _, part := range strings.Split(prelude, ",") {
		lName := strings.TrimSpace(part)
		if lName != "" {
			if _, ok := graph.LayerOrder[lName]; !ok {
				graph.LayerOrder[lName] = len(graph.LayerOrder) + 1
			}
		}
	}
}

func processAtRule(
	node *theme.AtRule,
	currentScope Scope,
	graph *TokenGraph,
	scopeProps map[string][]string,
	sourceOrder *int,
	walk func([]theme.Rule, Scope),
) {
	if strings.EqualFold(node.Name, "@layer") && node.Prelude != "" {
		registerLayerOrder(node.Prelude, graph)
	}

	childScope := currentScope
	childScope.AtRules = append(childScope.AtRules, AtRule{
		Name:       node.Name,
		Prelude:    node.Prelude,
		Conditions: parseConditions(node.Name, node.Prelude),
	})
	if strings.EqualFold(node.Name, "@layer") && node.Prelude != "" {
		childScope.Layers = append(childScope.Layers, node.Prelude)
	}

	for _, decl := range node.Declarations {
		processDeclaration(decl, childScope, graph, scopeProps, sourceOrder)
	}

	if len(node.Rules) > 0 {
		walk(node.Rules, childScope)
	}
}

func processStyleRule(
	node *theme.StyleRule,
	currentScope Scope,
	graph *TokenGraph,
	scopeProps map[string][]string,
	sourceOrder *int,
	scopes *[]Scope,
	walk func([]theme.Rule, Scope),
) {
	resolvedSelector := resolveNestedSelector(currentScope.Selector, node.Selector)
	childScope := currentScope
	childScope.Selector = resolvedSelector
	childScope.Specificity = ComputeSpecificity(resolvedSelector)
	*sourceOrder++
	childScope.SourceOrder = *sourceOrder
	*scopes = append(*scopes, childScope)

	for _, decl := range node.Declarations {
		processDeclaration(decl, childScope, graph, scopeProps, sourceOrder)
	}

	if len(node.Rules) > 0 {
		walk(node.Rules, childScope)
	}
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
				processAtRule(node, currentScope, graph, scopeProps, &sourceOrder, walkRules)
			case *theme.StyleRule:
				processStyleRule(node, currentScope, graph, scopeProps, &sourceOrder, &scopes, walkRules)
			}
		}
	}

	initialScope := Scope{
		Selector:    "",
		AtRules:     make([]AtRule, 0),
		Layers:      make([]string, 0),
		SourceOrder: 0,
		Specificity: Specificity{},
	}

	for _, decl := range sheet.Declarations {
		processDeclaration(decl, initialScope, graph, scopeProps, &sourceOrder)
	}

	walkRules(sheet.Rules, initialScope)
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
		unescapedProp := theme.UnescapeCSS(prop)
		refs := extractVarReferences(val)
		graph.AddToken(unescapedProp, val, tokenScope, decl.Span, refs)
		return
	}

	// Properti non-custom (misal: color-scheme: light dark)
	scopeProps[strings.ToLower(prop)] = append(scopeProps[strings.ToLower(prop)], val)
}

func extractVarReferences(val string) []string {
	return theme.ExtractAllVarNames(val)
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
	if name != "@media" && name != "@supports" && name != "@container" {
		return nil
	}
	if prelude == "" {
		return nil
	}
	return []Condition{
		ParseCondition(atName, prelude),
	}
}
