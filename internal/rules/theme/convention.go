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

// DefaultCharitesConvention mengimplementasikan konvensi bawaan Charites & DTCG:
//   - Opacity 10 & 20 -> mencari varian "-light" (misal: "primary-light", "--color-primary-light")
//   - Opacity 5 & 8   -> mencari varian "-subtle" (misal: "primary-subtle", "--color-primary-subtle")
//   - Fallback secondary: jika base "secondary" dan "secondary-light" tidak ada di graph,
//     tetapi "muted-light" ada di graph, merekomendasikan "muted-light".
//
// Catatan: Adapter ini HANYA mengembalikan kandidat JIKA token tersebut benar-benar ADA
// di dalam TokenGraph proyek pengguna (fakta membuktikan eksistensi token).
type DefaultCharitesConvention struct{}

// NewDefaultCharitesConvention membuat instance konvensi bawaan Charites.
func NewDefaultCharitesConvention() *DefaultCharitesConvention {
	return &DefaultCharitesConvention{}
}

// FindOpacityReplacement mencari token pengganti semantik resmi untuk base color dan modifier opacity.
func (c *DefaultCharitesConvention) FindOpacityReplacement(
	base string,
	opacity string,
	ctx themeengine.Context,
) ([]TokenCandidate, bool) {
	if ctx == nil {
		return nil, false
	}

	var suffixes []string
	switch opacity {
	case "10", "20":
		suffixes = []string{"-light"}
	case "5", "8":
		suffixes = []string{"-subtle"}
	default:
		return nil, false
	}

	for _, suffix := range suffixes {
		// 1. Cek langsung <base><suffix>, misal: "primary-light"
		candidateName := base + suffix
		if cand, ok := c.lookupCandidate(candidateName, ctx); ok {
			return []TokenCandidate{cand}, true
		}

		// 2. Fallback secondary -> muted
		if base == "secondary" {
			mutedCand := "muted" + suffix
			if cand, ok := c.lookupCandidate(mutedCand, ctx); ok {
				return []TokenCandidate{cand}, true
			}
		}
	}

	return nil, false
}

func (c *DefaultCharitesConvention) lookupCandidate(shortName string, ctx themeengine.Context) (TokenCandidate, bool) {
	// Cek bentuk custom property: "--color-<name>", "--<name>"
	variants := []string{
		"--color-" + shortName,
		"--" + shortName,
	}

	for _, v := range variants {
		tokens := ctx.ByName(v)
		if len(tokens) > 0 {
			// Token eksis di graph SSOT pengguna!
			return TokenCandidate{
				Name:     shortName,
				RawValue: tokens[0].RawValue,
			}, true
		}
	}

	return TokenCandidate{}, false
}
