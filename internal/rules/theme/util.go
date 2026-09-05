package theme

import (
	"strings"
)

// TailwindPrimitiveColors adalah daftar nama rumpun palet warna primitif bawaan Tailwind CSS.
var TailwindPrimitiveColors = map[string]bool{
	"slate":   true,
	"gray":    true,
	"zinc":    true,
	"neutral": true,
	"stone":   true,
	"red":     true,
	"orange":  true,
	"amber":   true,
	"yellow":  true,
	"lime":    true,
	"green":   true,
	"emerald": true,
	"teal":    true,
	"cyan":    true,
	"sky":     true,
	"blue":    true,
	"indigo":  true,
	"violet":  true,
	"purple":  true,
	"fuchsia": true,
	"pink":    true,
	"rose":    true,
}

// StripVariants memisahkan rantai varian Tailwind (misal: "md:hover:focus:bg-blue-500")
// menjadi daftar varian ["md", "hover", "focus"] dan utility dasar "bg-blue-500".
// Mendukung arbitrary variant bertanda kurung siku (seperti "[&>svg]:text-white").
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

// SplitAlphaModifier memisahkan utilitas dari modifier slash alpha (seperti "text-white/[0.06]").
func SplitAlphaModifier(token string) (base string, alpha string, hasAlpha bool) {
	// Abaikan slash yang berada di dalam arbitrary bracket (misal: "[color:rgb(0/0/0)]")
	lastSlash := -1
	inBracket := false
	for i := 0; i < len(token); i++ {
		b := token[i]
		switch b {
		case '[':
			inBracket = true
		case ']':
			inBracket = false
		case '/':
			if !inBracket {
				lastSlash = i
			}
		}
	}

	if lastSlash == -1 {
		return token, "", false
	}

	return token[:lastSlash], token[lastSlash+1:], true
}

// IsHexColor memeriksa apakah sebuah string adalah literal warna heksadesimal valid (#rgb, #rgba, #rrggbb, #rrggbbaa).
func IsHexColor(s string) bool {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, "#") {
		return false
	}
	hexPart := s[1:]
	n := len(hexPart)
	if n != 3 && n != 4 && n != 6 && n != 8 {
		return false
	}
	for i := 0; i < n; i++ {
		b := hexPart[i]
		isHex := (b >= '0' && b <= '9') || (b >= 'a' && b <= 'f') || (b >= 'A' && b <= 'F')
		if !isHex {
			return false
		}
	}
	return true
}

// IsColorFunction memeriksa apakah string diawali fungsi warna CSS seperti rgb(), rgba(), hsl(), hsla(), oklch(), oklab().
func IsColorFunction(s string) bool {
	s = strings.ToLower(strings.TrimSpace(s))
	prefixes := []string{"rgb(", "rgba(", "hsl(", "hsla(", "oklch(", "oklab("}
	for _, p := range prefixes {
		if strings.HasPrefix(s, p) && strings.HasSuffix(s, ")") {
			return true
		}
	}
	return false
}

// IsTailwindPrimitiveColor memeriksa apakah segmen nama warna adalah palet primitif bawaan Tailwind (misal: "blue-600").
func IsTailwindPrimitiveColor(colorName string) bool {
	dashIdx := strings.LastIndexByte(colorName, '-')
	if dashIdx == -1 {
		return false
	}
	palette := colorName[:dashIdx]
	scale := colorName[dashIdx+1:]

	if !TailwindPrimitiveColors[palette] {
		return false
	}

	switch scale {
	case "50", "100", "200", "300", "400", "500", "600", "700", "800", "900", "950":
		return true
	default:
		return false
	}
}

// ParseArbitraryProperty mengekstrak properti dan nilai dari sintaks arbitrary property Tailwind (misal: "[color:#fff]").
func ParseArbitraryProperty(token string) (prop, val string, ok bool) {
	if len(token) < 4 || token[0] != '[' || token[len(token)-1] != ']' {
		return "", "", false
	}
	inner := token[1 : len(token)-1]
	colonIdx := strings.IndexByte(inner, ':')
	if colonIdx == -1 || colonIdx == 0 || colonIdx == len(inner)-1 {
		return "", "", false
	}

	p := strings.TrimSpace(inner[:colonIdx])
	v := strings.TrimSpace(inner[colonIdx+1:])
	return p, v, true
}

// IsColorProperty memeriksa apakah sebuah nama properti CSS berhubungan dengan pewarnaan.
func IsColorProperty(prop string) bool {
	p := strings.ToLower(strings.TrimSpace(prop))
	switch p {
	case "color", "background-color", "background", "border-color",
		"outline-color", "fill", "stroke", "stop-color", "ring-color",
		"accent-color", "caret-color":
		return true
	default:
		return false
	}
}

