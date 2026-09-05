package wiki

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/will2469/charites/internal/rules"
)

// Generator mengelola perakitan dan pembuatan dokumentasi wiki markdown berbasis rules.Registry.
type Generator struct {
	reg *rules.Registry
}

// NewGenerator membuat instance Generator baru. Jika reg bernilai nil, menggunakan rules.DefaultRegistry().
func NewGenerator(reg *rules.Registry) *Generator {
	if reg == nil {
		reg = rules.DefaultRegistry()
	}
	return &Generator{reg: reg}
}

// Generate menghasilkan berkas Home.md dan <category>.md ke direktori target secara atomik menggunakan staging.
func (g *Generator) Generate(outputDir string) error {
	tmpDir, err := os.MkdirTemp("", "charites-wiki-staging-*")
	if err != nil {
		return fmt.Errorf("failed to create staging directory: %w", err)
	}
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	allRules := g.reg.All()

	// 1. Kumpulkan kategori unik terurut leksikografis
	catMap := make(map[string][]rules.Rule)
	for _, r := range allRules {
		cat := r.Category()
		catMap[cat] = append(catMap[cat], r)
	}

	categories := make([]string, 0, len(catMap))
	for cat := range catMap {
		categories = append(categories, cat)
	}
	slices.Sort(categories)

	// 2. Render Home.md
	homeContent := g.renderHome(categories, catMap, allRules)
	homePath := filepath.Join(tmpDir, "Home.md")
	if wErr := os.WriteFile(homePath, []byte(homeContent), 0o600); wErr != nil {
		return fmt.Errorf("failed to write Home.md: %w", wErr)
	}

	// 3. Render <category>.md untuk setiap domain kategori
	for _, cat := range categories {
		cRules := catMap[cat]
		catContent := g.renderCategory(cat, cRules)
		catPath := filepath.Join(tmpDir, cat+".md")
		if wErr := os.WriteFile(catPath, []byte(catContent), 0o600); wErr != nil {
			return fmt.Errorf("failed to write %s.md: %w", cat, wErr)
		}
	}

	// 4. Pastikan direktori target ada
	if mErr := os.MkdirAll(outputDir, 0o750); mErr != nil {
		return fmt.Errorf("failed to create output directory %s: %w", outputDir, mErr)
	}

	// 5. Salin berkas terverifikasi dari staging ke target direktori secara deterministik
	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		return fmt.Errorf("failed to read staging directory: %w", err)
	}

	for _, entry := range entries {
		srcFile := filepath.Join(tmpDir, entry.Name())
		destFile := filepath.Join(outputDir, entry.Name())

		content, err := os.ReadFile(filepath.Clean(srcFile))
		if err != nil {
			return fmt.Errorf("failed to read staged file %s: %w", entry.Name(), err)
		}

		if err := os.WriteFile(destFile, content, 0o600); err != nil {
			return fmt.Errorf("failed to write destination file %s: %w", destFile, err)
		}
	}

	return nil
}

func (g *Generator) renderHome(categories []string, catMap map[string][]rules.Rule, allRules []rules.Rule) string {
	var b strings.Builder

	b.WriteString("# Charites Static Analysis Rule Catalog\n\n")
	b.WriteString("Welcome to the **Charites Static Analysis Rule Catalog**. Charites is an ultra-fast, zero-CGO, zero-Node.js static analysis compiler for Astro, React TSX, and Tailwind CSS design tokens.\n\n")
	b.WriteString("---\n\n")
	b.WriteString("## Categories\n\n")
	b.WriteString("| Category | Rules Count | Documentation |\n")
	b.WriteString("| :--- | :---: | :--- |\n")

	for _, cat := range categories {
		cRules := catMap[cat]
		fmt.Fprintf(&b, "| `%s` | %d | [`%s.md`](%s.md) |\n", cat, len(cRules), cat, cat)
	}

	b.WriteString("\n---\n\n")
	b.WriteString("## All Registered Rules\n\n")
	b.WriteString("| Rule ID | Category | Severity | Description | Documentation |\n")
	b.WriteString("| :--- | :---: | :---: | :--- | :--- |\n")

	for _, r := range allRules {
		anchor := strings.ReplaceAll(strings.ToLower(r.ID()), ".", "")
		fmt.Fprintf(&b, "| `%s` | `%s` | `%s` | %s | [`%s.md#%s`](%s.md#%s) |\n",
			r.ID(), r.Category(), strings.ToUpper(string(r.DefaultSeverity())), r.Description(), r.Category(), anchor, r.Category(), anchor)
	}

	b.WriteString("\n---\n\n")
	b.WriteString("## Architectural Principles\n\n")
	b.WriteString("1. **Deterministic Execution:** Pure-function AST visitors without file system or network I/O during evaluation.\n")
	b.WriteString("2. **1-SSOT Tri-Corpus Assurance:** Every rule is validated against a 3-part golden test corpus (`positive/`, `negative/`, `adversarial/`).\n")
	b.WriteString("3. **Canonical Semgrep Identifiers:** All rules follow the `<category>.<slug>` standard.\n")

	return b.String()
}

func (g *Generator) renderCategory(category string, categoryRules []rules.Rule) string {
	var b strings.Builder

	titleCat := strings.ToUpper(category[:1]) + category[1:]
	fmt.Fprintf(&b, "# %s Rules (`%s`)\n\n", titleCat, category)
	fmt.Fprintf(&b, "The `%s` category contains static analysis rules for code quality, architectural constraints, and design system governance.\n\n", category)
	b.WriteString("---\n\n")
	b.WriteString("## Category Rule Index\n\n")
	b.WriteString("| Rule ID | Severity | Summary | Status |\n")
	b.WriteString("| :--- | :---: | :--- | :---: |\n")

	for _, r := range categoryRules {
		anchor := strings.ReplaceAll(strings.ToLower(r.ID()), ".", "")
		fmt.Fprintf(&b, "| [`%s`](#%s) | `%s` | %s | `enabled` |\n",
			r.ID(), anchor, strings.ToUpper(string(r.DefaultSeverity())), r.Description())
	}

	b.WriteString("\n---\n\n")

	for _, r := range categoryRules {
		fmt.Fprintf(&b, "## `%s`\n\n", r.ID())
		fmt.Fprintf(&b, "> **Rule ID:** `%s`  \n", r.ID())
		fmt.Fprintf(&b, "> **Severity:** `%s`  \n", strings.ToUpper(string(r.DefaultSeverity())))
		fmt.Fprintf(&b, "> **Category:** `%s`  \n", r.Category())
		b.WriteString("> **Target Standards:** W3C Design Tokens Community Group (DTCG), Tailwind CSS Standards  \n\n")

		b.WriteString("### 1. Overview\n")
		fmt.Fprintf(&b, "%s.\n\n", r.Description())

		b.WriteString("### 2. How to Suppress (Ignore Directives)\n\n")
		b.WriteString("Suppress this rule via canonical directive:\n\n")
		fmt.Fprintf(&b, "```astro\n<!-- charites:ignore %s intentional exception -->\n```\n\n", r.ID())
		fmt.Fprintf(&b, "```tsx\n// charites:ignore %s intentional exception\n```\n\n", r.ID())

		b.WriteString("### 3. Configuration Reference (`charites.yaml`)\n\n")
		b.WriteString("```yaml\nrules:\n")
		fmt.Fprintf(&b, "  %s:\n    severity: %s\n```\n\n", r.ID(), r.DefaultSeverity())
		b.WriteString("---\n\n")
	}

	return b.String()
}
