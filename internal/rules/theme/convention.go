package theme

import (
	themeengine "github.com/will2469/charites/internal/token"
)

// TokenCandidate merepresentasikan kandidat token pengganti yang ditemukan dalam graph.
type TokenCandidate struct {
	Name     string // e.g. "primary-light", "--color-primary-light"
	RawValue string
}

// TokenConvention mendefinisikan antarmuka adapter inferensi semantik (Layer 4)
// untuk menemukan padanan token resmi dari modifier opacity.
type TokenConvention interface {
	FindOpacityReplacement(
		base string,
		opacity string,
		ctx themeengine.Context,
	) ([]TokenCandidate, bool)
}

// ConfigurableConvention mengimplementasikan TokenConvention yang dapat dikonfigurasi penuh
// (design-agnostic) melalui charites.yaml atau adapter proyek khusus.
type ConfigurableConvention struct {
	opacityMappings map[string][]string
	fallbacks       map[string][]string
	prefixes        []string
}

// ConventionOption mendefinisikan fungsi opsi konfigurasi untuk ConfigurableConvention.
type ConventionOption func(*ConfigurableConvention)

// WithOpacityMapping mendaftarkan pemetaan nilai modifier opacity ke daftar suffix kandidat token.
func WithOpacityMapping(opacity string, suffixes ...string) ConventionOption {
	return func(c *ConfigurableConvention) {
		c.opacityMappings[opacity] = append(c.opacityMappings[opacity], suffixes...)
	}
}

// WithFallback mendaftarkan alias dasar fallback untuk nama warna dasar tertentu.
func WithFallback(base string, fallbackBases ...string) ConventionOption {
	return func(c *ConfigurableConvention) {
		c.fallbacks[base] = append(c.fallbacks[base], fallbackBases...)
	}
}

// WithPrefixes menetapkan daftar prefix custom property CSS yang dicari di token graph.
func WithPrefixes(prefixes ...string) ConventionOption {
	return func(c *ConfigurableConvention) {
		c.prefixes = prefixes
	}
}

// NewConfigurableConvention membuat instance baru ConfigurableConvention dengan opsi kustom.
func NewConfigurableConvention(opts ...ConventionOption) *ConfigurableConvention {
	c := &ConfigurableConvention{
		opacityMappings: make(map[string][]string),
		fallbacks:       make(map[string][]string),
		prefixes:        []string{"--color-", "--"},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// NewDefaultCharitesConvention membuat instance konvensi bawaan Charites & DTCG.
// Catatan: Ini adalah konfigurasi awal ramah-pengembang dan BUKAN kebenaran semantik absolut;
// dapat dioverride sepenuhnya via konfigurasi proyek charites.yaml.
func NewDefaultCharitesConvention() *ConfigurableConvention {
	return NewConfigurableConvention(
		WithOpacityMapping("10", "-light"),
		WithOpacityMapping("20", "-light"),
		WithOpacityMapping("5", "-subtle"),
		WithOpacityMapping("8", "-subtle"),
		WithFallback("secondary", "muted"),
		WithPrefixes("--color-", "--"),
	)
}

// NewConventionFromConfig membuat instance TokenConvention berdasarkan konfigurasi charites.yaml.
func NewConventionFromConfig(cfg themeengine.ConventionConfig) *ConfigurableConvention {
	c := NewDefaultCharitesConvention()

	if len(cfg.OpacityMappings) > 0 {
		c.opacityMappings = make(map[string][]string, len(cfg.OpacityMappings))
		for op, suffixes := range cfg.OpacityMappings {
			c.opacityMappings[op] = suffixes
		}
	}

	if len(cfg.Fallbacks) > 0 {
		c.fallbacks = make(map[string][]string, len(cfg.Fallbacks))
		for b, fbs := range cfg.Fallbacks {
			c.fallbacks[b] = fbs
		}
	}

	if len(cfg.Prefixes) > 0 {
		c.prefixes = cfg.Prefixes
	}

	return c
}

// AddOpacityMapping menambahkan pemetaan opacity baru secara dinamis.
func (c *ConfigurableConvention) AddOpacityMapping(opacity string, suffixes ...string) {
	c.opacityMappings[opacity] = append(c.opacityMappings[opacity], suffixes...)
}

// AddFallback menambahkan relasi fallback baru secara dinamis.
func (c *ConfigurableConvention) AddFallback(base string, fallbackBases ...string) {
	c.fallbacks[base] = append(c.fallbacks[base], fallbackBases...)
}

// FindOpacityReplacement mencari token pengganti semantik resmi untuk base color dan modifier opacity
// berdasarkan aturan pemetaan yang dikonfigurasi.
func (c *ConfigurableConvention) FindOpacityReplacement(
	base string,
	opacity string,
	ctx themeengine.Context,
) ([]TokenCandidate, bool) {
	if ctx == nil {
		return nil, false
	}

	suffixes, exists := c.opacityMappings[opacity]
	if !exists || len(suffixes) == 0 {
		return nil, false
	}

	for _, suffix := range suffixes {
		// 1. Cek langsung <base><suffix>, misal: "primary-light"
		candidateName := base + suffix
		if cand, ok := c.lookupCandidate(candidateName, ctx); ok {
			return []TokenCandidate{cand}, true
		}

		// 2. Cek basis fallback yang terkonfigurasi (misal: secondary -> muted, banana -> kuning)
		if fbs, ok := c.fallbacks[base]; ok {
			for _, fb := range fbs {
				fbCand := fb + suffix
				if cand, ok := c.lookupCandidate(fbCand, ctx); ok {
					return []TokenCandidate{cand}, true
				}
			}
		}
	}

	return nil, false
}

func (c *ConfigurableConvention) lookupCandidate(shortName string, ctx themeengine.Context) (TokenCandidate, bool) {
	for _, p := range c.prefixes {
		fullName := p + shortName
		if tok, ok := ctx.LookupToken(fullName); ok {
			return TokenCandidate{
				Name:     shortName,
				RawValue: tok.RawValue,
			}, true
		}
	}

	// Cek tanpa prefix (jika token dideklarasikan langsung sebagai shortName)
	if tok, ok := ctx.LookupToken(shortName); ok {
		return TokenCandidate{
			Name:     shortName,
			RawValue: tok.RawValue,
		}, true
	}

	return TokenCandidate{}, false
}
