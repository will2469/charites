package theme_test

import (
	"testing"

	"github.com/will2469/charites/internal/ir"
	"github.com/will2469/charites/internal/rules/theme"
)

func makeNode(classes ...string) *ir.Node {
	return &ir.Node{
		Classes: classes,
		Span:    ir.Span{Line: 1, Column: 1},
	}
}

func makeNodeWithAttr(attrs map[string]string, classes ...string) *ir.Node {
	return &ir.Node{
		Classes:    classes,
		Attributes: attrs,
		Span:       ir.Span{Line: 1, Column: 1},
	}
}

func TestHardcodeColorRule(t *testing.T) {
	rule := theme.NewHardcodeColorRule()

	if rule.ID() != "theme.hardcode-color" {
		t.Fatalf("unexpected ID: %s", rule.ID())
	}
	if rule.Category() != "theme" {
		t.Fatalf("unexpected category: %s", rule.Category())
	}
	if rule.DefaultSeverity() != ir.SeverityWarn {
		t.Fatalf("unexpected severity: %v", rule.DefaultSeverity())
	}

	doc := rule.Doc()
	if len(doc.TargetStandards) == 0 {
		t.Fatalf("empty doc target standards")
	}

	tests := []struct {
		name    string
		classes []string
		want    int
	}{
		{"ArbitraryHex", []string{"bg-[#2563eb]"}, 1},
		{"ArbitraryProp", []string{"[color:#fff]"}, 1},
		{"ArbitraryPropBg", []string{"[background-color:#1e293b]"}, 1},
		{"ArbitraryRGB", []string{"text-[rgb(255,0,0)]"}, 1},
		{"VariantArbitrary", []string{"hover:bg-[#123456]"}, 1},
		{"CleanToken", []string{"bg-primary"}, 0},
		{"CleanVar", []string{"bg-[var(--primary)]"}, 0},
		{"NonColorArbitrary", []string{"p-[13px]", "w-[200px]"}, 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			node := makeNode(tc.classes...)
			diags := rule.Evaluate(node)
			if len(diags) != tc.want {
				t.Fatalf("expected %d diagnostics, got %d: %+v", tc.want, len(diags), diags)
			}
		})
	}

	// Edge cases
	if diags := rule.Evaluate(nil); len(diags) != 0 {
		t.Errorf("expected 0 diags on nil node")
	}
	if diags := rule.Evaluate(makeNode()); len(diags) != 0 {
		t.Errorf("expected 0 diags on empty node")
	}
}

func TestPrimitiveInComponentRule(t *testing.T) {
	rule := theme.NewPrimitiveInComponentRule()

	if rule.ID() != "theme.primitive-in-component" {
		t.Fatalf("unexpected ID: %s", rule.ID())
	}
	if rule.DefaultSeverity() != ir.SeverityError {
		t.Fatalf("unexpected severity: %v", rule.DefaultSeverity())
	}

	tests := []struct {
		name    string
		classes []string
		want    int
	}{
		{"PrimitiveBg", []string{"bg-blue-600"}, 1},
		{"PrimitiveText", []string{"text-slate-800"}, 1},
		{"PrimitiveBorder", []string{"border-gray-200"}, 1},
		{"PrimitiveRing", []string{"ring-emerald-500"}, 1},
		{"VariantPrimitive", []string{"hover:bg-blue-700"}, 1},
		{"ArbitraryPrimitiveVar", []string{"bg-[var(--blue-500)]"}, 1},
		{"SemanticTokens", []string{"bg-primary", "text-card-foreground", "border-border"}, 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			node := makeNode(tc.classes...)
			diags := rule.Evaluate(node)
			if len(diags) != tc.want {
				t.Fatalf("expected %d diagnostics, got %d: %+v", tc.want, len(diags), diags)
			}
		})
	}

	if diags := rule.Evaluate(nil); len(diags) != 0 {
		t.Errorf("expected 0 diags on nil node")
	}
}

func TestHardcodeMonochromeRule(t *testing.T) {
	rule := theme.NewHardcodeMonochromeRule()

	if rule.ID() != "theme.hardcode-monochrome" {
		t.Fatalf("unexpected ID: %s", rule.ID())
	}
	if rule.DefaultSeverity() != ir.SeverityWarn {
		t.Fatalf("unexpected severity: %v", rule.DefaultSeverity())
	}

	tests := []struct {
		name    string
		classes []string
		want    int
	}{
		{"StaticWhite", []string{"bg-white"}, 1},
		{"StaticBlack", []string{"text-black"}, 1},
		{"AlphaBlack", []string{"bg-black/50"}, 1},
		{"AlphaWhite", []string{"text-white/[0.06]"}, 1},
		{"DirectionalWhite", []string{"border-t-white"}, 1},
		{"SemanticTokens", []string{"bg-background", "text-foreground", "bg-card"}, 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			node := makeNode(tc.classes...)
			diags := rule.Evaluate(node)
			if len(diags) != tc.want {
				t.Fatalf("expected %d diagnostics, got %d: %+v", tc.want, len(diags), diags)
			}
		})
	}

	if diags := rule.Evaluate(nil); len(diags) != 0 {
		t.Errorf("expected 0 diags on nil node")
	}
}

