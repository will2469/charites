package theme

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Standard CSS search paths relative to project root for automatic discovery.
var standardCSSPaths = []string{
	"src/style/global.css",
	"src/styles/global.css",
	"styles/global.css",
	"src/global.css",
	"src/index.css",
	"global.css",
	"index.css",
	"tests/fixtures/global.css",
}

// DiscoverAndLoad mencari berkas CSS tema berdasarkan discovery fallback engine Charites:
//  1. Jika customPath tidak kosong, periksa customPath.
//  2. Jika customPath kosong, cari di daftar standardCSSPaths dengan melakukan upward walk
//     mulai dari projectRoot (atau direktori kerja saat ini) hingga akar filesystem.
//  3. Jika berkas ditemukan, parse dan kembalikan Context terstruktur.
//  4. Jika tidak ditemukan berkas CSS manapun, kembalikan Context kosong tanpa error
//     (mematuhi invarian Zero-Config Default: YES).
func DiscoverAndLoad(projectRoot, customPath string) (*Context, error) {
	var targetPath string

	if customPath != "" {
		resolved := customPath
		if !filepath.IsAbs(resolved) && projectRoot != "" {
			resolved = filepath.Join(projectRoot, resolved)
		}
		if _, err := os.Stat(resolved); err == nil {
			targetPath = resolved
		} else {
			return nil, fmt.Errorf("specified theme css file not found: %s: %w", customPath, err)
		}
	} else {
		startDir := projectRoot
		if startDir == "" {
			if cwd, err := os.Getwd(); err == nil {
				startDir = cwd
			}
		}
		if abs, err := filepath.Abs(startDir); err == nil {
			startDir = abs
		}

		curr := startDir
		for {
			for _, rel := range standardCSSPaths {
				cand := filepath.Join(curr, rel)
				if _, err := os.Stat(cand); err == nil {
					targetPath = cand
					break
				}
			}
			if targetPath != "" {
				break
			}
			parent := filepath.Dir(curr)
			if parent == curr {
				break
			}
			curr = parent
		}
	}

	if targetPath == "" {
		// Zero-config operability: kembalikan context kosong yang aman
		return NewContext(), nil
	}

	data, err := os.ReadFile(filepath.Clean(targetPath))
	if err != nil {
		return nil, fmt.Errorf("failed to read theme css file %s: %w", targetPath, err)
	}

	ctx, err := ParseCSS(data)
	if err != nil {
		return nil, err
	}
	ctx.Path = targetPath
	return ctx, nil
}

// ParseCSS mem-parse buffer CSS menjadi Context terstruktur.
func ParseCSS(src []byte) (*Context, error) {
	ctx := NewContext()
	cleaned := stripCSSComments(src)

	scope := scopeContext{}
	parseBlocksAndDecls(cleaned, scope, ctx)

	return ctx, nil
}

type scopeContext struct {
	inLayer      bool
	isDark       bool
	isRoot       bool
	isThemeBlock bool
}

// stripCSSComments menghapus seluruh komentar /* ... */ dari buffer CSS.
func stripCSSComments(src []byte) []byte {
	var buf bytes.Buffer
	buf.Grow(len(src))

	i := 0
	for i < len(src) {
		if i+1 < len(src) && src[i] == '/' && src[i+1] == '*' {
			end := bytes.Index(src[i+2:], []byte("*/"))
			if end == -1 {
				break
			}
			i += end + 4
			continue
		}
		buf.WriteByte(src[i])
		i++
	}

	return buf.Bytes()
}

