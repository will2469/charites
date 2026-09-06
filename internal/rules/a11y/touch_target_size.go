package a11y

import (
	"fmt"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// TouchTargetSizeRule memastikan elemen interaktif memenuhi ukuran sentuh minimum
// setidaknya 44x44px sesuai Apple HIG dan WCAG 2.2 SC 2.5.8.
type TouchTargetSizeRule struct{}

// NewTouchTargetSizeRule membuat instance baru TouchTargetSizeRule.
func NewTouchTargetSizeRule() *TouchTargetSizeRule {
	return &TouchTargetSizeRule{}
}

// ID mengembalikan Charites Rule ID kanonikal berformat a11y.touch-target-size.
func (r *TouchTargetSizeRule) ID() string {
	return "a11y.touch-target-size"
}

// Description mengembalikan deskripsi ringkas aturan.
func (r *TouchTargetSizeRule) Description() string {
	return "Enforces minimum 44x44px physical touch target size on interactive controls (Apple HIG / WCAG 2.5.8)"
}

// Category mengembalikan nama kategori rule.
func (r *TouchTargetSizeRule) Category() string {
	return "a11y"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warning).
func (r *TouchTargetSizeRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *TouchTargetSizeRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"WCAG 2.2 Success Criterion 2.5.8 (Target Size - Minimum)",
			"Apple Human Interface Guidelines (Touch Controls: 44x44pt minimum)",
			"Google Material Design (Touch Targets: 48x48dp baseline)",
		},
		CoreInvariant: "Interactive controls (<button>, <a>, <input>, <select>) must provide an effective physical touch target area of at least 44x44px on mobile devices.",
		Grounding: "Touch screens rely on human fingertip interaction, which has an average contact area of 10-14mm (approx. 44-48 CSS pixels).\n\n" +
			"When interactive elements (like icon buttons or close triggers) are sized with explicit small dimensions such as `h-6 w-6` (24px) without hit-box padding compensation or minimum size constraints:\n" +
			"1. High Error Rate: Users frequently miss the intended button or trigger adjacent elements.\n" +
			"2. Motor Accessibility Failures: Users with tremors or motor impairments cannot reliably activate the control.\n" +
			"3. Frustrating Mobile UX: Tap feedback does not register, causing repeated tapping and accidental rage clicks.\n\n" +
			"Charites calculates the effective box model dimensions and flags elements falling below 44px.",
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Small 24px icon button without hit-box compensation",
				Code:     `<button className="h-6 w-6"><TrashIcon className="h-4 w-4" /></button>`,
			},
			{
				Language: "astro",
				Comment:  "Close button with 32px height and width",
				Code:     `<button class="h-8 w-8 p-1"><CloseIcon /></button>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "44px (2.75rem) target with centered inner icon",
				Code:     `<button className="h-11 w-11 flex items-center justify-center"><TrashIcon className="h-5 w-5" /></button>`,
			},
			{
				Language: "astro",
				Comment:  "Explicit minimum dimensions to guarantee 44px (2.75rem) tap zone",
				Code:     `<button class="h-8 w-8 min-h-11 min-w-11 flex items-center justify-center"><CloseIcon /></button>`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "WCAG 2.5.8 Non-Compliance",
				Severity: "HIGH",
				Impact:   "Users with motor disabilities cannot reliably tap interactive elements on mobile devices.",
			},
			{
				Vector:   "Rage Clicks & Accidental Actions",
				Severity: "MEDIUM",
				Impact:   "Users accidentally tap sibling destructive actions due to cramped bounding boxes.",
			},
		},
	}
}

// Evaluate memeriksa ukuran fisik elemen interaktif.
// Mematuhi 0 B/op, 0 allocs/op pada node bersih (QUAL-03).
func (r *TouchTargetSizeRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || len(node.Classes) == 0 {
		return nil
	}

	if !IsInteractiveElement(node) {
		return nil
	}

	// Pengecualian inline link di dalam paragraf teks
	if node.Parent != nil && strings.ToLower(node.Parent.Tag) == "p" && strings.ToLower(node.Tag) == "a" {
		return nil
	}

	// Pengecualian elemen di dalam label (label memperluas touch hit target)
	if node.Parent != nil && strings.ToLower(node.Parent.Tag) == "label" {
		return nil
	}

	// Pengecualian input seleksi native (checkbox, radio, hidden) sesuai WCAG 2.5.8
	if strings.ToLower(node.Tag) == "input" && node.Attributes != nil {
		inputType := strings.Trim(strings.TrimSpace(strings.ToLower(node.Attributes["type"])), "\"'`")
		if inputType == "checkbox" || inputType == "radio" || inputType == "hidden" {
			return nil
		}
	}

	var (
		explicitW float64
		explicitH float64
		hasW      bool
		hasH      bool
		hasMinW44 bool
		hasMinH44 bool
		padX      float64
		padY      float64
		hasPadX   bool
		hasPadY   bool
	)

	for _, class := range node.Classes {
		base := StripVariantsOnlyBase(class)

		// Periksa apakah ada kompensasi min-w atau min-h >= 44px
		if checkMinSizeCompensation(base, &hasMinW44, &hasMinH44) {
			continue
		}

		// Periksa width / height eksplisit
		checkExplicitDimensions(base, &explicitW, &explicitH, &hasW, &hasH)

		// Periksa padding
		checkPadding(base, &padX, &padY, &hasPadX, &hasPadY)
	}

	// Jika sudah ada min-w >= 44px dan min-h >= 44px, ukuran sentuh aman
	if hasMinW44 && hasMinH44 {
		return nil
	}

	// Hitung apakah dimensi melanggar ambang 44px
	if violatesTouchTarget(hasW, hasH, explicitW, explicitH, hasPadX, hasPadY, padX, padY, hasMinW44, hasMinH44) {
		dimDesc := describeDimensions(hasW, hasH, explicitW, explicitH)
		return []ir.Diagnostic{
			{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  fmt.Sprintf("Interactive <%s> touch target size is too small (%s < 44px)", node.Tag, dimDesc),
				Hint:     "Ensure at least 44x44px (2.75rem) target area using 'size-11', 'min-h-11 min-w-11', or padding compensation (Apple HIG / WCAG 2.5.8).",
			},
		}
	}

	return nil
}

