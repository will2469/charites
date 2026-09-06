package ux

import (
	"fmt"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// DestructiveActionUnconfirmedRule memastikan aksi destruktif (seperti delete, remove, destroy, purge, revoke)
// dilindungi mekanisme konfirmasi (dialog modal, AlertDialog, atau confirm()) untuk mencegah kesalahan slip tidak disengaja.
type DestructiveActionUnconfirmedRule struct{}

// NewDestructiveActionUnconfirmedRule membuat instance baru dari DestructiveActionUnconfirmedRule.
func NewDestructiveActionUnconfirmedRule() *DestructiveActionUnconfirmedRule {
	return &DestructiveActionUnconfirmedRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *DestructiveActionUnconfirmedRule) ID() string {
	return "ux.destructive-action-unconfirmed"
}

// Description mengembalikan ringkasan aturan.
func (r *DestructiveActionUnconfirmedRule) Description() string {
	return "Enforces confirmation gating for destructive actions to prevent accidental data loss from slips"
}

// Category mengembalikan nama kategori rule.
func (r *DestructiveActionUnconfirmedRule) Category() string {
	return "ux"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (error).
func (r *DestructiveActionUnconfirmedRule) DefaultSeverity() ir.Severity {
	return ir.SeverityError
}

// Doc mengembalikan dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *DestructiveActionUnconfirmedRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Nielsen Heuristic #5: Error Prevention (Slips and Lapses)",
			"ISO 9241-110 Ergonomics of Human-System Interaction (Suitability for the Task)",
			"Material Design & WCAG Defensive Action Guidelines",
		},
		CoreInvariant: "Destructive user operations ('delete', 'remove', 'destroy', 'purge', 'revoke') must be gated by a confirmation dialog or 2-step verification.",
		Grounding: "Destructive actions such as deleting user accounts, clearing billing databases, or revoking credentials cause permanent, " +
			"often irreversible data loss.\n\n" +
			"Executing these operations on a single click without confirmation exposes users to motor slips, touchscreen taps during scrolling, " +
			"and mistaken identity clicks. Gating destructive actions behind an explicit confirmation dialog (e.g. '<AlertDialog>' or 'window.confirm') " +
			"provides a cognitive pause and protects against catastrophic slips.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Irreversible Data Destruction from Motor Slips",
				Severity: "CRITICAL",
				Impact:   "Accidental single-click actions immediately wipe critical business data or terminate accounts without user consent.",
			},
			{
				Vector:   "User Anxiety & Hesitation",
				Severity: "MEDIUM",
				Impact:   "Users fear interacting with danger-styled buttons when no confirmation boundary protects them from permanent loss.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Destructive button triggering account deletion directly on single click",
				Code: `<button
  onClick={() => deleteAccount(user.id)}
  className="bg-destructive text-destructive-foreground px-4 py-2 rounded"
>
  Hapus Akun Permanen
</button>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Destructive button safely wrapped in AlertDialogTrigger confirmation modal",
				Code: `<AlertDialogTrigger asChild>
  <button className="bg-destructive text-destructive-foreground px-4 py-2 rounded">
    Hapus Akun Permanen
  </button>
</AlertDialogTrigger>`,
			},
		},
	}
}

// Evaluate memeriksa apakah aksi destruktif dilindungi dialog konfirmasi.
func (r *DestructiveActionUnconfirmedRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	if !isInteractiveElement(node) {
		return nil
	}

	act, isDestructive := detectDestructiveAction(node)
	if !isDestructive {
		return nil
	}

	for _, attrVal := range node.Attributes {
		if hasConfirmationGating(node, attrVal) {
			return nil
		}
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message: fmt.Sprintf(
				"Destructive action %q on <%s> lacks confirmation gating (missing AlertDialog, modal, or confirm dialog).",
				act,
				node.Tag,
			),
			Hint: "Gate destructive actions behind a confirmation dialog (e.g. <AlertDialog>, confirm(), or a two-step confirmation state) to prevent accidental data loss.",
		},
	}
}

func isInteractiveElement(node *ir.Node) bool {
	tagLower := strings.ToLower(node.Tag)
	if tagLower == "button" || tagLower == "a" {
		return true
	}
	if role, ok := getAttrCaseInsensitive(node, "role"); ok {
		r := cleanAttrValue(role)
		if r == "button" || r == "link" || r == "menuitem" {
			return true
		}
	}
	return false
}