// parseBlocksAndDecls memindai blok CSS dan deklarasi secara rekursif.
func parseBlocksAndDecls(data []byte, scope scopeContext, ctx *Context) {
	i := 0
	declStart := 0

	for i < len(data) {
		if data[i] == '{' {
			// Header selector / at-rule adalah data dari declStart hingga i
			rawHeader := string(data[declStart:i])
			header := strings.TrimSpace(rawHeader)

			// Cari kurung kurawal penutup '}' penyeimbang
			braceDepth := 1
			blockStart := i + 1
			blockEnd := -1
			for j := blockStart; j < len(data); j++ {
				if data[j] == '{' {
					braceDepth++
				} else if data[j] == '}' {
					braceDepth--
					if braceDepth == 0 {
						blockEnd = j
						break
					}
				}
			}

			if blockEnd == -1 {
				// Blok tidak ditutup sempurna, parsing hingga akhir
				blockEnd = len(data)
			}

			inner := data[blockStart:blockEnd]
			childScope := resolveChildScope(header, scope, ctx)

			// Jika blok memiliki nested child '{', rekursi; jika tidak, parse deklarasi
			if bytes.IndexByte(inner, '{') != -1 {
				parseBlocksAndDecls(inner, childScope, ctx)
			} else {
				parseDeclarations(inner, childScope, ctx)
			}

			if blockEnd < len(data) {
				i = blockEnd + 1
			} else {
				i = len(data)
			}
			declStart = i
			continue
		}
		i++
	}

	// Sisa deklarasi di luar kurung kurawal jika ada
	if declStart < len(data) {
		rem := data[declStart:]
		if bytes.IndexByte(rem, ':') != -1 {
			parseDeclarations(rem, scope, ctx)
		}
	}
}

func resolveChildScope(header string, current scopeContext, ctx *Context) scopeContext {
	child := current
	lower := strings.ToLower(header)

	if strings.Contains(lower, "@layer") {
		child.inLayer = true
		if strings.Contains(lower, "theme") {
			ctx.HasLayerTheme = true
		}
	}

	if strings.HasPrefix(lower, "@theme") {
		child.isThemeBlock = true
	}

	if strings.Contains(lower, ".dark") ||
		strings.Contains(lower, "data-theme=\"dark\"") ||
		strings.Contains(lower, "data-theme='dark'") ||
		strings.Contains(lower, "prefers-color-scheme: dark") {
		child.isDark = true
		child.isRoot = false
	} else if strings.Contains(lower, ":root") || strings.Contains(lower, "html") {
		child.isRoot = true
		child.isDark = false
	}

	return child
}

// parseDeclarations mem-parse pernyataan CSS (prop: val;) di dalam sebuah blok.
func parseDeclarations(block []byte, scope scopeContext, ctx *Context) {
	statements := bytes.Split(block, []byte(";"))
	for _, rawStmt := range statements {
		stmt := strings.TrimSpace(string(rawStmt))
		if stmt == "" {
			continue
		}

		colonIdx := strings.IndexByte(stmt, ':')
		if colonIdx == -1 {
			continue
		}

		prop := strings.TrimSpace(stmt[:colonIdx])
		val := strings.TrimSpace(stmt[colonIdx+1:])

		// Tangani properti standar color-scheme
		if strings.EqualFold(prop, "color-scheme") {
			if scope.isDark {
				ctx.HasDarkColorScheme = true
				ctx.DarkColorScheme = val
			} else {
				ctx.HasRootColorScheme = true
				ctx.RootColorScheme = val
			}
			continue
		}

		// Hanya proses custom properties (--*)
		if !strings.HasPrefix(prop, "--") || val == "" {
			continue
		}

		// Simpan raw variable
		ctx.RawVariables[prop] = val

		// Simpan berdasarkan scope deklarasi
		if scope.isDark {
			ctx.DarkTokens[prop] = val
		} else {
			ctx.LightTokens[prop] = val
		}

		// Kategorisasi token
		categorizeToken(prop, val, ctx)
	}
}

