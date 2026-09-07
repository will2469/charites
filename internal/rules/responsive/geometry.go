package responsive

import (
	"strconv"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// TierCount adalah jumlah breakpoint tier standar (baseline s/d 2xl).
const TierCount = 6

// BreakpointTier merepresentasikan tingkatan responsif mobile-first.
type BreakpointTier int

// Daftar konstanta breakpoint tier responsif mobile-first.
const (
	TierBaseline BreakpointTier = iota
	TierSm
	TierMd
	TierLg
	TierXl
	Tier2Xl
)

func (t BreakpointTier) String() string {
	switch t {
	case TierSm:
		return "sm"
	case TierMd:
		return "md"
	case TierLg:
		return "lg"
	case TierXl:
		return "xl"
	case Tier2Xl:
		return "2xl"
	default:
		return "baseline"
	}
}

// allTiers berisi seluruh urutan breakpoint tier dari mobile ke desktop.
var allTiers = [...]BreakpointTier{
	TierBaseline,
	TierSm,
	TierMd,
	TierLg,
	TierXl,
	Tier2Xl,
}

// DisplayMode merepresentasikan nilai CSS display pada suatu tier.
type DisplayMode int

// Nilai mode tampilan display CSS.
const (
	DisplayDefault DisplayMode = iota // tidak didefinisikan / normal flow
	DisplayBlock
	DisplayInlineBlock
	DisplayFlex
	DisplayInlineFlex
	DisplayGrid
	DisplayInlineGrid
	DisplayNone
)

// IsFlex memeriksa apakah display mode merupakan flex atau inline-flex.
func (d DisplayMode) IsFlex() bool {
	return d == DisplayFlex || d == DisplayInlineFlex
}

// FlexDirection merepresentasikan arah flexbox pada suatu tier.
type FlexDirection int

// Nilai arah aliran flexbox (row vs column).
const (
	FlexDirectionRow FlexDirection = iota // default CSS flex-direction: row
	FlexDirectionRowReverse
	FlexDirectionColumn
	FlexDirectionColumnReverse
)

// IsHorizontal memeriksa apakah arah flexbox mengalir secara horizontal (row / row-reverse).
func (d FlexDirection) IsHorizontal() bool {
	return d == FlexDirectionRow || d == FlexDirectionRowReverse
}

// FlexWrap merepresentasikan perilaku pematahan baris flexbox pada suatu tier.
type FlexWrap int

// Nilai pematahan baris flexbox.
const (
	FlexWrapNoWrap FlexWrap = iota // default CSS flex-wrap: nowrap
	FlexWrapWrap
	FlexWrapReverse
)

// ShrinkState merepresentasikan kemampuan susut flex item pada suatu tier.
// Sesuai CSS Flexbox, default tanpa utilitas adalah ShrinkEnabled (flex-shrink: 1).
type ShrinkState int

// Status kemampuan penyusutan flex item.
const (
	ShrinkUnknown  ShrinkState = iota
	ShrinkEnabled              // flex-shrink: 1 (default flexbox atau explicit shrink)
	ShrinkDisabled             // flex-shrink: 0 (shrink-0 / flex-none)
)

// GapState menyimpan informasi jarak celah horizontal dan vertikal.
type GapState struct {
	HasHorizontal bool
	HorizontalVal string
	HasVertical   bool
	VerticalVal   string
}

// ContainerTierState merangkum status geometri sebuah container pada suatu tier.
type ContainerTierState struct {
	Display   DisplayMode
	Direction FlexDirection
	Wrap      FlexWrap
	Gap       GapState
	HasScroll bool // true jika memiliki overflow-x-auto / overflow-x-scroll
}

// ChildTierWidth menyimpan informasi deklarasi lebar anak pada suatu tier.
type ChildTierWidth struct {
	IsFractional bool
	Ratio        float64
	RawClass     string
	IsFlex1      bool
}

// ChildGeometry merangkum geometri dari sebuah elemen anak di seluruh tier.
type ChildGeometry struct {
	NodeID int
	Node   *ir.Node
	Width  [TierCount]ChildTierWidth
	Shrink [TierCount]ShrinkState
}

// ResolvedGeometry adalah hasil komputasi geometri lengkap untuk sebuah kontainer beserta anak-anaknya.
type ResolvedGeometry struct {
	Container [TierCount]ContainerTierState
	Children  []ChildGeometry
}

// parseBreakpointPrefix memecah utility class menjadi BreakpointTier dan base class.
func parseBreakpointPrefix(cls string) (BreakpointTier, string) {
	if strings.HasPrefix(cls, "sm:") {
		return TierSm, cls[3:]
	}
	if strings.HasPrefix(cls, "md:") {
		return TierMd, cls[3:]
	}
	if strings.HasPrefix(cls, "lg:") {
		return TierLg, cls[3:]
	}
	if strings.HasPrefix(cls, "xl:") {
		return TierXl, cls[3:]
	}
	if strings.HasPrefix(cls, "2xl:") {
		return Tier2Xl, cls[4:]
	}
	return TierBaseline, cls
}

// parseFractionalWidth mengekstrak nilai rasio numerik dari utility lebar Tailwind (misal: w-1/2 -> 0.5).
func parseFractionalWidth(base string) (float64, bool) {
	if !strings.HasPrefix(base, "w-") {
		return 0, false
	}
	val := base[2:]

	if idx := strings.IndexByte(val, '/'); idx != -1 {
		numStr := val[:idx]
		denStr := val[idx+1:]
		num, err1 := strconv.Atoi(numStr)
		den, err2 := strconv.Atoi(denStr)
		if err1 == nil && err2 == nil && den > 0 && num > 0 && num < den {
			return float64(num) / float64(den), true
		}
	}

	if strings.HasPrefix(val, "[") && strings.HasSuffix(val, "%]") {
		pctStr := val[1 : len(val)-2]
		pct, err := strconv.ParseFloat(pctStr, 64)
		if err == nil && pct > 0 && pct < 100 {
			return pct / 100.0, true
		}
	}

	return 0, false
}

const (
	precedenceGapShorthand = 1
	precedenceGapSpecific  = 2
)

type gapDeclaration struct {
	horizontal bool
	rawVal     string
	hasH       bool
	vertical   bool
	hasV       bool
	precedence int
}

func parseGapClass(base string) (gapDeclaration, bool) {
	if strings.HasPrefix(base, "gap-x-") {
		val := base[len("gap-x-"):]
		if val == "0" || val == "[0]" || val == "[0px]" {
			return gapDeclaration{horizontal: true, hasH: false, precedence: precedenceGapSpecific}, true
		}
		return gapDeclaration{horizontal: true, hasH: true, rawVal: base, precedence: precedenceGapSpecific}, true
	}
	if strings.HasPrefix(base, "gap-y-") {
		val := base[len("gap-y-"):]
		if val == "0" || val == "[0]" || val == "[0px]" {
			return gapDeclaration{vertical: true, hasV: false, precedence: precedenceGapSpecific}, true
		}
		return gapDeclaration{vertical: true, hasV: true, rawVal: base, precedence: precedenceGapSpecific}, true
	}
	if strings.HasPrefix(base, "gap-") {
		val := base[len("gap-"):]
		if val == "0" || val == "[0]" || val == "[0px]" {
			return gapDeclaration{horizontal: true, hasH: false, vertical: true, hasV: false, precedence: precedenceGapShorthand}, true
		}
		return gapDeclaration{horizontal: true, hasH: true, vertical: true, hasV: true, rawVal: base, precedence: precedenceGapShorthand}, true
	}
	return gapDeclaration{}, false
}

type rawContainerTier struct {
	displayDecl   DisplayMode
	displayPrec   int
	directionDecl FlexDirection
	directionPrec int
	wrapDecl      FlexWrap
	wrapPrec      int
	hGap          bool
	hGapVal       string
	hGapPrec      int
	vGap          bool
	vGapVal       string
	vGapPrec      int
	hasScrollDecl bool
	hasScrollPrec int
}

func applyContainerDisplayClass(raw *rawContainerTier, base string) {
	modes := map[string]DisplayMode{
		"flex":         DisplayFlex,
		"inline-flex":  DisplayInlineFlex,
		"block":        DisplayBlock,
		"inline-block": DisplayInlineBlock,
		"grid":         DisplayGrid,
		"inline-grid":  DisplayInlineGrid,
		"hidden":       DisplayNone,
	}
	if m, ok := modes[base]; ok && raw.displayPrec <= 1 {
		raw.displayDecl = m
		raw.displayPrec = 1
	}
}

func applyContainerDirectionClass(raw *rawContainerTier, base string) {
	dirs := map[string]FlexDirection{
		"flex-row":         FlexDirectionRow,
		"flex-row-reverse": FlexDirectionRowReverse,
		"flex-col":         FlexDirectionColumn,
		"flex-col-reverse": FlexDirectionColumnReverse,
	}
	if d, ok := dirs[base]; ok && raw.directionPrec <= 1 {
		raw.directionDecl = d
		raw.directionPrec = 1
	}
}

func applyContainerWrapClass(raw *rawContainerTier, base string) {
	wraps := map[string]FlexWrap{
		"flex-wrap":         FlexWrapWrap,
		"flex-wrap-reverse": FlexWrapReverse,
		"flex-nowrap":       FlexWrapNoWrap,
	}
	if w, ok := wraps[base]; ok && raw.wrapPrec <= 1 {
		raw.wrapDecl = w
		raw.wrapPrec = 1
	}
}

func applyContainerScrollClass(raw *rawContainerTier, base string) {
	switch base {
	case "overflow-x-auto", "overflow-x-scroll":
		raw.hasScrollDecl = true
		raw.hasScrollPrec = 1
	case "overflow-x-hidden", "overflow-x-visible":
		raw.hasScrollDecl = false
		raw.hasScrollPrec = 2
	}
}

func extractContainerTierDecls(classes []string) [TierCount]rawContainerTier {
	var rawTiers [TierCount]rawContainerTier
	for _, cls := range classes {
		tier, base := parseBreakpointPrefix(cls)
		idx := int(tier)

		applyContainerDisplayClass(&rawTiers[idx], base)
		applyContainerDirectionClass(&rawTiers[idx], base)
		applyContainerWrapClass(&rawTiers[idx], base)
		applyContainerScrollClass(&rawTiers[idx], base)

		if gDecl, ok := parseGapClass(base); ok {
			if gDecl.horizontal && gDecl.precedence >= rawTiers[idx].hGapPrec {
				rawTiers[idx].hGap = gDecl.hasH
				rawTiers[idx].hGapVal = cls
				rawTiers[idx].hGapPrec = gDecl.precedence
			}
			if gDecl.vertical && gDecl.precedence >= rawTiers[idx].vGapPrec {
				rawTiers[idx].vGap = gDecl.hasV
				rawTiers[idx].vGapVal = cls
				rawTiers[idx].vGapPrec = gDecl.precedence
			}
		}
	}
	return rawTiers
}

func cascadeContainerTiers(rawTiers [TierCount]rawContainerTier) [TierCount]ContainerTierState {
	var states [TierCount]ContainerTierState
	currDisplay := DisplayDefault
	currDirection := FlexDirectionRow
	currWrap := FlexWrapNoWrap
	currHGap := false
	currHGapVal := ""
	currVGap := false
	currVGapVal := ""
	currScroll := false

	for i := 0; i < TierCount; i++ {
		raw := rawTiers[i]
		if raw.displayPrec > 0 {
			currDisplay = raw.displayDecl
		}
		if raw.directionPrec > 0 {
			currDirection = raw.directionDecl
		}
		if raw.wrapPrec > 0 {
			currWrap = raw.wrapDecl
		}
		if raw.hGapPrec > 0 {
			currHGap = raw.hGap
			currHGapVal = raw.hGapVal
		}
		if raw.vGapPrec > 0 {
			currVGap = raw.vGap
			currVGapVal = raw.vGapVal
		}
		if raw.hasScrollPrec > 0 {
			currScroll = raw.hasScrollDecl
		}

		states[i] = ContainerTierState{
			Display:   currDisplay,
			Direction: currDirection,
			Wrap:      currWrap,
			Gap: GapState{
				HasHorizontal: currHGap,
				HorizontalVal: currHGapVal,
				HasVertical:   currVGap,
				VerticalVal:   currVGapVal,
			},
			HasScroll: currScroll,
		}
	}
	return states
}

type rawChildTier struct {
	widthDecl  ChildTierWidth
	widthPrec  int
	shrinkDecl ShrinkState
	shrinkPrec int
}

func applyChildWidthClass(raw *rawChildTier, cls string, base string) {
	if ratio, ok := parseFractionalWidth(base); ok {
		if raw.widthPrec <= 2 {
			raw.widthDecl = ChildTierWidth{IsFractional: true, Ratio: ratio, RawClass: cls}
			raw.widthPrec = 2
		}
	} else if base == "w-full" || base == "w-[100%]" {
		if raw.widthPrec <= 2 {
			raw.widthDecl = ChildTierWidth{IsFractional: false, Ratio: 1.0, RawClass: cls}
			raw.widthPrec = 2
		}
	} else if base == "flex-1" {
		if raw.widthPrec <= 3 {
			raw.widthDecl = ChildTierWidth{IsFractional: false, IsFlex1: true, RawClass: cls}
			raw.widthPrec = 3
		}
	}
}

func applyChildShrinkClass(raw *rawChildTier, base string) {
	switch base {
	case "shrink-0", "flex-shrink-0", "flex-none":
		if raw.shrinkPrec <= 2 {
			raw.shrinkDecl = ShrinkDisabled
			raw.shrinkPrec = 2
		}
	case "shrink", "flex-shrink":
		if raw.shrinkPrec <= 2 {
			raw.shrinkDecl = ShrinkEnabled
			raw.shrinkPrec = 2
		}
	}
}

func cascadeChildTiers(rawChildTiers [TierCount]rawChildTier) ([TierCount]ChildTierWidth, [TierCount]ShrinkState) {
	var widths [TierCount]ChildTierWidth
	var shrinks [TierCount]ShrinkState
	currWidth := ChildTierWidth{}
	currShrink := ShrinkEnabled

	for i := 0; i < TierCount; i++ {
		raw := rawChildTiers[i]
		if raw.widthPrec > 0 {
			currWidth = raw.widthDecl
		}
		if raw.shrinkPrec > 0 {
			currShrink = raw.shrinkDecl
		}
		widths[i] = currWidth
		shrinks[i] = currShrink
	}
	return widths, shrinks
}

func extractChildGeometry(child *ir.Node, childID int) ChildGeometry {
	var rawChildTiers [TierCount]rawChildTier

	for _, cls := range child.Classes {
		tier, base := parseBreakpointPrefix(cls)
		idx := int(tier)
		applyChildWidthClass(&rawChildTiers[idx], cls, base)
		applyChildShrinkClass(&rawChildTiers[idx], base)
	}

	widths, shrinks := cascadeChildTiers(rawChildTiers)
	return ChildGeometry{
		NodeID: childID,
		Node:   child,
		Width:  widths,
		Shrink: shrinks,
	}
}

// ResolveGeometry mengomputasi geometri tata letak kontainer dan anak-anaknya
// melalui cascade per-tier dan aturan spesifisitas utilitas.
func ResolveGeometry(container *ir.Node) *ResolvedGeometry {
	if container == nil || container.Type != ir.NodeElement {
		return nil
	}

	geom := &ResolvedGeometry{}
	rawTiers := extractContainerTierDecls(container.Classes)
	geom.Container = cascadeContainerTiers(rawTiers)

	childGeoms := make([]ChildGeometry, 0, len(container.Children))
	childID := 0
	for _, child := range container.Children {
		if child == nil || child.Type != ir.NodeElement {
			continue
		}
		childID++
		childGeoms = append(childGeoms, extractChildGeometry(child, childID))
	}

	geom.Children = childGeoms
	return geom
}

// GeometryDriftFinding merangkum bukti kegagalan tata letak responsif pada sebuah flex container.
type GeometryDriftFinding struct {
	Tier        BreakpointTier
	FractionSum float64
	GapVal      string
	Reason      string
}

func evaluateTierDrift(tier BreakpointTier, tierIdx int, geom *ResolvedGeometry) *GeometryDriftFinding {
	cState := geom.Container[tierIdx]
	if !cState.Display.IsFlex() || !cState.Direction.IsHorizontal() || !cState.Gap.HasHorizontal {
		return nil
	}

	var (
		fractionCount         int
		fractionSum           float64
		hasFractionalNoShrink bool
		shrink0Count          int
	)

	for _, child := range geom.Children {
		w := child.Width[tierIdx]
		s := child.Shrink[tierIdx]

		if w.IsFractional {
			fractionCount++
			fractionSum += w.Ratio
			if s == ShrinkDisabled {
				hasFractionalNoShrink = true
				shrink0Count++
			}
		}
	}

	if fractionCount < 2 || fractionSum < 1.0-0.001 {
		return nil
	}

	if cState.HasScroll && shrink0Count >= 2 && fractionCount >= 2 {
		return nil
	}

	if cState.Wrap == FlexWrapWrap || cState.Wrap == FlexWrapReverse {
		return &GeometryDriftFinding{
			Tier:        tier,
			FractionSum: fractionSum,
			GapVal:      cState.Gap.HorizontalVal,
			Reason:      "unexpected line wrapping: fractional widths total 100% or more with active horizontal gap under 'flex-wrap'",
		}
	}

	if hasFractionalNoShrink {
		return &GeometryDriftFinding{
			Tier:        tier,
			FractionSum: fractionSum,
			GapVal:      cState.Gap.HorizontalVal,
			Reason:      "container blowout: fractional widths total 100% or more with active horizontal gap while flex shrinking is disabled ('shrink-0' / 'flex-none')",
		}
	}

	return nil
}

// ClassifyFlexGapDrift mengevaluasi ResolvedGeometry untuk membuktikan apakah terjadi Flexbox Gap Drift.
func ClassifyFlexGapDrift(geom *ResolvedGeometry) *GeometryDriftFinding {
	if geom == nil {
		return nil
	}

	for tierIdx, tier := range allTiers {
		if finding := evaluateTierDrift(tier, tierIdx, geom); finding != nil {
			return finding
		}
	}

	return nil
}
