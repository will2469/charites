package drift

import (
	"math"
	"strconv"
	"strings"
)

// TailwindPaletteSrgb mendefinisikan kamus referensi statis nilai sRGB hex untuk palet primitif Tailwind CSS.
// Berfungsi sebagai acuan komputasi kontras saat fase analisis statis tanpa melakukan evaluasi CSS dinamis.
var TailwindPaletteSrgb = map[string]string{
	// Monokrom dasar
	"white": "#ffffff",
	"black": "#000000",

	// Yellow
	"yellow-50":  "#fefce8",
	"yellow-100": "#fef9c3",
	"yellow-200": "#fef08a",
	"yellow-300": "#fde047",
	"yellow-400": "#facc15",
	"yellow-500": "#eab308",
	"yellow-600": "#ca8a04",
	"yellow-700": "#a16207",
	"yellow-800": "#854d0e",
	"yellow-900": "#713f12",
	"yellow-950": "#422006",

	// Amber
	"amber-50":  "#fffbeb",
	"amber-100": "#fef3c7",
	"amber-200": "#fde68a",
	"amber-300": "#fcd34d",
	"amber-400": "#fbbf24",
	"amber-500": "#f59e0b",
	"amber-600": "#d97706",
	"amber-700": "#b45309",
	"amber-800": "#92400e",
	"amber-900": "#78350f",
	"amber-950": "#451a03",

	// Orange
	"orange-50":  "#fff7ed",
	"orange-100": "#ffedd5",
	"orange-200": "#fed7aa",
	"orange-300": "#fdba74",
	"orange-400": "#fb923c",
	"orange-500": "#f97316",
	"orange-600": "#ea580c",
	"orange-700": "#c2410c",
	"orange-800": "#9a3412",
	"orange-900": "#7c2d12",
	"orange-950": "#431407",

	// Red
	"red-50":  "#fef2f2",
	"red-100": "#fee2e2",
	"red-200": "#fecaca",
	"red-300": "#fca5a5",
	"red-400": "#f87171",
	"red-500": "#ef4444",
	"red-600": "#dc2626",
	"red-700": "#b91c1c",
	"red-800": "#991b1b",
	"red-900": "#7f1d1d",
	"red-950": "#450a0a",

	// Rose
	"rose-50":  "#fff1f2",
	"rose-100": "#ffe4e6",
	"rose-200": "#fecdd3",
	"rose-300": "#fda4af",
	"rose-400": "#fb7185",
	"rose-500": "#f43f5e",
	"rose-600": "#e11d48",
	"rose-700": "#be123c",
	"rose-800": "#9f1239",
	"rose-900": "#881337",
	"rose-950": "#4c0519",

	// Green
	"green-50":  "#f0fdf4",
	"green-100": "#dcfce7",
	"green-200": "#bbf7d0",
	"green-300": "#86efac",
	"green-400": "#4ade80",
	"green-500": "#22c55e",
	"green-600": "#16a34a",
	"green-700": "#15803d",
	"green-800": "#166534",
	"green-900": "#14532d",
	"green-950": "#052e16",

	// Emerald
	"emerald-50":  "#ecfdf5",
	"emerald-100": "#d1fae5",
	"emerald-200": "#a7f3d0",
	"emerald-300": "#6ee7b7",
	"emerald-400": "#34d399",
	"emerald-500": "#10b981",
	"emerald-600": "#059669",
	"emerald-700": "#047857",
	"emerald-800": "#065f46",
	"emerald-900": "#064e3b",
	"emerald-950": "#022c22",

	// Teal
	"teal-50":  "#f0fdfa",
	"teal-100": "#ccfbf1",
	"teal-200": "#99f6e4",
	"teal-300": "#5eead4",
	"teal-400": "#2dd4bf",
	"teal-500": "#14b8a6",
	"teal-600": "#0d9488",
	"teal-700": "#0f766e",
	"teal-800": "#115e59",
	"teal-900": "#134e4a",
	"teal-950": "#042f2e",

	// Blue
	"blue-50":  "#eff6ff",
	"blue-100": "#dbeafe",
	"blue-200": "#bfdbfe",
	"blue-300": "#93c5fd",
	"blue-400": "#60a5fa",
	"blue-500": "#3b82f6",
	"blue-600": "#2563eb",
	"blue-700": "#1d4ed8",
	"blue-800": "#1e40af",
	"blue-900": "#1e3a8a",
	"blue-950": "#172554",

	// Sky
	"sky-50":  "#f0f9ff",
	"sky-100": "#e0f2fe",
	"sky-200": "#bae6fd",
	"sky-300": "#7dd3fc",
	"sky-400": "#38bdf8",
	"sky-500": "#0ea5e9",
	"sky-600": "#0284c7",
	"sky-700": "#0369a1",
	"sky-800": "#075985",
	"sky-900": "#0c4a6e",
	"sky-950": "#082f49",

	// Cyan
	"cyan-50":  "#ecfeff",
	"cyan-100": "#cffafe",
	"cyan-200": "#a5f3fc",
	"cyan-300": "#67e8f9",
	"cyan-400": "#22d3ee",
	"cyan-500": "#06b6d4",
	"cyan-600": "#0891b2",
	"cyan-700": "#0e7490",
	"cyan-800": "#155e75",
	"cyan-900": "#164e63",
	"cyan-950": "#083344",
}

