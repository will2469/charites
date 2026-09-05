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

func TestHardcodeSizeRule(t *testing.T) {
	rule := theme.NewHardcodeSizeRule()

	if rule.ID() != "theme.hardcode-size" {
		t.Fatalf("unexpected ID: %s", rule.ID())
	}
	if rule.Category() != "theme" {
		t.Fatalf("unexpected category: %s", rule.Category())
	}
	if rule.DefaultSeverity() != ir.SeverityWarn {
		t.Fatalf("unexpected severity: %v", rule.DefaultSeverity())
	}

	tests := []struct {
		name    string
		classes []string
		want    int
	}{
		{"ArbitraryPadding", []string{"p-[19px]"}, 1},
		{"ArbitraryWidth", []string{"w-[320px]"}, 1},
		{"ArbitraryTextSize", []string{"text-[15px]"}, 1},
		{"ArbitraryGap", []string{"gap-[13px]"}, 1},
		{"ArbitraryPropertyWidth", []string{"[width:100px]"}, 1},
		{"ArbitraryPropertyPadding", []string{"[padding:19px]"}, 1},
		{"StandardSpacing", []string{"p-5", "w-80", "text-base", "gap-3"}, 0},
		{"StandardFractions", []string{"w-full", "w-1/2"}, 0},
		{"CSSVariableSize", []string{"p-[var(--spacing-custom)]", "w-[var(--container-max)]"}, 0},
		{"BananaNonSpatial", []string{"scale-[1.5]", "duration-[300ms]", "rotate-[45deg]"}, 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			node := makeNode(tc.classes...)
			diags := rule.Evaluate(node)
			if len(diags) != tc.want {
				t.Fatalf("got %d diagnostics, want %d", len(diags), tc.want)
			}
		})
	}
}

func TestHardcodeBorderRadiusRule(t *testing.T) {
	rule := theme.NewHardcodeBorderRadiusRule()

	if rule.ID() != "theme.hardcode-border-radius" {
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
		{"ArbitraryRadiusPx", []string{"rounded-[7px]"}, 1},
		{"ArbitraryRadiusRem", []string{"rounded-[0.5rem]"}, 1},
		{"DirectionalRadius", []string{"rounded-t-[10px]"}, 1},
		{"ArbitraryPropertyRadius", []string{"[border-radius:8px]"}, 1},
		{"StandardTokens", []string{"rounded-none", "rounded-sm", "rounded", "rounded-md", "rounded-lg", "rounded-full"}, 0},
		{"StandardDirectional", []string{"rounded-t-xl", "rounded-b-none"}, 0},
		{"CSSVariableRadius", []string{"rounded-[var(--radius)]"}, 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			node := makeNode(tc.classes...)
			diags := rule.Evaluate(node)
			if len(diags) != tc.want {
				t.Fatalf("got %d diagnostics, want %d", len(diags), tc.want)
			}
		})
	}
}

func TestHardcodeZIndexRule(t *testing.T) {
	rule := theme.NewHardcodeZIndexRule()

	if rule.ID() != "theme.hardcode-z-index" {
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
		{"ArbitraryLargeZ", []string{"z-[9999]"}, 1},
		{"ArbitraryZHundred", []string{"z-[100]"}, 1},
		{"ArbitraryPropertyZ", []string{"[z-index:1000]"}, 1},
		{"StandardZScale", []string{"z-0", "z-10", "z-20", "z-30", "z-40", "z-50"}, 0},
		{"ZAuto", []string{"z-auto"}, 0},
		{"CSSVariableZ", []string{"z-[var(--z-modal)]"}, 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			node := makeNode(tc.classes...)
			diags := rule.Evaluate(node)
			if len(diags) != tc.want {
				t.Fatalf("got %d diagnostics, want %d", len(diags), tc.want)
			}
		})
	}
}

func TestHardcodeShadowColorRule(t *testing.T) {
	rule := theme.NewHardcodeShadowColorRule()

	if rule.ID() != "theme.hardcode-shadow-color" {
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
		{"ArbitraryShadowHex", []string{"shadow-[0_4px_10px_#00000040]"}, 1},
		{"ArbitraryShadowRGBA", []string{"shadow-[0_10px_15px_rgba(0,0,0,0.1)]"}, 1},
		{"ArbitraryPropertyShadow", []string{"[box-shadow:0_4px_6px_#000]"}, 1},
		{"StandardShadows", []string{"shadow-sm", "shadow", "shadow-md", "shadow-lg", "shadow-xl", "shadow-none"}, 0},
		{"CSSVariableShadowColor", []string{"shadow-[0_4px_6px_var(--shadow-color)]"}, 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			node := makeNode(tc.classes...)
			diags := rule.Evaluate(node)
			if len(diags) != tc.want {
				t.Fatalf("got %d diagnostics, want %d", len(diags), tc.want)
			}
		})
	}
}