// categorizeToken menganalisis nama dan nilai custom property dan memperbarui registry token.
func categorizeToken(prop, val string, ctx *Context) {
	// 1. Color Tokens (--color-*)
	if strings.HasPrefix(prop, "--color-") {
		colorName := strings.TrimPrefix(prop, "--color-")
		if colorName == "*" {
			return
		}
		registerColorToken(colorName, val, ctx)
		return
	}

	// 2. Fallback color format (--primary, dsb.) jika nilainya adalah fungsi warna CSS atau hex
	if isColorValue(val) {
		shortName := strings.TrimPrefix(prop, "--")
		registerColorToken(shortName, val, ctx)
		return
	}

	// 3. Radius Tokens (--radius-*)
	if strings.HasPrefix(prop, "--radius-") {
		radName := strings.TrimPrefix(prop, "--radius-")
		ctx.RadiusScale[radName] = val
		return
	}

	// 4. Z-Index Tokens (--z-*)
	if strings.HasPrefix(prop, "--z-") {
		zName := strings.TrimPrefix(prop, "--z-")
		ctx.ZIndexScale[zName] = val
		return
	}

	// 5. Elevation / Shadow Tokens (--shadow-*)
	if strings.HasPrefix(prop, "--shadow-") {
		shadowName := strings.TrimPrefix(prop, "--shadow-")
		ctx.ShadowScale[shadowName] = val
		return
	}

	// 6. Backdrop Blur Tokens (--blur-*)
	if strings.HasPrefix(prop, "--blur-") {
		blurName := strings.TrimPrefix(prop, "--blur-")
		ctx.BlurScale[blurName] = val
		return
	}

	// 7. Spacing Tokens (--spacing-*)
	if strings.HasPrefix(prop, "--spacing-") {
		spacingName := strings.TrimPrefix(prop, "--spacing-")
		ctx.SpacingScale[spacingName] = val
		return
	}
}

// registerColorToken memetakan token warna dan varian opacity kalibrasinya.
func registerColorToken(name, val string, ctx *Context) {
	// Cek apakah token merupakan varian -light
	if strings.HasSuffix(name, "-light") {
		base := strings.TrimSuffix(name, "-light")
		ctx.CalibratedTokens[name] = true
		ctx.BaseColors[base] = val

		// Daftarkan ekuivalensi standar DTCG untuk -light: /10, /20
		registerOpacity(ctx, base, "10", name)
		registerOpacity(ctx, base, "20", name)

		// Jika muted-light terdaftar, sediakan fallback untuk secondary jika secondary-light belum ada
		if base == "muted" {
			if _, hasSec := ctx.OpacityVariants["secondary/10"]; !hasSec {
				registerOpacity(ctx, "secondary", "10", name)
				registerOpacity(ctx, "secondary", "20", name)
			}
		}
		return
	}

	// Cek apakah token merupakan varian -subtle
	if strings.HasSuffix(name, "-subtle") {
		base := strings.TrimSuffix(name, "-subtle")
		ctx.CalibratedTokens[name] = true
		ctx.BaseColors[base] = val

		// Daftarkan ekuivalensi standar DTCG untuk -subtle: /5, /8
		registerOpacity(ctx, base, "5", name)
		registerOpacity(ctx, base, "8", name)

		// Jika muted-subtle terdaftar, sediakan fallback untuk secondary jika secondary-subtle belum ada
		if base == "muted" {
			if _, hasSec := ctx.OpacityVariants["secondary/5"]; !hasSec {
				registerOpacity(ctx, "secondary", "5", name)
				registerOpacity(ctx, "secondary", "8", name)
			}
		}
		return
	}

	// Token warna dasar (misal primary, secondary, destructive, background)
	ctx.BaseColors[name] = val
}

func registerOpacity(ctx *Context, base, opacity, replacement string) {
	key := base + "/" + opacity
	ctx.OpacityVariants[key] = replacement
}

// isColorValue memeriksa apakah sebuah string nilai CSS merepresentasikan warna.
func isColorValue(val string) bool {
	lower := strings.ToLower(strings.TrimSpace(val))
	return strings.HasPrefix(lower, "oklch(") ||
		strings.HasPrefix(lower, "oklab(") ||
		strings.HasPrefix(lower, "rgb(") ||
		strings.HasPrefix(lower, "rgba(") ||
		strings.HasPrefix(lower, "hsl(") ||
		strings.HasPrefix(lower, "hsla(") ||
		strings.HasPrefix(lower, "#")
}
