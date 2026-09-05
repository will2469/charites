package analyzer_test

import (
	"testing"

	"github.com/will2469/charites/internal/analyzer"
	"github.com/will2469/charites/internal/ir"
)

func TestParseDirectives(t *testing.T) {
	src := []byte(`
// Baris 2: komentar biasa
// charites:ignore  theme.hardcode-opacity-color  ,  theme.hardcode-opacity-color  ,  a11y.alt
const a = 1;
/* charites:ignore perf.bundle */
<!-- charites:ignore * -->
// charites:ignore
`)

	directives := analyzer.ParseDirectives(src)

	// Baris 3: theme.hardcode-opacity-color dan a11y.alt (spasi ditrim dan deduplikasi)
	rulesLine3, ok := directives[3]
	if !ok {
		t.Fatalf("expected directive on line 3, got none")
	}
	if len(rulesLine3) != 2 || rulesLine3[0] != "theme.hardcode-opacity-color" || rulesLine3[1] != "a11y.alt" {
		t.Errorf("line 3 rules mismatch: %+v", rulesLine3)
	}

	// Baris 5: block comment /* ... */
	rulesLine5, ok := directives[5]
	if !ok || len(rulesLine5) != 1 || rulesLine5[0] != "perf.bundle" {
		t.Errorf("line 5 rules mismatch: %+v", rulesLine5)
	}

	// Baris 6: HTML comment <!-- ... -->
	rulesLine6, ok := directives[6]
	if !ok || len(rulesLine6) != 1 || rulesLine6[0] != "*" {
		t.Errorf("line 6 rules mismatch: %+v", rulesLine6)
	}

	// Baris 7: Empty directive -> diabaikan
	if _, ok := directives[7]; ok {
		t.Errorf("expected line 7 empty directive to be ignored")
	}
}

func TestContext_IsIgnored_Scopes(t *testing.T) {
	// Skenario 1: Same-line trailing comment
	ctxSameLine := analyzer.NewContext("a.tsx", map[int][]string{
		10: {"theme.hardcode-opacity-color"},
	})
	diagSameLine := ir.Diagnostic{Line: 10, Rule: "theme.hardcode-opacity-color"}
	if !ctxSameLine.IsIgnored(diagSameLine, nil) {
		t.Error("expected same-line comment to suppress diagnostic")
	}

	// Skenario 2: Next-line preceding comment (komentar di baris N, diag di baris N+1)
	ctxNextLine := analyzer.NewContext("a.tsx", map[int][]string{
		9: {"theme.hardcode-opacity-color"},
	})
	diagNextLine := ir.Diagnostic{Line: 10, Rule: "theme.hardcode-opacity-color"}
	if !ctxNextLine.IsIgnored(diagNextLine, nil) {
		t.Error("expected next-line comment to suppress diagnostic")
	}

	// Skenario 3: Wildcard suppression (*)
	ctxWildcard := analyzer.NewContext("a.tsx", map[int][]string{
		10: {"*"},
	})
	diagWildcard := ir.Diagnostic{Line: 10, Rule: "any.arbitrary.rule"}
	if !ctxWildcard.IsIgnored(diagWildcard, nil) {
		t.Error("expected wildcard to suppress any diagnostic")
	}

	// Skenario 4: Node Span Scope (Multi-line JSX element)
	// Komentar di baris 10
	// Node membentang baris 11 s/d 15
	// Diagnostic terjadi di baris 13
	ctxSpan := analyzer.NewContext("a.tsx", map[int][]string{
		10: {"theme.hardcode-opacity-color"},
	})
	multiLineNode := &ir.Node{
		Type: ir.NodeElement,
		Span: ir.Span{
			Line:      11,
			Column:    1,
			EndLine:   15,
			EndColumn: 3,
		},
	}
	diagSpan := ir.Diagnostic{Line: 13, Rule: "theme.hardcode-opacity-color"}
	if !ctxSpan.IsIgnored(diagSpan, multiLineNode) {
		t.Error("expected node span scope (N-1 preceding node) to suppress diagnostic on line 13")
	}

	// Diagnostic di luar span node (baris 16) tidak boleh ditekan
	diagOutOfSpan := ir.Diagnostic{Line: 16, Rule: "theme.hardcode-opacity-color"}
	if ctxSpan.IsIgnored(diagOutOfSpan, multiLineNode) {
		t.Error("diagnostic outside node span must NOT be suppressed")
	}
}
