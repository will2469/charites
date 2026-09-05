package analyzer_test

import (
	"testing"

	"github.com/will2469/charites/internal/analyzer"
	"github.com/will2469/charites/internal/ir"
)

func TestParseDirectives_LineComments(t *testing.T) {
	src := []byte(`
// Baris 2: komentar biasa
// charites:ignore  theme.hardcode-opacity-color  ,  theme.hardcode-opacity-color  ,  a11y.alt
const a = 1;
// charites:ignore
`)
	directives := analyzer.ParseDirectives(src)

	rulesLine3, ok := directives[3]
	if !ok {
		t.Fatalf("expected directive on line 3, got none")
	}
	if len(rulesLine3) != 2 || rulesLine3[0] != "theme.hardcode-opacity-color" || rulesLine3[1] != "a11y.alt" {
		t.Errorf("line 3 rules mismatch: %+v", rulesLine3)
	}

	if _, exists := directives[5]; exists {
		t.Errorf("expected line 5 empty directive to be ignored")
	}
}

func TestParseDirectives_BlockAndHTMLComments(t *testing.T) {
	src := []byte(`
/* charites:ignore perf.bundle */
<!-- charites:ignore * -->
`)
	directives := analyzer.ParseDirectives(src)

	rulesLine2, ok := directives[2]
	if !ok || len(rulesLine2) != 1 || rulesLine2[0] != "perf.bundle" {
		t.Errorf("line 2 rules mismatch: %+v", rulesLine2)
	}

	rulesLine3, ok := directives[3]
	if !ok || len(rulesLine3) != 1 || rulesLine3[0] != "*" {
		t.Errorf("line 3 rules mismatch: %+v", rulesLine3)
	}
}

func TestParseDirectives_TemplateInterpolation(t *testing.T) {
	src := []byte("const b = `before ${ /* charites:ignore theme.inline-template */ 42 } after`;\n")
	directives := analyzer.ParseDirectives(src)

	rulesLine1, ok := directives[1]
	if !ok || len(rulesLine1) != 1 || rulesLine1[0] != "theme.inline-template" {
		t.Errorf("line 1 rules mismatch: %+v", rulesLine1)
	}
}

func TestParseDirectives_AdversarialNonComments(t *testing.T) {
	// Seluruh penanda charites:ignore di bawah ini BUKAN komentar yang sah
	// dan WAJIB ditolak (tidak menghasilkan direktif apa pun).
	src := []byte(`
const strDouble = "charites:ignore theme.hardcode-opacity-color";
const strSingle = 'charites:ignore theme.hardcode-opacity-color';
const tplLit = ` + "`" + `charites:ignore theme.hardcode-opacity-color` + "`" + `;
const jsxAttr = <input placeholder="charites:ignore theme.hardcode-opacity-color" />;
const htmlText = <p>charites:ignore theme.hardcode-opacity-color</p>;
const jsExpr = (code === "charites:ignore");
const fakeLineCommentInStr = "// charites:ignore theme.fake";
const fakeHTMLCommentInStr = "<!-- charites:ignore theme.fake -->";
const fakeEscapedQuote = "escaped \" charites:ignore theme.fake \" still in string";
`)

	directives := analyzer.ParseDirectives(src)

	if len(directives) != 0 {
		t.Fatalf("expected 0 directives extracted from non-comment adversarial cases, got %d: %+v", len(directives), directives)
	}
}

func TestContext_DeepCopy_CallerIsolation(t *testing.T) {
	// Invarian P1: caller map -> deep copy -> Context-owned map
	callerMap := map[int][]string{
		10: {"theme.color", "a11y.alt"},
		20: {"perf.bundle"},
	}

	ctx := analyzer.NewContext("test.tsx", callerMap)

	// Mutasi caller map setelah konstruksi Context
	callerMap[10][0] = "MUTATED_RULE"
	callerMap[10] = append(callerMap[10], "EXTRA_RULE")
	callerMap[99] = []string{"LEAKED_RULE"}
	delete(callerMap, 20)

	// Buktikan behavior Context tetap murni dan tidak terpengaruh oleh mutasi caller
	diag10Original := ir.Diagnostic{Line: 10, Rule: "theme.color"}
	if !ctx.IsIgnored(diag10Original, nil) {
		t.Error("expected theme.color on line 10 to still be ignored despite caller mutating slice element")
	}

	diag10Mutated := ir.Diagnostic{Line: 10, Rule: "MUTATED_RULE"}
	if ctx.IsIgnored(diag10Mutated, nil) {
		t.Error("MUTATED_RULE must NOT be recognized by Context")
	}

	diag10Extra := ir.Diagnostic{Line: 10, Rule: "EXTRA_RULE"}
	if ctx.IsIgnored(diag10Extra, nil) {
		t.Error("EXTRA_RULE must NOT be recognized by Context")
	}

	diag20 := ir.Diagnostic{Line: 20, Rule: "perf.bundle"}
	if !ctx.IsIgnored(diag20, nil) {
		t.Error("expected perf.bundle on line 20 to still be ignored despite caller deleting key from map")
	}

	diag99 := ir.Diagnostic{Line: 99, Rule: "LEAKED_RULE"}
	if ctx.IsIgnored(diag99, nil) {
		t.Error("LEAKED_RULE on line 99 must NOT be recognized by Context")
	}
}

