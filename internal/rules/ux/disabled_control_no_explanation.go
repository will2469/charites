package ux

import (
	"fmt"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// DisabledControlNoExplanationRule memastikan kontrol interaktif yang dinonaktifkan (disabled)
// memiliki petunjuk feedforward (aria-describedby, tooltip, atau teks penjelas) untuk menerangkan prasyarat yang belum terpenuhi.
type DisabledControlNoExplanationRule struct{}

// NewDisabledControlNoExplanationRule membuat instance baru dari DisabledControlNoExplanationRule.
func NewDisabledControlNoExplanationRule() *DisabledControlNoExplanationRule {
	return &DisabledControlNoExplanationRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *DisabledControlNoExplanationRule) ID() string {
	return "ux.disabled-control-no-explanation"
}

// Description mengembalikan ringkasan aturan.
func (r *DisabledControlNoExplanationRule) Description() string {
	return "Enforces feedforward explanation for disabled interactive controls to prevent user dead ends"
}

// Category mengembalikan nama kategori rule.
func (r *DisabledControlNoExplanationRule) Category() string {
	return "ux"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warning).
func (r *DisabledControlNoExplanationRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *DisabledControlNoExplanationRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Nielsen Heuristic #1: Visibility of System Status",
			"Feedforward Principle & Gulf of Evaluation (Don Norman)",
			"WCAG 2.2 Success Criterion 3.3.2 (Labels or Instructions)",
		},
		CoreInvariant: "Disabled interactive controls ('disabled', 'aria-disabled=\"true\"') must provide a feedforward explanation via 'aria-describedby', tooltip, or visible helper text.",
		Grounding: "When users encounter a disabled button or locked form control without any visible reason, " +
			"they reach a cognitive dead end: they cannot proceed and have no information on how to unlock the action.\n\n" +
			"Providing contextual feedforward (e.g. 'Minimum belanja Rp 50.000 untuk checkout' or an explanatory tooltip) " +
			"clarifies system constraints and transforms frustration into actionable user guidance.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Cognitive Dead Ends & Workflow Abandonment",
				Severity: "MEDIUM",
				Impact:   "Users are unable to diagnose why an action is blocked and abandon conversion or submission workflows.",
			},
			{
				Vector:   "Accessibility Barrier (Screen Reader Confusion)",
				Severity: "MEDIUM",
				Impact:   "Assisted technology users hear 'button disabled' without auditory explanation of required missing inputs.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Disabled checkout button without explanation of minimum order constraint",
				Code: `<div className="mt-4">
  <button disabled={cartTotal < 50000} className="bg-primary text-white px-4 py-2 rounded">
    Checkout
  </button>
</div>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Disabled button accompanied by aria-describedby linking to explanatory hint text",
				Code: `<div className="mt-4">
  <button
    disabled={cartTotal < 50000}
    aria-describedby="min-order-hint"
    className="bg-primary text-white px-4 py-2 rounded disabled:opacity-50"
  >
    Checkout
  </button>
  {cartTotal < 50000 && (
    <p id="min-order-hint" className="text-xs text-muted-foreground mt-1">
      Minimum belanja Rp 50.000 untuk melanjutkan checkout.
    </p>
  )}
</div>`,
			},
		},
	}
}

// Evaluate memeriksa apakah kontrol disabled memiliki penjelasan feedforward.
func (r *DisabledControlNoExplanationRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	if !isInteractiveControl(node) {
		return nil
	}

	if !isDisabledControl(node) {
		return nil
	}

	if hasFeedforwardExplanation(node) {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message: fmt.Sprintf(
				"Interactive control <%s> is disabled without feedforward explanation (missing aria-describedby, tooltip, or visible reason).",
				node.Tag,
			),
			Hint: "Provide an explanation using aria-describedby, title/tooltip, or nearby helper text explaining why this control is disabled.",
		},
	}
}

func isInteractiveControl(node *ir.Node) bool {
	tagLower := strings.ToLower(node.Tag)
	switch tagLower {
	case "button", "select", "textarea":
		return true
	case "input":
		if t, ok := getAttrCaseInsensitive(node, "type"); ok {
			typ := cleanAttrValue(t)
			if typ == "hidden" {
				return false
			}
		}
		return true
	}

	if role, ok := getAttrCaseInsensitive(node, "role"); ok {
		r := cleanAttrValue(role)
		if r == "button" || r == "link" || r == "menuitem" || r == "tab" {
			return true
		}
	}

	return false
}

func isDisabledControl(node *ir.Node) bool {
	if val, ok := getAttrCaseInsensitive(node, "disabled"); ok {
		cleaned := cleanAttrValue(val)
		if cleaned == "false" || val == "{false}" {
			return false
		}
		return true
	}

	if val, ok := getAttrCaseInsensitive(node, "aria-disabled"); ok {
		cleaned := cleanAttrValue(val)
		if cleaned == "true" || val == "{true}" {
			return true
		}
	}

	return false
}