func TestBackdropBlurHardcodeRule(t *testing.T) {
	rule := theme.NewBackdropBlurHardcodeRule()

	if rule.ID() != "theme.backdrop-blur-hardcode" {
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
		{"ArbitraryBackdropBlur", []string{"backdrop-blur-[5px]"}, 1},
		{"ArbitraryFilterBlur", []string{"blur-[12px]"}, 1},
		{"ArbitraryPropertyBackdropBlur", []string{"[backdrop-filter:blur(7px)]"}, 1},
		{"StandardBackdropBlur", []string{"backdrop-blur-sm", "backdrop-blur-md", "backdrop-blur-lg", "backdrop-blur-none"}, 0},
		{"StandardBlur", []string{"blur-sm", "blur-md", "blur-none"}, 0},
		{"CSSVariableBlur", []string{"backdrop-blur-[var(--blur-glass)]"}, 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			node := makeNode(tc.classes...)
			diags := rule.Evaluate(node)
			if len(diags) != tc.want {
				t.Fatalf("got %d diagnostics, want %d", len(diags), tc.want)
			}
		})
	}
}

func TestFocusRingHardcodeRule(t *testing.T) {
	rule := theme.NewFocusRingHardcodeRule()

	if rule.ID() != "theme.focus-ring-hardcode" {
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
		{"ArbitraryHexRing", []string{"focus:ring-[#3b82f6]"}, 1},
		{"PrimitiveColorRing", []string{"ring-blue-500"}, 1},
		{"PrimitiveOffsetMonochrome", []string{"ring-offset-white"}, 1},
		{"PrimitiveOutline", []string{"outline-blue-500"}, 1},
		{"ArbitraryOutlineHex", []string{"outline-[#3b82f6]"}, 1},
		{"SemanticRing", []string{"focus-visible:ring-ring", "focus:ring-ring", "ring-ring", "ring-primary"}, 0},
		{"SemanticOffset", []string{"ring-offset-background"}, 0},
		{"RingWidthsAndKeywords", []string{"ring", "ring-1", "ring-2", "ring-4", "ring-offset-2", "outline-none"}, 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			node := makeNode(tc.classes...)
			diags := rule.Evaluate(node)
			if len(diags) != tc.want {
				t.Fatalf("got %d diagnostics, want %d", len(diags), tc.want)
			}
		})
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

func BenchmarkKelompok2_ZeroAllocClean(b *testing.B) {
	node := makeNode("p-4", "w-80", "text-base", "gap-3", "rounded-md", "z-10", "shadow-md", "backdrop-blur-sm", "focus-visible:ring-2", "focus-visible:ring-ring")
	rulesK2 := []struct {
		name string
		eval func(*ir.Node) []ir.Diagnostic
	}{
		{"HardcodeSize", theme.NewHardcodeSizeRule().Evaluate},
		{"HardcodeBorderRadius", theme.NewHardcodeBorderRadiusRule().Evaluate},
		{"HardcodeZIndex", theme.NewHardcodeZIndexRule().Evaluate},
		{"HardcodeShadowColor", theme.NewHardcodeShadowColorRule().Evaluate},
		{"BackdropBlurHardcode", theme.NewBackdropBlurHardcodeRule().Evaluate},
		{"FocusRingHardcode", theme.NewFocusRingHardcodeRule().Evaluate},
	}

	for _, tc := range rulesK2 {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = tc.eval(node)
			}
		})
	}
}