// RelativeLuminance menghitung luminansi relatif sebuah warna sRGB sesuai formula resmi W3C WCAG 2.1.
func RelativeLuminance(r, g, b uint8) float64 {
	transform := func(val uint8) float64 {
		s := float64(val) / 255.0
		if s <= 0.04045 {
			return s / 12.92
		}
		return math.Pow((s+0.055)/1.055, 2.4)
	}

	rLinear := transform(r)
	gLinear := transform(g)
	bLinear := transform(b)

	return 0.2126*rLinear + 0.7152*gLinear + 0.0722*bLinear
}

// ParseHexColor mem-parse string hex (#ffffff atau #fff) menjadi komponen R, G, B.
func ParseHexColor(hex string) (uint8, uint8, uint8, bool) {
	clean := strings.TrimPrefix(strings.TrimSpace(hex), "#")
	if len(clean) == 3 {
		r, err1 := strconv.ParseUint(string([]byte{clean[0], clean[0]}), 16, 8)
		g, err2 := strconv.ParseUint(string([]byte{clean[1], clean[1]}), 16, 8)
		b, err3 := strconv.ParseUint(string([]byte{clean[2], clean[2]}), 16, 8)
		if err1 != nil || err2 != nil || err3 != nil {
			return 0, 0, 0, false
		}
		return uint8(r), uint8(g), uint8(b), true
	}
	if len(clean) == 6 {
		r, err1 := strconv.ParseUint(clean[0:2], 16, 8)
		g, err2 := strconv.ParseUint(clean[2:4], 16, 8)
		b, err3 := strconv.ParseUint(clean[4:6], 16, 8)
		if err1 != nil || err2 != nil || err3 != nil {
			return 0, 0, 0, false
		}
		return uint8(r), uint8(g), uint8(b), true
	}
	return 0, 0, 0, false
}

// ContrastRatio menghitung rasio kontras WCAG antara dua warna hex.
func ContrastRatio(hex1, hex2 string) (float64, bool) {
	r1, g1, b1, ok1 := ParseHexColor(hex1)
	r2, g2, b2, ok2 := ParseHexColor(hex2)
	if !ok1 || !ok2 {
		return 0, false
	}

	l1 := RelativeLuminance(r1, g1, b1)
	l2 := RelativeLuminance(r2, g2, b2)

	lighter := math.Max(l1, l2)
	darker := math.Min(l1, l2)

	ratio := (lighter + 0.05) / (darker + 0.05)
	return math.Round(ratio*100) / 100, true
}

// ResolveTailwindHex mengembalikan representasi sRGB hex dari nama token Tailwind jika tersedia di kamus statis.
func ResolveTailwindHex(token string) (string, bool) {
	// Buang awalan utilitas: bg-, text-, border-
	cleaned := token
	for _, prefix := range []string{"bg-", "text-", "border-"} {
		cleaned = strings.TrimPrefix(cleaned, prefix)
	}

	// Buang modifier opasitas: amber-500/80 -> amber-500
	if idx := strings.IndexByte(cleaned, '/'); idx != -1 {
		cleaned = cleaned[:idx]
	}

	hex, found := TailwindPaletteSrgb[cleaned]
	return hex, found
}

// CheckContrastHazard mengevaluasi apakah kombinasi kelas bg dan text menghasilkan pelanggaran kontras WCAG AA (< 4.5:1).
// Mengembalikan (ratio, failsAA, true) jika kedua warna dikenali secara statis.
// Mengabaikan token dinamis (var(--...), semantic token) dengan mengembalikan ok = false.
func CheckContrastHazard(bgClass, textClass string) (ratio float64, failsAA bool, ok bool) {
	bgHex, bgOk := ResolveTailwindHex(bgClass)
	textHex, textOk := ResolveTailwindHex(textClass)

	if !bgOk || !textOk {
		return 0, false, false
	}

	r, computed := ContrastRatio(bgHex, textHex)
	if !computed {
		return 0, false, false
	}

	return r, r < 4.5, true
}
