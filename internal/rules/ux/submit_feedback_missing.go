package ux

import (
	"fmt"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// SubmitFeedbackMissingRule memastikan setiap pemicu mutasi asinkron memenuhi dua kontrak:
// R1 (Reentry Guard: penguncian interaksi via disabled) dan R2 (Perceivable Feedback: umpan balik visual via aria-busy/spinner),
// guna mematuhi Doherty Threshold (< 400ms) dan mencegah aksi ganda (duplicate submission).
type SubmitFeedbackMissingRule struct{}

// NewSubmitFeedbackMissingRule membuat instance baru dari SubmitFeedbackMissingRule.
func NewSubmitFeedbackMissingRule() *SubmitFeedbackMissingRule {
	return &SubmitFeedbackMissingRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *SubmitFeedbackMissingRule) ID() string {
	return "ux.submit-feedback-missing"
}

// Description mengembalikan ringkasan aturan.
func (r *SubmitFeedbackMissingRule) Description() string {
	return "Enforces reentry guard (disabled) and perceivable feedback (aria-busy/spinner) on async mutation triggers"
}

// Category mengembalikan nama kategori rule.
func (r *SubmitFeedbackMissingRule) Category() string {
	return "ux"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warning).
func (r *SubmitFeedbackMissingRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *SubmitFeedbackMissingRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Doherty Threshold (< 400ms Interaction Feedback)",
			"Nielsen Heuristic #1: Visibility of System Status",
			"ISO 9241-110 Ergonomics of Human-System Interaction (Self-Descriptiveness)",
		},
		CoreInvariant: "Async mutation triggers must provide both Reentry Guard (R1: 'disabled={isPending}') and Perceivable Feedback (R2: 'aria-busy' or spinner).",
		Grounding: "When users trigger a mutation (such as submitting an order or charging a payment) without immediate feedback, " +
			"they perceive the system as unresponsive within the 400ms Doherty Threshold.\n\n" +
			"In the absence of a reentry lock, users repeatedly click the submit button, triggering duplicate requests, " +
			"double financial charges, and server-side race conditions. Enforcing both R1 (reentry lockout) and R2 (visual feedback) " +
			"guarantees transactional safety and cognitive assurance.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Duplicate Submissions & Double Charging",
				Severity: "HIGH",
				Impact:   "Rapid repeated clicks by uncertain users cause duplicate database inserts and multiple payment gateway captures.",
			},
			{
				Vector:   "Cognitive Disorientation & Interaction Rage",
				Severity: "MEDIUM",
				Impact:   "Users doubt whether their input registered, leading to rage clicks and workflow cancellation.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Async submit button without disabling control or showing pending feedback",
				Code: `<button
  onClick={async () => {
    await api.post("/orders", orderData);
  }}
  className="bg-primary text-white px-4 py-2"
>
  Bayar Sekarang
</button>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Compliant button satisfying R1 (disabled={isPending}) and R2 (aria-busy & dynamic text)",
				Code: `<button
  onClick={handlePayment}
  disabled={isPending}
  aria-busy={isPending}
  className="bg-primary text-white px-4 py-2 disabled:opacity-50"
>
  {isPending ? "Memproses Pembayaran..." : "Bayar Sekarang"}
</button>`,
			},
		},
	}
}

// Evaluate memeriksa apakah pemicu mutasi asinkron memenuhi kontrak R1 dan R2.
func (r *SubmitFeedbackMissingRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	if !isMutationTriggerElement(node) {
		return nil
	}

	if !hasAsyncMutationHandlerInNode(node) {
		return nil
	}

	r1 := hasReentryGuard(node)
	r2 := hasPerceivableFeedback(node)

	if !r1 && !r2 {
		return []ir.Diagnostic{
			{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: ir.SeverityError,
				Message: fmt.Sprintf(
					"Async mutation trigger <%s> lacks both Reentry Guard (R1: disabled) and Perceivable Feedback (R2: aria-busy/spinner), risking duplicate submissions.",
					node.Tag,
				),
				Hint: "Add 'disabled={isPending}' to prevent duplicate submissions and 'aria-busy={isPending}' or a <Spinner /> for immediate visual feedback.",
			},
		}
	}

	if !r1 {
		return []ir.Diagnostic{
			{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: ir.SeverityWarn,
				Message: fmt.Sprintf(
					"Async mutation trigger <%s> lacks Reentry Guard (R1: disabled={isPending}) to prevent duplicate submissions during in-flight network requests.",
					node.Tag,
				),
				Hint: "Lock the control with 'disabled={isPending}' while the async mutation is executing.",
			},
		}
	}

	if !r2 {
		return []ir.Diagnostic{
			{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: ir.SeverityWarn,
				Message: fmt.Sprintf(
					"Async mutation trigger <%s> lacks Perceivable Feedback (R2: aria-busy/spinner) to acknowledge user action within Doherty Threshold (< 400ms).",
					node.Tag,
				),
				Hint: "Provide immediate visual feedback using 'aria-busy={isPending}' or an animated spinner.",
			},
		}
	}

	return nil
}

func isMutationTriggerElement(node *ir.Node) bool {
	tagLower := strings.ToLower(node.Tag)
	if tagLower == "button" || tagLower == "form" {
		return true
	}
	if role, ok := getAttrCaseInsensitive(node, "role"); ok {
		r := cleanAttrValue(role)
		if r == "button" || r == "form" {
			return true
		}
	}
	return false
}

func hasAsyncMutationHandlerInNode(node *ir.Node) bool {
	for attrName, attrVal := range node.Attributes {
		if !isEventHandlerOrActionAttr(attrName) {
			continue
		}
		if detectAsyncMutationHandler(attrVal) {
			return true
		}
	}
	return false
}
