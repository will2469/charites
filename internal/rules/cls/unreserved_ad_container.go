package cls

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// UnreservedAdContainerRule mendeteksi kontainer slot iklan dinamis yang tidak memiliki
// reservasi dimensi vertikal (min-height atau aspect-ratio) atau kerangka skeleton awal.
type UnreservedAdContainerRule struct{}

// NewUnreservedAdContainerRule membuat instance baru UnreservedAdContainerRule.
func NewUnreservedAdContainerRule() *UnreservedAdContainerRule {
	return &UnreservedAdContainerRule{}
}

// ID mengembalikan identifier kanonikal Semgrep rule.
func (r *UnreservedAdContainerRule) ID() string {
	return "cls.unreserved-ad-container"
}

// Description mengembalikan ringkasan aturan.
func (r *UnreservedAdContainerRule) Description() string {
	return "Warns when dynamic ad containers lack reserved vertical dimensions or initial skeleton placeholders"
}

// Category mengembalikan nama kategori rule.
func (r *UnreservedAdContainerRule) Category() string {
	return "cls"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *UnreservedAdContainerRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *UnreservedAdContainerRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Cumulative Layout Shift (CLS) Metric Specification",
			"Google Publisher Tag / AdSense CLS Best Practices",
			"Interactive Advertising Bureau (IAB) Standard Ad Unit Specifications",
		},
		CoreInvariant: "Dynamic ad slot containers must define a reserved bounding box (using 'min-h-*', 'h-*', or 'aspect-*') or contain an initial placeholder skeleton before third-party ad scripts inject payloads.",
		Grounding: "Ad tags and third-party advertising SDKs (such as Google AdSense, Google Publisher Tag, or Carbon Ads) execute client-side bidding and late script downloads.\n\n" +
			"When ad containers start with an empty 0px height in the normal document flow, the loaded advertisement abruptly expands the container, shoving the main article or page content downward. This sudden shift is one of the leading contributors to poor Core Web Vitals.\n\n" +
			"Declaring a minimum height corresponding to standard IAB ad dimensions (e.g. 'min-h-[90px]' for leaderboard banners or 'min-h-[250px]' for medium rectangles) reserves the necessary vertical space in advance.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Severe Downward Content Shift",
				Severity: "HIGH",
				Impact:   "Late ad injections jolt the reader's viewport position, frustrating users and ruining reading continuity.",
			},
			{
				Vector:   "Core Web Vitals Penalty",
				Severity: "HIGH",
				Impact:   "Ad insertion shifts contribute heavily to high session CLS scores in Google Search Console / CrUX.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Ad container without reserved height or skeleton placeholder",
				Code:     `<div id="ad-leaderboard" data-ad-slot="12345" className="w-full text-center" />`,
			},
			{
				Language: "astro",
				Comment:  "AdBanner component without dimension constraints",
				Code:     `<AdBanner slotId="banner-top" class="my-4" />`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Ad container with reserved IAB leaderboard min-height",
				Code:     `<div id="ad-leaderboard" data-ad-slot="12345" className="w-full min-h-[90px] md:min-h-[250px] bg-muted/20" />`,
			},
			{
				Language: "tsx",
				Comment:  "Ad slot containing an initial skeleton placeholder",
				Code: `<div id="ad-sidebar" data-ad-slot="67890" className="w-full">
  <Skeleton className="w-full h-[250px]" />
</div>`,
			},
		},
	}
}

// Evaluate memeriksa apakah kontainer iklan dinamis memiliki reservasi dimensi minimum.
func (r *UnreservedAdContainerRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if !isAdContainer(node) {
		return nil
	}

	// Jangan beri alarm palsu pada JSX spread attributes ({...props})
	if hasSpreadProps(node.Attributes) {
		return nil
	}

	// 1. Cek pembatas tinggi vertikal langsung pada kontainer (min-h-*, h-*, aspect-*, inline style)
	if hasBoundedHeight(node) {
		return nil
	}

	// 2. Cek apakah kontainer sudah menyediakan skeleton atau placeholder visual
	if hasSkeletonOrFallback(node) {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("Dynamic ad container <%s> lacks reserved layout dimensions ('min-h-*', 'h-*', or 'aspect-*'). Ad payloads loaded asynchronously will cause sudden downward Cumulative Layout Shift (CLS).", node.Tag),
			Hint:     "Reserve vertical space for the expected ad dimensions (e.g. 'min-h-[90px]' for leaderboards or 'min-h-[250px]' for rects) or provide an initial skeleton placeholder.",
		},
	}
}
