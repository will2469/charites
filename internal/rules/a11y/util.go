package a11y

import (
	"strconv"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// StripVariants memisahkan rantai varian Tailwind (misal: "sm:focus:outline-none")
// menjadi daftar varian ["sm", "focus"] dan utility dasar "outline-none".
func StripVariants(token string) ([]string, string) {
	if !strings.Contains(token, ":") {
		return nil, token
	}

	var variants []string
	curr := token
	inBracket := false

	for i := 0; i < len(curr); i++ {
		b := curr[i]
		switch b {
		case '[':
			inBracket = true
		case ']':
			inBracket = false
		case ':':
			if !inBracket {
				variant := curr[:i]
				variants = append(variants, variant)
				curr = curr[i+1:]
				i = -1
			}
		}
	}

	return variants, curr
}

// StripVariantsOnlyBase mengembalikan utility dasar tanpa mengalokasikan slice varian.
// Mengembalikan substring langsung tanpa alokasi memori (0 B/op, 0 allocs/op).
func StripVariantsOnlyBase(token string) string {
	if !strings.Contains(token, ":") {
		return token
	}

	lastColon := -1
	inBracket := false
	for i := 0; i < len(token); i++ {
		b := token[i]
		switch b {
		case '[':
			inBracket = true
		case ']':
			inBracket = false
		case ':':
			if !inBracket {
				lastColon = i
			}
		}
	}

	if lastColon != -1 && lastColon+1 < len(token) {
		return token[lastColon+1:]
	}
	return token
}

// HasResponsivePrefix memeriksa apakah token memiliki prefix breakpoint responsif (sm:, md:, lg:, xl:, 2xl:).
func HasResponsivePrefix(token string) bool {
	variants, _ := StripVariants(token)
	for _, v := range variants {
		switch v {
		case "sm", "md", "lg", "xl", "2xl":
			return true
		}
	}
	return false
}

// HasMinBreakpointPrefix memeriksa apakah token memiliki modifier min-width responsive breakpoint (sm:, md:, lg:, xl:, 2xl:).
func HasMinBreakpointPrefix(token string) bool {
	if !strings.Contains(token, ":") {
		return false
	}
	variants, _ := StripVariants(token)
	for _, v := range variants {
		switch v {
		case "sm", "md", "lg", "xl", "2xl":
			return true
		}
	}
	return false
}

// IsInteractiveElement memeriksa apakah node merupakan elemen interaktif
// (<button>, <a>, <input>, <select>, <textarea>, <summary>, atau elemen dengan role="button" / onClick / tabIndex non-negatif).
func IsInteractiveElement(node *ir.Node) bool {
	if node == nil || node.Type != ir.NodeElement {
		return false
	}

	tag := strings.ToLower(node.Tag)
	switch tag {
	case "button", "a", "input", "select", "textarea", "summary":
		return true
	}

	if node.Attributes != nil {
		if role, ok := node.Attributes["role"]; ok && strings.EqualFold(strings.Trim(strings.TrimSpace(role), "\"'`"), "button") {
			return true
		}
		if _, ok := node.Attributes["onClick"]; ok {
			return true
		}
		if _, ok := node.Attributes["onclick"]; ok {
			return true
		}
		if ti, ok := node.Attributes["tabIndex"]; ok && strings.Trim(strings.TrimSpace(ti), "\"'`") != "-1" {
			return true
		}
		if ti, ok := node.Attributes["tabindex"]; ok && strings.Trim(strings.TrimSpace(ti), "\"'`") != "-1" {
			return true
		}
	}

	return false
}

// IsTextualInput memeriksa apakah node merupakan kontrol teks formulir (<input>, <select>, <textarea>)
// yang rentan terhadap auto-zoom viewport atau text clipping.
func IsTextualInput(node *ir.Node) bool {
	if node == nil || node.Type != ir.NodeElement {
		return false
	}

	tag := strings.ToLower(node.Tag)
	switch tag {
	case "select", "textarea":
		return true
	case "input":
		if node.Attributes != nil {
			inputType := strings.Trim(strings.TrimSpace(strings.ToLower(node.Attributes["type"])), "\"'`")
			switch inputType {
			case "hidden", "checkbox", "radio", "range", "button", "submit", "reset", "image", "color", "file":
				return false
			}
		}
		return true
	}

	return false
}

// ParseTailwindSizeToPx mengonversi token numerik Tailwind (misal: "6" dari "w-6" atau "11" dari "h-11")
// menjadi nilai piksel. Mengembalikan false jika token bukan ukuran yang dikenal.
func ParseTailwindSizeToPx(sizeStr string) (float64, bool) {
	// Dukungan arbitrary value: "[24px]" atau "[2.5rem]"
	if strings.HasPrefix(sizeStr, "[") && strings.HasSuffix(sizeStr, "]") {
		raw := sizeStr[1 : len(sizeStr)-1]
		if strings.HasSuffix(raw, "px") {
			val, err := strconv.ParseFloat(strings.TrimSuffix(raw, "px"), 64)
			if err == nil {
				return val, true
			}
		}
		if strings.HasSuffix(raw, "rem") {
			val, err := strconv.ParseFloat(strings.TrimSuffix(raw, "rem"), 64)
			if err == nil {
				return val * 16.0, true // 1rem = 16px default
			}
		}
		return 0, false
	}

	switch sizeStr {
	case "0":
		return 0, true
	case "0.5":
		return 2, true
	case "1":
		return 4, true
	case "1.5":
		return 6, true
	case "2":
		return 8, true
	case "2.5":
		return 10, true
	case "3":
		return 12, true
	case "3.5":
		return 14, true
	case "4":
		return 16, true
	case "5":
		return 20, true
	case "6":
		return 24, true
	case "7":
		return 28, true
	case "8":
		return 32, true
	case "9":
		return 36, true
	case "10":
		return 40, true
	case "11":
		return 44, true
	case "12":
		return 48, true
	case "14":
		return 56, true
	case "16":
		return 64, true
	case "20":
		return 80, true
	case "24":
		return 96, true
	}

	// Coba parse angka numerik reguler * 4
	val, err := strconv.ParseFloat(sizeStr, 64)
	if err == nil {
		return val * 4.0, true
	}

	return 0, false
}

// ParseFontSizeToPx mengekstrak estimasi ukuran font piksel dari kelas Tailwind (misal: text-sm -> 14px).
func ParseFontSizeToPx(baseClass string) (float64, bool) {
	if !strings.HasPrefix(baseClass, "text-") {
		return 0, false
	}
	sub := strings.TrimPrefix(baseClass, "text-")

	// Arbitrary value: text-[14px]
	if strings.HasPrefix(sub, "[") && strings.HasSuffix(sub, "]") {
		raw := sub[1 : len(sub)-1]
		if strings.HasSuffix(raw, "px") {
			val, err := strconv.ParseFloat(strings.TrimSuffix(raw, "px"), 64)
			if err == nil {
				return val, true
			}
		}
		if strings.HasSuffix(raw, "rem") {
			val, err := strconv.ParseFloat(strings.TrimSuffix(raw, "rem"), 64)
			if err == nil {
				return val * 16.0, true
			}
		}
		return 0, false
	}

	switch sub {
	case "xs":
		return 12.0, true
	case "sm":
		return 14.0, true
	case "base":
		return 16.0, true
	case "lg":
		return 18.0, true
	case "xl":
		return 20.0, true
	case "2xl":
		return 24.0, true
	case "3xl":
		return 30.0, true
	case "4xl":
		return 36.0, true
	case "5xl":
		return 48.0, true
	case "6xl":
		return 60.0, true
	case "7xl":
		return 72.0, true
	case "8xl":
		return 96.0, true
	case "9xl":
		return 128.0, true
	}

	return 0, false
}

// CleanAttr membersihkan pembungkus kutip, kurung kurawal ekspresi JSX, dan spasi dari nilai atribut mentah.
func CleanAttr(val string) string {
	return strings.Trim(strings.TrimSpace(val), "\"'`{}")
}

// FindRoot menelusuri pointer Parent ke atas hingga menemukan node akar dokumen.
func FindRoot(node *ir.Node) *ir.Node {
	if node == nil {
		return nil
	}
	curr := node
	for curr.Parent != nil {
		curr = curr.Parent
	}
	return curr
}

// HasEnclosingLabel memeriksa apakah node dibungkus di dalam elemen <label> pada rantai ancestor.
func HasEnclosingLabel(node *ir.Node) bool {
	if node == nil {
		return false
	}
	curr := node.Parent
	for curr != nil {
		if curr.Type == ir.NodeElement && strings.EqualFold(curr.Tag, "label") {
			return true
		}
		curr = curr.Parent
	}
	return false
}

// HasDocumentElementWithID memeriksa apakah terdapat elemen dalam dokumen dengan id yang cocok.
func HasDocumentElementWithID(root *ir.Node, targetID string) bool {
	if root == nil || targetID == "" {
		return false
	}

	cleanTarget := CleanAttr(targetID)
	// Jika target id adalah ekspresi dinamis (mengandung operasi JS/interpolasi), jangan beri false-positive
	if strings.ContainsAny(cleanTarget, "${}+()?:") {
		return true
	}

	for n := range root.Walk() {
		if n.Type != ir.NodeElement || n.Attributes == nil {
			continue
		}
		if rawID, ok := n.Attributes["id"]; ok {
			if CleanAttr(rawID) == cleanTarget {
				return true
			}
		}
	}

	return false
}

// HasAssociatedLabel memeriksa apakah kontrol input memiliki label asosiasi,
// baik melalui elemen pembungkus <label> maupun deklarasi <label htmlFor="id"> di dokumen.
func HasAssociatedLabel(root, inputNode *ir.Node, inputID string) bool {
	if HasEnclosingLabel(inputNode) {
		return true
	}

	cleanID := CleanAttr(inputID)
	if cleanID == "" || root == nil {
		return false
	}

	for n := range root.Walk() {
		if n.Type != ir.NodeElement || !strings.EqualFold(n.Tag, "label") || n.Attributes == nil {
			continue
		}

		if forVal, ok := n.Attributes["htmlFor"]; ok && CleanAttr(forVal) == cleanID {
			return true
		}
		if forVal, ok := n.Attributes["for"]; ok && CleanAttr(forVal) == cleanID {
			return true
		}
	}

	return false
}
