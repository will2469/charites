package ux

import (
	"fmt"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// LayoutOwnerKind mengklasifikasikan jenis kepemilikan hubungan tata letak spasial.
type LayoutOwnerKind int

// Konstanta jenis ownership tata letak.
const (
	LayoutOwnerNone           LayoutOwnerKind = iota // Tidak memiliki ownership tata letak
	LayoutOwnerContainerGap                          // flex/grid dengan utility gap-*
	LayoutOwnerContainerSpace                        // container dengan utility space-x-* atau space-y-*
	LayoutOwnerMarginSequence                        // block parent terverifikasi dengan sekuens margin sibling homogen tanpa spasi saingan
	LayoutOwnerPeerContainers                        // parent yang membawahi peer containers homogen dengan gap-*
)

// SpacingRole mengklasifikasikan peran struktural spasial.
type SpacingRole int

// Konstanta peran spasial.
const (
	SpacingRoleUnknown   SpacingRole = iota // Peran tidak diketahui
	SpacingRoleInline                       // Spasi inline (icon <-> teks)
	SpacingRoleComponent                    // Spasi level komponen (label <-> input)
	SpacingRoleGroup                        // Spasi level grup (field <-> field, card <-> card, peer <-> peer)
	SpacingRoleSection                      // Spasi level seksi makro (header <-> konten, seksi <-> seksi)
)

// SpacingAxis mengklasifikasikan sumbu spasi (horizontal vs vertikal).
type SpacingAxis int

// Konstanta sumbu spasi.
const (
	AxisUnknown    SpacingAxis = iota // Sumbu tidak diketahui
	AxisHorizontal                    // Sumbu horizontal (x)
	AxisVertical                      // Sumbu vertikal (y)
)

// SpacingSource mengklasifikasikan mekanisme CSS yang mendasari pengukuran spasi.
type SpacingSource int

// Konstanta mekanisme CSS spasi.
const (
	SpacingSourceUnknown SpacingSource = iota // Mekanisme tidak diketahui
	SpacingSourceGap                          // CSS gap (flex/grid)
	SpacingSourceSpace                        // Tailwind space-between sibling selector
	SpacingSourceMargin                       // CSS margin (mb/mt/mr/ml)
)

// ResponsiveState menyimpan varian breakpoint/media/container yang aktif.
type ResponsiveState struct {
	Raw string // "base", "sm", "md", "lg", "xl", "2xl", "max-sm", "@md", etc.
}

// SpatialRelationship merepresentasikan hubungan struktural spasial antara dua node atau container.
type SpatialRelationship struct {
	From   *ir.Node
	To     *ir.Node
	Parent *ir.Node
	Axis   SpacingAxis
	Role   SpacingRole
}

// SpacingOccurrence merepresentasikan satu titik observasi spasi.
type SpacingOccurrence struct {
	Node         *ir.Node // AST node yang menyandang sintaks / target diagnostik
	Relationship SpatialRelationship
	Value        float64
	RawToken     string
	Source       SpacingSource
	Responsive   ResponsiveState
}

// RhythmPartitionKey adalah kunci partisi lengkap agar spasi berbeda sumbu/mekanisme/peran tidak bercampur.
type RhythmPartitionKey struct {
	OwnerID    string
	Axis       SpacingAxis
	Source     SpacingSource
	Role       SpacingRole
	Responsive ResponsiveState
}

// RhythmGroup adalah kumpulan observasi spasi homogen yang dapat dibandingkan.
type RhythmGroup struct {
	Key         RhythmPartitionKey
	Occurrences []SpacingOccurrence
}

// RhythmConfig menyimpan ambang batas analisis ritme majority-based.
type RhythmConfig struct {
	MinOccurrences int     // Bawaan: 3
	MinDominance   float64 // Bawaan: 0.60 (60%)
}

// DefaultRhythmConfig mengembalikan konfigurasi default sesuai spesifikasi v1.
func DefaultRhythmConfig() RhythmConfig {
	return RhythmConfig{
		MinOccurrences: 3,
		MinDominance:   0.60,
	}
}

// isResponsiveVariant memeriksa apakah nama varian merepresentasikan kondisi media/viewport/container.
func isResponsiveVariant(v string) bool {
	switch v {
	case "sm", "md", "lg", "xl", "2xl":
		return true
	}
	if strings.HasPrefix(v, "max-") || strings.HasPrefix(v, "min-") || strings.HasPrefix(v, "@") {
		return true
	}
	return false
}

// extractResponsiveState memecah rantai varian Tailwind menggunakan StripVariants
// untuk mengekstrak modifier breakpoint/media/container yang aktif.
func extractResponsiveState(token string) (ResponsiveState, string) {
	variants, base := StripVariants(token)
	if len(variants) == 0 {
		return ResponsiveState{Raw: "base"}, base
	}

	for _, v := range variants {
		if isResponsiveVariant(v) {
			return ResponsiveState{Raw: v}, base
		}
	}

	return ResponsiveState{Raw: "base"}, base
}

// isBlockContainerTag mengecek apakah tag HTML/JSX merupakan block container yang valid untuk sekuens margin.
func isBlockContainerTag(tag string) bool {
	switch strings.ToLower(tag) {
	case "div", "form", "fieldset", "section", "ul", "ol", "nav", "article", "aside", "main":
		return true
	}
	if len(tag) > 0 && tag[0] >= 'A' && tag[0] <= 'Z' {
		return true
	}
	return false
}

// getElementChildren mengembalikan slice anak yang merupakan NodeElement.
func getElementChildren(node *ir.Node) []*ir.Node {
	var elements []*ir.Node
	for _, child := range node.Children {
		if child != nil && child.Type == ir.NodeElement {
			elements = append(elements, child)
		}
	}
	return elements
}

// hasCompetingContainerSpacing mengecek apakah node memiliki kelas gap atau space yang bersaing.
func hasCompetingContainerSpacing(node *ir.Node) bool {
	for _, cls := range node.Classes {
		base := StripVariantsOnlyBase(cls)
		if strings.HasPrefix(base, "gap-") ||
			strings.HasPrefix(base, "space-y-") ||
			strings.HasPrefix(base, "space-x-") {
			return true
		}
	}
	return false
}

// extractChildMarginToken mengekstrak token margin (misal: "mb-4") dari elemen child.
func extractChildMarginToken(child *ir.Node, prefix string) (string, float64, ResponsiveState, bool) {
	for _, cls := range child.Classes {
		resp, base := extractResponsiveState(cls)
		cleanBase := StripVariantsOnlyBase(base)
		if strings.HasPrefix(cleanBase, prefix) {
			valStr := cleanBase[len(prefix):]
			if val, ok := parseTailwindSpacingNumber(valStr); ok {
				return cleanBase, val, resp, true
			}
		}
	}
	return "", 0, ResponsiveState{}, false
}

// checkMarginSequenceProof memverifikasi 5 syarat pembuktian struktural sekuens margin.
func checkMarginSequenceProof(node *ir.Node, children []*ir.Node) (string, bool) {
	if !isBlockContainerTag(node.Tag) || hasCompetingContainerSpacing(node) || len(children) < 3 {
		return "", false
	}

	firstTag := children[0].Tag
	for _, c := range children {
		if c.Tag != firstTag {
			return "", false
		}
	}

	for _, prefix := range []string{"mb-", "mt-"} {
		contiguous, ok := verifyContiguousMargins(children, prefix)
		if ok && contiguous >= 3 {
			return prefix, true
		}
	}

	return "", false
}

// verifyContiguousMargins memeriksa bahwa margin tidak diselingi oleh elemen tanpa margin.
func verifyContiguousMargins(children []*ir.Node, prefix string) (int, bool) {
	firstIdx := -1
	lastIdx := -1
	count := 0

	for i, c := range children {
		if _, _, _, ok := extractChildMarginToken(c, prefix); ok {
			if firstIdx == -1 {
				firstIdx = i
			}
			lastIdx = i
			count++
		}
	}

	if count < 3 {
		return 0, false
	}
	if lastIdx-firstIdx+1 != count {
		return 0, false // Ada elemen terselip tanpa margin (non-adjacent)
	}

	return count, true
}

// extractContainerGapToken mengekstrak token gap pada suatu container child peer.
func extractContainerGapToken(node *ir.Node) (string, float64, SpacingAxis, ResponsiveState, bool) {
	isCol := isFlexColumn(node)
	for _, cls := range node.Classes {
		resp, base := extractResponsiveState(cls)
		cleanBase := StripVariantsOnlyBase(base)

		axis, valStr, ok := parseGapUtility(cleanBase, isCol)
		if !ok {
			continue
		}

		val, ok := parseTailwindSpacingNumber(valStr)
		if ok {
			return cleanBase, val, axis, resp, true
		}
	}
	return "", 0, AxisUnknown, ResponsiveState{}, false
}

// checkPeerContainersProof memeriksa apakah anak-anak direct merupakan peer containers yang masing-masing mendeklarasikan gap.
func checkPeerContainersProof(children []*ir.Node) bool {
	if len(children) < 3 {
		return false
	}

	var firstAxis SpacingAxis
	count := 0

	for _, c := range children {
		if c == nil || c.Type != ir.NodeElement {
			return false
		}
		_, _, axis, _, ok := extractContainerGapToken(c)
		if ok {
			if firstAxis == AxisUnknown {
				firstAxis = axis
			} else if firstAxis != axis {
				return false
			}
			count++
		}
	}

	return count >= 3 && count == len(children)
}

// ClassifyLayoutOwner menentukan jenis ownership tata letak suatu node.
func ClassifyLayoutOwner(node *ir.Node) LayoutOwnerKind {
	if node == nil || node.Type != ir.NodeElement {
		return LayoutOwnerNone
	}

	children := getElementChildren(node)

	// Cek apakah node membawahi sekelompok peer container dengan gap-*
	if checkPeerContainersProof(children) {
		return LayoutOwnerPeerContainers
	}

	hasGap := false
	hasSpace := false

	for _, cls := range node.Classes {
		base := StripVariantsOnlyBase(cls)
		if strings.HasPrefix(base, "gap-") {
			hasGap = true
		} else if strings.HasPrefix(base, "space-x-") || strings.HasPrefix(base, "space-y-") {
			hasSpace = true
		}
	}

	if hasGap {
		return LayoutOwnerContainerGap
	}
	if hasSpace {
		return LayoutOwnerContainerSpace
	}

	if _, ok := checkMarginSequenceProof(node, children); ok {
		return LayoutOwnerMarginSequence
	}

	return LayoutOwnerNone
}

func isMacroLandmarkTag(tag string) bool {
	switch strings.ToLower(tag) {
	case "header", "footer", "section", "article", "aside", "nav", "main", "h1", "h2", "h3":
		return true
	}
	switch tag {
	case "Header", "Footer", "Section", "Article", "Aside", "Nav", "Main", "SectionHeader", "ActionFooter":
		return true
	}
	return false
}

func isInputLikeTag(tag string) bool {
	switch strings.ToLower(tag) {
	case "input", "select", "textarea", "button":
		return true
	}
	switch tag {
	case "Input", "TextField", "Select", "Textarea", "Button":
		return true
	}
	return false
}

func isLabelTag(tag string) bool {
	return strings.EqualFold(tag, "label")
}

// ClassifyRelationship mengklasifikasikan peran spasial.
// Memprioritaskan keterbandingan struktural secara umum (SpacingRoleGroup),
// dan menggunakan semantik murni sebagai penghalusan (refinement).
func ClassifyRelationship(rel *SpatialRelationship) SpacingRole {
	if rel == nil || rel.From == nil {
		return SpacingRoleUnknown
	}

	tagFrom := rel.From.Tag
	tagTo := ""
	if rel.To != nil {
		tagTo = rel.To.Tag
	}

	// 1. Refinement: Komponen kontrol (label <-> input)
	if isLabelTag(tagFrom) && isInputLikeTag(tagTo) {
		return SpacingRoleComponent
	}

	// 2. Refinement: Batas makro struktural yang berbeda (header <-> footer, dll)
	if isMacroLandmarkTag(tagFrom) && isMacroLandmarkTag(tagTo) && tagFrom != tagTo {
		return SpacingRoleSection
	}

	// 3. Baseline: Semua sibling di bawah layout owner yang sama memiliki keterbandingan struktural
	return SpacingRoleGroup
}

// BuildRhythmGroups mengekstrak observasi spasi dan membaginya ke dalam partisi yang aman.
func BuildRhythmGroups(owner *ir.Node, kind LayoutOwnerKind) []RhythmGroup {
	children := getElementChildren(owner)
	if len(children) == 0 {
		return nil
	}

	ownerID := fmt.Sprintf("%p_%d_%d", owner, owner.Span.Line, owner.Span.Column)
	groupMap := make(map[string]*RhythmGroup)

	if kind == LayoutOwnerPeerContainers || checkPeerContainersProof(children) {
		buildPeerGapOccurrences(owner, children, ownerID, groupMap)
	}

	switch kind {
	case LayoutOwnerMarginSequence:
		prefix, ok := checkMarginSequenceProof(owner, children)
		if ok {
			buildMarginOccurrences(owner, children, prefix, ownerID, groupMap)
		}
	case LayoutOwnerContainerGap:
		buildContainerGapOccurrences(owner, children, ownerID, groupMap)
	case LayoutOwnerContainerSpace:
		buildContainerSpaceOccurrences(owner, children, ownerID, groupMap)
	}

	var result []RhythmGroup
	for _, g := range groupMap {
		if g != nil && len(g.Occurrences) > 0 {
			result = append(result, *g)
		}
	}
	return result
}

func buildPeerGapOccurrences(owner *ir.Node, children []*ir.Node, ownerID string, groupMap map[string]*RhythmGroup) {
	for i, c := range children {
		token, val, axis, resp, ok := extractContainerGapToken(c)
		if !ok {
			continue
		}

		var nextPeer *ir.Node
		if i+1 < len(children) {
			nextPeer = children[i+1]
		}

		rel := SpatialRelationship{
			From:   c,
			To:     nextPeer,
			Parent: owner,
			Axis:   axis,
			Role:   SpacingRoleGroup,
		}

		key := RhythmPartitionKey{
			OwnerID:    ownerID + "_peers",
			Axis:       axis,
			Source:     SpacingSourceGap,
			Role:       SpacingRoleGroup,
			Responsive: resp,
		}
		keyStr := fmt.Sprintf("%s_%d_%d_%d_%s", key.OwnerID, key.Axis, key.Source, key.Role, key.Responsive.Raw)

		if _, exists := groupMap[keyStr]; !exists {
			groupMap[keyStr] = &RhythmGroup{Key: key}
		}
		groupMap[keyStr].Occurrences = append(groupMap[keyStr].Occurrences, SpacingOccurrence{
			Node:         c,
			Relationship: rel,
			Value:        val,
			RawToken:     token,
			Source:       SpacingSourceGap,
			Responsive:   resp,
		})
	}
}

func buildMarginOccurrences(owner *ir.Node, children []*ir.Node, prefix string, ownerID string, groupMap map[string]*RhythmGroup) {
	// Occurrence margin hanya valid jika memiliki next sibling di dalam layout group (child terakhir diabaikan).
	for i := 0; i < len(children)-1; i++ {
		c := children[i]
		token, val, resp, ok := extractChildMarginToken(c, prefix)
		if !ok {
			continue
		}

		nextChild := children[i+1]
		rel := SpatialRelationship{
			From:   c,
			To:     nextChild,
			Parent: owner,
			Axis:   AxisVertical,
		}
		rel.Role = ClassifyRelationship(&rel)
		if rel.Role == SpacingRoleUnknown {
			continue
		}

		key := RhythmPartitionKey{
			OwnerID:    ownerID,
			Axis:       AxisVertical,
			Source:     SpacingSourceMargin,
			Role:       rel.Role,
			Responsive: resp,
		}
		keyStr := fmt.Sprintf("%s_%d_%d_%d_%s", key.OwnerID, key.Axis, key.Source, key.Role, key.Responsive.Raw)

		if _, exists := groupMap[keyStr]; !exists {
			groupMap[keyStr] = &RhythmGroup{Key: key}
		}
		groupMap[keyStr].Occurrences = append(groupMap[keyStr].Occurrences, SpacingOccurrence{
			Node:         c,
			Relationship: rel,
			Value:        val,
			RawToken:     token,
			Source:       SpacingSourceMargin,
			Responsive:   resp,
		})
	}
}

func parseGapUtility(cleanBase string, isCol bool) (SpacingAxis, string, bool) {
	switch {
	case strings.HasPrefix(cleanBase, "gap-y-"):
		return AxisVertical, cleanBase[len("gap-y-"):], true
	case strings.HasPrefix(cleanBase, "gap-x-"):
		return AxisHorizontal, cleanBase[len("gap-x-"):], true
	case strings.HasPrefix(cleanBase, "gap-"):
		if isCol {
			return AxisVertical, cleanBase[len("gap-"):], true
		}
		return AxisHorizontal, cleanBase[len("gap-"):], true
	default:
		return AxisUnknown, "", false
	}
}

func parseSpaceUtility(cleanBase string) (SpacingAxis, string, bool) {
	switch {
	case strings.HasPrefix(cleanBase, "space-y-") && !strings.HasPrefix(cleanBase, "space-y-reverse"):
		return AxisVertical, cleanBase[len("space-y-"):], true
	case strings.HasPrefix(cleanBase, "space-x-") && !strings.HasPrefix(cleanBase, "space-x-reverse"):
		return AxisHorizontal, cleanBase[len("space-x-"):], true
	default:
		return AxisUnknown, "", false
	}
}

func isFlexColumn(node *ir.Node) bool {
	for _, cls := range node.Classes {
		if StripVariantsOnlyBase(cls) == "flex-col" {
			return true
		}
	}
	return false
}

func appendPairOccurrences(children []*ir.Node, owner *ir.Node, ownerID string, axis SpacingAxis, source SpacingSource, val float64, token string, resp ResponsiveState, groupMap map[string]*RhythmGroup) {
	for i := 0; i < len(children)-1; i++ {
		rel := SpatialRelationship{
			From:   children[i],
			To:     children[i+1],
			Parent: owner,
			Axis:   axis,
		}
		rel.Role = ClassifyRelationship(&rel)
		if rel.Role == SpacingRoleUnknown {
			continue
		}

		key := RhythmPartitionKey{
			OwnerID:    ownerID,
			Axis:       axis,
			Source:     source,
			Role:       rel.Role,
			Responsive: resp,
		}
		keyStr := fmt.Sprintf("%s_%d_%d_%d_%s", key.OwnerID, key.Axis, key.Source, key.Role, key.Responsive.Raw)

		if _, exists := groupMap[keyStr]; !exists {
			groupMap[keyStr] = &RhythmGroup{Key: key}
		}
		groupMap[keyStr].Occurrences = append(groupMap[keyStr].Occurrences, SpacingOccurrence{
			Node:         children[i],
			Relationship: rel,
			Value:        val,
			RawToken:     token,
			Source:       source,
			Responsive:   resp,
		})
	}
}

func buildContainerGapOccurrences(owner *ir.Node, children []*ir.Node, ownerID string, groupMap map[string]*RhythmGroup) {
	if len(children) < 2 {
		return
	}

	isCol := isFlexColumn(owner)
	for _, cls := range owner.Classes {
		resp, base := extractResponsiveState(cls)
		cleanBase := StripVariantsOnlyBase(base)

		axis, valStr, ok := parseGapUtility(cleanBase, isCol)
		if !ok {
			continue
		}

		val, ok := parseTailwindSpacingNumber(valStr)
		if !ok {
			continue
		}

		appendPairOccurrences(children, owner, ownerID, axis, SpacingSourceGap, val, cleanBase, resp, groupMap)
	}
}

func buildContainerSpaceOccurrences(owner *ir.Node, children []*ir.Node, ownerID string, groupMap map[string]*RhythmGroup) {
	if len(children) < 2 {
		return
	}

	for _, cls := range owner.Classes {
		resp, base := extractResponsiveState(cls)
		cleanBase := StripVariantsOnlyBase(base)

		axis, valStr, ok := parseSpaceUtility(cleanBase)
		if !ok {
			continue
		}

		val, ok := parseTailwindSpacingNumber(valStr)
		if !ok {
			continue
		}

		appendPairOccurrences(children, owner, ownerID, axis, SpacingSourceSpace, val, cleanBase, resp, groupMap)
	}
}

// findDominantValue mencari nilai spasi dominan dan menghitung rasio dominansinya jika memenuhi syarat.
func findDominantValue(occurrences []SpacingOccurrence, cfg RhythmConfig) (float64, float64, bool) {
	if len(occurrences) < cfg.MinOccurrences {
		return 0, 0, false
	}

	counts := make(map[float64]int)
	for _, occ := range occurrences {
		counts[occ.Value]++
	}

	maxCount := 0
	var dominantVal float64
	for val, count := range counts {
		if count > maxCount {
			maxCount = count
			dominantVal = val
		}
	}

	tieCount := 0
	for _, count := range counts {
		if count == maxCount {
			tieCount++
		}
	}
	if tieCount > 1 {
		return 0, 0, false
	}

	total := len(occurrences)
	dominance := float64(maxCount) / float64(total)
	if dominance < cfg.MinDominance || maxCount < 2 {
		return 0, 0, false
	}

	return dominantVal, dominance, true
}

func evaluateRhythmGroup(g RhythmGroup, cfg RhythmConfig) []ir.Diagnostic {
	dominantVal, dominance, ok := findDominantValue(g.Occurrences, cfg)
	if !ok {
		return nil
	}

	var diags []ir.Diagnostic
	for _, occ := range g.Occurrences {
		if occ.Value != dominantVal {
			diags = append(diags, ir.Diagnostic{
				Line:     occ.Node.Span.Line,
				Column:   occ.Node.Span.Column,
				Rule:     "ux.spacing-rhythm-drift",
				Severity: ir.SeverityWarn,
				Message: fmt.Sprintf(
					"Spacing token '%s' (%.2g rem) deviates from local dominant rhythm (%.2g rem, %.0f%% dominance). Align with sibling cadence or document intentional visual boundary.",
					occ.RawToken, occ.Value, dominantVal, dominance*100,
				),
				Hint: "Align spacing with the dominant local family or use an explicit container/role boundary.",
			})
		}
	}
	return diags
}

// DetectRhythmDrift menganalisis setiap grup untuk mendeteksi outlier dari ritme dominan.
func DetectRhythmDrift(groups []RhythmGroup, cfg RhythmConfig) []ir.Diagnostic {
	var diags []ir.Diagnostic
	for _, g := range groups {
		if d := evaluateRhythmGroup(g, cfg); len(d) > 0 {
			diags = append(diags, d...)
		}
	}
	return diags
}