func TestUnpairedDarkVariantRule(t *testing.T) {
	rule := theme.NewUnpairedDarkVariantRule()

	if rule.ID() != "theme.unpaired-dark-variant" {
		t.Fatalf("unexpected ID: %s", rule.ID())
	}
	if rule.DefaultSeverity() != ir.SeverityWarn {
		t.Fatalf("unexpected severity: %v", rule.DefaultSeverity())
	}

	t.Run("UnpairedDarkBgWithoutBase", func(t *testing.T) {
		node := makeNode("dark:bg-zinc-900")
		diags := rule.Evaluate(node)
		if len(diags) != 1 {
			t.Fatalf("expected 1 diag, got %d: %+v", len(diags), diags)
		}
	})

	t.Run("InvertedContainerWithChildTextBlind", func(t *testing.T) {
		child := makeNode("text-zinc-900")
		parent := makeNode("bg-white", "dark:bg-zinc-900")
		parent.Children = []*ir.Node{child}
		child.Parent = parent

		diags := rule.Evaluate(parent)
		if len(diags) != 1 {
			t.Fatalf("expected 1 diag on parent evaluating child, got %d: %+v", len(diags), diags)
		}
	})

	t.Run("CompliantSemanticTokens", func(t *testing.T) {
		child := makeNode("text-card-foreground")
		parent := makeNode("bg-card")
		parent.Children = []*ir.Node{child}
		child.Parent = parent

		diags := rule.Evaluate(parent)
		if len(diags) != 0 {
			t.Fatalf("expected 0 diags, got %d: %+v", len(diags), diags)
		}
	})

	t.Run("CompliantPairedVariants", func(t *testing.T) {
		child := makeNode("text-zinc-900", "dark:text-zinc-100")
		parent := makeNode("bg-white", "dark:bg-zinc-900")
		parent.Children = []*ir.Node{child}
		child.Parent = parent

		diags := rule.Evaluate(parent)
		if len(diags) != 0 {
			t.Fatalf("expected 0 diags, got %d: %+v", len(diags), diags)
		}
	})
}

func TestShadowWithoutBorderDarkRule(t *testing.T) {
	rule := theme.NewShadowWithoutBorderDarkRule()

	if rule.ID() != "theme.shadow-without-border-dark" {
		t.Fatalf("unexpected ID: %s", rule.ID())
	}

	t.Run("ElevatedShadowWithoutBorder", func(t *testing.T) {
		node := makeNode("bg-card", "shadow-xl", "rounded-xl", "p-6")
		diags := rule.Evaluate(node)
		if len(diags) != 1 {
			t.Fatalf("expected 1 diag, got %d: %+v", len(diags), diags)
		}
	})

	t.Run("ElevatedShadowWithBorder", func(t *testing.T) {
		node := makeNode("bg-card", "border", "border-border", "shadow-xl", "rounded-xl", "p-6")
		diags := rule.Evaluate(node)
		if len(diags) != 0 {
			t.Fatalf("expected 0 diags, got %d: %+v", len(diags), diags)
		}
	})

	t.Run("ElevatedShadowWithRing", func(t *testing.T) {
		node := makeNode("bg-zinc-900", "ring-1", "ring-border", "shadow-lg", "p-4")
		diags := rule.Evaluate(node)
		if len(diags) != 0 {
			t.Fatalf("expected 0 diags, got %d: %+v", len(diags), diags)
		}
	})

	t.Run("LowElevationShadowSmClean", func(t *testing.T) {
		node := makeNode("bg-card", "shadow-sm", "p-4")
		diags := rule.Evaluate(node)
		if len(diags) != 0 {
			t.Fatalf("expected 0 diags, got %d: %+v", len(diags), diags)
		}
	})
}

func TestNestedOpacityContrastRule(t *testing.T) {
	rule := theme.NewNestedOpacityContrastRule()

	if rule.ID() != "theme.nested-opacity-contrast" {
		t.Fatalf("unexpected ID: %s", rule.ID())
	}

	t.Run("CompoundedOpacityViolation", func(t *testing.T) {
		child := makeNode("text-foreground/50")
		parent := makeNode("bg-muted/40", "opacity-80")
		parent.Children = []*ir.Node{child}
		child.Parent = parent

		diags := rule.Evaluate(parent)
		if len(diags) != 1 {
			t.Fatalf("expected 1 diag, got %d: %+v", len(diags), diags)
		}
	})

	t.Run("CompliantSolidTokens", func(t *testing.T) {
		child := makeNode("text-muted-foreground")
		parent := makeNode("bg-muted")
		parent.Children = []*ir.Node{child}
		child.Parent = parent

		diags := rule.Evaluate(parent)
		if len(diags) != 0 {
			t.Fatalf("expected 0 diags, got %d: %+v", len(diags), diags)
		}
	})
}