func TestContext_DiagnosticsOwnership_Cloning(t *testing.T) {
	// Invarian P1: DiagnosticsList() mengembalikan clone independen
	ctx := analyzer.NewContext("a.tsx", nil)
	ctx.AddDiagnostic(ir.Diagnostic{
		File:     "a.tsx",
		Line:     10,
		Column:   5,
		Rule:     "theme.hardcode-color",
		Severity: ir.SeverityWarn,
		Message:  "original message",
	})

	list1 := ctx.DiagnosticsList()
	if len(list1) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(list1))
	}

	// Mutasi eksternal pada hasil return
	list1[0].Message = "MUTATED MESSAGE"
	list1[0].Rule = "MUTATED_RULE"
	list1 = append(list1, ir.Diagnostic{Rule: "EXTRA_DIAG"})
	if len(list1) != 2 {
		t.Fatalf("expected mutated external slice to have 2 items")
	}

	// Panggil DiagnosticsList() lagi: verifikasi data internal tidak terkorupsi
	list2 := ctx.DiagnosticsList()
	if len(list2) != 1 {
		t.Fatalf("expected 1 diagnostic in internal list, got %d", len(list2))
	}
	if list2[0].Message != "original message" || list2[0].Rule != "theme.hardcode-color" {
		t.Errorf("internal diagnostics were corrupted by external mutation: %+v", list2[0])
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

	// Skenario 5: Komentar berada pada opening tag baris pertama node (node.Span.Line)
	ctxOpeningTag := analyzer.NewContext("a.tsx", map[int][]string{
		11: {"theme.hardcode-opacity-color"},
	})
	if !ctxOpeningTag.IsIgnored(diagSpan, multiLineNode) {
		t.Error("expected comment on opening tag (node.Span.Line) to suppress diagnostic inside node span")
	}

	// Skenario 6: Boundary check line 1 (node.Span.Line == 1 dan diag.Line == 1)
	ctxLine1 := analyzer.NewContext("a.tsx", map[int][]string{
		1: {"theme.color"},
	})
	nodeLine1 := &ir.Node{
		Type: ir.NodeElement,
		Span: ir.Span{Line: 1, EndLine: 2},
	}
	diagLine1 := ir.Diagnostic{Line: 1, Rule: "theme.color"}
	if !ctxLine1.IsIgnored(diagLine1, nodeLine1) {
		t.Error("expected line 1 diagnostic to be suppressed")
	}
	diagLine1Unmatched := ir.Diagnostic{Line: 1, Rule: "other.rule"}
	if ctxLine1.IsIgnored(diagLine1Unmatched, nodeLine1) {
		t.Error("unmatched rule on line 1 must not be ignored")
	}
}

func TestParseDirectives_MultilineComments(t *testing.T) {
	srcBlock := []byte(`
/*
 * line 2
 * charites:ignore theme.multiline-block
 * line 4
 */
`)
	dirBlock := analyzer.ParseDirectives(srcBlock)
	if rules, ok := dirBlock[4]; !ok || len(rules) != 1 || rules[0] != "theme.multiline-block" {
		t.Errorf("expected directive on line 4 for multiline block comment, got: %+v", dirBlock)
	}

	srcHTML := []byte(`
<!--
  line 2
  charites:ignore theme.multiline-html
  line 4
-->
`)
	dirHTML := analyzer.ParseDirectives(srcHTML)
	if rules, ok := dirHTML[4]; !ok || len(rules) != 1 || rules[0] != "theme.multiline-html" {
		t.Errorf("expected directive on line 4 for multiline HTML comment, got: %+v", dirHTML)
	}
}

func TestParseDirectives_UnclosedCommentsAtEOF(t *testing.T) {
	srcUnclosedBlock := []byte(`/* charites:ignore theme.unclosed-block`)
	dirUnclosedBlock := analyzer.ParseDirectives(srcUnclosedBlock)
	if rules, ok := dirUnclosedBlock[1]; !ok || len(rules) != 1 || rules[0] != "theme.unclosed-block" {
		t.Errorf("expected unclosed block directive at EOF, got: %+v", dirUnclosedBlock)
	}

	srcUnclosedHTML := []byte(`<!-- charites:ignore theme.unclosed-html`)
	dirUnclosedHTML := analyzer.ParseDirectives(srcUnclosedHTML)
	if rules, ok := dirUnclosedHTML[1]; !ok || len(rules) != 1 || rules[0] != "theme.unclosed-html" {
		t.Errorf("expected unclosed HTML directive at EOF, got: %+v", dirUnclosedHTML)
	}
}

func TestParseDirectives_TemplateExpressionsAndEscapes(t *testing.T) {
	src := []byte("const str = \"line1\\nline2 \\\" escaped quote\";\n" +
		"const strSingle = 'single \\' quote';\n" +
		"const tpl = `line1\\`\n" +
		"line2 ${\n" +
		"  // charites:ignore theme.template-line-comment\n" +
		"  (() => {\n" +
		"    const nested = { a: { b: 1 } };\n" +
		"    const s1 = 'str';\n" +
		"    const s2 = \"str2\";\n" +
		"    const t2 = `nested ${2}`;\n" +
		"    /* charites:ignore theme.template-block-comment */\n" +
		"    return 42;\n" +
		"  })()\n" +
		"}`;\n")

	dirs := analyzer.ParseDirectives(src)

	// Verifikasi line comment di dalam ${ ... }
	if rules, ok := dirs[5]; !ok || len(rules) != 1 || rules[0] != "theme.template-line-comment" {
		t.Errorf("expected template line comment directive on line 5, got: %+v", dirs)
	}

	// Verifikasi block comment di dalam ${ ... }
	if rules, ok := dirs[11]; !ok || len(rules) != 1 || rules[0] != "theme.template-block-comment" {
		t.Errorf("expected template block comment directive on line 11, got: %+v", dirs)
	}
}
