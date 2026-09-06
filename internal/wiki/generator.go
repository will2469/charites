package wiki

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"text/template"

	"github.com/will2469/charites/internal/ir"
	"github.com/will2469/charites/internal/rules"
)

//go:embed templates/*.tmpl
var templateFS embed.FS

// Generator mengelola perakitan dan pembuatan dokumentasi wiki markdown berbasis template dan rules.Registry.
type Generator struct {
	reg          *rules.Registry
	tmplHome     *template.Template
	tmplCategory *template.Template
	tmplRule     *template.Template
	tmplSidebar  *template.Template
}

// NewGenerator membuat instance Generator baru berbasis rules.Registry dan memuat embedded templates.
func NewGenerator(reg *rules.Registry) *Generator {
	if reg == nil {
		reg = rules.DefaultRegistry()
	}

	funcMap := template.FuncMap{
		"ToUpper": strings.ToUpper,
	}

	tmplHome := template.Must(template.New("home.md.tmpl").Funcs(funcMap).ParseFS(templateFS, "templates/home.md.tmpl"))
	tmplCategory := template.Must(template.New("category.md.tmpl").Funcs(funcMap).ParseFS(templateFS, "templates/category.md.tmpl"))
	tmplRule := template.Must(template.New("rule.md.tmpl").Funcs(funcMap).ParseFS(templateFS, "templates/rule.md.tmpl"))
	tmplSidebar := template.Must(template.New("sidebar.md.tmpl").Funcs(funcMap).ParseFS(templateFS, "templates/sidebar.md.tmpl"))

	return &Generator{
		reg:          reg,
		tmplHome:     tmplHome,
		tmplCategory: tmplCategory,
		tmplRule:     tmplRule,
		tmplSidebar:  tmplSidebar,
	}
}

type homeCategoryEntry struct {
	Name  string
	Count int
}

type homeRuleEntry struct {
	ID          string
	Category    string
	Slug        string
	Severity    string
	Description string
}

type homeTemplateData struct {
	Categories []homeCategoryEntry
	Rules      []homeRuleEntry
}

type categoryTemplateData struct {
	Title    string
	Category string
	Rules    []homeRuleEntry
}

type sidebarCategoryEntry struct {
	Name  string
	Title string
	Count int
	Rules []homeRuleEntry
}

type sidebarTemplateData struct {
	Categories []sidebarCategoryEntry
}

type ruleTemplateData struct {
	ID            string
	Severity      string
	RawSeverity   string
	Category      string
	Description   string
	Standards     string
	CoreInvariant string
	Grounding     string
	BadExamples   []ir.CodeExample
	GoodExamples  []ir.CodeExample
	Risks         []ir.RiskItem
}

// Generate menghasilkan berkas Home.md, <category>.md, dan <category>/<slug>.md ke direktori target.
func (g *Generator) Generate(outputDir string) error {
	tmpDir, err := os.MkdirTemp("", ".charites-wiki-staging-*")
	if err != nil {
		return fmt.Errorf("failed to create staging directory: %w", err)
	}
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	allRules := g.reg.All()

	// 1. Kelompokkan kategori unik terurut leksikografis
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
	homeContent, err := g.renderHome(categories, catMap, allRules)
	if err != nil {
		return fmt.Errorf("failed to render Home.md: %w", err)
	}
	homePath := filepath.Join(tmpDir, "Home.md")
	if wErr := os.WriteFile(homePath, []byte(homeContent), 0o600); wErr != nil {
		return fmt.Errorf("failed to write Home.md: %w", wErr)
	}

	// 2b. Render _Sidebar.md for GitHub Wiki hierarchical navigation
	sidebarContent, sErr := g.renderSidebar(categories, catMap)
	if sErr != nil {
		return fmt.Errorf("failed to render _Sidebar.md: %w", sErr)
	}
	sidebarPath := filepath.Join(tmpDir, "_Sidebar.md")
	if wErr := os.WriteFile(sidebarPath, []byte(sidebarContent), 0o600); wErr != nil {
		return fmt.Errorf("failed to write _Sidebar.md: %w", wErr)
	}

	// 3. Render <category>.md dan <category>/<slug>.md untuk setiap domain

	for _, cat := range categories {
		if err := g.renderCategoryDocs(tmpDir, cat, catMap[cat]); err != nil {
			return err
		}
	}

	// 4. Pastikan direktori target ada dan bersihkan subdirektori lawas
	if mErr := os.MkdirAll(outputDir, 0o750); mErr != nil {
		return fmt.Errorf("failed to create output directory %s: %w", outputDir, mErr)
	}
	for _, cat := range categories {
		_ = os.RemoveAll(filepath.Join(outputDir, cat))
	}

	// 5. Salin struktur dari staging ke target direktori secara rekursif
	return copyTree(tmpDir, outputDir)
}

func (g *Generator) renderRuleDoc(tmpDir, cat string, r rules.Rule) error {
	slug := strings.TrimPrefix(r.ID(), cat+".")
	ruleContent, rErr := g.renderRule(r, cat, slug)
	if rErr != nil {
		return fmt.Errorf("failed to render rule %s: %w", r.ID(), rErr)
	}
	rulePath := filepath.Join(tmpDir, r.ID()+".md")
	if wErr := os.WriteFile(rulePath, []byte(ruleContent), 0o600); wErr != nil {
		return fmt.Errorf("failed to write rule doc %s: %w", rulePath, wErr)
	}
	return nil
}

