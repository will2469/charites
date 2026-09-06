package ux

import (
	"fmt"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// OrphanedErrorStateRule memastikan setter state error yang dipanggil dalam event handler
// memiliki elemen presentasi error di dalam komponen sehingga error dapat dilihat oleh pengguna.
type OrphanedErrorStateRule struct{}

// NewOrphanedErrorStateRule membuat instance baru dari OrphanedErrorStateRule.
func NewOrphanedErrorStateRule() *OrphanedErrorStateRule {
	return &OrphanedErrorStateRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *OrphanedErrorStateRule) ID() string {
	return "ux.orphaned-error-state"
}

// Description mengembalikan ringkasan aturan.
func (r *OrphanedErrorStateRule) Description() string {
	return "Flags error state updates in event handlers that lack corresponding UI error presentation elements"
}

// Category mengembalikan nama kategori rule.
func (r *OrphanedErrorStateRule) Category() string {
	return "ux"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warning).
func (r *OrphanedErrorStateRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *OrphanedErrorStateRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Nielsen Heuristic #9: Help Users Recognize, Diagnose, and Recover from Errors",
			"ISO 9241-110 Ergonomics of Human-System Interaction (Error Tolerance)",
			"WCAG 2.2 Success Criterion 3.3.1 (Error Identification)",
		},
		CoreInvariant: "Validation error setters invoked in component handlers must have a corresponding error presentation element in the UI.",
		Grounding: "When client-side validation logic flags invalid input and updates internal component state (e.g. 'setEmailError(\"Format salah\")'), " +
			"that state must be surfaced to the user.\n\n" +
			"If the error state is never rendered in JSX or communicated via accessible error indicators, " +
			"the form silently blocks submission while the user remains completely unaware of what went wrong.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Silent Submission Failure & Ghost Validation",
				Severity: "HIGH",
				Impact:   "Users submit forms that silently fail validation without displaying any error messages, causing confusion and frustration.",
			},
			{
				Vector:   "Inaccessible Error Notification",
				Severity: "MEDIUM",
				Impact:   "Screen reader users and keyboard navigators receive no feedback when form constraints are violated.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Event handler updates error state but no error display exists in the JSX",
				Code: `export function LoginForm() {
  const [email, setEmail] = useState("");
  const [emailError, setEmailError] = useState("");

  const handleSubmit = (e) => {
    e.preventDefault();
    if (!email.includes("@")) {
      setEmailError("Format email tidak valid");
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <input value={email} onChange={e => setEmail(e.target.value)} />
      <button type="submit">Masuk</button>
    </form>
  );
}`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Error state rendered visibly and accessibly with role='alert' and destructive color",
				Code: `export function LoginForm() {
  const [email, setEmail] = useState("");
  const [emailError, setEmailError] = useState("");

  const handleSubmit = (e) => {
    e.preventDefault();
    if (!email.includes("@")) {
      setEmailError("Format email tidak valid");
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <input value={email} onChange={e => setEmail(e.target.value)} />
      {emailError && (
        <p role="alert" className="text-sm text-destructive font-medium">
          {emailError}
        </p>
      )}
      <button type="submit">Masuk</button>
    </form>
  );
}`,
			},
		},
	}
}

// Evaluate memeriksa apakah pemanggilan setter error di handler memiliki presentasi UI.
func (r *OrphanedErrorStateRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	setterName, found := detectErrorSetterInNode(node)
	if !found {
		return nil
	}

	root := getComponentRoot(node)
	if hasErrorPresentationInSubtree(root) {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message: fmt.Sprintf(
				"Event handler in <%s> updates error state (%s) but no corresponding error presentation element exists in the component.",
				node.Tag,
				setterName,
			),
			Hint: "Render the error state in your JSX (e.g., '{error && <p role=\"alert\" className=\"text-destructive\">{error}</p>}') so users can recognize and recover from errors.",
		},
	}
}

func detectErrorSetterInNode(node *ir.Node) (string, bool) {
	for attrName, attrVal := range node.Attributes {
		if !isEventHandlerOrActionAttr(attrName) {
			continue
		}
		if setter, ok := extractErrorSetterFromText(attrVal); ok {
			return setter, true
		}
	}
	return "", false
}

func extractErrorSetterFromText(val string) (string, bool) {
	lower := strings.ToLower(val)
	errorKeywords := [...]string{
		"seterror(", "setemailerror(", "setformerror(", "setvalidationerror(",
		"seterr(", "setfielderror(", "setpassworderror(", "setmessageerror(",
	}
	for _, kw := range errorKeywords {
		if strings.Contains(lower, kw) {
			return strings.TrimSuffix(kw, "("), true
		}
	}
	return "", false
}

func getComponentRoot(node *ir.Node) *ir.Node {
	curr := node
	for curr.Parent != nil {
		curr = curr.Parent
	}
	return curr
}
