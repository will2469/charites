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

	if lastColon == -1 {
		return token
	}
	return token[lastColon+1:]
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

// HasHardcodedScalarUnit memeriksa apakah nilai string memuat unit skalar (px, rem, em, pt) tanpa var(--...).
func HasHardcodedScalarUnit(val string) bool {
	if strings.Contains(val, "var(--") {
		return false
	}
	for _, unit := range []string{"px", "rem", "em", "pt"} {
		if strings.HasSuffix(val, unit) {
			numPart := val[:len(val)-len(unit)]
			if len(numPart) > 0 && isNumericString(numPart) {
				return true
			}
		}
	}
	return false
}

func isNumericString(s string) bool {
	hasDigit := false
	for i := 0; i < len(s); i++ {
		b := s[i]
		if b >= '0' && b <= '9' {
			hasDigit = true
		} else if b != '.' && b != '-' {
			return false
		}
	}
	return hasDigit
}

// SpatialSizePrefixes adalah daftar prefix Tailwind yang mengatur ukuran spasial dan tipografi.
var SpatialSizePrefixes = []string{
	"min-w-[", "max-w-[", "min-h-[", "max-h-[",
	"gap-x-[", "gap-y-[", "space-x-[", "space-y-[",
	"inset-x-[", "inset-y-[",
	"size-[", "gap-[", "inset-[",
	"px-[", "py-[", "pt-[", "pb-[", "pl-[", "pr-[", "ps-[", "pe-[",
	"mx-[", "my-[", "mt-[", "mb-[", "ml-[", "mr-[", "ms-[", "me-[",
	"w-[", "h-[", "p-[", "m-[",
	"top-[", "bottom-[", "left-[", "right-[",
	"leading-[", "tracking-[",
}

// SpatialPropertyNames adalah daftar nama properti CSS yang mengatur dimensi dan spacing.
var SpatialPropertyNames = map[string]bool{
	"width": true, "height": true, "min-width": true, "max-width": true,
	"min-height": true, "max-height": true,
	"padding": true, "padding-top": true, "padding-bottom": true, "padding-left": true, "padding-right": true,
	"margin": true, "margin-top": true, "margin-bottom": true, "margin-left": true, "margin-right": true,
	"gap": true, "column-gap": true, "row-gap": true,
	"top": true, "bottom": true, "left": true, "right": true,
	"font-size": true, "line-height": true, "letter-spacing": true,
}

// IsHardcodedSizeUtility mendeteksi apakah token menggunakan utility dimensi, spacing, posisi, atau tipografi hardcode.
func IsHardcodedSizeUtility(base string) bool {
	if prop, val, ok := ParseArbitraryProperty(base); ok {
		if SpatialPropertyNames[strings.ToLower(prop)] {
			return HasHardcodedScalarUnit(val)
		}
		return false
	}

	// Tangani typography text-[...px] / text-[...rem]
	if strings.HasPrefix(base, "text-[") && strings.HasSuffix(base, "]") {
		inner := base[6 : len(base)-1]
		if !IsHexColor(inner) && !IsColorFunction(inner) && HasHardcodedScalarUnit(inner) {
			return true
		}
		return false
	}

	for _, prefix := range SpatialSizePrefixes {
		if strings.HasPrefix(base, prefix) && strings.HasSuffix(base, "]") {
			inner := base[len(prefix) : len(base)-1]
			if HasHardcodedScalarUnit(inner) {
				return true
			}
		}
	}
	return false
}

// BorderRadiusPrefixes adalah daftar prefix border radius Tailwind.
var BorderRadiusPrefixes = []string{
	"rounded-tl-[", "rounded-tr-[", "rounded-bl-[", "rounded-br-[",
	"rounded-ss-[", "rounded-se-[", "rounded-es-[", "rounded-ee-[",
	"rounded-t-[", "rounded-b-[", "rounded-l-[", "rounded-r-[",
	"rounded-s-[", "rounded-e-[",
	"rounded-[",
}

// IsHardcodedBorderRadius mendeteksi radius sudut arbitrer.
func IsHardcodedBorderRadius(base string) bool {
	if prop, val, ok := ParseArbitraryProperty(base); ok {
		p := strings.ToLower(prop)
		if p == "border-radius" || (strings.HasPrefix(p, "border-") && strings.HasSuffix(p, "-radius")) {
			return HasHardcodedScalarUnit(val)
		}
		return false
	}

	for _, prefix := range BorderRadiusPrefixes {
		if strings.HasPrefix(base, prefix) && strings.HasSuffix(base, "]") {
			inner := base[len(prefix) : len(base)-1]
			if strings.Contains(inner, "var(--") {
				return false
			}
			return HasHardcodedScalarUnit(inner) || isNumericString(inner)
		}
	}
	return false
}