func TestHardcodeBorderColorRule(t *testing.T) {
	rule := theme.NewHardcodeBorderColorRule()

	if rule.ID() != "theme.hardcode-border-color" {
		t.Fatalf("unexpected ID: %s", rule.ID())
	}
	if rule.DefaultSeverity() != ir.SeverityWarn {
		t.Fatalf("unexpected severity: %v", rule.DefaultSeverity())
	}

	tests := []struct {
		name    string
		classes []string
		want    int
	}{
		{"PrimitiveBorder", []string{"border-gray-200"}, 1},
		{"ArbitraryBorderHex", []string{"border-[#e5e5e5]"}, 1},
		{"DirectionalBorderPrimitive", []string{"border-t-slate-300"}, 1},
		{"DividePrimitive", []string{"divide-gray-200"}, 1},
		{"MonochromeBorder", []string{"border-white"}, 1},
		{"SemanticBorder", []string{"border-border", "border-input", "divide-border"}, 0},
		{"BorderWidths", []string{"border", "border-0", "border-2", "border-4"}, 0},
		{"BorderStyles", []string{"border-solid", "border-dashed", "border-none"}, 0},
		{"BorderCollapse", []string{"border-collapse"}, 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			node := makeNode(tc.classes...)
			diags := rule.Evaluate(node)
			if len(diags) != tc.want {
				t.Fatalf("expected %d diagnostics, got %d: %+v", tc.want, len(diags), diags)
			}
		})
	}

	if diags := rule.Evaluate(nil); len(diags) != 0 {
		t.Errorf("expected 0 diags on nil node")
	}
}

func TestGradientHardcodeRule(t *testing.T) {
	rule := theme.NewGradientHardcodeRule()

	if rule.ID() != "theme.gradient-hardcode" {
		t.Fatalf("unexpected ID: %s", rule.ID())
	}
	if rule.DefaultSeverity() != ir.SeverityWarn {
		t.Fatalf("unexpected severity: %v", rule.DefaultSeverity())
	}

	tests := []struct {
		name    string
		classes []string
		want    int
	}{
		{"ArbitraryHexStop", []string{"from-[#3b82f6]"}, 1},
		{"PrimitiveStop", []string{"to-blue-500"}, 1},
		{"ViaPrimitive", []string{"via-purple-600"}, 1},
		{"MonochromeStop", []string{"from-white"}, 1},
		{"SemanticStops", []string{"from-primary", "to-secondary", "via-accent"}, 0},
		{"TransparentStops", []string{"from-transparent", "to-transparent"}, 0},
		{"PercentageStops", []string{"from-10%", "via-50%", "to-90%"}, 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			node := makeNode(tc.classes...)
			diags := rule.Evaluate(node)
			if len(diags) != tc.want {
				t.Fatalf("expected %d diagnostics, got %d: %+v", tc.want, len(diags), diags)
			}
		})
	}

	if diags := rule.Evaluate(nil); len(diags) != 0 {
		t.Errorf("expected 0 diags on nil node")
	}
}

func TestInlineStyleHardcodeRule(t *testing.T) {
	rule := theme.NewInlineStyleHardcodeRule()

	if rule.ID() != "theme.inline-style-hardcode" {
		t.Fatalf("unexpected ID: %s", rule.ID())
	}
	if rule.DefaultSeverity() != ir.SeverityError {
		t.Fatalf("unexpected severity: %v", rule.DefaultSeverity())
	}

	tests := []struct {
		name  string
		attrs map[string]string
		want  int
	}{
		{"HexInHTMLStyle", map[string]string{"style": "color: #2563eb; background: #ffffff;"}, 1},
		{"RGBInJSXStyle", map[string]string{"style": "{{ color: '#2563eb', backgroundColor: 'rgb(255, 0, 0)' }}"}, 1},
		{"NamedColorInStyle", map[string]string{"style": "color: red;"}, 1},
		{"SafeDynamicTransform", map[string]string{"style": "top: 0; transform: translateY(-50%);"}, 0},
		{"SafeCSSVar", map[string]string{"style": "color: var(--primary); background: var(--background);"}, 0},
		{"SafeInherit", map[string]string{"style": "color: inherit; background: transparent;"}, 0},
		{"EmptyStyle", map[string]string{"style": ""}, 0},
		{"NoStyle", map[string]string{"class": "text-primary"}, 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			node := makeNodeWithAttr(tc.attrs)
			diags := rule.Evaluate(node)
			if len(diags) != tc.want {
				t.Fatalf("expected %d diagnostics, got %d: %+v", tc.want, len(diags), diags)
			}
		})
	}

	if diags := rule.Evaluate(nil); len(diags) != 0 {
		t.Errorf("expected 0 diags on nil node")
	}
}