func checkMinSizeCompensation(base string, hasMinW44, hasMinH44 *bool) bool {
	if strings.HasPrefix(base, "min-w-") {
		sub := strings.TrimPrefix(base, "min-w-")
		if px, ok := ParseTailwindSizeToPx(sub); ok && px >= 44.0 {
			*hasMinW44 = true
			return true
		}
	}
	if strings.HasPrefix(base, "min-h-") {
		sub := strings.TrimPrefix(base, "min-h-")
		if px, ok := ParseTailwindSizeToPx(sub); ok && px >= 44.0 {
			*hasMinH44 = true
			return true
		}
	}
	return false
}

func checkExplicitDimensions(base string, explicitW, explicitH *float64, hasW, hasH *bool) {
	if strings.HasPrefix(base, "size-") {
		sub := strings.TrimPrefix(base, "size-")
		if px, ok := ParseTailwindSizeToPx(sub); ok {
			*explicitW = px
			*explicitH = px
			*hasW = true
			*hasH = true
		}
		return
	}

	if strings.HasPrefix(base, "w-") {
		sub := strings.TrimPrefix(base, "w-")
		if px, ok := ParseTailwindSizeToPx(sub); ok {
			*explicitW = px
			*hasW = true
		}
		return
	}

	if strings.HasPrefix(base, "h-") {
		sub := strings.TrimPrefix(base, "h-")
		if px, ok := ParseTailwindSizeToPx(sub); ok {
			*explicitH = px
			*hasH = true
		}
	}
}

func checkPadding(base string, padX, padY *float64, hasPadX, hasPadY *bool) {
	if strings.HasPrefix(base, "p-") {
		sub := strings.TrimPrefix(base, "p-")
		if px, ok := ParseTailwindSizeToPx(sub); ok {
			*padX = px * 2
			*padY = px * 2
			*hasPadX = true
			*hasPadY = true
		}
		return
	}
	if strings.HasPrefix(base, "px-") {
		sub := strings.TrimPrefix(base, "px-")
		if px, ok := ParseTailwindSizeToPx(sub); ok {
			*padX = px * 2
			*hasPadX = true
		}
		return
	}
	if strings.HasPrefix(base, "py-") {
		sub := strings.TrimPrefix(base, "py-")
		if px, ok := ParseTailwindSizeToPx(sub); ok {
			*padY = px * 2
			*hasPadY = true
		}
	}
}

func violatesTouchTarget(hasW, hasH bool, explicitW, explicitH float64, hasPadX, hasPadY bool, padX, padY float64, hasMinW44, hasMinH44 bool) bool {
	// Jika ada ukuran eksplisit yang lebih kecil dari 44px
	if hasW && explicitW < 44.0 && !hasMinW44 {
		// Jika padding tidak cukup mengompensasi
		if !hasPadX || (explicitW+padX < 44.0) {
			return true
		}
	}
	if hasH && explicitH < 44.0 && !hasMinH44 {
		if !hasPadY || (explicitH+padY < 44.0) {
			return true
		}
	}
	return false
}

func describeDimensions(hasW, hasH bool, explicitW, explicitH float64) string {
	if hasW && hasH {
		return fmt.Sprintf("%.0fx%.0fpx", explicitW, explicitH)
	}
	if hasW {
		return fmt.Sprintf("width %.0fpx", explicitW)
	}
	if hasH {
		return fmt.Sprintf("height %.0fpx", explicitH)
	}
	return "dimensions"
}