// IsHardcodedZIndex mendeteksi z-index arbitrer skalar.
func IsHardcodedZIndex(base string) bool {
	if prop, val, ok := ParseArbitraryProperty(base); ok {
		if strings.ToLower(prop) == "z-index" {
			return !strings.Contains(val, "var(--") && isNumericString(val)
		}
		return false
	}

	if strings.HasPrefix(base, "z-[") && strings.HasSuffix(base, "]") {
		inner := base[3 : len(base)-1]
		if strings.Contains(inner, "var(--") {
			return false
		}
		return isNumericString(inner)
	}
	return false
}

// IsHardcodedShadowColor mendeteksi apakah bayangan arbitrer menyematkan warna hardcode.
func IsHardcodedShadowColor(base string) bool {
	if prop, val, ok := ParseArbitraryProperty(base); ok {
		if strings.ToLower(prop) == "box-shadow" {
			return CheckHardcodedColorValue(val)
		}
		return false
	}

	if strings.HasPrefix(base, "shadow-[") && strings.HasSuffix(base, "]") {
		inner := base[8 : len(base)-1]
		if strings.Contains(inner, "#") {
			return true
		}
		if strings.Contains(inner, "rgb(") || strings.Contains(inner, "rgba(") ||
			strings.Contains(inner, "hsl(") || strings.Contains(inner, "hsla(") ||
			strings.Contains(inner, "oklch(") || strings.Contains(inner, "oklab(") {
			return true
		}
	}
	return false
}

// IsHardcodedBackdropBlur mendeteksi blur atau backdrop-blur arbitrer.
func IsHardcodedBackdropBlur(base string) bool {
	if prop, val, ok := ParseArbitraryProperty(base); ok {
		p := strings.ToLower(prop)
		if p == "backdrop-filter" || p == "filter" {
			if idx := strings.Index(val, "blur("); idx != -1 {
				inner := val[idx+5:]
				if closeIdx := strings.IndexByte(inner, ')'); closeIdx != -1 {
					inner = strings.TrimSpace(inner[:closeIdx])
					return HasHardcodedScalarUnit(inner) || isNumericString(inner)
				}
			}
		}
		return false
	}

	if strings.HasPrefix(base, "backdrop-blur-[") && strings.HasSuffix(base, "]") {
		inner := base[15 : len(base)-1]
		return HasHardcodedScalarUnit(inner) || isNumericString(inner)
	}
	if strings.HasPrefix(base, "blur-[") && strings.HasSuffix(base, "]") {
		inner := base[6 : len(base)-1]
		return HasHardcodedScalarUnit(inner) || isNumericString(inner)
	}
	return false
}

// SafeRingTokens adalah daftar token ring atau outline bawaan yang aman.
var SafeRingTokens = map[string]bool{
	"ring": true, "ring-1": true, "ring-2": true, "ring-4": true, "ring-8": true,
	"ring-inset":    true,
	"ring-offset-0": true, "ring-offset-1": true, "ring-offset-2": true, "ring-offset-4": true, "ring-offset-8": true,
	"outline-none": true, "outline-0": true, "outline-1": true, "outline-2": true, "outline-4": true, "outline-8": true,
	"outline-dashed": true, "outline-dotted": true, "outline-double": true,
	"ring-ring": true, "ring-primary": true, "ring-border": true, "ring-background": true,
	"ring-offset-background": true, "ring-offset-ring": true,
}

// FocusRingBracketPrefixes adalah daftar prefix kurung siku arbitrer untuk ring dan outline fokus.
var FocusRingBracketPrefixes = []string{"ring-[", "ring-offset-[", "outline-["}

func isArbitraryFocusRingBracket(base string) bool {
	for _, prefix := range FocusRingBracketPrefixes {
		if strings.HasPrefix(base, prefix) && strings.HasSuffix(base, "]") {
			inner := base[len(prefix) : len(base)-1]
			if strings.Contains(inner, "var(--") {
				return false
			}
			return IsHexColor(inner) || IsColorFunction(inner)
		}
	}
	return false
}

func isPrimitiveFocusRing(base string) bool {
	if strings.HasPrefix(base, "ring-offset-") {
		rem := base[12:]
		return IsMonochromeColor(rem) || IsTailwindPrimitiveColor(rem)
	}
	if strings.HasPrefix(base, "ring-") {
		return IsTailwindPrimitiveColor(base[5:])
	}
	if strings.HasPrefix(base, "outline-") {
		return IsTailwindPrimitiveColor(base[8:])
	}
	return false
}

func isArbitraryFocusRingProperty(base string) bool {
	prop, val, ok := ParseArbitraryProperty(base)
	if !ok {
		return false
	}
	p := strings.ToLower(prop)
	if !strings.HasPrefix(p, "ring-") && !strings.HasPrefix(p, "outline-") {
		return false
	}
	return CheckHardcodedColorValue(val)
}

