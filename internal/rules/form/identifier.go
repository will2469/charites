package form

import (
	"strings"
	"unicode"
)

// IdentifierClass merepresentasikan klasifikasi semantik dari identitas sebuah field input.
type IdentifierClass int

const (
	// IdentifierUnknown menandakan field tidak memiliki sinyal identitas maupun kuantitas yang jelas.
	IdentifierUnknown IdentifierClass = iota
	// IdentifierQuantity menandakan field ditujukan untuk nilai kuantitas/matematika (misal: count, qty, total, amount).
	IdentifierQuantity
	// IdentifierIdentity menandakan field merupakan kode identitas/serial diskrit (misal: postal_code, account_number, NIK).
	IdentifierIdentity
)

// IdentifierSource mendefinisikan sumber atribut DOM dari mana bukti identitas diekstraksi.
type IdentifierSource uint8

const (
	// IdentifierSourceNone menandakan tidak ada atribut sumber yang cocok.
	IdentifierSourceNone IdentifierSource = iota
	// IdentifierSourceName menandakan bukti diekstraksi dari atribut name.
	IdentifierSourceName
	// IdentifierSourceID menandakan bukti diekstraksi dari atribut id.
	IdentifierSourceID
	// IdentifierSourceAutocomplete menandakan bukti diekstraksi dari atribut autocomplete.
	IdentifierSourceAutocomplete
	// IdentifierSourceAriaLabel menandakan bukti diekstraksi dari atribut aria-label atau aria-labelledby.
	IdentifierSourceAriaLabel
	// IdentifierSourcePlaceholder menandakan bukti diekstraksi dari atribut placeholder.
	IdentifierSourcePlaceholder
	// IdentifierSourceTestID menandakan bukti diekstraksi dari atribut data-testid.
	IdentifierSourceTestID
)

// String mengembalikan representasi string dari nama atribut sumber.
func (s IdentifierSource) String() string {
	switch s {
	case IdentifierSourceName:
		return "name"
	case IdentifierSourceID:
		return "id"
	case IdentifierSourceAutocomplete:
		return "autocomplete"
	case IdentifierSourceAriaLabel:
		return "aria-label"
	case IdentifierSourcePlaceholder:
		return "placeholder"
	case IdentifierSourceTestID:
		return "data-testid"
	default:
		return ""
	}
}

// IdentifierEvidence menyimpan informasi rinci bukti semantik yang ditemukan pada elemen input.
type IdentifierEvidence struct {
	Value   string
	Source  IdentifierSource
	Class   IdentifierClass
	Matched string
}

// strongExactIdentityTokens adalah token identitas tunggal yang berdiri sendiri tanpa perlu kualifikasi.
var strongExactIdentityTokens = map[string]struct{}{
	"nik":      {},
	"kk":       {},
	"ktp":      {},
	"npwp":     {},
	"bpjs":     {},
	"pin":      {},
	"otp":      {},
	"passport": {},
	"ssn":      {},
	"serial":   {},
	"zipcode":  {},
	"postcode": {},
}

// strongQualifiedSequences adalah urutan token semantik yang secara definitif menyatakan kode identitas.
var strongQualifiedSequences = [][]string{
	{"postal", "code"},
	{"zip", "code"},
	{"phone", "number"},
	{"phone", "no"},
	{"mobile", "number"},
	{"mobile", "no"},
	{"telepon", "nomor"},
	{"account", "number"},
	{"account", "no"},
	{"rekening", "nomor"},
	{"no", "rekening"},
	{"order", "id"},
	{"order", "number"},
	{"order", "no"},
	{"tracking", "number"},
	{"tracking", "id"},
	{"invoice", "number"},
	{"invoice", "id"},
	{"invoice", "no"},
	{"security", "pin"},
	{"security", "code"},
	{"credit", "card"},
	{"card", "number"},
	{"card", "no"},
	{"nomor", "kk"},
	{"no", "kk"},
	{"license", "number"},
	{"license", "id"},
}

// contextualIdentityTokens adalah token identitas lokal/singkat yang memerlukan resolusi kuantitas.
var contextualIdentityTokens = map[string]struct{}{
	"rt":      {},
	"rw":      {},
	"hp":      {},
	"wa":      {},
	"phone":   {},
	"postal":  {},
	"telepon": {},
}

// quantityModifierTokens adalah kata penunjuk kuantitas atau agregasi numerik.
var quantityModifierTokens = map[string]struct{}{
	"count":      {},
	"total":      {},
	"jumlah":     {},
	"banyak":     {},
	"qty":        {},
	"quantity":   {},
	"amount":     {},
	"volume":     {},
	"capacity":   {},
	"panjang":    {},
	"lebar":      {},
	"tinggi":     {},
	"berat":      {},
	"weight":     {},
	"height":     {},
	"width":      {},
	"length":     {},
	"duration":   {},
	"durasi":     {},
	"age":        {},
	"usia":       {},
	"year":       {},
	"tahun":      {},
	"month":      {},
	"bulan":      {},
	"day":        {},
	"hari":       {},
	"price":      {},
	"harga":      {},
	"rate":       {},
	"score":      {},
	"percent":    {},
	"percentage": {},
	"quota":      {},
	"kuota":      {},
	"size":       {},
	"batch":      {},
}

