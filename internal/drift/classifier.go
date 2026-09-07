package drift

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// ResolveScope menentukan ScopeKey dari node AST berdasarkan tag dan atribut ARIA role.
func ResolveScope(node *ir.Node) (ScopeKey, bool) {
	if node == nil || node.Type != ir.NodeElement {
		return ScopeKey{}, false
	}

	tag := node.Tag

	// 1. Exact Components
	switch tag {
	case "Button", "TabsTrigger":
		return ScopeKey{Kind: KindButton, Confidence: ConfidenceExact}, true
	case "Input", "SelectTrigger":
		return ScopeKey{Kind: KindInput, Confidence: ConfidenceExact}, true
	case "Card":
		return ScopeKey{Kind: KindCard, Confidence: ConfidenceExact}, true
	case "Dialog", "DialogContent":
		return ScopeKey{Kind: KindDialog, Confidence: ConfidenceExact}, true
	case "Sheet", "SheetContent":
		return ScopeKey{Kind: KindSheet, Confidence: ConfidenceExact}, true
	case "Badge":
		return ScopeKey{Kind: KindBadge, Confidence: ConfidenceExact}, true
	case "PopoverContent", "DropdownMenuContent", "TooltipContent":
		return ScopeKey{Kind: KindMenu, Confidence: ConfidenceExact}, true
	case "Alert":
		return ScopeKey{Kind: KindAlert, Confidence: ConfidenceExact}, true
	}

	// 2. Native HTML Tags
	switch tag {
	case "button":
		return ScopeKey{Kind: KindButton, Confidence: ConfidenceNative}, true
	case "input", "textarea", "select":
		return ScopeKey{Kind: KindInput, Confidence: ConfidenceNative}, true
	}

	// 3. Semantic ARIA Roles
	if role, ok := node.GetAttr("role"); ok {
		switch role {
		case "button":
			return ScopeKey{Kind: KindButton, Confidence: ConfidenceSemantic}, true
		case "textbox":
			return ScopeKey{Kind: KindInput, Confidence: ConfidenceSemantic}, true
		case "dialog":
			return ScopeKey{Kind: KindDialog, Confidence: ConfidenceSemantic}, true
		case "status":
			return ScopeKey{Kind: KindBadge, Confidence: ConfidenceSemantic}, true
		case "alert":
			return ScopeKey{Kind: KindAlert, Confidence: ConfidenceSemantic}, true
		}
	}

	return ScopeKey{}, false
}

// ExtractOccurrences mengekstrak daftar StyleOccurrence dari sebuah node AST IR jika dikenali.
func ExtractOccurrences(filePath string, node *ir.Node) []StyleOccurrence {
	scope, ok := ResolveScope(node)
	if !ok || len(node.Classes) == 0 {
		return nil
	}

	// Salin props relevan untuk evaluasi pengecualian kontekstual
	props := make(map[string]string)
	if node.Attributes != nil {
		for k, v := range node.Attributes {
			if k != "class" && k != "className" {
				props[k] = v
			}
		}
	}

	occurrences := make([]StyleOccurrence, 0, len(node.Classes))

	for _, classToken := range node.Classes {
		cat, normToken, variant, matched := classifyClassToken(classToken)
		if !matched {
			continue
		}

		// Filter dimensional property jika relevan untuk scope
		if !isCategoryApplicable(scope.Kind, cat) {
			continue
		}

		occurrences = append(occurrences, StyleOccurrence{
			FilePath:   filePath,
			Span:       node.Span,
			Scope:      scope,
			Category:   cat,
			Token:      normToken,
			Variant:    variant,
			RawClasses: node.RawClasses,
			Tag:        node.Tag,
			Props:      props,
		})
	}

	return occurrences
}

// ExtractTreeOccurrences melakukan traversal pohon AST IR dan mengekstrak seluruh StyleOccurrence dari setiap node.
func ExtractTreeOccurrences(filePath string, root *ir.Node) []StyleOccurrence {
	if root == nil {
		return nil
	}
	var occurrences []StyleOccurrence
	for node := range root.Walk() {
		occs := ExtractOccurrences(filePath, node)
		if len(occs) > 0 {
			occurrences = append(occurrences, occs...)
		}
	}
	return occurrences
}

// isCategoryApplicable memfilter apakah kategori properti relevan untuk komponen tersebut.
func isCategoryApplicable(kind ComponentKind, cat PropertyCategory) bool {
	switch kind {
	case KindButton:
		return cat == CatRounded || cat == CatActiveScale || cat == CatDisabledOpacity ||
			cat == CatDisabledCursor || cat == CatFocusRingWidth || cat == CatFocusRingColor ||
			cat == CatFocusRingOffset || cat == CatShadow || cat == CatBorderWidth ||
			cat == CatBorderColor || cat == CatBackground || cat == CatHeight
	case KindInput:
		return cat == CatRounded || cat == CatBorderWidth || cat == CatBorderColor ||
			cat == CatBackground || cat == CatHeight || cat == CatFocusRingWidth || cat == CatFocusRingColor
	case KindCard, KindDialog, KindSheet:
		return cat == CatRounded || cat == CatShadow || cat == CatBorderWidth || cat == CatBorderColor || cat == CatBackground
	case KindBadge:
		return cat == CatRounded || cat == CatBackground || cat == CatBorderWidth || cat == CatBorderColor
	case KindMenu:
		return cat == CatRounded || cat == CatShadow || cat == CatBorderWidth || cat == CatBorderColor || cat == CatBackground
	case KindAlert:
		return cat == CatRounded || cat == CatBorderWidth || cat == CatBorderColor || cat == CatBackground
	default:
		return true
	}
}