// IsHardcodedFocusRing mendeteksi ring atau outline fokus dengan warna mentah / primitif.
func IsHardcodedFocusRing(base string) bool {
	if SafeRingTokens[base] {
		return false
	}
	return isArbitraryFocusRingProperty(base) ||
		isArbitraryFocusRingBracket(base) ||
		isPrimitiveFocusRing(base)
}

// IsElevatedShadow memeriksa apakah suatu class adalah bayangan elevasi kontainer (shadow, shadow-md, shadow-lg, shadow-xl, shadow-2xl).
func IsElevatedShadow(base string) bool {
	switch base {
	case "shadow", "shadow-md", "shadow-lg", "shadow-xl", "shadow-2xl":
		return true
	default:
		return false
	}
}

// HasBorderOrRing memeriksa apakah dalam daftar class terdapat deklarasi border atau ring aktif.
func HasBorderOrRing(classes []string) bool {
	for _, class := range classes {
		base := StripVariantsOnlyBase(class)
		if isBorderClass(base) || isRingClass(base) {
			return true
		}
	}
	return false
}

func isBorderClass(base string) bool {
	if base == "border" {
		return true
	}
	if strings.HasPrefix(base, "border-") {
		rem := base[7:]
		switch rem {
		case "0", "none", "transparent":
			return false
		default:
			return true
		}
	}
	return false
}

func isRingClass(base string) bool {
	if base == "ring" {
		return true
	}
	if strings.HasPrefix(base, "ring-") {
		rem := base[5:]
		switch rem {
		case "0", "none", "transparent", "inset":
			return false
		default:
			return true
		}
	}
	return false
}

// HasDarkVariant memeriksa apakah dalam daftar class terdapat class yang diawali varian dark:
func HasDarkVariant(classes []string) bool {
	for _, class := range classes {
		if strings.HasPrefix(class, "dark:") || strings.Contains(class, ":dark:") {
			return true
		}
	}
	return false
}

// HasContainerOpacity memeriksa apakah class kontainer memuat modifier opacity atau background slash opacity.
func HasContainerOpacity(classes []string) (string, bool) {
	for _, class := range classes {
		base := StripVariantsOnlyBase(class)
		if strings.HasPrefix(base, "opacity-") && len(base) > 8 {
			return class, true
		}
		if strings.HasPrefix(base, "bg-") && strings.Contains(base, "/") {
			return class, true
		}
	}
	return "", false
}

// HasTextOrInnerOpacity memeriksa apakah elemen anak memuat teks dengan slash opacity atau modifier opacity.
func HasTextOrInnerOpacity(classes []string) (string, bool) {
	for _, class := range classes {
		base := StripVariantsOnlyBase(class)
		if strings.HasPrefix(base, "opacity-") && len(base) > 8 {
			return class, true
		}
		if strings.HasPrefix(base, "text-") && strings.Contains(base, "/") {
			return class, true
		}
	}
	return "", false
}

// IsThemeGraphicAsset mendeteksi apakah path/URL gambar merujuk pada aset visual seperti SVG, logo, ilustrasi, diagram, atau chart.
func IsThemeGraphicAsset(src string) bool {
	lower := strings.ToLower(strings.Trim(strings.TrimSpace(src), "\"'`"))
	if lower == "" {
		return false
	}
	if strings.HasSuffix(lower, ".svg") {
		return true
	}
	graphicKeywords := []string{"logo", "diagram", "illustration", "chart", "graphic", "scheme"}
	for _, kw := range graphicKeywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}

// SVGElementNames adalah himpunan nama elemen SVG umum yang mendukung atribut pewarnaan vektor.
var SVGElementNames = map[string]bool{
	"svg": true, "path": true, "rect": true, "circle": true,
	"polygon": true, "ellipse": true, "line": true, "polyline": true,
	"stop": true, "g": true, "mask": true, "pattern": true, "text": true,
}

// IsSVGElement memeriksa apakah tag HTML/JSX merupakan elemen SVG.
func IsSVGElement(tag string) bool {
	return SVGElementNames[strings.ToLower(tag)]
}

// IsHardcodedSVGAttribute memeriksa apakah pasangan nama atribut dan nilai SVG memuat warna heksadesimal atau primitif hardcoded.
func IsHardcodedSVGAttribute(attrName, attrVal string) bool {
	name := strings.ToLower(strings.TrimSpace(attrName))
	switch name {
	case "fill", "stroke", "stop-color", "stopcolor":
		val := strings.Trim(strings.TrimSpace(attrVal), "\"'`")
		if val == "" || IsSafeColorValue(val) || strings.HasPrefix(val, "url(#") {
			return false
		}
		if IsHexColor(val) || IsColorFunction(val) || IsMonochromeColor(val) || TailwindPrimitiveColors[val] {
			return true
		}
	}
	return false
}