// authoritativeAutocompleteIdentity memetakan nilai autocomplete standar yang mewakili data identitas diskrit.
var authoritativeAutocompleteIdentity = map[string]string{
	"tel":              "tel",
	"tel-national":     "tel-national",
	"tel-country-code": "tel-country-code",
	"tel-area-code":    "tel-area-code",
	"tel-local":        "tel-local",
	"postal-code":      "postal-code",
	"cc-number":        "cc-number",
	"cc-csc":           "cc-csc",
	"cc-exp":           "cc-exp",
	"one-time-code":    "one-time-code",
}

// TokenizeIdentifier memecah string identifier menjadi token-token kata lowercase berdasarkan
// pemisah tanda baca (_ - . / : + spasi) dan transisi camelCase.
func TokenizeIdentifier(s string) []string {
	if s == "" {
		return nil
	}

	var tokens []string
	var current strings.Builder
	runes := []rune(s)
	n := len(runes)

	flush := func() {
		if current.Len() > 0 {
			tokens = append(tokens, strings.ToLower(current.String()))
			current.Reset()
		}
	}

	for i := 0; i < n; i++ {
		r := runes[i]

		// Pemisah eksplisit
		if r == '_' || r == '-' || r == '.' || r == '/' || r == ':' || r == '+' || unicode.IsSpace(r) {
			flush()
			continue
		}

		if unicode.IsUpper(r) {
			// Transisi lowercase -> Uppercase (e.g. postalCode)
			if i > 0 && unicode.IsLower(runes[i-1]) {
				flush()
			} else if i > 0 && unicode.IsUpper(runes[i-1]) && i+1 < n && unicode.IsLower(runes[i+1]) {
				// Transisi multi-uppercase -> Uppercase+lowercase (e.g. XMLHttp -> XML, Http)
				flush()
			}
			current.WriteRune(r)
		} else {
			current.WriteRune(r)
		}
	}
	flush()

	return tokens
}

// containsSequence memeriksa apakah urutan sub-slice seq muncul secara berurutan di dalam tokens.
func containsSequence(tokens []string, seq ...string) bool {
	if len(seq) == 0 || len(tokens) < len(seq) {
		return false
	}
	for i := 0; i <= len(tokens)-len(seq); i++ {
		match := true
		for j := 0; j < len(seq); j++ {
			if tokens[i+j] != seq[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

// ClassifyIdentifier menganalisis daftar token dan mengklasifikasikan kelas semantik identifier.
// Menggunakan hierarki presedensi ketat:
// 1. Strong Exact Token & Strong Qualified Sequence -> selalu IdentifierIdentity.
// 2. Contextual Identity:
//   - Jika ada token kuantitas -> IdentifierQuantity (misal: rt_count, total_rt, phone_count).
//   - Jika berdiri sendiri/tanpa kuantitas -> IdentifierIdentity (misal: no_rt, rt, phone).
//
// 3. Kuantitas murni -> IdentifierQuantity (misal: item_count, discount_percent).
// 4. Default -> IdentifierUnknown.
func ClassifyIdentifier(tokens []string) (IdentifierClass, string) {
	if len(tokens) == 0 {
		return IdentifierUnknown, ""
	}

	// 1a. Periksa Strong Exact Identity Tokens
	for _, tok := range tokens {
		if _, ok := strongExactIdentityTokens[tok]; ok {
			return IdentifierIdentity, tok
		}
	}

	// 1b. Periksa Strong Qualified Sequences
	for _, seq := range strongQualifiedSequences {
		if containsSequence(tokens, seq...) {
			return IdentifierIdentity, strings.Join(seq, "_")
		}
	}

	// Periksa keberadaan Quantity Modifier
	var hasQuantityModifier bool
	var matchedQuantity string
	for _, tok := range tokens {
		if _, ok := quantityModifierTokens[tok]; ok {
			hasQuantityModifier = true
			matchedQuantity = tok
			break
		}
	}

	// 2. Periksa Contextual Identity Tokens
	for _, tok := range tokens {
		if _, ok := contextualIdentityTokens[tok]; ok {
			if hasQuantityModifier {
				return IdentifierQuantity, matchedQuantity
			}
			return IdentifierIdentity, tok
		}
	}

	// 3. Kuantitas murni
	if hasQuantityModifier {
		return IdentifierQuantity, matchedQuantity
	}

	return IdentifierUnknown, ""
}
