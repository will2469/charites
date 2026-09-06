package ux

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// NavOverflowChunkingRule mendeteksi landmark navigasi yang memuat lebih dari 7 tautan datar
// tanpa mekanisme chunking (dropdown, disclosure, atau sub-menu), melanggar Hukum Miller (7 ± 2).
type NavOverflowChunkingRule struct{}

// NewNavOverflowChunkingRule membuat instance baru dari NavOverflowChunkingRule.
func NewNavOverflowChunkingRule() *NavOverflowChunkingRule {
	return &NavOverflowChunkingRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *NavOverflowChunkingRule) ID() string {
	return "ux.nav-overflow-chunking"
}

// Description mengembalikan ringkasan aturan.
func (r *NavOverflowChunkingRule) Description() string {
	return "Warns when a navigation landmark contains more than 7 direct navigation links without chunking mechanisms"
}

// Category mengembalikan nama kategori rule.
func (r *NavOverflowChunkingRule) Category() string {
	return "ux"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *NavOverflowChunkingRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *NavOverflowChunkingRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Miller's Law (Information Processing Capacity: 7 ± 2 Chunks)",
			"Information Architecture Chunking & Category Hierarchy (Rosenfeld & Morville)",
			"W3C WAI-ARIA Authoring Practices Guide 1.2 (Navigation Menubars)",
		},
		CoreInvariant: "Navigation landmarks ('<nav>' or 'role=\"navigation\"') must not present more than 7 flat direct links without grouping into disclosures, dropdown menus, or category drawers.",
		Grounding: "Miller's Law dictates that human working memory can reliably retain only 7 ± 2 distinct chunks of information at any single time.\n\n" +
			"When a main navigation bar presents 8 or more flat links in a single row or list without visual or hierarchical chunking, users experience choice paralysis and elevated cognitive scan latency.\n\n" +
			"To maintain optimal information architecture, high-density menus should group secondary destinations into nested dropdowns, accordions, or an overflow 'More...' disclosure container.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Cognitive Overload & Choice Paralysis",
				Severity: "MEDIUM",
				Impact:   "Users take significantly longer to locate key navigation targets and frequently miss secondary features.",
			},
			{
				Vector:   "Visual Clutter on Narrow Viewports",
				Severity: "MEDIUM",
				Impact:   "Flat multi-link navigation rows wrap awkwardly or cause accidental taps on mobile touchscreens.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Ten flat navigation links inside <nav> without grouping or overflow chunking",
				Code: `<nav className="flex gap-4">
  <a href="/home">Beranda</a>
  <a href="/profil">Profil</a>
  <a href="/layanan">Layanan</a>
  <a href="/berita">Berita</a>
  <a href="/transparansi">Transparansi</a>
  <a href="/anggaran">Anggaran</a>
  <a href="/regulasi">Regulasi</a>
  <a href="/galeri">Galeri</a>
  <a href="/kontak">Kontak</a>
  <a href="/bantuan">Bantuan</a>
</nav>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Primary destinations kept to 4 links, with remaining links grouped into a DropdownMenu",
				Code: `<nav className="flex gap-4 items-center">
  <a href="/home">Beranda</a>
  <a href="/profil">Profil</a>
  <a href="/layanan">Layanan</a>
  <a href="/berita">Berita</a>
  <DropdownMenu>
    <button type="button" aria-expanded="false">Lainnya</button>
  </DropdownMenu>
</nav>`,
			},
		},
	}
}

// Evaluate memeriksa apakah kontainer navigasi memuat > 7 tautan tanpa chunking.
func (r *NavOverflowChunkingRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if !isNavigationLandmark(node) {
		return nil
	}

	// Jika sudah ada mekanisme chunking, struktur navigasi sudah compliant
	if hasChunkingMechanism(node) {
		return nil
	}

	// Hitung link navigasi
	linkCount := 0
	countLinks(node, &linkCount)

	if linkCount > 7 {
		return []ir.Diagnostic{
			{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message: fmt.Sprintf(
					"Navigation landmark contains %d flat links exceeding Miller's Law threshold (7±2) without chunking mechanisms (dropdown, accordion, or overflow menu).",
					linkCount,
				),
				Hint: "Group secondary navigation items into a dropdown menu, disclosure, or overflow drawer.",
			},
		}
	}

	return nil
}

func countLinks(node *ir.Node, count *int) {
	if node == nil {
		return
	}

	for _, child := range node.Children {
		if child.Type != ir.NodeElement {
			continue
		}

		if isNavLinkNode(child) {
			*count++
		} else if child.Tag == "ul" || child.Tag == "ol" || child.Tag == "li" || child.Tag == "div" {
			// Telusuri struktur list / wrapper navigasi
			countLinks(child, count)
		}
	}
}
