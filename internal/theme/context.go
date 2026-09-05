package theme

import (
	"sort"
	"strings"
)

// Context merepresentasikan Single Source of Truth (SSOT) dari token tema proyek.
// Diekstrak langsung dari berkas CSS (seperti global.css) secara deterministik dan thread-safe.
// Menyediakan metode query berkecepatan tinggi dengan alokasi memori nol (zero-alloc)
// pada jalur evaluasi aturan statis (hot-path).
type Context struct {
	// Path berkas CSS sumber (jika dibaca dari filesystem).
	Path string

	// Flag struktural arsitektur CSS
	HasLayerTheme      bool
	HasRootColorScheme bool
	RootColorScheme    string // "light", "dark", "light dark"
	HasDarkColorScheme bool
	DarkColorScheme    string

	// Raw Custom Properties (Key: "--color-primary", Val: "oklch(0.55 0.22 260)")
	RawVariables map[string]string

	// Tokens berdasarkan scope deklarasi
	LightTokens map[string]string // Custom properties di dalam :root / html
	DarkTokens  map[string]string // Custom properties di dalam .dark / [data-theme="dark"]

	// Normalized Base Colors: "primary", "secondary", "destructive", "accent", dsb.
	// Memetakan nama dasar ke representasi CSS aslinya.
	BaseColors map[string]string

	// Calibrated Opacity Variants:
	// Memetakan pasangan "color/opacity" (misal "primary/10", "primary/20", "primary/5")
	// ke token semantik resmi pengganti (misal "primary-light", "primary-subtle").
	OpacityVariants map[string]string

	// Calibrated Tokens Set: "primary-light", "primary-subtle", dsb.
	CalibratedTokens map[string]bool

	// Scale & Shape Tokens
	RadiusScale  map[string]string // "sm" -> "0.25rem", "md" -> "0.5rem", etc.
	ZIndexScale  map[string]string // "dropdown" -> "1000", "modal" -> "1300", etc.
	ShadowScale  map[string]string // "sm", "md", "lg", "xl"
	BlurScale    map[string]string // "sm", "md", "lg"
	SpacingScale map[string]string // "1", "2", "3", "4", "6", "8", "container"
}

// NewContext menginisialisasi instance Context baru dengan seluruh map teralokasi.
func NewContext() *Context {
	return &Context{
		RawVariables:     make(map[string]string),
		LightTokens:      make(map[string]string),
		DarkTokens:       make(map[string]string),
		BaseColors:       make(map[string]string),
		OpacityVariants:  make(map[string]string),
		CalibratedTokens: make(map[string]bool),
		RadiusScale:      make(map[string]string),
		ZIndexScale:      make(map[string]string),
		ShadowScale:      make(map[string]string),
		BlurScale:        make(map[string]string),
		SpacingScale:     make(map[string]string),
	}
}

// IsKnownColor memeriksa apakah nama warna tertentu terdaftar dalam token semantik tema.
// Bersifat zero-alloc pada hot-path.
func (c *Context) IsKnownColor(name string) bool {
	if c == nil || len(c.BaseColors) == 0 {
		return false
	}
	_, ok := c.BaseColors[name]
	return ok
}

// ReplacementForSlash mengembalikan token semantik pengganti untuk pasangan "color/opacity"
// (misalnya "primary/10" -> "primary-light", "primary/5" -> "primary-subtle").
// Menerima langsung slice token slash tanpa alokasi string baru.
func (c *Context) ReplacementForSlash(colorSlash string) (string, bool) {
	if c == nil || len(c.OpacityVariants) == 0 {
		return "", false
	}
	rep, ok := c.OpacityVariants[colorSlash]
	return rep, ok
}

// ReplacementForOpacity mengembalikan token pengganti untuk nama dasar dan nilai opacity terpisah.
func (c *Context) ReplacementForOpacity(baseColor, opacity string) (string, bool) {
	if c == nil || len(c.OpacityVariants) == 0 {
		return "", false
	}
	key := baseColor + "/" + opacity
	rep, ok := c.OpacityVariants[key]
	return rep, ok
}

// IsCalibratedToken memeriksa apakah nama token merupakan token varian kalibrasi
// (seperti "primary-light", "primary-subtle").
func (c *Context) IsCalibratedToken(token string) bool {
	if c == nil || len(c.CalibratedTokens) == 0 {
		return false
	}
	return c.CalibratedTokens[token]
}

// IsKnownRadius memeriksa apakah nama radius terdaftar dalam skala radius tema.
// Contoh: "sm", "md", "lg", "xl", "full".
func (c *Context) IsKnownRadius(name string) bool {
	if c == nil || len(c.RadiusScale) == 0 {
		return false
	}
	_, ok := c.RadiusScale[name]
	return ok
}

// IsKnownZIndex memeriksa apakah nama z-index terdaftar dalam skala z-index tema.
// Contoh: "dropdown", "sticky", "modal", "popover", "tooltip", "toast".
func (c *Context) IsKnownZIndex(name string) bool {
	if c == nil || len(c.ZIndexScale) == 0 {
		return false
	}
	_, ok := c.ZIndexScale[name]
	return ok
}

// IsKnownShadow memeriksa apakah nama shadow terdaftar dalam skala elevasi tema.
func (c *Context) IsKnownShadow(name string) bool {
	if c == nil || len(c.ShadowScale) == 0 {
		return false
	}
	_, ok := c.ShadowScale[name]
	return ok
}

// IsKnownBlur memeriksa apakah nama blur terdaftar dalam skala backdrop-filter tema.
func (c *Context) IsKnownBlur(name string) bool {
	if c == nil || len(c.BlurScale) == 0 {
		return false
	}
	_, ok := c.BlurScale[name]
	return ok
}

// HasColorScheme mengembalikan true jika berkas tema mendeklarasikan properti `color-scheme`.
func (c *Context) HasColorScheme() bool {
	if c == nil {
		return false
	}
	return c.HasRootColorScheme
}

// IsLayered mengembalikan true jika tema membungkus deklarasi di dalam `@layer theme`.
func (c *Context) IsLayered() bool {
	if c == nil {
		return false
	}
	return c.HasLayerTheme
}

// CheckDarkParity memeriksa apakah token warna tertentu (misal "--color-primary" atau "primary")
// didefinisikan ulang secara semantik pada dark mode scope.
func (c *Context) CheckDarkParity(token string) bool {
	if c == nil || len(c.DarkTokens) == 0 {
		return false
	}
	if !strings.HasPrefix(token, "--") {
		token = "--color-" + token
	}
	_, ok := c.DarkTokens[token]
	return ok
}

// MissingDarkParityTokens mengembalikan daftar token light mode yang tidak memiliki padanan
// inversi pada dark mode scope, diurutkan secara alfabetis deterministik.
func (c *Context) MissingDarkParityTokens() []string {
	if c == nil || len(c.LightTokens) == 0 {
		return nil
	}

	var missing []string
	for token := range c.LightTokens {
		// Hanya periksa token semantik warna (--color-*)
		if !strings.HasPrefix(token, "--color-") {
			continue
		}
		// Abaikan token turunan/light/subtle/hover jika dark mode menggunakan strategi kalkulasi mandiri
		// tetapi periksa token permukaan dan peran inti: primary, secondary, card, popover, background, foreground, border
		if _, exists := c.DarkTokens[token]; !exists {
			missing = append(missing, token)
		}
	}

	sort.Strings(missing)
	return missing
}