func TestPseudoHardcodeColorRule(t *testing.T) {
	rule := theme.NewPseudoHardcodeColorRule()

	if rule.ID() != "theme.pseudo-hardcode-color" {
		t.Fatalf("unexpected ID: %s", rule.ID())
	}
	if rule.DefaultSeverity() != ir.SeverityWarn {
		t.Fatalf("unexpected severity: %v", rule.DefaultSeverity())
	}

	tests := []struct {
		name    string
		classes []string
		want    int
	}{
		{"PlaceholderPrimitive", []string{"placeholder:text-gray-400"}, 1},
		{"SelectionPrimitive", []string{"selection:bg-blue-200"}, 1},
		{"FilePrimitive", []string{"file:bg-blue-50"}, 1},
		{"MarkerPrimitive", []string{"marker:text-slate-500"}, 1},
		{"SelectionHex", []string{"selection:text-[#2563eb]"}, 1},
		{"PlaceholderMonochrome", []string{"placeholder:text-white"}, 1},
		{"PlaceholderSemantic", []string{"placeholder:text-muted-foreground"}, 0},
		{"SelectionSemantic", []string{"selection:bg-primary-light"}, 0},
		{"NonPseudoVariant", []string{"hover:bg-blue-600"}, 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			node := makeNode(tc.classes...)
			diags := rule.Evaluate(node)
			if len(diags) != tc.want {
				t.Fatalf("expected %d diagnostics, got %d: %+v", tc.want, len(diags), diags)
			}
		})
	}

	if diags := rule.Evaluate(nil); len(diags) != 0 {
		t.Errorf("expected 0 diags on nil node")
	}
}

func TestImportantOverrideRule(t *testing.T) {
	rule := theme.NewImportantOverrideRule()

	if rule.ID() != "theme.important-override" {
		t.Fatalf("unexpected ID: %s", rule.ID())
	}
	if rule.DefaultSeverity() != ir.SeverityError {
		t.Fatalf("unexpected severity: %v", rule.DefaultSeverity())
	}

	tests := []struct {
		name    string
		classes []string
		want    int
	}{
		{"ImportantBgPrimitive", []string{"!bg-red-500"}, 1},
		{"ImportantTextMonochrome", []string{"!text-white"}, 1},
		{"ImportantBgSemantic", []string{"!bg-primary"}, 1},
		{"ImportantBorderSemantic", []string{"!border-border"}, 1},
		{"VariantImportantBg", []string{"hover:!bg-blue-500"}, 1},
		{"ImportantLayoutBlock", []string{"!block"}, 0},
		{"ImportantDisplayHidden", []string{"!hidden"}, 0},
		{"ImportantWidth", []string{"!w-full"}, 0},
		{"ImportantPadding", []string{"!p-4"}, 0},
		{"CleanColorClass", []string{"bg-primary", "text-foreground"}, 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			node := makeNode(tc.classes...)
			diags := rule.Evaluate(node)
			if len(diags) != tc.want {
				t.Fatalf("expected %d diagnostics, got %d: %+v", tc.want, len(diags), diags)
			}
		})
	}

	if diags := rule.Evaluate(nil); len(diags) != 0 {
		t.Errorf("expected 0 diags on nil node")
	}
}

func BenchmarkKelompok1_ZeroAllocClean(b *testing.B) {
	node := makeNode("p-4", "flex", "items-center", "justify-between", "rounded-lg", "shadow-sm")
	allRules := []struct {
		name string
		eval func(*ir.Node) []ir.Diagnostic
	}{
		{"HardcodeColor", theme.NewHardcodeColorRule().Evaluate},
		{"PrimitiveInComponent", theme.NewPrimitiveInComponentRule().Evaluate},
		{"HardcodeMonochrome", theme.NewHardcodeMonochromeRule().Evaluate},
		{"HardcodeBorderColor", theme.NewHardcodeBorderColorRule().Evaluate},
		{"GradientHardcode", theme.NewGradientHardcodeRule().Evaluate},
		{"InlineStyleHardcode", theme.NewInlineStyleHardcodeRule().Evaluate},
		{"PseudoHardcodeColor", theme.NewPseudoHardcodeColorRule().Evaluate},
		{"ImportantOverride", theme.NewImportantOverrideRule().Evaluate},
	}

	for _, tc := range allRules {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = tc.eval(node)
			}
		})
	}
}
