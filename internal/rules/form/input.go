package form

import (
	"strconv"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// InputTypeClass merepresentasikan klasifikasi jenis atribut type pada elemen input.
type InputTypeClass int

const (
	// InputTypeUnknown menandakan tipe input tidak dapat ditentukan secara statis (misal: binding dinamis type={foo}).
	InputTypeUnknown InputTypeClass = iota
	// InputTypeText menandakan generic single-line text input (type="text", type="", atau type tidak dideklarasikan).
	InputTypeText
	// InputTypeOther menandakan kontrol input dengan tipe khusus (misal: number, password, email, tel, checkbox, radio, dll.).
	InputTypeOther
)

func (t InputTypeClass) String() string {
	switch t {
	case InputTypeText:
		return "text"
	case InputTypeOther:
		return "other"
	default:
		return "unknown"
	}
}

// InputFacts merangkum fakta semantik dan interaksi dari sebuah elemen form input.
type InputFacts struct {
	IsInputTag      bool
	IsNumberType    bool // Dipertahankan untuk kompatibilitas penuh dengan rule Issue #5
	TypeClass       InputTypeClass
	IsDisabled      bool
	IsReadOnly      bool
	HasWheelHandler bool
	HasDeclaredMin  bool

	Identifier IdentifierEvidence
	Content    ContentEvidence
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

// classifyInputType mengklasifikasikan atribut type ke dalam InputTypeClass dan mengembalikan flag isNumberType.
func classifyInputType(rawVal string, present bool) (InputTypeClass, bool) {
	if !present {
		return InputTypeText, false
	}

	raw := strings.TrimSpace(rawVal)
	if raw == "" {
		return InputTypeText, false
	}

	// Jika JSX brace expression: type={"text"} vs type={dynamicVar}
	if strings.HasPrefix(raw, "{") && strings.HasSuffix(raw, "}") {
		inner := strings.TrimSpace(raw[1 : len(raw)-1])
		if (strings.HasPrefix(inner, "\"") && strings.HasSuffix(inner, "\"")) ||
			(strings.HasPrefix(inner, "'") && strings.HasSuffix(inner, "'")) ||
			(strings.HasPrefix(inner, "`") && strings.HasSuffix(inner, "`")) {
			clean := strings.ToLower(cleanAttrVal(inner))
			if clean == "" || clean == "text" {
				return InputTypeText, false
			}
			if clean == "number" {
				return InputTypeOther, true
			}
			return InputTypeOther, false
		}
		// Ekspresi dinamis non-literal (type={dynamicVar})
		return InputTypeUnknown, false
	}

	clean := strings.ToLower(cleanAttrVal(raw))
	if clean == "" || clean == "text" {
		return InputTypeText, false
	}
	if clean == "number" {
		return InputTypeOther, true
	}
	return InputTypeOther, false
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

type contentChannelDef struct {
	source IdentifierSource
	key    string
	isText bool
}

var contentChannels = []contentChannelDef{
	{source: IdentifierSourceName, key: "name", isText: false},
	{source: IdentifierSourceID, key: "id", isText: false},
	{source: IdentifierSourcePlaceholder, key: "placeholder", isText: true},
	{source: IdentifierSourceAriaLabel, key: "aria-label", isText: true},
}

type channelFinding struct {
	source  IdentifierSource
	value   string
	tokens  []string
	intent  ContentIntent
	matched string
}

// extractContentEvidence mengevaluasi seluruh channel atribut secara independen untuk menentukan intensi bentuk konten.
// Precedence intent:
// Strong Single-Line Qualifier > Strong Multiline Token > Contextual Single-Line Qualifier > Unknown.
func extractContentEvidence(attrs map[string]string) ContentEvidence {
	if attrs == nil {
		return ContentEvidence{}
	}

	findings := make([]channelFinding, 0, len(contentChannels))
	for _, ch := range contentChannels {
		_, rawVal, ok := getAttrCI(attrs, ch.key)
		if !ok {
			continue
		}
		cleanVal := cleanAttrVal(rawVal)
		if cleanVal == "" {
			continue
		}

		var tokens []string
		if ch.isText {
			tokens = TokenizeTextWords(cleanVal)
		} else {
			tokens = TokenizeIdentifier(cleanVal)
		}

		intent, matched := ClassifyContentIntent(tokens)
		findings = append(findings, channelFinding{
			source:  ch.source,
			value:   cleanVal,
			tokens:  tokens,
			intent:  intent,
			matched: matched,
		})
	}

	if len(findings) == 0 {
		return ContentEvidence{}
	}

	// 1. Evaluasi apakah ada channel yang menghasilkan SingleLine (qualifier menang atas multiline)
	for _, f := range findings {
		if f.intent == ContentIntentSingleLine {
			return ContentEvidence{
				Value:   f.value,
				Source:  f.source,
				Intent:  ContentIntentSingleLine,
				Matched: f.matched,
			}
		}
	}

	// 2. Evaluasi apakah ada channel yang menghasilkan Multiline
	for _, f := range findings {
		if f.intent == ContentIntentMultiline {
			return ContentEvidence{
				Value:   f.value,
				Source:  f.source,
				Intent:  ContentIntentMultiline,
				Matched: f.matched,
			}
		}
	}

	// 3. Fallback ke temuan pertama dengan intent Unknown
	return ContentEvidence{
		Value:  findings[0].value,
		Source: findings[0].source,
		Intent: ContentIntentUnknown,
	}
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

	// 2. Evaluasi Tipe Input
	_, typeVal, typePresent := getAttrCI(attrs, "type")
	facts.TypeClass, facts.IsNumberType = classifyInputType(typeVal, typePresent)

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

	// 7. Evaluasi Intensi Bentuk Konten (Multiline vs SingleLine)
	facts.Content = extractContentEvidence(attrs)

	return facts
}