func TestImageThemeHardcodeRule(t *testing.T) {
	rule := theme.NewImageThemeHardcodeRule()

	if rule.ID() != "theme.image-theme-hardcode" {
		t.Fatalf("unexpected ID: %s", rule.ID())
	}

	t.Run("StaticSVGLogoViolation", func(t *testing.T) {
		node := makeNodeWithAttr(map[string]string{"src": "/images/logo-black.svg", "alt": "Logo"})
		node.Tag = "img"
		diags := rule.Evaluate(node)
		if len(diags) != 1 {
			t.Fatalf("expected 1 diag, got %d: %+v", len(diags), diags)
		}
	})

	t.Run("ThemePairedImgClean", func(t *testing.T) {
		node := makeNodeWithAttr(map[string]string{"src": "/images/logo-light.svg", "alt": "Logo"}, "dark:hidden")
		node.Tag = "img"
		diags := rule.Evaluate(node)
		if len(diags) != 0 {
			t.Fatalf("expected 0 diags, got %d: %+v", len(diags), diags)
		}
	})

	t.Run("InsidePictureClean", func(t *testing.T) {
		picture := &ir.Node{Tag: "picture"}
		img := makeNodeWithAttr(map[string]string{"src": "/images/logo.svg"})
		img.Tag = "img"
		img.Parent = picture
		picture.Children = []*ir.Node{img}

		diags := rule.Evaluate(img)
		if len(diags) != 0 {
			t.Fatalf("expected 0 diags, got %d: %+v", len(diags), diags)
		}
	})

	t.Run("RegularPhotoClean", func(t *testing.T) {
		node := makeNodeWithAttr(map[string]string{"src": "/photos/avatar.jpg", "alt": "Avatar"})
		node.Tag = "img"
		diags := rule.Evaluate(node)
		if len(diags) != 0 {
			t.Fatalf("expected 0 diags, got %d: %+v", len(diags), diags)
		}
	})
}

func TestSVGHardcodeFillRule(t *testing.T) {
	rule := theme.NewSVGHardcodeFillRule()

	if rule.ID() != "theme.svg-hardcode-fill" {
		t.Fatalf("unexpected ID: %s", rule.ID())
	}

	t.Run("HardcodedHexFill", func(t *testing.T) {
		node := makeNodeWithAttr(map[string]string{"fill": "#000000", "d": "M10 10"})
		node.Tag = "path"
		diags := rule.Evaluate(node)
		if len(diags) != 1 {
			t.Fatalf("expected 1 diag, got %d: %+v", len(diags), diags)
		}
	})

	t.Run("HardcodedStopColor", func(t *testing.T) {
		node := makeNodeWithAttr(map[string]string{"stop-color": "#3b82f6", "offset": "100%"})
		node.Tag = "stop"
		diags := rule.Evaluate(node)
		if len(diags) != 1 {
			t.Fatalf("expected 1 diag, got %d: %+v", len(diags), diags)
		}
	})

	t.Run("CurrentColorClean", func(t *testing.T) {
		node := makeNodeWithAttr(map[string]string{"fill": "currentColor", "stroke": "none"})
		node.Tag = "path"
		diags := rule.Evaluate(node)
		if len(diags) != 0 {
			t.Fatalf("expected 0 diags, got %d: %+v", len(diags), diags)
		}
	})

	t.Run("CSSVarClean", func(t *testing.T) {
		node := makeNodeWithAttr(map[string]string{"stop-color": "var(--primary)", "offset": "100%"})
		node.Tag = "stop"
		diags := rule.Evaluate(node)
		if len(diags) != 0 {
			t.Fatalf("expected 0 diags, got %d: %+v", len(diags), diags)
		}
	})
}

func BenchmarkKelompok3_ZeroAllocClean(b *testing.B) {
	node := makeNode("bg-card", "border", "border-border", "text-card-foreground", "p-6", "rounded-xl")
	rulesK3 := []struct {
		name string
		eval func(*ir.Node) []ir.Diagnostic
	}{
		{"UnpairedDarkVariant", theme.NewUnpairedDarkVariantRule().Evaluate},
		{"ShadowWithoutBorderDark", theme.NewShadowWithoutBorderDarkRule().Evaluate},
		{"NestedOpacityContrast", theme.NewNestedOpacityContrastRule().Evaluate},
		{"ImageThemeHardcode", theme.NewImageThemeHardcodeRule().Evaluate},
		{"SVGHardcodeFill", theme.NewSVGHardcodeFillRule().Evaluate},
	}

	for _, tc := range rulesK3 {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = tc.eval(node)
			}
		})
	}
}