func (g *Generator) renderCategoryDocs(tmpDir, cat string, cRules []rules.Rule) error {
	catContent, cErr := g.renderCategory(cat, cRules)
	if cErr != nil {
		return fmt.Errorf("failed to render %s.md: %w", cat, cErr)
	}
	catPath := filepath.Join(tmpDir, cat+".md")
	if wErr := os.WriteFile(catPath, []byte(catContent), 0o600); wErr != nil {
		return fmt.Errorf("failed to write %s.md: %w", cat, wErr)
	}

	for _, r := range cRules {
		if err := g.renderRuleDoc(tmpDir, cat, r); err != nil {
			return err
		}
	}
	return nil
}

func copyTree(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, relErr := filepath.Rel(src, path)
		if relErr != nil {
			return relErr
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o750)
		}
		// #nosec G122 -- controlled staging path
		data, readErr := os.ReadFile(filepath.Clean(path)) //nolint:gosec
		if readErr != nil {
			return readErr
		}
		// #nosec G703 -- controlled destination path
		return os.WriteFile(filepath.Clean(target), data, 0o600) //nolint:gosec
	})
}

func (g *Generator) renderHome(categories []string, catMap map[string][]rules.Rule, allRules []rules.Rule) (string, error) {
	data := homeTemplateData{
		Categories: make([]homeCategoryEntry, 0, len(categories)),
		Rules:      make([]homeRuleEntry, 0, len(allRules)),
	}

	for _, cat := range categories {
		data.Categories = append(data.Categories, homeCategoryEntry{
			Name:  cat,
			Count: len(catMap[cat]),
		})
	}

	for _, r := range allRules {
		slug := strings.TrimPrefix(r.ID(), r.Category()+".")
		data.Rules = append(data.Rules, homeRuleEntry{
			ID:          r.ID(),
			Category:    r.Category(),
			Slug:        slug,
			Severity:    strings.ToUpper(string(r.DefaultSeverity())),
			Description: r.Description(),
		})
	}

	var buf bytes.Buffer
	if err := g.tmplHome.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (g *Generator) renderCategory(category string, categoryRules []rules.Rule) (string, error) {
	titleCat := strings.ToUpper(category[:1]) + category[1:]
	data := categoryTemplateData{
		Title:    titleCat,
		Category: category,
		Rules:    make([]homeRuleEntry, 0, len(categoryRules)),
	}

	for _, r := range categoryRules {
		slug := strings.TrimPrefix(r.ID(), category+".")
		data.Rules = append(data.Rules, homeRuleEntry{
			ID:          r.ID(),
			Category:    category,
			Slug:        slug,
			Severity:    strings.ToUpper(string(r.DefaultSeverity())),
			Description: r.Description(),
		})
	}

	var buf bytes.Buffer
	if err := g.tmplCategory.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (g *Generator) renderSidebar(categories []string, catMap map[string][]rules.Rule) (string, error) {
	data := sidebarTemplateData{
		Categories: make([]sidebarCategoryEntry, 0, len(categories)),
	}

	for _, cat := range categories {
		cRules := catMap[cat]
		titleCat := strings.ToUpper(cat[:1]) + cat[1:]
		entry := sidebarCategoryEntry{
			Name:  cat,
			Title: titleCat,
			Count: len(cRules),
			Rules: make([]homeRuleEntry, 0, len(cRules)),
		}
		for _, r := range cRules {
			slug := strings.TrimPrefix(r.ID(), cat+".")
			entry.Rules = append(entry.Rules, homeRuleEntry{
				ID:          r.ID(),
				Category:    cat,
				Slug:        slug,
				Severity:    strings.ToUpper(string(r.DefaultSeverity())),
				Description: r.Description(),
			})
		}
		data.Categories = append(data.Categories, entry)
	}

	var buf bytes.Buffer
	if err := g.tmplSidebar.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (g *Generator) renderRule(r rules.Rule, category, slug string) (string, error) {
	_ = slug
	data := ruleTemplateData{
		ID:          r.ID(),
		Severity:    strings.ToUpper(string(r.DefaultSeverity())),
		RawSeverity: string(r.DefaultSeverity()),
		Category:    category,
		Description: r.Description(),
	}

	if docRule, ok := r.(rules.DocumentedRule); ok {
		doc := docRule.Doc()
		if len(doc.TargetStandards) > 0 {
			data.Standards = strings.Join(doc.TargetStandards, ", ")
		}
		data.CoreInvariant = doc.CoreInvariant
		data.Grounding = doc.Grounding
		data.BadExamples = doc.BadExamples
		data.GoodExamples = doc.GoodExamples
		data.Risks = doc.Risks
	}

	var buf bytes.Buffer
	if err := g.tmplRule.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// RenderRuleDoc me-render dokumentasi 8-Pillars lengkap untuk sebuah rule ke dalam format Markdown.
func RenderRuleDoc(r rules.Rule) (string, error) {
	if r == nil {
		return "", fmt.Errorf("cannot render nil rule")
	}
	gen := NewGenerator(nil)
	slug := strings.TrimPrefix(r.ID(), r.Category()+".")
	return gen.renderRule(r, r.Category(), slug)
}
