package a11y

import (
	"fmt"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// InputCrampedPaddingRule mendeteksi kontrol input dengan tinggi atau bantalan vertikal kerdil (< 42px)
// yang memotong aksen tipografi dan mempersulit sentuhan jari pengguna.
type InputCrampedPaddingRule struct{}

// NewInputCrampedPaddingRule membuat instance baru InputCrampedPaddingRule.
func NewInputCrampedPaddingRule() *InputCrampedPaddingRule {
	return &InputCrampedPaddingRule{}
}

// ID mengembalikan Charites Rule ID kanonikal berformat a11y.input-cramped-padding.
func (r *InputCrampedPaddingRule) ID() string {
	return "a11y.input-cramped-padding"
}

// Description mengembalikan deskripsi ringkas aturan.
func (r *InputCrampedPaddingRule) Description() string {
	return "Flags input controls with cramped vertical padding or height under 42px that clip text and impede touch targeting"
}

// Category mengembalikan nama kategori rule.
func (r *InputCrampedPaddingRule) Category() string {
	return "a11y"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warning).
func (r *InputCrampedPaddingRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *InputCrampedPaddingRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"WCAG 2.2 Success Criterion 2.5.8 (Target Size - Minimum)",
			"Apple Human Interface Guidelines (Form Controls & Hit Targets: 44pt minimum)",
			"Material Design 3 (Text Fields: 48dp baseline, 40-42dp dense minimum)",
		},
		CoreInvariant: "Form input controls (<input>, <select>, <textarea>) must maintain an effective height of at least 42px with adequate vertical padding ('py-2.5' or 'h-11') to ensure touch comfort and prevent text clipping.",
		Grounding: "Text inputs require sufficient spatial breathing room for cursor positioning, focus states, and internationalized typography ascenders and descenders (e.g. g, j, y, and diacritics).\n\n" +
			"When form controls specify undersized fixed heights like `h-7` (28px) or `h-8` (32px), or cramped vertical padding such as `py-1` (4px) or `p-1`:\n" +
			"1. Motor Targeting Errors: Users frequently miss the input tap area, accidentally tapping adjacent elements or the form background.\n" +
			"2. Glyph Clipping: Font descenders collide with or are cut off by the field's bottom border.\n" +
			"3. Inconsistent Visual Hierarchy: The input appears visually shrunk and uninviting compared to system controls.\n\n" +
			"Charites checks for undersized explicit heights (< 42px) and cramped padding without minimum height compensation.",
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Cramped 32px height input with tight padding",
				Code:     `<input className="h-8 px-2 text-sm border rounded" placeholder="Nama lengkap" />`,
			},
			{
				Language: "astro",
				Comment:  "Input with cramped py-1 vertical padding",
				Code:     `<input class="py-1 px-2 text-xs border rounded" placeholder="Kode pos" />`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Ergonomic 44px height input with balanced padding",
				Code:     `<input className="h-11 px-3.5 py-2.5 text-base sm:text-sm border rounded-lg" placeholder="Nama lengkap" />`,
			},
			{
				Language: "astro",
				Comment:  "Spacious 10px vertical padding input ensuring > 44px total height",
				Code:     `<input class="py-2.5 px-3.5 text-base border rounded-md" placeholder="Kode pos" />`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "Touch Targeting Failure",
				Severity: "HIGH",
				Impact:   "Users experience difficulty tapping inputs on mobile touchscreens due to cramped hit targets.",
			},
			{
				Vector:   "Text & Glyph Clipping",
				Severity: "MEDIUM",
				Impact:   "Letters with descenders or accents are visually truncated by field borders.",
			},
		},
	}
}

// Evaluate memeriksa tinggi dan bantalan vertikal pada kontrol isian formulir.
// Mematuhi 0 B/op, 0 allocs/op pada node bersih (QUAL-03).
func (r *InputCrampedPaddingRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || len(node.Classes) == 0 {
		return nil
	}

	if !IsTextualInput(node) {
		return nil
	}

	state := inspectInputDimensions(node.Classes)

	// Jika ada kompensasi min-h >= 42px, elemen terlindungi
	if state.hasMinHCompensation {
		return nil
	}

	// Kasus 1: Tinggi eksplisit di bawah 42px
	if state.hasHeight && state.explicitHeight < 42.0 {
		return []ir.Diagnostic{
			{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  fmt.Sprintf("Input <%s> has cramped height '%s' (%.0fpx < 42px), impeding touch targeting and text readability", node.Tag, state.heightClass, state.explicitHeight),
				Hint:     "Use at least 'h-11' (2.75rem / 44px) or 'py-2.5' for accessible form controls (WCAG 2.5.8 / Apple HIG).",
			},
		}
	}

	// Kasus 2: Tanpa tinggi eksplisit, tetapi padding vertikal terlalu sempit (< 14px total padding)
	// Default line-height (24px) + vertical padding (< 14px) = total tinggi < 38-42px
	if !state.hasHeight && state.hasVertPad && state.totalVertPad < 14.0 {
		return []ir.Diagnostic{
			{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  fmt.Sprintf("Input <%s> has cramped vertical padding '%s' (total padding %.0fpx), resulting in height under 42px", node.Tag, state.padClass, state.totalVertPad),
				Hint:     "Provide at least 'py-2.5' or 'h-11' to guarantee comfortable touch targets and prevent glyph clipping.",
			},
		}
	}

	return nil
}

type inputDimState struct {
	explicitHeight      float64
	hasHeight           bool
	heightClass         string
	totalVertPad        float64
	hasVertPad          bool
	padClass            string
	hasMinHCompensation bool
}

func inspectInputDimensions(classes []string) inputDimState {
	var state inputDimState

	for _, class := range classes {
		base := StripVariantsOnlyBase(class)
		if checkInputMinHeight(base, &state) {
			continue
		}
		if checkInputHeight(base, class, &state) {
			continue
		}
		checkInputPadding(base, class, &state)
	}

	return state
}

func checkInputMinHeight(base string, state *inputDimState) bool {
	if strings.HasPrefix(base, "min-h-") {
		sub := strings.TrimPrefix(base, "min-h-")
		if px, ok := ParseTailwindSizeToPx(sub); ok && px >= 42.0 {
			state.hasMinHCompensation = true
			return true
		}
	}
	return false
}

func checkInputHeight(base, class string, state *inputDimState) bool {
	if strings.HasPrefix(base, "h-") {
		sub := strings.TrimPrefix(base, "h-")
		if px, ok := ParseTailwindSizeToPx(sub); ok {
			state.explicitHeight = px
			state.hasHeight = true
			state.heightClass = class
			return true
		}
	}
	return false
}

func checkInputPadding(base, class string, state *inputDimState) {
	if strings.HasPrefix(base, "py-") {
		sub := strings.TrimPrefix(base, "py-")
		if px, ok := ParseTailwindSizeToPx(sub); ok {
			state.totalVertPad = px * 2
			state.hasVertPad = true
			state.padClass = class
		}
		return
	}

	if strings.HasPrefix(base, "p-") && !strings.HasPrefix(base, "px-") &&
		!strings.HasPrefix(base, "pt-") && !strings.HasPrefix(base, "pb-") &&
		!strings.HasPrefix(base, "pl-") && !strings.HasPrefix(base, "pr-") {
		sub := strings.TrimPrefix(base, "p-")
		if px, ok := ParseTailwindSizeToPx(sub); ok {
			state.totalVertPad = px * 2
			state.hasVertPad = true
			state.padClass = class
		}
	}
}