func classifyRadius(base, variant string) (PropertyCategory, string, string, bool) {
	if base == "rounded" || strings.HasPrefix(base, "rounded-") {
		return CatRounded, base, variant, true
	}
	return CatUnknown, "", "", false
}

func classifyScale(base, variant string) (PropertyCategory, string, string, bool) {
	if (strings.HasPrefix(base, "scale-") || strings.HasPrefix(base, "-scale-")) &&
		(variant == "active" || strings.Contains(variant, "active")) {
		return CatActiveScale, base, variant, true
	}
	return CatUnknown, "", "", false
}

func classifyDisabled(base, variant string) (PropertyCategory, string, string, bool) {
	if !strings.Contains(variant, "disabled") {
		return CatUnknown, "", "", false
	}
	if strings.HasPrefix(base, "opacity-") {
		return CatDisabledOpacity, base, variant, true
	}
	if strings.HasPrefix(base, "cursor-") {
		return CatDisabledCursor, base, variant, true
	}
	return CatUnknown, "", "", false
}

func classifyFocus(base, variant string) (PropertyCategory, string, string, bool) {
	if !strings.Contains(variant, "focus") && !strings.Contains(variant, "focus-visible") {
		return CatUnknown, "", "", false
	}
	if base != "ring" && !strings.HasPrefix(base, "ring-") {
		return CatUnknown, "", "", false
	}
	if strings.HasPrefix(base, "ring-offset-") {
		return CatFocusRingOffset, base, variant, true
	}
	if isRingWidth(base) {
		return CatFocusRingWidth, base, variant, true
	}
	return CatFocusRingColor, base, variant, true
}

func classifySurfaces(base, variant string) (PropertyCategory, string, string, bool) {
	if (base == "shadow" || strings.HasPrefix(base, "shadow-")) && isElevationShadow(base) {
		return CatShadow, base, variant, true
	}
	if base == "border" || strings.HasPrefix(base, "border-") {
		if isBorderWidth(base) {
			return CatBorderWidth, base, variant, true
		}
		if isBorderColor(base) {
			return CatBorderColor, base, variant, true
		}
	}
	return CatUnknown, "", "", false
}

func classifyDimensions(base, variant string) (PropertyCategory, string, string, bool) {
	if variant != "" {
		return CatUnknown, "", "", false
	}
	if strings.HasPrefix(base, "bg-") {
		return CatBackground, base, variant, true
	}
	if strings.HasPrefix(base, "h-") {
		return CatHeight, base, variant, true
	}
	if strings.HasPrefix(base, "px-") {
		return CatPaddingX, base, variant, true
	}
	if strings.HasPrefix(base, "py-") {
		return CatPaddingY, base, variant, true
	}
	return CatUnknown, "", "", false
}

// classifyClassToken mengklasifikasi sebuah class token ke PropertyCategory, token yang dinormalisasi, dan varian.
func classifyClassToken(raw string) (PropertyCategory, string, string, bool) {
	variant, base := splitVariantAndBase(raw)

	if cat, tok, v, ok := classifyRadius(base, variant); ok {
		return cat, tok, v, true
	}
	if cat, tok, v, ok := classifyScale(base, variant); ok {
		return cat, tok, v, true
	}
	if cat, tok, v, ok := classifyDisabled(base, variant); ok {
		return cat, tok, v, true
	}
	if cat, tok, v, ok := classifyFocus(base, variant); ok {
		return cat, tok, v, true
	}
	if cat, tok, v, ok := classifySurfaces(base, variant); ok {
		return cat, tok, v, true
	}
	if cat, tok, v, ok := classifyDimensions(base, variant); ok {
		return cat, tok, v, true
	}

	return CatUnknown, "", "", false
}

func splitVariantAndBase(token string) (string, string) {
	colonIdx := strings.LastIndexByte(token, ':')
	if colonIdx == -1 {
		return "", token
	}
	return token[:colonIdx], token[colonIdx+1:]
}

func isRingWidth(base string) bool {
	switch base {
	case "ring", "ring-0", "ring-1", "ring-2", "ring-4", "ring-8", "ring-inset":
		return true
	default:
		return false
	}
}

func isElevationShadow(base string) bool {
	switch base {
	case "shadow", "shadow-none", "shadow-sm", "shadow-md", "shadow-lg", "shadow-xl", "shadow-2xl", "shadow-inner":
		return true
	default:
		return false
	}
}

func isBorderWidth(base string) bool {
	switch base {
	case "border", "border-0", "border-2", "border-4", "border-8",
		"border-t", "border-r", "border-b", "border-l",
		"border-t-0", "border-t-2", "border-t-4", "border-t-8",
		"border-b-0", "border-b-2", "border-b-4", "border-b-8",
		"border-x", "border-y", "border-x-0", "border-y-0", "border-x-2", "border-y-2":
		return true
	default:
		return false
	}
}

func isBorderColor(base string) bool {
	if isBorderWidth(base) {
		return false
	}
	// Hindari border-collapse, border-solid, border-dashed, border-dotted, border-double, border-none
	switch base {
	case "border-solid", "border-dashed", "border-dotted", "border-double", "border-none",
		"border-collapse", "border-separate":
		return false
	}
	return strings.HasPrefix(base, "border-")
}
