package ux

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// SilentCatchSwallowRule mendeteksi blok catch pada handler interaksi pengguna yang menelan error secara senyap
// (hanya console.log atau kosong) tanpa memberi umpan balik visual (toast, alert, banner error) atau me-rethrow error.
type SilentCatchSwallowRule struct{}

// NewSilentCatchSwallowRule membuat instance baru dari SilentCatchSwallowRule.
func NewSilentCatchSwallowRule() *SilentCatchSwallowRule {
	return &SilentCatchSwallowRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *SilentCatchSwallowRule) ID() string {
	return "ux.silent-catch-swallow"
}

// Description mengembalikan ringkasan aturan.
func (r *SilentCatchSwallowRule) Description() string {
	return "Detects swallowed catch blocks in event handlers that lack user feedback (toast/alert) or re-throw"
}

// Category mengembalikan nama kategori rule.
func (r *SilentCatchSwallowRule) Category() string {
	return "ux"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (error).
func (r *SilentCatchSwallowRule) DefaultSeverity() ir.Severity {
	return ir.SeverityError
}

// Doc mengembalikan dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *SilentCatchSwallowRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Nielsen Heuristic #9: Help Users Recognize, Diagnose, and Recover from Errors",
			"ISO 9241-110 Ergonomics of Human-System Interaction (Error Management)",
			"Zero-Trust Error Transparency Guidelines",
		},
		CoreInvariant: "Catch blocks in user interaction handlers must provide visible UI feedback ('toast', error state, alert) or re-throw the error.",
		Grounding: "When user interaction handlers catch exceptions and only log them to 'console.log' or discard them entirely, " +
			"the failure is silently swallowed.\n\n" +
			"The user receives no feedback, falsely assumes their changes were saved, and navigates away, only to discover later " +
			"that critical data was lost. Surfacing visible feedback (e.g. 'toast.error(...)', 'setError(...)', or banner notifications) " +
			"guarantees that errors are transparent, allowing users to understand the problem and re-attempt the action.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Silent Data Loss & False Sense of Completion",
				Severity: "HIGH",
				Impact:   "Users believe changes succeeded when they actually failed on the network, leading to unrecoverable data loss.",
			},
			{
				Vector:   "Lack of Failure Diagnostics",
				Severity: "MEDIUM",
				Impact:   "Users cannot self-diagnose network errors or invalid parameters, resulting in confusion and support tickets.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Catch block silently logging to console without notifying the user",
				Code: `<button
  onClick={async () => {
    try {
      await api.updateProfile(formData);
    } catch (e) {
      console.error(e); // Pengguna tidak tahu aksinya gagal!
    }
  }}
>
  Simpan Profil
</button>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Catch block notifying user with a toast error notification",
				Code: `<button
  onClick={async () => {
    try {
      await api.updateProfile(formData);
    } catch (e) {
      toast.error("Gagal memperbarui profil. Silakan coba lagi.");
    }
  }}
>
  Simpan Profil
</button>`,
			},
		},
	}
}

// Evaluate memeriksa apakah ada blok catch di event handler yang menelan error tanpa umpan balik antarmuka.
func (r *SilentCatchSwallowRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	for attrName, attrVal := range node.Attributes {
		if !isEventHandlerOrActionAttr(attrName) {
			continue
		}

		if detectSilentCatchSwallow(attrVal) {
			return []ir.Diagnostic{
				{
					Line:     node.Span.Line,
					Column:   node.Span.Column,
					Rule:     r.ID(),
					Severity: r.DefaultSeverity(),
					Message: fmt.Sprintf(
						"Event handler %q contains a catch block that silently swallows errors without surfacing UI feedback (toast, alert, or error state) or re-throwing.",
						attrName,
					),
					Hint: "Provide visible feedback to the user upon failure (e.g. 'toast.error(\"Gagal menyimpan\");' or 'setError(err.message);') so they are aware of the error.",
				},
			}
		}
	}

	return nil
}