// ExtractRawColorFromArbitrary mengekstrak nilai warna dari class kurung siku Tailwind,
// misalnya "bg-[#1e293b]" -> ("#1e293b", true) atau "text-[rgb(255,0,0)]" -> ("rgb(255,0,0)", true).
func ExtractRawColorFromArbitrary(base string) (string, bool) {
	bracketStart := strings.IndexByte(base, '[')
	bracketEnd := strings.LastIndexByte(base, ']')
	if bracketStart == -1 || bracketEnd <= bracketStart+1 {
		return "", false
	}

	content := base[bracketStart+1 : bracketEnd]
	content = strings.TrimSpace(content)

	// Jika berformat properti "[color:#fff]"
	if colonIdx := strings.IndexByte(content, ':'); colonIdx != -1 {
		prop := strings.TrimSpace(content[:colonIdx])
		val := strings.TrimSpace(content[colonIdx+1:])
		if IsColorProperty(prop) && (IsHexColor(val) || IsColorFunction(val)) {
			return val, true
		}
		return "", false
	}

	if IsHexColor(content) || IsColorFunction(content) {
		return content, true
	}

	return "", false
}

// IsMonochromeColor memeriksa apakah nama warna adalah monokrom statis (white / black).
func IsMonochromeColor(s string) bool {
	return s == "white" || s == "black"
}

// OrderedColorPrefixes adalah daftar prefix utility pewarnaan Tailwind terurut dari yang terpanjang.
var OrderedColorPrefixes = []string{
	"border-t-", "border-r-", "border-b-", "border-l-", "border-x-", "border-y-",
	"divide-x-", "divide-y-", "divide-",
	"ring-offset-", "ring-",
	"outline-",
	"border-",
	"bg-",
	"text-",
	"fill-",
	"stroke-",
	"accent-",
	"caret-",
	"decoration-",
	"from-", "via-", "to-",
}

// SplitColorPrefix memisahkan prefix pewarnaan dari sisa token warna.
func SplitColorPrefix(token string) (prefix, remainder string, ok bool) {
	for _, p := range OrderedColorPrefixes {
		if strings.HasPrefix(token, p) {
			return p, token[len(p):], true
		}
	}
	return "", token, false
}

// IsPseudoVariant memeriksa apakah sebuah identifier varian adalah pseudo element atau pseudo class bertarget.
func IsPseudoVariant(variant string) bool {
	switch variant {
	case "placeholder", "selection", "file", "marker", "backdrop", "before", "after":
		return true
	default:
		return false
	}
}

// HasPseudoVariant memeriksa apakah salah satu varian dalam rantai adalah pseudo variant.
func HasPseudoVariant(variants []string) bool {
	for _, v := range variants {
		if IsPseudoVariant(v) {
			return true
		}
	}
	return false
}

// IsNonColorBorderKeyword mengidentifikasi keyword border non-warna (seperti width, style, collapse).
func IsNonColorBorderKeyword(s string) bool {
	switch s {
	case "0", "2", "4", "8",
		"solid", "dashed", "dotted", "double", "hidden", "none",
		"collapse", "separate":
		return true
	default:
		return false
	}
}

// IsColorStyleProperty memeriksa apakah nama properti inline style CSS/JSX berhubungan dengan warna.
func IsColorStyleProperty(p string) bool {
	switch p {
	case "color", "background", "background-color", "backgroundcolor",
		"border-color", "bordercolor", "fill", "stroke", "border",
		"outline-color", "outlinecolor":
		return true
	default:
		return false
	}
}

// IsSafeColorValue memeriksa apakah nilai warna aman (menggunakan CSS variable atau keyword sistem).
func IsSafeColorValue(val string) bool {
	return strings.Contains(val, "var(--") ||
		val == "none" || val == "transparent" || val == "currentcolor" ||
		val == "inherit" || val == "initial" || val == "unset"
}

// CheckHardcodedColorValue memeriksa apakah suatu nilai berisi literal warna heksadesimal, fungsi warna, atau nama primitif.
func CheckHardcodedColorValue(val string) bool {
	if IsSafeColorValue(val) {
		return false
	}
	if strings.Contains(val, "#") {
		return true
	}
	if IsColorFunction(val) || strings.Contains(val, "rgb(") || strings.Contains(val, "rgba(") ||
		strings.Contains(val, "hsl(") || strings.Contains(val, "hsla(") ||
		strings.Contains(val, "oklch(") || strings.Contains(val, "oklab(") {
		return true
	}
	if IsMonochromeColor(val) || TailwindPrimitiveColors[val] {
		return true
	}
	return false
}

// ExtractInlineStyleHardcodedColor memeriksa apakah nilai atribut style memuat deklarasi warna yang di-hardcode.
func ExtractInlineStyleHardcodedColor(styleVal string) (string, bool) {
	// Normalisasi pemisah deklarasi (; untuk CSS string, , untuk JSX object)
	clean := strings.Trim(styleVal, " {}\"'`")
	decls := strings.FieldsFunc(clean, func(r rune) bool {
		return r == ';' || r == ','
	})

	for _, decl := range decls {
		decl = strings.TrimSpace(decl)
		colonIdx := strings.IndexByte(decl, ':')
		if colonIdx == -1 {
			continue
		}
		prop := strings.ToLower(strings.Trim(decl[:colonIdx], " \"'"))
		val := strings.ToLower(strings.Trim(decl[colonIdx+1:], " \"'"))

		if !IsColorStyleProperty(prop) {
			continue
		}
		if CheckHardcodedColorValue(val) {
			return prop + ": " + val, true
		}
	}

	return "", false
}
