package form

import (
	"strconv"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// InputFacts merangkum fakta semantik dan interaksi dari sebuah elemen form input.
type InputFacts struct {
	IsInputTag      bool
	IsNumberType    bool
	IsDisabled      bool
	IsReadOnly      bool
	HasWheelHandler bool
	HasDeclaredMin  bool

	Identifier IdentifierEvidence
}

// getAttrCI mengambil nilai atribut dari map node.Attributes secara case-insensitive.
func getAttrCI(attrs map[string]string, keys ...string) (string, string, bool) {
	if attrs == nil {
		return "", "", false
	}
	for _, target := range keys {
		for k, v := range attrs {
			if strings.EqualFold(k, target) {
				return k, v, true
			}
		}
	}
	return "", "", false
}

// cleanAttrVal membersihkan tanda kutip (" ' `), kurung kurawal ({ }), dan spasi.
func cleanAttrVal(val string) string {
	return strings.Trim(strings.TrimSpace(val), "\"'`{}")
}

// isNumberTypeAttr memeriksa apakah atribut type bernilai static "number".
// Mengabaikan ekspresi dinamis (misal: type={inputType} atau type={cond ? "number" : "text"})
// untuk menjaga kepastian statis (static certainty).
func isNumberTypeAttr(rawVal string) bool {
	raw := strings.TrimSpace(rawVal)
	if raw == "" {
		return false
	}

	// Jika JSX brace expression: type={"number"} vs type={dynamicVar}
	if strings.HasPrefix(raw, "{") && strings.HasSuffix(raw, "}") {
		inner := strings.TrimSpace(raw[1 : len(raw)-1])
		// Periksa apakah inner adalah string literal yang diapit tanda kutip
		if (strings.HasPrefix(inner, "\"") && strings.HasSuffix(inner, "\"")) ||
			(strings.HasPrefix(inner, "'") && strings.HasSuffix(inner, "'")) ||
			(strings.HasPrefix(inner, "`") && strings.HasSuffix(inner, "`")) {
			clean := cleanAttrVal(inner)
			return strings.EqualFold(clean, "number")
		}
		// Ekspresi dinamis diabaikan
		return false
	}

	clean := cleanAttrVal(raw)
	return strings.EqualFold(clean, "number")
}

// hasDeclaredMinAttr memeriksa apakah developer telah mendeklarasikan batas bawah domain yang valid.
// min="0", min="1", min="-50", min={0}, min={domainMin} -> true
// min="", min={undefined}, min={null}, atau tidak ada min -> false
func hasDeclaredMinAttr(rawVal string, present bool) bool {
	if !present {
		return false
	}
	raw := strings.TrimSpace(rawVal)
	if raw == "" {
		return false
	}

	// Jika JSX brace expression: min={0}, min={domainMin} vs min={undefined}
	if strings.HasPrefix(raw, "{") && strings.HasSuffix(raw, "}") {
		inner := strings.TrimSpace(raw[1 : len(raw)-1])
		if inner == "" || inner == "undefined" || inner == "null" {
			return false
		}
		return true
	}

	clean := cleanAttrVal(raw)
	if clean == "" {
		return false
	}

	// Validasi parsing numerik jika literal
	if _, err := strconv.ParseFloat(clean, 64); err == nil {
		return true
	}

	// String non-empty non-numeric tetap dianggap deklarasi eksplisit (misal template string atau konstanta)
	return true
}

// extractInteractionFacts memeriksa status disabled dan readOnly dari atribut elemen.
func extractInteractionFacts(attrs map[string]string) (disabled, readOnly bool) {
	if _, disVal, ok := getAttrCI(attrs, "disabled", "aria-disabled"); ok {
		clean := cleanAttrVal(disVal)
		if clean != "false" {
			disabled = true
		}
	}
	if _, roVal, ok := getAttrCI(attrs, "readonly", "readOnly"); ok {
		clean := cleanAttrVal(roVal)
		if clean != "false" {
			readOnly = true
		}
	}
	return disabled, readOnly
}

type evidenceSourceDef struct {
	keys   []string
	source IdentifierSource
}

var evidenceSources = []evidenceSourceDef{
	{[]string{"name"}, IdentifierSourceName},
	{[]string{"id"}, IdentifierSourceID},
	{[]string{"autocomplete", "autoComplete"}, IdentifierSourceAutocomplete},
	{[]string{"aria-label", "aria-labelledby"}, IdentifierSourceAriaLabel},
	{[]string{"placeholder"}, IdentifierSourcePlaceholder},
	{[]string{"data-testid", "data-test-id", "data-test"}, IdentifierSourceTestID},
}

// findEvidenceByClass mencari atribut pertama yang cocok dengan targetClass.
func findEvidenceByClass(attrs map[string]string, targetClass IdentifierClass) (IdentifierEvidence, bool) {
	for _, es := range evidenceSources {
		_, val, ok := getAttrCI(attrs, es.keys...)
		if !ok {
			continue
		}
		cleanVal := cleanAttrVal(val)
		if cleanVal == "" {
			continue
		}

		if targetClass == IdentifierIdentity && es.source == IdentifierSourceAutocomplete {
			if matched, isAuto := authoritativeAutocompleteIdentity[cleanVal]; isAuto {
				return IdentifierEvidence{
					Value:   cleanVal,
					Source:  es.source,
					Class:   IdentifierIdentity,
					Matched: matched,
				}, true
			}
		}

		tokens := TokenizeIdentifier(cleanVal)
		class, matched := ClassifyIdentifier(tokens)
		if class == targetClass {
			return IdentifierEvidence{
				Value:   cleanVal,
				Source:  es.source,
				Class:   targetClass,
				Matched: matched,
			}, true
		}
	}
	return IdentifierEvidence{}, false
}

// extractIdentifierEvidence mengevaluasi atribut berdasarkan hierarki otoritas (ranked evidence).
func extractIdentifierEvidence(attrs map[string]string) IdentifierEvidence {
	// Pass 1: cari bukti identitas (IdentifierIdentity)
	if ev, ok := findEvidenceByClass(attrs, IdentifierIdentity); ok {
		return ev
	}

	// Pass 2: cari bukti kuantitas (IdentifierQuantity)
	if ev, ok := findEvidenceByClass(attrs, IdentifierQuantity); ok {
		return ev
	}

	// Pass 3: fallback identifier primer (IdentifierUnknown)
	for _, es := range evidenceSources {
		_, val, ok := getAttrCI(attrs, es.keys...)
		if !ok {
			continue
		}
		cleanVal := cleanAttrVal(val)
		if cleanVal != "" {
			return IdentifierEvidence{
				Value:  cleanVal,
				Source: es.source,
				Class:  IdentifierUnknown,
			}
		}
	}

	return IdentifierEvidence{}
}

// ExtractInputFacts mengekstraksi metadata semantik dan interaksi dari sebuah node IR secara murni.
func ExtractInputFacts(node *ir.Node) InputFacts {
	var facts InputFacts
	if node == nil || node.Type != ir.NodeElement {
		return facts
	}

	// 1. Validasi Tag Input (native <input> atau konvensi komponen <Input>)
	if node.Tag != "input" && node.Tag != "Input" {
		return facts
	}
	facts.IsInputTag = true

	attrs := node.Attributes

	// 2. Evaluasi Tipe Input ("number")
	if _, typeVal, ok := getAttrCI(attrs, "type"); ok {
		facts.IsNumberType = isNumberTypeAttr(typeVal)
	}

	// 3. Evaluasi Pengecualian Interaksi (disabled & readOnly)
	facts.IsDisabled, facts.IsReadOnly = extractInteractionFacts(attrs)

	// 4. Evaluasi Wheel Protection Handler
	if _, _, ok := getAttrCI(attrs, "onwheel", "onWheel", "onwheelcapture", "onWheelCapture"); ok {
		facts.HasWheelHandler = true
	}

	// 5. Evaluasi Deklarasi Min (Lower Bound)
	_, minVal, minPresent := getAttrCI(attrs, "min")
	facts.HasDeclaredMin = hasDeclaredMinAttr(minVal, minPresent)

	// 6. Evaluasi Identifier Berdasarkan Bobot Otoritas (Ranked Evidence)
	facts.Identifier = extractIdentifierEvidence(attrs)

	return facts
}
