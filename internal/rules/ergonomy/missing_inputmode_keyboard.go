package ergonomy

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// MissingInputmodeKeyboardRule menegakkan penggunaan atribut inputmode atau type kontekstual pada kontrol form mobile.
type MissingInputmodeKeyboardRule struct{}

// NewMissingInputmodeKeyboardRule membuat instance baru dari MissingInputmodeKeyboardRule.
func NewMissingInputmodeKeyboardRule() *MissingInputmodeKeyboardRule {
	return &MissingInputmodeKeyboardRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *MissingInputmodeKeyboardRule) ID() string {
	return "ergonomy.missing-inputmode-keyboard"
}

// Description mengembalikan ringkasan aturan.
func (r *MissingInputmodeKeyboardRule) Description() string {
	return "Enforces contextual virtual keyboard inputmode and type attributes on mobile form inputs (Tesler's Law)"
}

// Category mengembalikan nama kategori rule.
func (r *MissingInputmodeKeyboardRule) Category() string {
	return "ergonomy"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (info).
func (r *MissingInputmodeKeyboardRule) DefaultSeverity() ir.Severity {
	return ir.SeverityInfo
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *MissingInputmodeKeyboardRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"HTML Living Standard Section 4.10.5.3 (The inputmode attribute)",
			"Tesler's Law (Conservation of Complexity in Virtual Keyboards)",
			"Apple iOS & Android Mobile Virtual Keyboard Guidelines",
		},
		CoreInvariant: "Form text inputs collecting numeric, phone, or email values must declare contextual 'inputmode' or specialized 'type' attributes to directly open the optimized mobile virtual keypad.",
		Grounding: "On mobile devices, focusing an input without specialized type or inputmode opens the full standard QWERTY keyboard.\n\n" +
			"For numeric, telephone, or OTP fields, this forces the user to repeatedly toggle keyboard layers to find digits. " +
			"According to Tesler's Law, complexity must be absorbed by software rather than offloaded to user manual effort. " +
			"Declaring 'inputmode=\"numeric\"' or 'type=\"tel\"' instantly summons large, thumb-friendly numeric keypads.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Mobile Cognitive Friction",
				Severity: "LOW",
				Impact:   "Users are forced to manually switch keyboard layers on small touchscreens to enter digits or phone numbers.",
			},
			{
				Vector:   "High Form Abandonment & Typing Errors",
				Severity: "LOW",
				Impact:   "Entering OTP or financial amounts on dense QWERTY keys leads to frequent miss-taps and delayed checkout flows.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Phone input missing type or inputMode",
				Code: `<input
  name="nomor_hp"
  placeholder="08123456789"
  className="h-11 px-3.5 py-2.5 border rounded-xl"
/>`,
			},
			{
				Language: "astro",
				Comment:  "OTP field defaulting to QWERTY keyboard",
				Code: `<input
  id="otp_code"
  placeholder="123456"
  class="h-11 px-3.5 py-2.5 border rounded-xl"
/>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Explicit telephone keypad and autocomplete",
				Code: `<input
  name="nomor_hp"
  type="tel"
  inputMode="tel"
  autoComplete="tel"
  placeholder="08123456789"
  className="h-11 px-3.5 py-2.5 border rounded-xl"
/>`,
			},
			{
				Language: "astro",
				Comment:  "Numeric keypad for OTP verification",
				Code: `<input
  id="otp_code"
  type="text"
  inputmode="numeric"
  pattern="[0-9]*"
  placeholder="123456"
  class="h-11 px-3.5 py-2.5 border rounded-xl"
/>`,
			},
		},
	}
}

// Evaluate memeriksa apakah kontrol input form memiliki inputmode kontekstual.
func (r *MissingInputmodeKeyboardRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || strings.ToLower(node.Tag) != "input" || node.Attributes == nil {
		return nil
	}

	typeAttr := cleanAttrValue(node.Attributes["type"])
	if isNonTextualInputType(typeAttr) {
		return nil
	}

	// Jika sudah memiliki inputmode eksplisit, dianggap Compliant
	if hasInputMode(node.Attributes) || typeAttr == "tel" || typeAttr == "email" || typeAttr == "number" {
		return nil
	}

	identifier := getSemanticIdentifier(node.Attributes)
	if identifier == "" {
		return nil
	}

	expectedType, suggestedMode := detectExpectedKeyboard(identifier)
	if expectedType == "" {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  "Input field '" + identifier + "' expects " + expectedType + " data but lacks contextual virtual keyboard hints. Mobile users will be shown the default QWERTY keyboard.",
			Hint:     "Add 'inputmode=\"" + suggestedMode + "\"' or 'type=\"" + suggestedMode + "\"' to directly display the dedicated mobile keypad (Tesler's Law).",
		},
	}
}

func cleanAttrValue(v string) string {
	return strings.Trim(strings.TrimSpace(strings.ToLower(v)), "\"'`{}")
}

func isNonTextualInputType(t string) bool {
	switch t {
	case "hidden", "submit", "button", "reset", "checkbox", "radio", "file", "image", "color", "range", "date", "datetime-local", "time", "month", "week":
		return true
	default:
		return false
	}
}

func hasInputMode(attrs map[string]string) bool {
	for k, v := range attrs {
		if strings.EqualFold(k, "inputmode") && cleanAttrValue(v) != "" {
			return true
		}
	}
	return false
}

func getSemanticIdentifier(attrs map[string]string) string {
	for _, key := range [...]string{"name", "id", "placeholder", "aria-label", "autocomplete"} {
		for k, val := range attrs {
			if strings.EqualFold(k, key) {
				cleanVal := cleanAttrValue(val)
				if cleanVal != "" {
					return cleanVal
				}
			}
		}
	}
	return ""
}

func detectExpectedKeyboard(id string) (expectedType, suggestedMode string) {
	// 1. Phone triggers
	for _, keyword := range [...]string{"phone", "telp", "telepon", "hp", "wa", "whatsapp", "mobile_number", "nohp", "no_hp", "handphone"} {
		if strings.Contains(id, keyword) {
			return "telephone", "tel"
		}
	}

	// 2. Numeric / OTP / Financial / Zip triggers
	for _, keyword := range [...]string{"otp", "pin", "kode_otp", "verification", "cvv", "cvc", "nominal", "harga", "price", "amount", "kodepos", "postal", "zip", "rekening", "nik", "ktp", "npwp"} {
		if strings.Contains(id, keyword) {
			return "numeric", "numeric"
		}
	}

	// 3. Email triggers
	for _, keyword := range [...]string{"email", "surel"} {
		if strings.Contains(id, keyword) {
			return "email", "email"
		}
	}

	return "", ""
}
